package archive

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type ZipArchiver struct {
	file          *os.File
	writer        *zip.Writer
	archivedCount uint
}

func NewZip() (*ZipArchiver, error) {
	var err error
	archiver := &ZipArchiver{}
	archiver.file, err = os.CreateTemp("", ".zip")
	if err != nil {
		return nil, err
	}
	archiver.writer = zip.NewWriter(archiver.file)
	return archiver, nil
}

func (archiver *ZipArchiver) Write(key string, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	var name string
	if key != "" {
		name = fmt.Sprintf("%s/%s", key, stat.Name())
	} else {
		name = stat.Name()
	}
	writer, err := archiver.writer.Create(name)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	if err != nil {
		return err
	}

	archiver.archivedCount += 1
	return nil
}

func (archiver *ZipArchiver) WriteWithCtx(
	ctx context.Context,
	key string,
	path string,
) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return archiver.Write(key, path)
	}
}

func (archiver *ZipArchiver) Save(path string) (*SaveResult, error) {
	if err := archiver.Close(false); err != nil {
		return nil, err
	}

	archivePath := filepath.Join(path, GenArchiveName())
	err := os.Rename(archiver.file.Name(), archivePath)
	if err != nil {
		_ = os.Remove(archiver.file.Name())
		return nil, err
	}

	stat, err := os.Stat(archivePath)
	if err != nil {
		return nil, err
	}
	return &SaveResult{
		archivePath:   archivePath,
		archivedCount: archiver.archivedCount,
		archiveSize:   stat.Size(),
	}, nil
}

func (archiver *ZipArchiver) Close(removeAnyway bool) error {
	errWriterClose := archiver.writer.Close()
	errFileClose := archiver.file.Close()
	if errWriterClose != nil {
		_ = os.Remove(archiver.file.Name())
		return errWriterClose
	}
	if errFileClose != nil {
		_ = os.Remove(archiver.file.Name())
		return errFileClose
	}

	if removeAnyway {
		_ = os.Remove(archiver.file.Name())
	}
	return nil
}
