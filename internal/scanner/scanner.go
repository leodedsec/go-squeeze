package scanner

import (
	"os"
	"path/filepath"
	"strings"
)

func arrayToMap(strArray []string) map[string]struct{} {
	result := make(map[string]struct{})
	for _, str := range strArray {
		result[strings.ToLower(str)] = struct{}{}
	}
	return result
}

func Scan(
	path string,
	recursion bool,
	excludeExtensions []string,
) ([]string, error) {
	var paths []string
	excludeExtensionsMap := arrayToMap(excludeExtensions)

	err := filepath.WalkDir(
		path,
		func(p string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				// apply to root dir
				if p == path {
					return nil
				}
				if !recursion {
					return filepath.SkipDir
				}
				return nil
			}

			if _, ok := excludeExtensionsMap[filepath.Ext(p)]; ok {
				return nil
			}

			paths = append(paths, p)

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return paths, nil
}
