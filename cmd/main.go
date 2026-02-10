package main

import (
	"fmt"
	"go-squeeze/internal/appinfo"
	"go-squeeze/internal/archive"
	"go-squeeze/internal/argparser"
	"go-squeeze/internal/console"
	"go-squeeze/internal/grouper"
	"go-squeeze/internal/scanner"
)

func main() {
	ap, err := argparser.New()
	if err != nil {
		console.Error(err.Error())
		return
	}

	paths, err := scanner.Scan(
		ap.InputPath(),
		ap.IsRecursion(),
		ap.ExcludeExtensions(),
	)
	if err != nil {
		console.Error(err.Error())
		return
	}

	zip, err := archive.NewZip()
	if err != nil {
		console.Error(err.Error())
		return
	}

	byExtension := grouper.ByExtension(paths)
	for key, pathsArray := range byExtension {
		for _, p := range pathsArray {
			err = zip.Write(key, p)
			if err != nil {
				console.Error(err.Error())
				return
			}
		}
	}

	saveResult, err := zip.Save(ap.OutputPath())
	if err != nil {
		console.Error(err.Error())
		return
	}

	tableMap := map[string]string{
		"Created":                 saveResult.ArchivePath(),
		"Size":                    fmt.Sprintf("%.2f Mb", float64(saveResult.ArchiveSize())/(1024*1024)),
		"Count of archived files": fmt.Sprintf("%d", saveResult.ArchivedCount()),
	}

	console.Table(
		fmt.Sprintf(
			"\n%s%s%s\n%sGithub: %s%s",
			console.Color.Green, appinfo.Name, console.Color.SysReset,
			console.Color.Green, appinfo.Github, console.Color.SysReset,
		),
		tableMap,
	)
	console.PressEnterToExit()
}
