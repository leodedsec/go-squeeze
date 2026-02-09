package argparser

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type listFlag []string

func (f *listFlag) String() string {
	b, _ := json.Marshal(*f)
	return string(b)
}

func (f *listFlag) Set(value string) error {
	for _, str := range strings.Split(value, ",") {
		*f = append(*f, str)
	}
	return nil
}

type ArgParser struct {
	inputPath         string
	outputPath        string
	recursion         bool
	excludeExtensions []string
}

func New() (*ArgParser, error) {
	currentDir, err := filepath.Abs(".")
	if err != nil {
		return nil, err
	}

	var excludeExtensions listFlag

	inputPath := flag.String(
		"ipath",
		currentDir,
		"Input path - path to directory.\n"+
			"Type: string\n"+
			"Default: current directory\n"+
			"Example: -ipath(--ipath) ./test",
	)
	outputPath := flag.String(
		"opath",
		currentDir,
		"Output path - path to save.\n"+
			"Type: string\n"+
			"Default: current directory\n"+
			"Example: -opath(--opath) ./test",
	)
	recursion := flag.Bool(
		"recursion",
		true,
		"Recursive scan of subdirectories.\n"+
			"Type: boolean\n"+
			"Default: true\n"+
			"Example: -recursion(--recursion) false",
	)
	flag.Var(
		&excludeExtensions,
		"exclude",
		"Exclude extensions.\n"+
			"Type: sequence of strings separated by ','\n"+
			"Default: empty\n"+
			"Example: -exclude(--exclude) txt,exe,csv",
	)

	flag.Parse()

	iPath, err := filepath.Abs(*inputPath)
	if err != nil {
		return nil, fmt.Errorf("invalid input path: %s", *inputPath)
	}
	oPath, err := filepath.Abs(*outputPath)
	if err != nil {
		return nil, fmt.Errorf("invalid output path: %s", *outputPath)
	}

	if err = os.MkdirAll(oPath, 0755); err != nil {
		return nil, err
	}

	for i, ext := range excludeExtensions {
		excludeExtensions[i] = "." + strings.TrimLeft(ext, ".")
	}

	return &ArgParser{
		inputPath:         iPath,
		outputPath:        oPath,
		recursion:         *recursion,
		excludeExtensions: excludeExtensions,
	}, nil
}

func (ap *ArgParser) InputPath() string {
	return ap.inputPath
}

func (ap *ArgParser) OutputPath() string {
	return ap.outputPath
}

func (ap *ArgParser) IsRecursion() bool {
	return ap.recursion
}

func (ap *ArgParser) ExcludeExtensions() []string {
	return ap.excludeExtensions
}
