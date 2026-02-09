package archive

import (
	"fmt"
	"time"
)

type saveResult struct {
	archivePath   string
	archivedCount uint
	archiveSize   int64
}

func (sv *saveResult) ArchivePath() string {
	return sv.archivePath
}

func (sv *saveResult) ArchivedCount() uint {
	return sv.archivedCount
}

func (sv *saveResult) ArchiveSize() int64 {
	return sv.archiveSize
}

type Archiver interface {
	Write(key string, path string) error
	Save(path string) (*saveResult, error)
}

func GenArchiveName() string {
	return fmt.Sprintf("%s.zip", time.Now().Format("20060102_150405"))
}
