package archive

import (
	"context"
	"fmt"
	"time"
)

type SaveResult struct {
	archivePath   string
	archivedCount uint
	archiveSize   int64
}

func (sv *SaveResult) ArchivePath() string {
	return sv.archivePath
}

func (sv *SaveResult) ArchivedCount() uint {
	return sv.archivedCount
}

func (sv *SaveResult) ArchiveSize() int64 {
	return sv.archiveSize
}

type Archiver interface {
	Write(key string, path string) error
	WriteWithCtx(ctx context.Context, key string, path string) error
	Save(path string) (*SaveResult, error)
	Close(removeAnyway bool) error
}

func GenArchiveName() string {
	return fmt.Sprintf("%s.zip", time.Now().Format("20060102_150405"))
}
