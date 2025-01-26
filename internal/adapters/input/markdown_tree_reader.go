package input

import (
	core "github.com/Theo-Hafsaoui/Anemon/internal/core"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type MarkdownSource struct{}

var hashtagRegex = regexp.MustCompile(`^#+`)

// `GetCVFrom` takes a root directory and extracts Markdown documents to generate a list of CV type.
//
// This function assumes the directory structure is a tree with a depth of 3, where each leaf
// is a Markdown (.md) document. Each document may contain multiple paragraphs, but the headers
// should not repeat within the same document.
func (*MarkdownSource) GetCVsFrom(root string) ([]core.CV, error) {
	cvsPath := root + "/cv"
	cvs := make([]core.CV, 0)

	current_lang := ""
	current_sections := make([]core.Section, 0)
	has_been_inside_dir := false

	err := filepath.Walk(cvsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if has_been_inside_dir {
				cvs[len(cvs)-1].Sections = current_sections
				current_sections = make([]core.Section, 0)
				has_been_inside_dir = false
			}
			current_lang = info.Name()
                        if current_lang != "cv"{
                            cvs = append(cvs, core.CV{Lang: current_lang})
                        }
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			has_been_inside_dir = true
			if current_lang == "" {
				return errors.New("markdown file found before lang directory")
			}
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			new_section := core.Section{Title: strings.TrimRight(info.Name(), ".md"),
                                                    Paragraphes: getParagrapheFrom(string(content))}
			current_sections = append(current_sections, new_section)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	cvs[len(cvs)-1].Sections = current_sections

	return cvs, nil
}

// Take string in the md format and return a slice of core.Paragraphe for each paragraph inside of it
// We define a paragraph by text block where \n\n indicate a separation
func getParagrapheFrom(s_section string) []core.Paragraphe {
	paragraphs := make([]core.Paragraphe, 0)
	for _, paragraphe := range strings.Split(s_section, "\n\n") {
		n_paragraphe, err := parse_paragraphe(paragraphe)
		if is_empty(n_paragraphe) {
			continue
		}
		if err != nil {
			fmt.Println("Failed to parse paragraphe ")
		} else {
			paragraphs = append(paragraphs, n_paragraphe)
		}
	}
	if len(paragraphs) == 0 {
		return nil
	}
	return paragraphs
}

// Return if a paragraphs has no header and no items
func is_empty(p core.Paragraphe) bool {
	no_header := p.H1 == "" && p.H2 == "" && p.H3 == "" && p.H4 == ""
	no_items := len(p.Items) == 0
	return no_header && no_items
}

/*
Parse parses a Markdown-like `paragraph` into a `Paragraph`,
extracting headings and descriptions based on the number of leading hashtags or stars.
Returns an error if the format is invalid.
*/
func parse_paragraphe(paragraph string) (core.Paragraphe, error) {
	var (
		n_paragraphe       core.Paragraphe
		bulletPrefix       = "- "
		skillAsteriskCount = 4 // Number of asterisks that signify a skill block
	)
	if len(strings.Split(paragraph, "\n\n")) > 1 {
		return n_paragraphe, errors.New("Tried to parse multiple paragraphs into a single section")
	}
	wasASkill := false
	lines := strings.Split(strings.TrimRight(paragraph, "\n"), "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		nbHashtags := len(hashtagRegex.FindString(line))
		if wasASkill {
			wasASkill = false
			n_paragraphe.H2 = strings.TrimLeft(line, bulletPrefix)
			continue
		}
		if nbHashtags == 0 && strings.HasPrefix(line, "*") && len(strings.Trim(line, "*")) == len(line)-skillAsteriskCount {
			n_paragraphe.H1 = strings.Trim(line, "*")
			wasASkill = true
			continue
		}
		if err := handleLine(&n_paragraphe, line, nbHashtags); err != nil {
			return n_paragraphe, err
		}
	}
	return n_paragraphe, nil
}

// handleLine processes a line based on the number of leading hashtags
func handleLine(n_paragraphe *core.Paragraphe, line string, nbHashtags int) error {
	if nbHashtags > 0 && line[nbHashtags] != ' ' {
		return fmt.Errorf("Err: cannot parse this md line {%s}, '#' should be followed by space", line)
	}
	switch nbHashtags {
	case 1:
		processesHeader(&n_paragraphe.H1, line, nbHashtags)
	case 2:
		processesHeader(&n_paragraphe.H2, line, nbHashtags)
	case 3:
		processesHeader(&n_paragraphe.H3, line, nbHashtags)
	case 4:
		processesHeader(&n_paragraphe.H4, line, nbHashtags)
	case 0:
		if strings.HasPrefix(line, "- ") && len(line) > 1 {
			n_paragraphe.Items = append(n_paragraphe.Items, strings.TrimLeft(line, "- "))
		}
	default:
		return fmt.Errorf("cannot parse this md line {%s}", line)
	}
	return nil
}

// Affect to the Header the line without the `-`
func processesHeader(pt_header *string, line string, nbHashtags int) {
	if *pt_header != "" {
		slog.Warn("Trying to overload Header", "oldHeader",
			*pt_header, "newHeader", line[nbHashtags+1:])
	}
	*pt_header = line[nbHashtags+1:]
}
