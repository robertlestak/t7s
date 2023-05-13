package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/robertlestak/t7s/pkg/t7s"
	log "github.com/sirupsen/logrus"
)

var (
	Version  = "dev"
	t7sflags = flag.NewFlagSet("t7s", flag.ExitOnError)
)

func init() {
	l, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		l = log.InfoLevel
	}
	log.SetLevel(l)
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options] [in] [out]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Options:\n")
	t7sflags.PrintDefaults()
}

func main() {
	t7sflags = flag.NewFlagSet("t7s", flag.ExitOnError)
	logLevel := t7sflags.String("log", log.GetLevel().String(), "Log level")
	leftdelim := t7sflags.String("left", "{{", "Left delimiter")
	rightdelim := t7sflags.String("right", "}}", "Right delimiter")
	indexOp := t7sflags.Bool("i", false, "Index templates")
	indexMerge := t7sflags.Bool("m", false, "Index merge with existing variables")
	varfile := t7sflags.String("v", "variables.yaml", "Variables file")
	require := t7sflags.Bool("r", false, "Require all variables to be set")
	version := t7sflags.Bool("version", false, "Print version")
	t7sflags.Usage = usage
	t7sflags.Parse(os.Args[1:])
	l, err := log.ParseLevel(*logLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(l)
	if *version {
		log.Info("Version: " + Version)
		os.Exit(0)
	}
	var inPath string = "-"
	var outPath string = "-"
	if len(t7sflags.Args()) > 0 {
		inPath = t7sflags.Args()[0]
	}
	if len(t7sflags.Args()) > 1 {
		outPath = t7sflags.Args()[1]
	}
	jt := t7s.JobTypeReplace
	if *indexOp {
		jt = t7s.JobTypeIndex
	}
	t := &t7s.T7s{
		LeftDelim:  *leftdelim,
		RightDelim: *rightdelim,
		InPath:     inPath,
		OutPath:    outPath,
		JobType:    jt,
		VarFile:    *varfile,
		Require:    *require,
		IndexMerge: *indexMerge,
	}
	if err := t.Run(); err != nil {
		log.Fatal(err)
	}
}
