package archive

import (
	"context"
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
	Write(ctx context.Context, prefix string, path string) error
	Save(path string) (*SaveResult, error)
	Close(removeAnyway bool) error
}
