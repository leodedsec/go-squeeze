package archive

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"os"
	"path/filepath"
)

type TarGzArchiver struct {
	file          *os.File
	tarWriter     *tar.Writer
	gzWriter      *gzip.Writer
	archivedCount uint
}

func NewTarGz() (*TarGzArchiver, error) {
	var err error
	archiver := &TarGzArchiver{}
	archiver.file, err = os.CreateTemp("", ".tar.gz")
	if err != nil {
		return nil, err
	}
	archiver.gzWriter = gzip.NewWriter(archiver.file)
	archiver.tarWriter = tar.NewWriter(archiver.gzWriter)
	return archiver, nil
}

func (archiver *TarGzArchiver) write(
	prefix string,
	path string,
) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	entryPath := buildEntryPath(prefix, info)
	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return err
	}
	header.Name = entryPath

	if err := archiver.tarWriter.WriteHeader(header); err != nil {
		return err
	}
	if err := writeFileTo(archiver.tarWriter, file); err != nil {
		return err
	}

	archiver.archivedCount += 1
	return nil
}

func (archiver *TarGzArchiver) Write(
	ctx context.Context,
	prefix string,
	path string,
) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return archiver.write(prefix, path)
	}
}

func (archiver *TarGzArchiver) Save(path string) (*SaveResult, error) {
	if err := archiver.Close(false); err != nil {
		return nil, err
	}

	archivePath := filepath.Join(path, genArchiveName("tar.gz"))
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

func (archiver *TarGzArchiver) Close(removeAnyway bool) error {
	errTarWriterClose := archiver.tarWriter.Close()
	errGzWriterClose := archiver.gzWriter.Close()
	errFileClose := archiver.file.Close()

	if errTarWriterClose != nil {
		_ = os.Remove(archiver.file.Name())
		return errTarWriterClose
	}
	if errGzWriterClose != nil {
		_ = os.Remove(archiver.file.Name())
		return errGzWriterClose
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
