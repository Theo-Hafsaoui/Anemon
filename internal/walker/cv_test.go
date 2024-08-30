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

	result, err := walkCV(rootDir)
	if err != nil {
		t.Fatalf("walkCV failed: %v", err)
	}

	expected := map[string]string{
		"eng/education": "Education",
		"eng/project":   "Project",
		"fr/education":  "Education",
		"fr/work":       "Work",
	}

	for key, expectedValue := range expected {
		if result[key] != expectedValue {
			t.Errorf("Expected %s to be %q, got %q", key, expectedValue, result[key])
		}
	}

	for key := range result {
		if _, found := expected[key]; !found {
			t.Errorf("Unexpected key found: %s", key)
		}
	}
}
