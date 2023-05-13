package t7s

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Variable struct {
	Name        string  `yaml:"name" json:"name"`
	Value       string  `yaml:"value" json:"value"`
	Description *string `yaml:"description,omitempty" json:"description,omitempty"`
	Required    *bool   `yaml:"required,omitempty" json:"required,omitempty"`
}

type VariablesCfg struct {
	Variables []Variable `yaml:"variables"`
}

type T7s struct {
	LeftDelim  string
	RightDelim string
	InPath     string
	OutPath    string
	VarFile    string
	IndexMerge bool
	Variables  []Variable
	JobType    JobType
	Require    bool
}

type JobType string

const (
	JobTypeReplace JobType = "replace"
	JobTypeIndex   JobType = "index"
)

type ProcessJob struct {
	Type  JobType
	In    string
	Out   string
	Cfg   *T7s
	Error error
}

func parseInPath(in string) (io.Reader, error) {
	l := log.WithFields(log.Fields{
		"fn": "parseInPath",
	})
	l.Debug("Starting")
	if in == "-" {
		l.Debug("Reading from stdin")
		return os.Stdin, nil
	}
	l.Debug("Reading from file " + in)
	// open file
	f, err := os.Open(in)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func parseOutPath(out string) (io.Writer, error) {
	l := log.WithFields(log.Fields{
		"fn": "parseOutPath",
	})
	l.Debug("Starting")
	if out == "-" {
		l.Debug("Writing to stdout")
		return os.Stdout, nil
	}
	l.Debug("Writing to file " + out)
	// ensure parent dir exists
	if err := os.MkdirAll(filepath.Dir(out), 0755); err != nil {
		return nil, err
	}
	// open file
	f, err := os.Create(out)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (t *T7s) LoadVariables() error {
	l := log.WithFields(log.Fields{
		"fn": "LoadVariables",
	})
	l.Debug("Starting")
	if t.VarFile == "" {
		return nil
	}
	if t.JobType == JobTypeIndex && !t.IndexMerge {
		return nil
	}
	l.Debug("Loading variables from " + t.VarFile)
	_, err := os.Stat(t.VarFile)
	if err != nil && !os.IsNotExist(err) {
		l.WithError(err).Error("Failed to stat file")
		return err
	} else if os.IsNotExist(err) {
		l.Debug("File does not exist")
		return nil
	}
	f, err := os.Open(t.VarFile)
	if err != nil {
		l.WithError(err).Error("Failed to open file")
		return err
	}
	defer f.Close()
	vc := VariablesCfg{}
	// parse yaml, load into t.Variables
	if err := yaml.NewDecoder(f).Decode(&vc); err != nil {
		// try json
		if err := json.NewDecoder(f).Decode(&vc); err != nil {
			l.WithError(err).Error("Failed to decode yaml")
			return err
		}
	}
	t.Variables = vc.Variables
	l.Debugf("Loaded %d variables", len(t.Variables))
	return nil
}

func (t *T7s) OutputVariables() error {
	l := log.WithFields(log.Fields{
		"fn": "OutputVariables",
	})
	l.Debug("Starting")
	of, err := parseOutPath(t.OutPath)
	if err != nil {
		l.WithError(err).Error("Failed to parse out path")
		return err
	}
	defer of.(io.WriteCloser).Close()
	// output yaml
	vc := VariablesCfg{
		Variables: t.Variables,
	}
	if err := yaml.NewEncoder(of).Encode(vc); err != nil {
		l.WithError(err).Error("Failed to encode yaml")
		return err
	}
	return nil
}

func (t *T7s) Run() error {
	l := log.WithFields(log.Fields{
		"fn": "Run",
	})
	l.Debug("Starting")
	if err := t.LoadVariables(); err != nil {
		return err
	}
	var isDir bool
	if t.InPath == "-" {
		isDir = false
	} else {
		instat, err := os.Stat(t.InPath)
		if err != nil {
			return err
		}
		isDir = instat.IsDir()
	}
	if isDir {
		return t.processDir()
	} else {
		return t.processSingle()
	}
}
