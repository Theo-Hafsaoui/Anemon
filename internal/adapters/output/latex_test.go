package output

import (
	"anemon/internal/core"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func removeWhitespace(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), "\n", "")
}

func TestApplySectionToTemplate(t *testing.T) {

	t.Run("happy path - should apply Professional section template with headers and items", func(t *testing.T) {
		template := "Start\n%EXPERIENCE_SECTIONS%\nEnd"
		headers := []string{"Company Name", "Position", "https://company.com", "Company"}
		items := []string{"Task 1", "Task 2"}
		keyword := []string{"1", "2"}

		want := strings.TrimSpace(`Start
%EXPERIENCE_SECTIONS%
\resumeSubheading
    {Company Name}{Position}
    {\href{https://company.com}{Company}}{ }
\resumeItemListStart
             \resumeItem{Task \textbf{1}}
\resumeItem{Task \textbf{2}}
\resumeItemListEnd
End`)

		processor := LatexProccesor{}
		got, err := processor.ApplySectionToTemplate(template, headers, items, "Professional",keyword)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if removeWhitespace(got) != removeWhitespace(want) {
			t.Errorf("ApplySectionToTemplate happy path failed;\nwant:\n%s\ngot:\n%s", want, got)
		}
	})

	t.Run("happy path - should apply Project section template with headers and items", func(t *testing.T) {
		template := "Start\n%PROJECTS_SECTIONS%\nEnd"
		headers := []string{"Project Name", "Project Description", "https://github.com/project"}
		items := []string{"Feature 1", "Feature 2"}
		keyword := []string{"1", "2"}

                want := `Start
%PROJECTS_SECTIONS%
\resumeProjectHeading
{\textbf{Project Name} | \emph{Project Description \href{https://github.com/project}{\faIcon{github}}}}{}
\resumeItemListStart
    \resumeItem{Feature \textbf{1}}
\resumeItem{Feature \textbf{2}}

\resumeItemListEnd

End`


		processor := LatexProccesor{}
		got, err := processor.ApplySectionToTemplate(template, headers, items, "Project",keyword)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if removeWhitespace(got) != removeWhitespace(want) {
			t.Errorf("ApplySectionToTemplate happy path failed;\nwant:\n%s\ngot:\n%s", want, got)
		}
	})

	t.Run("failure - should return error for unsupported section", func(t *testing.T) {
		template := "Start\n%EXPERIENCE_SECTIONS%\nEnd"
		headers := []string{"Company Name", "Position", "https://company.com", "Company"}
		items := []string{"Task 1", "Task 2"}
		keyword := []string{"1", "2"}

		processor := LatexProccesor{}
		_, err := processor.ApplySectionToTemplate(template, headers, items, "UnsupportedSection",keyword)
		if err == nil || err.Error() != "Don't know type UnsupportedSection" {
			t.Errorf("expected error 'Don't know type UnsupportedSection', got %v", err)
		}
	})

	t.Run("missing headers - should handle missing headers gracefully in Project section", func(t *testing.T) {
		template := "Start\n%PROJECTS_SECTIONS%\nEnd"
		headers := []string{"Project Name"}
		items := []string{"Feature 1", "Feature 2"}

		want := strings.TrimSpace(`Start
%PROJECTS_SECTIONS%
\resumeProjectHeading
{\textbf{Project Name} | \emph{$2 \href{$3}{\faIcon{github}}}}{}
\resumeItemListStart
    \resumeItem{Feature 1}
    \resumeItem{Feature 2}
\resumeItemListEnd
End`)

		processor := LatexProccesor{}
		got, err := processor.ApplySectionToTemplate(template, headers, items, "Project",nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if removeWhitespace(got) != removeWhitespace(want) {
			t.Errorf("ApplySectionToTemplate missing headers case failed;\nwant:\n%s\ngot:\n%s", want, got)
		}
	})
}

func TestReplaceWithSectionTemplate(t *testing.T) {

	t.Run("happy path - should replace hook with section template including headers and items", func(t *testing.T) {
		template := "Start\n%PROJECTS_SECTIONS%\nEnd"
		headers := []string{"Project Name", "Project Description", "https://github.com/project"}
		items := []string{"Feature 1", "Feature 2"}
		want := strings.TrimSpace(`Start
%PROJECTS_SECTIONS%
\resumeProjectHeading
{\textbf{Project Name} | \emph{Project Description \href{https://github.com/project}{\faIcon {github}}}}{}
\resumeItemListStart
    \resumeItem{Feature 1}
\resumeItem{Feature 2}

\resumeItemListEnd

End`)

		got := strings.TrimSpace(replaceWithSectionTemplate(template, ProjectTemplate, headers, items,nil))
		if removeWhitespace(got) != removeWhitespace(want) {
			t.Errorf("replaceWithSectionTemplate happy path failed;\nwant:\n%s\ngot:\n%s", want, got)
		}
	})

	t.Run("failure - should work even with missing headers", func(t *testing.T) {
		template := "Start\n%PROJECTS_SECTIONS%\nEnd"
		headers := []string{"Project Name", "Project Description"}
		items := []string{"Feature 1", "Feature 2"}
		want := strings.TrimSpace(`Start
%PROJECTS_SECTIONS%
\resumeProjectHeading
{\textbf{Project Name} | \emph{Project Description \href{$3}{\faIcon {github}}}}{}
\resumeItemListStart
    \resumeItem{Feature 1}
\resumeItem{Feature 2}

\resumeItemListEnd

End`)

		got := strings.TrimSpace(replaceWithSectionTemplate(template, ProjectTemplate, headers, items,nil))
		if removeWhitespace(got) != removeWhitespace(want) {
			t.Errorf("replaceWithSectionTemplate happy path failed;\nwant:\n%s\ngot:\n%s", want, got)
		}
	})
}

func TestReplaceHeaders(t *testing.T) {

	t.Run("happy path - should replace %ITEMS% with formatted items", func(t *testing.T) {
		template := "This is $1 and this is $2."
		section_items := []string{"ok", "ok"}
		want := "This is ok and this is ok."

		got := replace_headers(template, section_items)
		if got != want {
			t.Errorf("replace_headers happy path failed;\nwant:\n%s\ngot:\n%s", want, got)
		}
	})

	t.Run("happy path - should replace %ITEMS% with formatted items even with to much", func(t *testing.T) {
		template := "This is $1 and this is $2."
		section_items := []string{"ok", "ok", "ko"}
		want := "This is ok and this is ok."

		got := replace_headers(template, section_items)
		if got != want {
			t.Errorf("replace_headers happy path failed;\nwant:\n%s\ngot:\n%s", want, got)
		}
	})

	t.Run("failure - should handle malformed template with extra %ITEMS% tags", func(t *testing.T) {
		template := "This is $1 and this is $2 but not this $3."
		section_items := []string{"ok", "ok"}
		want := "This is ok and this is ok but not this $3."

		got := replace_headers(template, section_items)
		if got != want {
			t.Errorf("replace_headers failure case failed;\nwant:\n%s\ngot:\n%s", want, got)
		}
	})

	t.Run("nothing - should rm %ITEMS% if section_items is empty", func(t *testing.T) {
		template := "This is $1 and this is $2 but not this."
		section_items := []string{}
		want := "This is $1 and this is $2 but not this."
		got := replace_headers(template, section_items)
		if got != want {
			t.Errorf("replace_headers 'nothing' case failed;\nwant:\n%s\ngot:\n%s", want, got)
		}
	})
}

func TestReplaceItems(t *testing.T) {
	t.Run("happy path - should replace %ITEMS% with formatted items", func(t *testing.T) {
		template := "Start\n%ITEMS% End"
		section_items := []string{"Item 1", "Item 2", "Item 3"}
		keyword := []string{"1", "2"}
                want := `Start
\resumeItem{Item \textbf{1}}
\resumeItem{Item \textbf{2}}
\resumeItem{Item 3}
 End`


		got := replace_items(template, section_items,keyword)
		if got != want {
			t.Errorf("replace_items happy path failed;\nwant:\n%s\ngot:\n%s", want, got)
		}
	})

	t.Run("failure - should handle malformed template with extra %ITEMS% tags", func(t *testing.T) {
		template := "Start\n%ITEMS%Middle\n%ITEMS%\nEnd"
		section_items := []string{"Single Item"}
		want := "Start\n\\resumeItem{Single Item}\nMiddle\n%ITEMS%\nEnd"
		keyword := []string{"1", "2"}

		got := replace_items(template, section_items,keyword)
		if got != want {
			t.Errorf("replace_items failure case failed;\nwant:\n%s\ngot:\n%s", want, got)
		}
	})

	t.Run("nothing - should rm %ITEMS% if section_items is empty", func(t *testing.T) {
		template := "Start\n%ITEMS%\nEnd"
		section_items := []string{}
		keyword := []string{"1", "2"}
		want := "Start\n\nEnd"
		got := replace_items(template, section_items,keyword)
		if got != want {
			t.Errorf("replace_items 'nothing' case failed;\nwant:\n%s\ngot:\n%s", want, got)
		}
	})
}

func TestSanitize(t *testing.T) {
	t.Run("happy path - should sanitize special characters", func(t *testing.T) {
		got := "100% **bold text** and *italic text*"
		expected := "100\\% \\textbf{bold text} and \\emph{italic text}"
		result := sanitize(got)
		if result != expected {
			t.Errorf("Sanitize happy path failed; expected:\n%s\ngot:\n%s", expected, result)
		}
	})

	t.Run("failure - should handle unmatched patterns gracefully", func(t *testing.T) {
		got := "**bold text with *unmatched italic**"
		expected := "\\textbf{bold text with *unmatched italic}"
		result := sanitize(got)
		if result != expected {
			t.Errorf("Sanitize failure case failed; expected:\n%s\ngot:\n%s", expected, result)
		}
	})

	t.Run("nothing - should change nothing", func(t *testing.T) {
		got := "Just a plain string"
		expected := "Just a plain string"
		result := sanitize(got)
		if result != expected {
			t.Errorf("Sanitize 'nothing' case failed; expected:\n%s\ngot:\n%s", expected, result)
		}
	})
}

func TestApplyInfoToTemplate(t *testing.T) {
	template := "%VARS%"
	params := core.Params{
		Info: struct {
			Name      string `yaml:"name"`
			FirstName string `yaml:"firstname"`
			Number    string `yaml:"number"`
			Mail      string `yaml:"mail"`
			GitHub    string `yaml:"github"`
			LinkedIn  string `yaml:"linkedin"`
		}{
			Name:      "John Doe",
			FirstName: "John",
			Number:    "12345",
			Mail:      "john.doe@example.com",
			GitHub:    "https://github.com/johndoe",
			LinkedIn:  "https://linkedin.com/in/johndoe",
		},
		Variante: map[string][]string{},
	}
	want := `\def\Name{John Doe}
\def\FirstName{John}
\def\Number{12345}
\def\Mail{john.doe@example.com}
\def\GitHub{https://github.com/johndoe}
\def\LinkedIn{https://linkedin.com/in/johndoe}`
	got := ApplyInfoToTemplate(template, params)

	if removeWhitespace(got) != removeWhitespace(want) {
		t.Errorf("expected:\n%s\ngot:\n%s", want, got)
	}
}

func TestGetListOfTemplate(t *testing.T) {
	root := "testdata"
	templateDir := filepath.Join(root, "assets", "latex", "output")
	err := os.MkdirAll(templateDir, os.ModePerm)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	defer os.RemoveAll(root)

	files := []string{"foo.tex", "bar.tex", "garbage"}
	for _, file := range files {
		f, err := os.Create(filepath.Join(templateDir, file))
		if err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
		f.Close()
	}

	got, err := getListOfTemplate(root)
	if err != nil {
		t.Fatalf("getListOfTemplate returned an error: %v", err)
	}
	expected := files
	if len(got) != len(expected)-1 { //Should ommit garbage
		t.Errorf("expected %d files, got %d", len(expected), len(got))
	}
	for _, filePath := range got {
		_, err := os.Stat(filePath)
		if err != nil {
			t.Errorf("expected path %s not found: %v", filePath, err)
		}
	}
}
