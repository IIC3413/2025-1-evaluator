package internal

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const copyLimit = 2 << 24

func copyZipFile(fd *zip.File, target string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("failed to copy zip file: %w", err)
		}
	}()

	fp := strings.Split(fd.Name, string(os.PathSeparator))
	// Early exit if source dir.
	if len(fp) == 0 || len(fp) == 1 {
		return nil
	}
	name := filepath.Join(target, filepath.Join(fp[1:]...))
	if !strings.HasPrefix(name, filepath.Clean(target)) {
		return fmt.Errorf("malicious path found: %s", fd.Name)
	}

	if fd.Mode().IsDir() {
		if err = os.MkdirAll(name, dirPermission); err != nil {
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

func copyFile(src, target string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("failed to copy file: %w", err)
		}
	}()

	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	nf, err := os.Create(target)
	if err != nil {
		return err
	}
	defer nf.Close()

	_, err = io.Copy(nf, f)
	return err
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

func sum(vs []int) int {
	sum := 0
	for _, v := range vs {
		sum += v
	}
	return sum
}
