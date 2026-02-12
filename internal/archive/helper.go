package archive

import (
	"fmt"
	"io"
	"os"
	"time"
)

func writeFileTo(w io.Writer, file *os.File) error {
	_, err := io.Copy(w, file)
	if err != nil {
		return err
	}
	return nil
}

func buildEntryPath(prefix string, info os.FileInfo) string {
	if prefix == "" {
		return info.Name()
	}
	return fmt.Sprintf("%s/%s", prefix, info.Name())
}

func genArchiveName(ext string) string {
	return fmt.Sprintf("%s.%s", time.Now().Format("20060102_150405"), ext)
}
