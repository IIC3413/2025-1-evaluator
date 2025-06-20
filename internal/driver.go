package internal

import (
	"archive/zip"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/schollz/progressbar/v3"
)

type Results struct {
	id  string
	pts []int
}

type Evaluator struct {
	ctx     *ExecContext
	results []Results
	logfile *os.File
	logger  *log.Logger
}

func NewEvaluator(ctx *ExecContext) (*Evaluator, error) {
	f, err := os.Create(filepath.Join(ioDir, resultDir, ctx.Label+".log"))
	if err != nil {
		return nil, fmt.Errorf("failed to setup logs: %w", err)
	}
	e := Evaluator{
		ctx:     ctx,
		results: make([]Results, len(ctx.Submissions)),
		logfile: f,
		logger:  log.New(f, "", log.LstdFlags),
	}
	return &e, nil
}

func (e *Evaluator) Eval() error {
	var err error
	// We run the tests before creating the output file so no submission can
	// search for it and make changes.
	prog := progressbar.Default(int64(len(e.ctx.Submissions)))
	for i := range len(e.ctx.Submissions) {
		if err = e.evalSubmission(i); err != nil {
			return err
		}
		if err = removeSubmission(); err != nil {
			return err
		}
		if err = prog.Add(1); err != nil {
			return err
		}
	}
	return e.writeResults()
}

func (e *Evaluator) FreeLogs() error {
	return e.logfile.Close()
}

func (e *Evaluator) evalSubmission(idx int) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf(
				"failed to eval %s: %w",
				filepath.Base(e.ctx.Submissions[idx]),
				err,
			)
		}
	}()

	e.addSubmissionEntry(idx)
	if err = e.unzipSubmission(idx); err != nil {
		return err
	}
	if err = e.insertDeps(); err != nil {
		return err
	}
	if err = e.ctx.Corrections.Apply(); err != nil {
		return err
	}
	if err = e.compileTests(); err != nil {
		// A compilation error means no points can be awarded.
		e.logger.Printf(
			"\nFailed to compile %s:\n%s====================================",
			filepath.Base(e.ctx.Submissions[idx]),
			err.Error(),
		)
		return nil
	}
	return e.runTests(idx)
}

func (e *Evaluator) addSubmissionEntry(idx int) {
	e.results[idx].id = e.ctx.Submissions[idx]
	e.results[idx].pts = make([]int, len(e.ctx.Tests))
}

func (e *Evaluator) unzipSubmission(idx int) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("failed to unzip submission: %w", err)
		}
	}()

	zr, err := zip.OpenReader(e.ctx.Submissions[idx])
	if err != nil {
		return err
	}
	defer zr.Close()

	target := filepath.Join(workingDir, srcDir)
	if err = os.Mkdir(target, 0o777); err != nil {
		return err
	}

	root, err := determineZipRoot(zr)
	if err != nil {
		root = ""
	}

	for _, fd := range zr.File {
		if err = copyZipFile(fd, root, target); err != nil {
			return err
		}
	}

	return nil
}

func (e *Evaluator) insertDeps() (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("failed to insert dependencies: %w", err)
		}
	}()

	outputs := filepath.Join(workingDir, outputsDir)
	if err = os.Mkdir(outputs, 0o777); err != nil {
		return err
	}
	// Hack: os.Mkdir perm argument are subject to applying umask so the
	// resulting directory is not writable by other users. We need
	// to explicitly change the permission so as to allow submission to write
	// to it.
	if err = os.Chmod(outputs, 0o777); err != nil {
		return err
	}
	if err = e.copyTests(); err != nil {
		return err
	}
	return e.copyData()
}

func (e *Evaluator) copyTests() (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("failed to copy tests: %w", err)
		}
	}()

	target := filepath.Join(workingDir, srcDir, testsDir)
	if err = os.RemoveAll(target); err != nil {
		return err
	}
	// Submissions should not be able to read the test directory.
	if err = os.Mkdir(target, 0o700); err != nil {
		return err
	}
	return os.CopyFS(
		target,
		os.DirFS(filepath.Join(ioDir, testsDir, e.ctx.Label)),
	)
}

func (e *Evaluator) copyData() (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("failed to copy databases: %w", err)
		}
	}()

	target := filepath.Join(workingDir, dataDir, dbDir)
	// Submissions should very much be able to read the database directory.
	if err = os.MkdirAll(target, 0o777); err != nil {
		return err
	}
	err = os.CopyFS(
		target,
		os.DirFS(filepath.Join(ioDir, dataDir, e.ctx.Label, dbDir)),
	)
	if err != nil {
		return err
	}
	return filepath.WalkDir(
		target,
		func(path string, _ fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			return os.Chmod(path, 0o777)
		},
	)
}

func (e *Evaluator) compileTests() (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("failed to compile: %w", err)
		}
	}()

	binPath := filepath.Join(workingDir, buildDir, e.ctx.Build)
	//nolint:gosec // no user provider paths.
	cmd1 := exec.Command(
		"cmake",
		"-B"+binPath,
		"-S"+workingDir,
		"-DCMAKE_BUILD_TYPE="+e.ctx.Build,
	)
	cmd2 := exec.Command("cmake", "--build", binPath, "-j", "8")
	if _, err = runCommand(cmd1); err != nil {
		return fmt.Errorf("cmake build tree: %w", err)
	}
	if _, err = runCommand(cmd2); err != nil {
		return fmt.Errorf("cmake build: %w", err)
	}
	// We change permission so submissions can not replace binaries.
	return os.Chmod(binPath, 0o701)
}

func (e *Evaluator) runTests(idx int) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("failed to run tests: %w", err)
		}
	}()

	binPath := filepath.Join(buildDir, e.ctx.Build, binDir)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for i, t := range e.ctx.Tests {
		//nolint:gosec // no user input used here.
		cmd := exec.CommandContext(ctx, filepath.Join(binPath, extlessBase(t)))
		if _, err = runCommandAsSubmission(cmd); err != nil {
			e.logger.Printf(
				"\nFailed to run test %d for submission %s:\n"+
					"%s================================",
				i,
				filepath.Base(e.ctx.Submissions[idx]),
				err.Error(),
			)
			if err = writeErrorOutput(t); err != nil {
				return err
			}
			continue // An error here means no points for this test.
		}
	}
	return e.computeScore(idx)
}

func (e *Evaluator) computeScore(idx int) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("failed to compute score: %w", err)
		}
	}()

	outputsPath := filepath.Join(ioDir, dataDir, e.ctx.Label, outputsDir)
	expected, err := os.ReadDir(outputsPath)
	if err != nil {
		return err
	}
	produced, err := os.ReadDir(filepath.Join(workingDir, outputsDir))
	if err != nil {
		return err
	}
	if len(expected) != len(produced) || len(produced) != len(e.ctx.Tests) {
		return errors.New("invalid number of expected & produced outputs")
	}

	var ok bool
	for i := range expected {
		ok, err = dirEntriesEq(
			DirEntry{
				expected[i],
				filepath.Join(outputsPath, expected[i].Name()),
			},
			DirEntry{
				produced[i],
				filepath.Join(workingDir, outputsDir, produced[i].Name()),
			},
		)
		if err != nil {
			return err
		}
		if ok {
			e.results[idx].pts[i] = 1
		}
	}

	return nil
}

func (e *Evaluator) writeResults() (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("failed to write results to csv: %w", err)
		}
	}()

	out := filepath.Join(ioDir, resultDir, e.ctx.Label) + ".csv"
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	if err = e.writeHeader(w); err != nil {
		return err
	}
	for i := range len(e.results) {
		if err = e.writeResult(w, i); err != nil {
			return err
		}
	}
	return nil
}

func (e *Evaluator) writeHeader(w *csv.Writer) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("failed to write results' header: %w", err)
		}
	}()
	header := make([]string, len(e.ctx.Tests)+2)
	header[0] = "submission"
	header[len(header)-1] = "total"
	for i, t := range e.ctx.Tests {
		header[i+1] = t
	}
	return w.Write(header)
}

func (e *Evaluator) writeResult(w *csv.Writer, idx int) (err error) {
	r := e.results[idx]

	defer func() {
		if err != nil {
			err = fmt.Errorf(
				"failed to write %s results: %w",
				extlessBase(r.id),
				err,
			)
		}
	}()

	line := make([]string, len(r.pts)+2)
	line[0] = extlessBase(r.id)
	line[len(line)-1] = strconv.Itoa(sum(r.pts))
	for i, p := range r.pts {
		line[i+1] = strconv.Itoa(p)
	}
	return w.Write(line)
}

func writeErrorOutput(t string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("failed to write error output: %w", err)
		}
	}()

	f, err := os.Create(filepath.Join(workingDir, outputsDir, t+".output"))
	if err != nil {
		return err
	}
	if _, err = f.WriteString("ERROR"); err != nil {
		return err
	}
	return f.Close()
}

func removeSubmission() error {
	src := filepath.Join(workingDir, srcDir)
	build := filepath.Join(workingDir, buildDir)
	data := filepath.Join(workingDir, dataDir)
	output := filepath.Join(workingDir, outputsDir)

	if err := os.RemoveAll(src); err != nil {
		return err
	}
	if err := os.RemoveAll(build); err != nil {
		return err
	}
	if err := os.RemoveAll(data); err != nil {
		return err
	}
	return os.RemoveAll(output)
}
