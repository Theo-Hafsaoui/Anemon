package core

import (
	"reflect"
	"testing"
)

type MockParamsSource struct {
	Params Params
	Err    error
}

func (s *MockParamsSource) GetParamsFrom(root string) (Params, error) {
	if s.Err != nil {
		return s.Params, s.Err
	}
	return s.Params, nil
}

type MockCompiler struct{}

func (c *MockCompiler) CompileTemplate(root string) (int, error) { return 0, nil }

type MockSource struct {
	CVs []CV
	Err error
}

func (s *MockSource) GetCVsFrom(root string) ([]CV, error) {
	if s.Err != nil {
		return nil, s.Err
	}
	return s.CVs, nil
}

type MockTemplateReader struct {
	Template string
	Err      error
}

func (r *MockTemplateReader) ReadCVTemplate(path string, params Params) (string, error) {
	if r.Err != nil {
		return "", r.Err
	}
	return r.Template, nil
}

type MockTemplateProcessor struct {
	AppliedTemplates []string
	GeneratedFiles   map[string]string
	ApplyErr         error
	MakeErr          error
}

func (p *MockTemplateProcessor) MakeNewTemplate(path string, template string, name string) error {
	if p.MakeErr != nil {
		return p.MakeErr
	}
	p.GeneratedFiles[name] = template
	return nil
}

func (p *MockTemplateProcessor) ApplySectionToTemplate(template string, headers []string, items []string, section string, keywords []string) (string, error) {
	if p.ApplyErr != nil {
		return "", p.ApplyErr
	}
	result := template + " | Section: " + section + " | Headers: " + headers[0] + ", " + headers[1] + ", " + headers[2] + ", " + headers[3]
	for _, item := range items {
		result += " | Item: " + item
	}
	p.AppliedTemplates = append(p.AppliedTemplates, result)
	return result, nil
}

func TestGenerateTemplates(t *testing.T) {
	root := "testRoot"
	baseTemplate := "base template content"
	params := Params{
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

	cv := CV{
		Lang: "EN",
		Sections: []Section{
			{
				Title: "Work Experience",
				Paragraphes: []Paragraphe{
					{
						H1:    "Job Title",
						H2:    "Company",
						H3:    "Location",
						H4:    "Date",
						Items: []string{"Managed projects", "Led team"},
					},
				},
			},
		},
	}

	t.Run("When giving one language and two variante should generate two cv in the language", func(t *testing.T) {
		source := &MockSource{CVs: []CV{cv}}
		paramsSource := &MockParamsSource{Params: params}
		templateReader := &MockTemplateReader{Template: baseTemplate}
		templateProcessor := &MockTemplateProcessor{GeneratedFiles: make(map[string]string)}
		compiler := &MockCompiler{}
		var builder BuilderService

		builder.SetRoot(root)
		builder.SetSource(source)
		builder.SetParamsSource(paramsSource)
		builder.SetTemplateReader(templateReader)
		builder.SetTemplateProcessor(templateProcessor)
		builder.SetCompiler(compiler)

		service := builder.GetService()

		err := service.GenerateTemplates()

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(templateProcessor.GeneratedFiles) != 2 {
			t.Fatalf("expected 2 generated file, got %d", len(templateProcessor.GeneratedFiles))
		}

		generatedTemplate, exists := templateProcessor.GeneratedFiles["CV-EN-optionA.tex"]
		if !exists {
			t.Fatalf("expected generated file 'CV-EN-optionA.tex' to exist")
		}

		expectedContent := baseTemplate + " | Section: Work Experience | Headers: Job Title, Company, Location, Date | Item: Managed projects | Item: Led team"
		if !reflect.DeepEqual(generatedTemplate, expectedContent) {
			t.Errorf("expected generated template content %v, got %v", expectedContent, generatedTemplate)
		}
	})

	t.Run("When giving one language and zero info or variante should generate one cv in the language", func(t *testing.T) {
		source := &MockSource{CVs: []CV{cv}}
		paramsSource := &MockParamsSource{Params: Params{}}
		templateReader := &MockTemplateReader{Template: baseTemplate}
		templateProcessor := &MockTemplateProcessor{GeneratedFiles: make(map[string]string)}
		compiler := &MockCompiler{}
		var builder BuilderService

		builder.SetRoot(root)
		builder.SetSource(source)
		builder.SetParamsSource(paramsSource)
		builder.SetTemplateReader(templateReader)
		builder.SetTemplateProcessor(templateProcessor)
		builder.SetCompiler(compiler)

		service := builder.GetService()
		err := service.GenerateTemplates()

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(templateProcessor.GeneratedFiles) != 1 {
			t.Fatalf("expected 1 generated file, got %d", len(templateProcessor.GeneratedFiles))
		}

		generatedTemplate, exists := templateProcessor.GeneratedFiles["CV-EN-simple.tex"]
		if !exists {
			t.Fatalf("expected generated file 'CV-EN-simple.tex' to exist")
		}

		expectedContent := baseTemplate + " | Section: Work Experience | Headers: Job Title, Company, Location, Date | Item: Managed projects | Item: Led team"
		if !reflect.DeepEqual(generatedTemplate, expectedContent) {
			t.Errorf("expected generated template content %v, got %v", expectedContent, generatedTemplate)
		}
	})

}

func TestGetScore(t *testing.T) {
	tests := []struct {
		item     string
		keywords []string
		expected int
	}{
		{item: "foo bar baz", keywords: []string{"foo", "bar", "baz"}, expected: 3},
		{item: "foo qux", keywords: []string{"foo", "bar", "baz"}, expected: 1},
		{item: "patate douce", keywords: []string{"foo", "bar", "baz"}, expected: 0},
		{item: "Foo Bar Baz", keywords: []string{"foo", "bar", "baz"}, expected: 0},
		{item: "foo bar baz", keywords: []string{}, expected: 0},
		{item: "", keywords: []string{"foo", "bar", "baz"}, expected: 0},
		{item: "foo bar foo baz", keywords: []string{"foo", "foo bar"}, expected: 3},
	}

	for _, tt := range tests {
		t.Run(tt.item, func(t *testing.T) {
			score := getScore(tt.item, tt.keywords)
			if score != tt.expected {
				t.Errorf("for item=%q and keywords=%v, expected %d, got %d", tt.item, tt.keywords, tt.expected, score)
			}
		})
	}
}

func TestSortByScore(t *testing.T) {
	keywords := []string{"foo", "bar"}
	tests := []struct {
		items    []string
		expected []string
	}{
		{
			items:    []string{" with foo", " with foo and bar", " with neither", " with bar"},
			expected: []string{" with foo and bar", " with foo", " with bar", " with neither"},
		},

		{
			items:    []string{" with foo bar foo bar", " with foo", " with bar bar bar", " with foo bar"},
			expected: []string{" with foo bar foo bar", " with bar bar bar", " with foo bar", " with foo"},
		},
		{
			items:    []string{" without keywords", " with bar", "Another  without keywords"},
			expected: []string{" with bar", " without keywords", "Another  without keywords"},
		},
	}

	for _, tt := range tests {
		items := make([]string, len(tt.items))
		copy(items, tt.items)
		sortByScore(items, keywords)

		for i := range items {
			if items[i] != tt.expected[i] {
				t.Errorf("after sorting, expected %v but got %v", tt.expected, items)
			}
		}
	}
}

func TestGetLowestParagraphe(t *testing.T) {
	keywords := []string{"skills", "experience", "project"}
	tests := []struct {
		section  Section
		expected int
	}{
		{
			section: Section{
				Title: "Work Experience",
				Paragraphes: []Paragraphe{
					{H1: "Experience 1", Items: []string{"skills", "project"}},
					{H1: "Experience 2", Items: []string{"experience", "skills"}},
					{H1: "Experience 3", Items: []string{"project"}},
				},
			},
			expected: 2,
		},
		{
			section: Section{
				Title: "Projects",
				Paragraphes: []Paragraphe{
					{H1: "Project A", Items: []string{"project", "skills"}},
					{H1: "Project B", Items: []string{"skills", "experience"}},
					{H1: "Project C", Items: []string{"project", "experience"}},
				},
			},
			expected: 0,
		},
		{
			section: Section{
				Title: "Skills",
				Paragraphes: []Paragraphe{
					{H1: "Skillset 1", Items: []string{"skills", "tools", "project"}},
					{H1: "Skillset 2", Items: []string{"experience", "skills"}},
				},
			},
			expected: 0,
		},
		{
			section: Section{
				Title: "Summary",
				Paragraphes: []Paragraphe{
					{H1: "Summary", Items: []string{}},
					{H1: "Overview", Items: []string{}},
				},
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		result := getLowestParagraphe(tt.section, keywords)
		if result != tt.expected {
			t.Errorf("expected index %v but got %v", tt.expected, result)
		}
	}
}
