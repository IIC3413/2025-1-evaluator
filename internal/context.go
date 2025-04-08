package internal

import (
	"crypto/rand"
	_ "embed"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/goccy/go-yaml"
)

const (
	ioDir              = "io"
	testsBaseDir       = "tests"
	submissionsBaseDir = "submissions"
	resultBaseDir      = "results"
	workingDir         = "wkdir"
	srcDir             = "src"
	cMakeFile          = "CMakeLists.txt"
	buildDir           = "build"
	releaseDir         = "Release"
	binDir             = "bin"
	dataDir            = "data"
	dbDir              = "dbs"
)

//go:embed data/CMakeLists.txt
var cmake []byte

type Config struct {
	OutputName    string `yaml:"output_name"`
	SumbissionDir string `yaml:"submission"`
	TestDir       string `yaml:"test"`
	VerCode       string `yaml:"verification_code"`
}

func OpenConfig(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	var config Config
	err = yaml.NewDecoder(f).Decode(&config)
	return config, err
}

type ExecContext struct {
	Config
	DBs           []string
	Submissions   []string
	Tests         []string
	TestTargetDir string
	VerOutputRGX  *regexp.Regexp
}

func SetUpContext(conf Config) (ExecContext, error) {
	subs, err := getDirFiles(
		filepath.Join(ioDir, submissionsBaseDir, conf.SumbissionDir),
	)
	if err != nil {
		return ExecContext{}, err
	}

	testsPath := filepath.Join(ioDir, testsBaseDir, conf.TestDir)
	tests, err := getDirFiles(testsPath, suffixFilter(".cc"))
	if err != nil {
		return ExecContext{}, err
	}
	dbs, err := getDirFiles(testsPath, suffixFilter(".db"))
	if err != nil {
		return ExecContext{}, err
	}

	ttd := rand.Text()
	if err = writeCMakeTargets(tests, ttd); err != nil {
		return ExecContext{}, err
	}

	vor, err := regexp.Compile(regexp.QuoteMeta(conf.VerCode) + `(\d+)\n`)
	if err != nil {
		return ExecContext{}, err
	}

	return ExecContext{conf, dbs, subs, tests, ttd, vor}, nil
}

func suffixFilter(sfx string) func(s string) bool {
	return func(s string) bool {
		return strings.HasSuffix(s, sfx)
	}
}

func getDirFiles(path string, filters ...func(string) bool) ([]string, error) {
	ents, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	fents := []string{}
	for _, e := range ents {
		if e.IsDir() {
			continue
		}
		for _, fl := range filters {
			if !fl(e.Name()) {
				continue
			}
		}
		fents = append(fents, filepath.Join(path, e.Name()))
	}
	return fents, nil
}

func writeCMakeTargets(tests []string, testsTarget string) error {
	testNames, err := buildTestNames(tests)
	if err != nil {
		return err
	}
	o := regexp.MustCompile(`TEST_REPLACE`).ReplaceAll(cmake, []byte(testNames))
	o = regexp.MustCompile(`TARGET_REPLACE`).ReplaceAll(o, []byte(testsTarget))
	f, err := os.Create(filepath.Join(workingDir, cMakeFile))
	if err != nil {
		return err
	}
	if _, err = f.Write(o); err != nil {
		return err
	}
	return f.Close()
}

func buildTestNames(tests []string) (string, error) {
	var (
		err error
		b   strings.Builder
	)
	for _, t := range tests {
		if err = b.WriteByte('\t'); err != nil {
			return "", err
		}
		if _, err = b.WriteString(extlessBase(t)); err != nil {
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
