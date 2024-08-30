package walker

import (
	"os"
	"path/filepath"
	"strings"
)

// walkCV traverses the directory tree starting at root and returns a map where
// the keys are the relative paths without the .md extension and the values are
// the contents of the markdown files.
func walkCV(root string) (map[string]string, error) {
	fileMap := make(map[string]string)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			relativePath := strings.TrimSuffix(strings.TrimPrefix(path, root+string(os.PathSeparator)), ".md")
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			fileMap[relativePath] = string(content)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return fileMap, nil
}
