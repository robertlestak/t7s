package t7s

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func (t *T7s) replaceVariables(f io.Reader, w io.Writer) error {
	l := log.WithFields(log.Fields{
		"fn": "replaceVariables",
	})
	l.Debug("Starting")
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		for _, v := range t.Variables {
			if strings.Contains(line, t.LeftDelim+v.Name+t.RightDelim) {
				if v.Required != nil && *v.Required && v.Value == "" {
					return fmt.Errorf("required variable %s is empty", v.Name)
				}
				line = strings.Replace(line, t.LeftDelim+v.Name+t.RightDelim, os.ExpandEnv(v.Value), -1)
			}
		}
		fmt.Fprintln(w, line)
		if strings.Contains(line, t.LeftDelim) && strings.Contains(line, t.RightDelim) {
			varName := strings.Split(line, t.LeftDelim)[1]
			varName = strings.Split(varName, t.RightDelim)[0]
			if t.Require {
				return fmt.Errorf("required variable %s is empty", varName)
			} else {
				l.WithFields(log.Fields{
					"variable": varName,
					"line":     line,
				}).Warn("Line contains unprocessed variables")
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (t *T7s) processReplaceFile(f io.Reader, w io.Writer) error {
	l := log.WithFields(log.Fields{
		"fn": "processReplaceFile",
	})
	l.Debug("Starting")
	if err := t.replaceVariables(f, w); err != nil {
		return err
	}
	l.Debug("Done")
	return nil
}
