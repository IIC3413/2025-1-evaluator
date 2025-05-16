package internal

import (
	"archive/zip"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

const copyLimit = 2 << 24

var ErrNoRoot = errors.New("zip file has not root")

func copyZipFile(fd *zip.File, root, target string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("failed to copy zip file: %w", err)
		}
	}()

	// Fuck. you. tim. cook.
	if strings.Contains(fd.Name, ".DS_Store") ||
		strings.Contains(fd.Name, "._") {
		return nil
	}

	fp := strings.Split(sourcefy(fd.Name, root), string(os.PathSeparator))
	// Early exit if source dir.
	if len(fp) == 0 || len(fp) == 1 {
		return nil
	}

	name := filepath.Join(target, filepath.Join(fp[1:]...))
	if !strings.HasPrefix(name, filepath.Clean(target)) {
		return fmt.Errorf("malicious path found: %s", fd.Name)
	}

	if fd.Mode().IsDir() {
		if err = os.MkdirAll(name, 0o777); err != nil {
			return err
		}
		return nil
	}

	f, err := fd.Open()
	if err != nil {
		return err
	}
	defer f.Close()
	nf, err := os.Create(name)
	if err != nil {
		return err
	}
	defer nf.Close()

	if _, err = io.CopyN(nf, f, copyLimit); err != nil {
		if !errors.Is(err, io.EOF) {
			return err
		}
	}
	return nil
}

func determineZipRoot(zr *zip.ReadCloser) (string, error) {
	root := strings.SplitN(zr.File[0].Name, string(os.PathSeparator), 2)[0]
	for _, fe := range zr.File {
		split := strings.SplitN(fe.Name, string(os.PathSeparator), 2)
		if len(split) == 1 {
			root = split[0]
		}
		if root != split[0] {
			return "", ErrNoRoot
		}
	}
	return root, nil
}

func runCommand(cmd *exec.Cmd) (string, error) {
	var stderr, stdout bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	err := cmd.Run()

	//nolint:errorlint // trust bruh.
	if _, ok := err.(*exec.ExitError); ok {
		return "", errors.New(stderr.String())
	} else if err != nil {
		return "", err
	}
	return stdout.String(), nil
}

func runCommandAsSubmission(cmd *exec.Cmd) (string, error) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{Uid: suid, Gid: sgid},
	}
	cmd.Dir = workingDir
	return runCommand(cmd)
}

func sum(vs []int) int {
	sum := 0
	for _, v := range vs {
		sum += v
	}
	return sum
}

type DirEntry struct {
	os.DirEntry
	path string
}

func dirEntriesEq(d1, d2 DirEntry) (bool, error) {
	if d1.IsDir() != d2.IsDir() {
		return false, fmt.Errorf(
			"unable to compare two directories (%s, %s)",
			d1.Name(),
			d2.Name(),
		)
	}

	f1, err := os.Open(d1.path)
	if err != nil {
		return false, err
	}
	defer f1.Close()
	f2, err := os.Open(d2.path)
	if err != nil {
		return false, err
	}
	defer f2.Close()
	return compReaders(bufio.NewReader(f1), bufio.NewReader(f2))
}

func compReaders(r1, r2 io.Reader) (bool, error) {
	var (
		err1, err2 error
		rb1, rb2   int
		buf1, buf2 = make([]byte, 1024), make([]byte, 1024)
	)
	for {
		rb1, err1 = r1.Read(buf1)
		rb2, err2 = r2.Read(buf2)
		if err1 != nil || err2 != nil {
			if (rb1 != 0 && errors.Is(err2, io.EOF)) ||
				(rb2 != 0 && errors.Is(err1, io.EOF)) {
				return false, nil
			}
			if errors.Is(err1, io.EOF) && errors.Is(err2, io.EOF) {
				break
			}
			return false, errors.Join(err1, err2)
		}
		if rb1 != rb2 {
			return false, nil
		}
		if !bytes.Equal(buf1[:rb1], buf2[:rb2]) {
			return false, nil
		}
	}
	return rb1 == rb2 && bytes.Equal(buf1[:rb1], buf2[:rb2]), nil
}

func sourcefy(path, root string) string {
	if !strings.Contains(path, "src") {
		if root == "" {
			return filepath.Join("src", path)
		}
		split := strings.SplitN(path, string(os.PathSeparator), 2)
		if len(split) == 1 {
			return filepath.Join("src", split[0])
		}
		return filepath.Join("src", split[1])
	}
	split := strings.SplitN(path, "src", 2)
	if len(split) == 1 {
		return "src" + split[0]
	}
	return "src" + split[1]
}
