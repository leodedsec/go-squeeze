package argparser

import (
	"encoding/json"
	"flag"
	"fmt"
	"go-squeeze/internal/appinfo"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

var availableModes []string = []string{
	"zip",
	"tar.gz",
}

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
	mode              string
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

	flag.Usage = func() {
		out := flag.CommandLine.Output()

		fmt.Fprintln(out, "Usage:")
		fmt.Fprintln(out, "  app [OPTIONS]")
		fmt.Fprintln(out)

		fmt.Fprintln(out, "Options:")
		fmt.Fprintln(out, "  -mode")
		fmt.Fprintln(out, "        Type of archive. Available values: zip, tar.gz (default zip).")
		fmt.Fprintln(out)

		fmt.Fprintln(out, "  -ipath PATH")
		fmt.Fprintf(out, "        Read input files from PATH (default: %s).", currentDir)
		fmt.Fprintln(out)

		fmt.Fprintln(out, "  -opath PATH")
		fmt.Fprintf(out, "        Write output archive file to PATH (default: %s).", currentDir)
		fmt.Fprintln(out)

		fmt.Fprintln(out, "  -recursion")
		fmt.Fprintln(out, "        Scan subdirectories recursively (default true).")
		fmt.Fprintln(out, "  -recursion=false")
		fmt.Fprintln(out, "        Disable recursive scan.")
		fmt.Fprintln(out, "  -exclude EXT[,EXT...]")
		fmt.Fprintln(out, "        Exclude files with specified extensions.")
		fmt.Fprintln(out)

		fmt.Fprintln(out, "  -info")
		fmt.Fprintln(out, "        Show application information and exit.")

		fmt.Fprintln(out, "Examples:")
		fmt.Fprintln(out, "  app -ipath ./src -opath ./out")
		fmt.Fprintln(out, "  app -exclude txt,csv,exe")
		fmt.Fprintln(out, "  app -recursion=false")
	}

	mode := flag.String("mode", "zip", "")
	inputPath := flag.String(
		"ipath",
		currentDir,
		"",
	)
	outputPath := flag.String(
		"opath",
		currentDir,
		"",
	)
	recursion := flag.Bool(
		"recursion",
		true,
		"",
	)
	flag.Var(
		&excludeExtensions,
		"exclude",
		"",
	)
	showInfo := flag.Bool("info", false, "")

	flag.Parse()

	if *showInfo {
		fmt.Printf(
			"%s — %s\nVersion: %s\nGithub: %s",
			appinfo.Name, appinfo.Description,
			appinfo.Version,
			appinfo.Github,
		)
		os.Exit(0)
	}

	if !slices.Contains(availableModes, *mode) {
		return nil, fmt.Errorf("invalid mode: %s", *mode)
	}

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
		mode:              *mode,
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

func (ap *ArgParser) Mode() string { return ap.mode }
