package internal

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const dirPermission = 0777 //nolint:gofumpt // it formatted is though.

type Results struct {
	id  string
	pts []int
}

func Run(ctx ExecContext) (err error) {
	// We run the tests before creating the output file so no submission can
	// search for it and make changes.
	rs := make([]Results, len(ctx.Submissions))
	for i, s := range ctx.Submissions {
		if rs[i], err = evalSubmission(s, ctx); err != nil {
			return err
		}
	}

	out := filepath.Join(ioDir, resultBaseDir, ctx.OutputName) + ".csv"
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	for _, r := range rs {
		if err = writeResult(w, r); err != nil {
			return err
		}
	}

	return nil
}

func evalSubmission(s string, ctx ExecContext) (r Results, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("failed evaluating %s: %w", filepath.Base(s), err)
		}
	}()

	r.id = extlessBase(s)
	if err = unzipSubmission(s); err != nil {
		return Results{}, err
	}
	if err = copyTests(ctx.TestTargetDir, ctx.Tests); err != nil {
		return Results{}, err
	}
	if err = compileTests(); err != nil {
		// A compilation error means no points can be awarded.
		log.Printf(
			"Failed compilation for submission %s: %s",
			r.id,
			err.Error(),
		)
		return Results{
			id:  r.id,
			pts: make([]int, len(ctx.Tests)),
		}, nil
	}

	if err = runTests(&r, ctx.Tests, ctx.VerCode); err != nil {
		return Results{}, err
	}
	return r, removeSubmission()
}

func unzipSubmission(s string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("failed to unzip submission: %w", err)
		}
	}()

	zr, err := zip.OpenReader(s)
	if err != nil {
		return err
	}
	defer zr.Close()

	target := filepath.Join(workingDir, srcDir)
	if err = os.Mkdir(target, dirPermission); err != nil {
		return err
	}

	for _, fd := range zr.File {
		if err = copyZipFile(fd, target); err != nil {
			return err
		}
	}

	return nil
}

func copyTests(tt string, tests []string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("failed to copy tests: %w", err)
		}
	}()

	target := filepath.Join(workingDir, srcDir, tt)
	if err = os.Mkdir(target, dirPermission); err != nil {
		return err
	}

	for _, t := range tests {
		err = copyFile(t, filepath.Join(target, filepath.Base(t)))
		if err != nil {
			return err
		}
	}
	return nil
}

func compileTests() (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("failed to compile: %w", err)
		}
	}()

	//nolint:gosec // no user provider paths.
	cmd1 := exec.Command(
		"cmake",
		"-B"+filepath.Join(workingDir, buildDir, releaseDir),
		"-S"+workingDir,
		"-DCMAKE_BUILD_TYPE=Release",
	)
	//nolint:gosec // no user provider paths.
	cmd2 := exec.Command(
		"cmake",
		"--build",
		filepath.Join(workingDir, buildDir, releaseDir),
		"-j",
		"8",
	)
	if _, err = runCommand(cmd1); err != nil {
		return fmt.Errorf("cmake build tree: %w", err)
	}
	if _, err = runCommand(cmd2); err != nil {
		return fmt.Errorf("cmake build: %w", err)
	}
	return nil
}

func runTests(r *Results, tests []string, vc string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("failed to run tests: %w", err)
		}
	}()

	r.pts = make([]int, len(tests))
	binPath := filepath.Join(buildDir, releaseDir, binDir)

	for i, t := range tests {
		//nolint:gosec // no user input used here.
		cmd := exec.Command(filepath.Join(workingDir, binPath, extlessBase(t)))
		//nolint:govet // no problem shadowing err here.
		result, err := runCommand(cmd)
		if err != nil {
			log.Printf(
				"Failed to run test %d for submission %s: %s",
				i,
				r.id,
				err.Error(),
			)
			continue // An error here means no points for this test.
		}
		rp, err := parseResult(result, vc)
		if err != nil {
			return err
		}
		r.pts[i] = rp
	}
	return nil
}

func parseResult(result string, vc string) (int, error) {
	if !strings.HasPrefix(result, vc) {
		return 0, fmt.Errorf("invalid result: %s", result)
	}

	info := strings.Replace(result, vc, "", 1)
	info = strings.ReplaceAll(info, " ", "")
	info = strings.ReplaceAll(info, "\n", "")

	pts, err := strconv.Atoi(info)
	if err != nil {
		return 0, fmt.Errorf("invalid result: %w", err)
	}

	return pts, nil
}

func writeResult(w *csv.Writer, r Results) error {
	line := make([]string, len(r.pts)+2)
	line[0] = r.id
	line[len(line)-1] = strconv.Itoa(sum(r.pts))
	for i, p := range r.pts {
		line[i+1] = strconv.Itoa(p)
	}
	return w.Write(line)
}

func removeSubmission() error {
	src := filepath.Join(workingDir, srcDir)
	data := filepath.Join(workingDir, dataDir)

	if err := os.RemoveAll(src); err != nil {
		return err
	}
	if err := os.RemoveAll(data); err != nil {
		return err
	}
	if err := os.Mkdir(data, dirPermission); err != nil {
		return err
	}
	return nil
}
