package input

import (
	core "anemon/internal/core"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var (
	yamlContent = `
info:
  name: Doe
  firstname: John
  number: "12345"
  mail: john.doe@example.com
  github: johndoe
  linkedin: john-doe-linkedin
variante:
  optionA:
    - "value1"
    - "value2"
  optionB:
    - "value3"
`
	work_input = `
# Back-End Intern
## February 2024 -- August 2024
### TechCorp
#### Internship
- Assisted in developing and optimizing a key business process for expanding into new markets, collaborating with various teams to ensure compliance and smooth integration.
- Participated in the migration of core backend logic from a monolithic application to a microservice architecture, improving system performance and scalability.
- Enhanced system monitoring and reliability by implementing tracing mechanisms and performance objectives.

# Back-End Intern
## February 2024 -- August 2024
### TechCorp
#### Internship
- Assisted in developing and optimizing a key business process for expanding into new markets, collaborating with various teams to ensure compliance and smooth integration.
- Participated in the migration of core backend logic from a monolithic application to a microservice architecture, improving system performance and scalability.
- Enhanced system monitoring and reliability by implementing tracing mechanisms and performance objectives.`

	skill_input = `
**Langage**
 - Langage A, Langage B, Langage C, Langage D, Langage E, Langage F, Langage G/H

**Langage**
 - Langage A, Langage B, Langage C, Langage D, Langage E, Langage F, Langage G/H`
	invalid_skill_input = `
 - Langage A, Langage B, Langage C, Langage D, Langage E, Langage F, Langage G/H

**Langag
 - Langage A, Langage B, Langage C, Langage D, Langage E, Langage F, Langage G/H`

	skill_paragraphe = core.Paragraphe{H1: "Langage", H2: "Langage A, Langage B, Langage C, Langage D, Langage E, Langage F, Langage G/H"}

	work_paragraphe = core.Paragraphe{H1: "Back-End Intern", H2: "February 2024 -- August 2024",
		H3: "TechCorp", H4: "Internship", Items: []string{
			"Assisted in developing and optimizing a key business process for expanding into new markets, collaborating with various teams to ensure compliance and smooth integration.",
			"Participated in the migration of core backend logic from a monolithic application to a microservice architecture, improving system performance and scalability.",
			"Enhanced system monitoring and reliability by implementing tracing mechanisms and performance objectives."}}

	invalid_input         = "ajsdlhsaeld##dafdbhkbhkjsd##"
	work_expected_result  = []core.Paragraphe{work_paragraphe, work_paragraphe}
	skill_expected_result = []core.Paragraphe{skill_paragraphe, skill_paragraphe}

	paths = []struct {
		relativePath string
		content      string
	}{
		{"cv/eng/education.md", work_input},
		{"cv/eng/project.md", work_input},
		{"cv/eng/work.md", work_input},
		{"cv/eng/skill.md", skill_input},
		{"cv/fr/education.md", work_input},
		{"cv/fr/work.md", work_input},
		{"cv/fr/skill.md", skill_input},
	}
)

func TestParagraphe(t *testing.T) {

	t.Run("Work Paragraphes should return a slice of valid Paragraphe", func(t *testing.T) {
		got := getParagrapheFrom(work_input)
		want := work_expected_result
		if !reflect.DeepEqual(got[0], want[0]) {
			t.Fatalf("the first Paragraphe should be :\n%s\n got :%s", want, got)
		}

		if !reflect.DeepEqual(got[1], want[1]) {
			t.Fatalf("the first Paragraphe should be :\n%s\n got :%s", want, got)
		}
	})

	t.Run("Invalid input should return nothing", func(t *testing.T) {
		result := getParagrapheFrom(invalid_input)
		if result != nil {
			t.Fatalf("Invalid input should return nil got %v", result)
		}
	})

	t.Run("Skill Paragraphe should be return from valid input", func(t *testing.T) {
		got := getParagrapheFrom(skill_input)
		want := skill_expected_result
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("the first Paragraphe should be :\n%s\n got :%s", want, got)
		}
	})

	t.Run("Invalid skill Paragraphe should return nil", func(t *testing.T) {
		got := getParagrapheFrom(invalid_skill_input)
		if got != nil {
			t.Fatalf("Invalid input should return nil got %v", got)
		}
	})
}

func TestSections(t *testing.T) {
	rootDir := t.TempDir()
	for _, p := range paths {
		fullPath := filepath.Join(rootDir, p.relativePath)
		if err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm); err != nil {
			t.Fatalf("Failed to create directories: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(p.content), os.ModePerm); err != nil {
			t.Fatalf("Failed to write file: %v", err)
		}
	}

	t.Run("Should return a valid list of cv", func(t *testing.T) {
		source := MarkdownSource{}
		got, err := source.GetCVsFrom(rootDir)
		if err != nil {
			t.Fatalf("Failed to getCV got %s", err)
		}

		lang_got := got[len(got)-1].Lang
		l_want := "fr"
		if lang_got != l_want {
			t.Fatalf("Should have %s got %s", l_want, lang_got)
		}

		sec_got := len(got[len(got)-1].Sections)
		s_want := 3
		if sec_got != s_want {
			t.Fatalf("Should have %d got %d", s_want, sec_got)
		}

		t_sec_got := got[len(got)-1].Sections[len(got[len(got)-1].Sections)-1].Title
		t_s_want := "work"
		if t_sec_got != t_s_want {
			t.Fatalf("Should have %s got %s", t_s_want, t_sec_got)
		}

		p_t_sec_got := got[len(got)-1].Sections[len(got[len(got)-1].Sections)-1].Paragraphes
		p_t_s_want := work_expected_result

		if len(p_t_sec_got) != len(p_t_s_want) {
			t.Fatalf("Should have len %d got %d", len(p_t_s_want), len(p_t_sec_got))
		}

		if p_t_sec_got[len(p_t_sec_got)-1].H1 != p_t_s_want[len(p_t_s_want)-1].H1 {
			t.Fatalf("Should have title %s got %s", p_t_s_want[len(p_t_s_want)-1].H1, p_t_sec_got[len(p_t_sec_got)-1].H1)
		}

	})
}

func TestGetParamsFrom(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	yamlFilePath := filepath.Join(tempDir, "params.yml")
	err = os.WriteFile(yamlFilePath, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write YAML file: %v", err)
	}

	source := &YamlSource{}
	params, err := source.GetParamsFrom(tempDir)
	if err != nil {
		t.Fatalf("GetParamsFrom returned an error: %v", err)
	}
	expectedParams := core.Params{
		Info: struct {
			Name      string `yaml:"name"`
			FirstName string `yaml:"firstname"`
			Number    string `yaml:"number"`
			Mail      string `yaml:"mail"`
			GitHub    string `yaml:"github"`
			LinkedIn  string `yaml:"linkedin"`
		}{
			Name:      "Doe",
			FirstName: "John",
			Number:    "12345",
			Mail:      "john.doe@example.com",
			GitHub:    "johndoe",
			LinkedIn:  "john-doe-linkedin",
		},
		Variante: map[string][]string{
			"optionA": {"value1", "value2"},
			"optionB": {"value3"},
		},
	}
	if !reflect.DeepEqual(params, expectedParams) {
		t.Errorf("Expected %+v, but got %+v", expectedParams, params)
	}
}
