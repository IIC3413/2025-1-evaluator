package internal

import (
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type correction struct {
	Path       string `yaml:"~"`
	FileName   string `yaml:"file_name"`
	TargetPath string `yaml:"target_path"`
}

type corrections struct {
	Entries []correction `yaml:"entries"`
}

func (cs corrections) Apply() error {
	for _, c := range cs.Entries {
		if err := c.Apply(); err != nil {
			return err
		}
	}
	return nil
}

func (c correction) Apply() error {
	f, err := os.Open(c.Path)
	if err != nil {
		return err
	}
	tf, err := os.Create(filepath.Join(workingDir, c.TargetPath))
	if err != nil {
		return err
	}
	_, err = io.Copy(tf, f)
	return err
}

func loadCorrections(lab string) (corrections, error) {
	var cors corrections
	corsPath := filepath.Join(ioDir, correctionsDir, lab)
	f, err := os.Open(filepath.Join(corsPath, correctionYaml))
	if err != nil {
		if os.IsNotExist(err) {
			return corrections{Entries: []correction{}}, nil
		}
		return corrections{}, err
	}
	if err = yaml.NewDecoder(f).Decode(&cors); err != nil {
		return corrections{}, err
	}
	for i := range cors.Entries {
		cors.Entries[i].Path = filepath.Join(corsPath, cors.Entries[i].FileName)
	}
	return cors, nil
}
