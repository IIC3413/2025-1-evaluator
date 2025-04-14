package internal

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	// Directories.
	ioDir          = "io-cpy"
	testsDir       = "tests"
	submissionsDir = "submissions"
	resultDir      = "results"
	workingDir     = "wkdir"
	srcDir         = "src"
	cMakeFile      = "CMakeLists.txt"
	buildDir       = "build"
	releaseDir     = "Release"
	binDir         = "bin"
	dataDir        = "data"
	dbDir          = "eval_dbs"
	outputsDir     = "outputs"
	// User and group for running tests on submissions.
	suid uint32 = 1001
	sgid uint32 = 1001
)

//go:embed data/CMakeLists.txt
var cmake []byte

type ExecContext struct {
	Submissions []string
	Tests       []string
	TestsPath   string
	InputsPath  string
	OutputsPath string
}

func SetUpContext(lab string) (*ExecContext, error) {
	subs, err := getDirFiles(
		filepath.Join(ioDir, submissionsDir, lab),
	)
	if err != nil {
		return nil, err
	}
	tests, testsPath, err := getTests(lab)
	if err != nil {
		return nil, err
	}
	inputsPath := filepath.Join(ioDir, dataDir, lab, dbDir)
	outputsPath := filepath.Join(ioDir, dataDir, lab, outputsDir)
	if err = writeCMakeTargets(tests); err != nil {
		return nil, err
	}

	return &ExecContext{subs, tests, testsPath, inputsPath, outputsPath}, nil
}

func getTests(lab string) ([]string, string, error) {
	testsDirPath := filepath.Join(ioDir, testsDir, lab)
	testCompTargets, err := getDirFiles(
		testsDirPath,
		func(s string) bool { return !strings.HasSuffix(s, ".h") },
	)
	if err != nil {
		return nil, "", err
	}
	for i := range testCompTargets {
		testCompTargets[i] = extlessBase(testCompTargets[i])
	}
	return testCompTargets, testsDirPath, nil
}

func getDirFiles(path string, filters ...func(string) bool) ([]string, error) {
	ents, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	fents := []string{}
ENTS_LOOP:
	for _, e := range ents {
		if e.IsDir() {
			continue
		}
		for _, fl := range filters {
			if !fl(e.Name()) {
				continue ENTS_LOOP
			}
		}
		fents = append(fents, filepath.Join(path, e.Name()))
	}
	return fents, nil
}

func writeCMakeTargets(targets []string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("failed to write cmake file: %w", err)
		}
	}()

	targetNames, err := buildTargetNames(targets)
	if err != nil {
		return err
	}
	o := regexp.MustCompile(`TEST_REPLACE`).ReplaceAll(
		cmake,
		[]byte(targetNames),
	)
	f, err := os.Create(filepath.Join(workingDir, cMakeFile))
	if err != nil {
		return err
	}
	if _, err = f.Write(o); err != nil {
		return err
	}
	return f.Close()
}

func buildTargetNames(targets []string) (string, error) {
	var (
		err error
		b   strings.Builder
	)
	for _, t := range targets {
		if err = b.WriteByte('\t'); err != nil {
			return "", err
		}
		if _, err = b.WriteString(t); err != nil {
			return "", err
		}
		if err = b.WriteByte('\n'); err != nil {
			return "", err
		}
	}
	return b.String(), nil
}

func extlessBase(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	return strings.ReplaceAll(base, ext, "")
}
