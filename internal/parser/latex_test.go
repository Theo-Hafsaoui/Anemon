package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
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

func TestApplySection(t *testing.T) {
    tests := []struct {
        name      string
        section   Section
        sectionType string
        want      string
    }{
        {
            name: "Professional Section",
            section: Section{
                first: "first",
                second: "second",
                third: "third",
                fourth: "fourth",
                description: []string{"item1", "item2"},
            },
            sectionType: "Professional",
            want: `
\resumeSubheading
    {first}{second}
    {\href{third}{fourth}}{ }
\resumeItemListStart
    \resumeItem{item1}
\resumeItem{item2}

\resumeItemListEnd
`,
        },
        {
            name: "Project Section",
            section: Section{
                first: "first",
                second: "second",
                third: "third",
                description: []string{"item1", "item2"},
            },
            sectionType: "Project",
            want: `
\resumeProjectHeading
{\textbf{first} $|$ \emph{second \href{third}{\faIcon{github}}}}{}
\resumeItemListStart
    \resumeItem{item1}
\resumeItem{item2}

\resumeItemListEnd
`,
        },
        {
            name: "Education Section",
            section: Section{
                first: "first",
                second: "second",
                third: "third",
                fourth: "fourth",
            },
            sectionType: "Education",
            want: `
\resumeSubheading
{\href{first}{second}}{}
{third}{fourth}
`,
        },
        {
            name: "Skill Section",
            section: Section{
                first: "first",
                second: "second",
                third: "third",
                fourth: "fourth",
            },
            sectionType: "Skill",
            want: `
\textbf{first}{: second} \\
`,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := applyToSection(tt.section, tt.sectionType)
            if err != nil {
                t.Fatalf("error when applying template: %v", err)
            }
            if diff := cmp.Diff(tt.want, got); diff != "" {
                t.Errorf("TestApplySection mismatch (-want +got):\n%s", diff)
            }
        })
    }
}

