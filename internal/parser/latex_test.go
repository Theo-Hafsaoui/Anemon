package parser

import (
    "os"
    "path/filepath"
    "testing"
)

func TestReadLatex(t *testing.T) {
    dir := filepath.Join("../../assets", "latex", "template")
    templateFile := filepath.Join(dir, "template.tex")
    backupFile := filepath.Join(dir, "save.tex")
    if _, err := os.Stat(templateFile); err == nil {
        err = os.Rename(templateFile, backupFile)
        if err != nil {
            t.Fatalf("Failed to rename template.tex to save.tex: %v", err)
        }
    }
    err := os.WriteFile(templateFile, []byte("Hello World"), 0644)
    if err != nil {
        t.Fatalf("Failed to write file: %v", err)
    }

    content, err := read_template()
    if err != nil {
        t.Fatalf("Failed to read file: %v", err)
    }

    if content != "Hello World" {
        t.Fatalf("Expected 'Hello World', got '%s'", content)
    }

    err = os.Remove(templateFile)
    if err != nil {
        t.Fatalf("Failed to remove file: %v", err)
    }

    if _, err := os.Stat(backupFile); err == nil {
        err = os.Rename(backupFile, templateFile)
        if err != nil {
            t.Fatalf("Failed to rename save.tex back to template.tex: %v", err)
        }
    }
}
func TestWriteLatex(t *testing.T) {
    err := writeTemplate("Hello, world", "hello")
    if err != nil {
        t.Fatalf("Failed to write file: %v", err)
    }

    err = os.Remove("../../assets/latex/output/hello.tex")
    if err != nil {
        t.Fatalf("Failed to remove file: %v", err)
    }
}

