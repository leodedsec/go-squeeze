package main

import (
	"context"
	"fmt"
	"go-squeeze/internal/appinfo"
	"go-squeeze/internal/archive"
	"go-squeeze/internal/argparser"
	"go-squeeze/internal/console"
	"go-squeeze/internal/grouper"
	"go-squeeze/internal/scanner"
	"os"
	"os/signal"
	"syscall"
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

	getArchiver := func(mode string) (archive.Archiver, error) {
		switch mode {
		case "zip":
			return archive.NewZip()
		case "tar.gz":
			return archive.NewTarGz()
		}
		return nil, fmt.Errorf("invalid mode: %s", mode)

	}

	archiver, err := getArchiver(ap.Mode())
	if err != nil {
		console.Error(err.Error())
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	stopCh := make(chan os.Signal, 1)
	errorCh := make(chan error, 1)
	doneCh := make(chan struct{})
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	byExtension := grouper.ByExtension(paths)
	go func() {
		defer close(doneCh)
		for key, pathsArray := range byExtension {
			for _, p := range pathsArray {
				if err := archiver.Write(ctx, key, p); err != nil {
					select {
					case errorCh <- err:
					default:
					}
					return
				}
			}
		}
	}()

	select {
	case <-stopCh:
		cancel()
		<-doneCh
		_ = archiver.Close(true)
		console.Info(
			"Operation stopped by the user. " +
				"The created archive file has been deleted.",
		)
		return
	case err := <-errorCh:
		cancel()
		<-doneCh
		_ = archiver.Close(true)
		console.Error(err.Error())
		return
	case <-doneCh:
		saveResult, err := archiver.Save(ap.OutputPath())
		if err != nil {
			console.Error(err.Error())
			return
		}
		tableMap := map[string]string{
			"Created": saveResult.ArchivePath(),
			"Size": fmt.Sprintf(
				"%.2f Mb",
				float64(saveResult.ArchiveSize())/(1024*1024),
			),
			"Count of archived files": fmt.Sprintf(
				"%d",
				saveResult.ArchivedCount(),
			),
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
		return
	}
}
