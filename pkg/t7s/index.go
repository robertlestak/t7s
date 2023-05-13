package t7s

import (
	"bufio"
	"io"
	"strings"

	log "github.com/sirupsen/logrus"
)

func findVarNames(line string, left string, right string) []string {
	l := log.WithFields(log.Fields{
		"fn": "findVarNames",
	})
	l.Debug("Starting")
	varNames := []string{}
	for {
		if !strings.Contains(line, left) {
			break
		}
		start := strings.Index(line, left)
		end := strings.Index(line, right)
		varName := line[start+len(left) : end]
		varNames = append(varNames, varName)
		line = line[end+len(right):]
	}
	l.WithFields(log.Fields{
		"varNames": varNames,
	}).Debug("Done")
	return varNames
}

func (t *T7s) indexVariables(f io.Reader, w io.Writer) error {
	l := log.WithFields(log.Fields{
		"fn": "indexVariables",
	})
	l.Debug("Starting")
	localVars := []Variable{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, t.LeftDelim) && strings.Contains(line, t.RightDelim) {
			varNames := findVarNames(line, t.LeftDelim, t.RightDelim)
			for _, varName := range varNames {
				// if variable is not already in localVars, add it
				var found bool
				for _, v := range localVars {
					if v.Name == varName {
						found = true
						break
					}
				}
				if !found {
					localVars = append(localVars, Variable{Name: varName})
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	l.WithFields(log.Fields{
		"localVars": len(localVars),
	}).Debug("Done")
	// for each variable in localVars, check if it is in t.Variables
	// if it is not, add it
	for _, v := range localVars {
		var found bool
		for _, tv := range t.Variables {
			if v.Name == tv.Name {
				found = true
				break
			}
		}
		if !found {
			t.Variables = append(t.Variables, v)
		}
	}
	l.WithFields(log.Fields{
		"t.Variables": len(t.Variables),
	}).Debug("Done")
	vc := VariablesCfg{
		Variables: t.Variables,
	}
	l.WithFields(log.Fields{
		"vc": vc,
	}).Debug("Done")
	return nil
}

func (t *T7s) processIndexFile(f io.Reader, w io.Writer) error {
	l := log.WithFields(log.Fields{
		"fn": "processIndexFile",
	})
	l.Debug("Starting")
	if err := t.indexVariables(f, w); err != nil {
		return err
	}
	l.Debug("Done")
	return nil
}
