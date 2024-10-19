package core

import (
	"fmt"
	"strings"
)

type ICVservice interface {
	generateTemplates(root string, source Source, templateReader TemplateReader, templateProcessor TemplateProcessor) error
}

type Compiler interface {
	CompileTemplate(root string) (int, error)
}

type SourceParams interface {
	GetParamsFrom(root string) (Params, error)
}

type Source interface {
	GetCVsFrom(root string) ([]CV, error)
}

type TemplateReader interface {
	ReadCVTemplate(root string, params Params) (string, error)
}

type TemplateProcessor interface {
	MakeNewTemplate(path string, template string, name string) error
	ApplySectionToTemplate(template string, headers []string, item []string, section string, keyword []string) (string, error)
}

// CV with Language and Sections
type CV struct {
	Lang     string
	Sections []Section
}

type Section struct {
	Title       string
	Paragraphes []Paragraphe
}

type Paragraphe struct {
	H1    string
	H2    string
	H3    string
	H4    string
	Items []string
}

type Params struct {
	Info struct {
		Name      string `yaml:"name"`
		FirstName string `yaml:"firstname"`
		Number    string `yaml:"number"`
		Mail      string `yaml:"mail"`
		GitHub    string `yaml:"github"`
		LinkedIn  string `yaml:"linkedin"`
	} `yaml:"info"`
	Variante map[string][]string `yaml:"variante"`
}

// Print Method for CV
func (cv *CV) Print() {
	fmt.Println("CV Language: " + cv.Lang)
	fmt.Println(strings.Repeat("=", 40))
	fmt.Printf("With %d Sections\n", len(cv.Sections))
	for _, section := range cv.Sections {
		fmt.Printf("Section: %s\n", section.Title)
		fmt.Println(strings.Repeat("-", 40))
		for _, p := range section.Paragraphes {
			if p.H1 != "" {
				fmt.Printf("H1: %s\n", p.H1)
			} else {
				fmt.Printf("No H1")
			}
			if p.H2 != "" {
				fmt.Printf("  H2: %s\n", p.H2)
			} else {
				fmt.Printf("No H2")
			}
			if p.H3 != "" {
				fmt.Printf("    H3: %s\n", p.H3)
			} else {
				fmt.Printf("No H3")
			}
			if p.H4 != "" {
				fmt.Printf("      H4: %s\n", p.H4)
			} else {
				fmt.Printf("No H4")
			}

			if len(p.Items) > 0 {
				fmt.Println("      Items:")
				for _, item := range p.Items {
					fmt.Printf("      - %s\n", item)
				}
			}
			fmt.Println()
		}
		fmt.Println()
	}
}
