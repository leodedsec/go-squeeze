package grouper

import (
	"path/filepath"
	"strings"
)

func extractExtension(path string) string {
	ext := filepath.Ext(path)
	if ext == "" {
		return "unknown"
	}
	return strings.ToLower(ext)
}

func ByExtension(paths []string) map[string][]string {
	result := make(map[string][]string)

	for _, path := range paths {
		ext := extractExtension(path)
		result[ext] = append(result[ext], path)
	}

	return result
}
