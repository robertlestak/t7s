package t7s

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func (t *T7s) processDir() error {
	l := log.WithFields(log.Fields{
		"fn": "processDir",
	})
	l.Debug("Starting")
	// walk dir
	// for each file
	//   process file
	jobs := make(chan ProcessJob, 1)
	results := make(chan ProcessJob, 1)
	for w := 1; w <= 10; w++ {
		go processWorker(jobs, results)
	}
	var totalCount int
	err := filepath.Walk(t.InPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			l.WithFields(log.Fields{
				"path": path,
			}).Debug("Skipping directory")
			return nil
		}
		// filepath relative to inpath
		relpath, err := filepath.Rel(t.InPath, path)
		if err != nil {
			return err
		}
		var outPath string
		if t.OutPath == "-" {
			outPath = "-"
		} else if t.JobType == JobTypeIndex {
			outPath = t.OutPath
		} else {
			outPath = filepath.Join(t.OutPath, relpath)
		}
		jobs <- ProcessJob{
			In:   path,
			Out:  outPath,
			Cfg:  t,
			Type: t.JobType,
		}
		totalCount++
		return nil
	})
	if err != nil {
		return err
	}
	close(jobs)
	for i := 0; i < totalCount; i++ {
		r := <-results
		if r.Error != nil {
			return r.Error
		}
	}
	if t.JobType == JobTypeIndex {
		if err := t.OutputVariables(); err != nil {
			return err
		}
	}
	return nil
}

func (t *T7s) processSingle() error {
	l := log.WithFields(log.Fields{
		"fn": "processSingle",
	})
	l.Debug("Starting")
	jobs := make(chan ProcessJob, 1)
	results := make(chan ProcessJob, 1)
	go processWorker(jobs, results)
	jobs <- ProcessJob{
		In:   t.InPath,
		Out:  t.OutPath,
		Cfg:  t,
		Type: t.JobType,
	}
	close(jobs)
	r := <-results
	if r.Error != nil {
		return r.Error
	}
	if t.JobType == JobTypeIndex {
		if err := t.OutputVariables(); err != nil {
			return err
		}
	}
	return nil
}

func processWorker(jobs <-chan ProcessJob, results chan<- ProcessJob) {
	l := log.WithFields(log.Fields{
		"fn": "processWorker",
	})
	l.Debug("Starting")
	for j := range jobs {
		l.WithFields(log.Fields{
			"in":  j.In,
			"out": j.Out,
		}).Debug("Processing")
		f, err := parseInPath(j.In)
		if err != nil {
			l.Debug("Error parsing in path")
			j.Error = err
			results <- j
			continue
		}
		w, err := parseOutPath(j.Out)
		if err != nil {
			l.Debug("Error parsing out path")
			j.Error = err
			results <- j
			continue
		}
		l.Debug("Processing file")
		switch j.Type {
		case JobTypeReplace:
			if err := j.Cfg.processReplaceFile(f, w); err != nil {
				l.Debug("Error processing file")
				j.Error = err
				results <- j
				continue
			}
		case JobTypeIndex:
			if err := j.Cfg.processIndexFile(f, w); err != nil {
				l.Debug("Error processing file")
				j.Error = err
				results <- j
				continue
			}
		}
		l.Debug("Done processing file")
		results <- j
	}
	l.Debug("Done")
}
