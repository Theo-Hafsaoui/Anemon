package walker

import (
	"os"
	"path/filepath"
	"strings"
)

// walkCV traverses the directory tree starting at root and returns a map of map where
// the keys are first the language folowed by the second key which is the section
// the values are the contents of the markdown files.
func WalkCV(root string) (map[string]map[string]string, error) {
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
        md_per_language := make(map[string]map[string]string)
        for k := range fileMap{
            k_split := strings.Split(k,"/")
            if md_per_language [k_split[0]] == nil {
                md_per_language[k_split[0]]= make(map[string]string)
            }
            md_per_language[k_split[0]][k_split[1]]=fileMap[k]
        }
	return md_per_language, nil
}
