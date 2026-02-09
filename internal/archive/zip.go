package archive

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type zipArchiver struct {
	file          *os.File
	writer        *zip.Writer
	archivedCount uint
}

func NewZip() (Archiver, error) {
	var err error
	archiver := &zipArchiver{}
	archiver.file, err = os.CreateTemp("", ".zip")
	if err != nil {
		return nil, err
	}
	archiver.writer = zip.NewWriter(archiver.file)
	return archiver, nil
}

func (archiver *zipArchiver) Write(key string, path string) error {
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

func (archiver *zipArchiver) Save(path string) (*saveResult, error) {
	if err := archiver.writer.Close(); err != nil {
		_ = archiver.file.Close()
		_ = os.Remove(archiver.file.Name())
		return nil, err
	}
	if err := archiver.file.Close(); err != nil {
		_ = os.Remove(archiver.file.Name())
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
	return &saveResult{
		archivePath:   archivePath,
		archivedCount: archiver.archivedCount,
		archiveSize:   stat.Size(),
	}, nil
}
