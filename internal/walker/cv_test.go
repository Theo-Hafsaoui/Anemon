package walker

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWalkCV(t *testing.T) {
	rootDir := t.TempDir()
	paths := []struct {
		relativePath string
		content      string
	}{
		{"eng/education.md", "Education"},
		{"eng/project.md", "Project"},
		{"fr/education.md", "Education"},
		{"fr/work.md", "Work"},
	}

	for _, p := range paths {
		fullPath := filepath.Join(rootDir, p.relativePath)
		if err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm); err != nil {
			t.Fatalf("Failed to create directories: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(p.content), os.ModePerm); err != nil {
			t.Fatalf("Failed to write file: %v", err)
		}
	}

	result, err := WalkCV(rootDir)
	if err != nil {
		t.Fatalf("walkCV failed: %v", err)
	}

	expected := map[string]map[string]string{
		"eng": {
			"education": "Education",
			"project":   "Project",
		},
		"fr": {
			"education": "Education",
			"work":      "Work",
		},
	}

	for lang, expectedFiles := range expected {
		if resultFiles, ok := result[lang]; ok {
			for file, expectedContent := range expectedFiles {
				if resultContent, ok := resultFiles[file]; ok {
					if resultContent != expectedContent {
						t.Errorf("Expected %s/%s to be %q, got %q", lang, file, expectedContent, resultContent)
					}
				} else {
					t.Errorf("Expected file %s/%s not found in result", lang, file)
				}
			}
		} else {
			t.Errorf("Expected language %s not found in result", lang)
		}
	}

	for lang := range result {
		if _, found := expected[lang]; !found {
			t.Errorf("Unexpected language found: %s", lang)
		}
	}
}
