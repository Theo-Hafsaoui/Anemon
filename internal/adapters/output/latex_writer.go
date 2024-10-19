package output

import (
	"anemon/internal/core"
	"errors"
	"log/slog"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
)

const NB_REPEAT = 1

type LatexReader struct{}
type LatexProccesor struct{}
var REG_EXTRACT_LANG_FROM_NAME = regexp.MustCompile(`-(.*?)-`)

// Write the template file in the assets directory
func (*LatexProccesor) MakeNewTemplate(path string, template string, name string) error {
        lang := REG_EXTRACT_LANG_FROM_NAME.FindStringSubmatch(name)
        san_template := sanitize_template(template,lang[1])
	err := os.WriteFile(path+"/assets/latex/output/"+name,
		[]byte(san_template), 0644)
	return err
}

// Read the template file in the assets directory from the root dir and apply the params given to it
func (*LatexReader) ReadCVTemplate(root string, params core.Params) (string, error) {
	file, err := os.ReadFile(root + "/assets/latex/template/template.tex")
	if err != nil {
		return "", err
	}
	return ApplyInfoToTemplate(string(file), params), nil
}

// Apply general information(name, mail...) to a template
func ApplyInfoToTemplate(template string, params core.Params) string {
	var varsBuilder strings.Builder
	infoValue := reflect.ValueOf(params.Info)
	for i := 0; i < infoValue.NumField(); i++ {
		field := infoValue.Type().Field(i)
		fieldValue := infoValue.Field(i)
		varsBuilder.WriteString("\\def\\" + field.Name + "{" + fieldValue.String() + "}\n")
	}
	return replaceVars(template, varsBuilder.String())
}

// Apply a section to a section type on a latex template
func (*LatexProccesor) ApplySectionToTemplate(template string, headers []string,
                                                item []string, section string, keyword []string) (string, error) {
	if len(section) < 2 {
		return "", errors.New("Don't know type " + section)
	}
	section = strings.ToUpper(string(section[0])) + section[1:]
	switch {
	case section == "Professional":
		template = replaceWithSectionTemplate(template, ProfessionalTemplate,
			headers, item,keyword)
	case section == "Project":
		template = replaceWithSectionTemplate(template, ProjectTemplate,
			headers, item,keyword)
	case section == "Education":
		template = replaceWithSectionTemplate(template, EducationTemplate,
			headers, nil,nil)
	case section == "Skill":
		template = replaceWithSectionTemplate(template, SkillTemplate,
			headers, nil,keyword)
	default:
		slog.Warn("Don't know type " + section)
		return "", errors.New("Don't know type " + section)
	}
	return template, nil
}

// Replace with the template defined in the template_sections.go const
func replaceWithSectionTemplate(template string, SectionTemplate TemplateStruct, headers []string, items []string,keywords []string) string {
	updated_template := strings.Replace(template, SectionTemplate.hook,
		SectionTemplate.hook+replace_headers(SectionTemplate.template, headers), 1)
	if items != nil {
		updated_template = replace_items(updated_template, items, keywords)
	}
	return updated_template
}

// Replace the %vars$ with the vars
func replaceVars(template string, vars string) string {
	updated_template := strings.Replace(template, "%VARS%", vars, 1)
	return updated_template
}

// Search and replace the headers in the template by their replacement
func replace_headers(sec_template string, replacements []string) string {
	for i := 0; i < len(replacements); i++ {
		position := fmt.Sprintf("$%d", i+1)
		sec_template = strings.Replace(sec_template,
			position, replacements[i], 1)
	}
	return sanitize(sec_template)
}

//replace the %item% keyword in the template with the item prepared for the CV
func replace_items(template string, o_section_items []string, keywords []string) string {
        section_items := emph_keyword(o_section_items,keywords)
	items := ""
	for _, item := range section_items {
		items += strings.Replace(single_item_template, "%ITEM%", item, 1)
	}
	template = strings.Replace(template,
		"%ITEMS%", items, 1)
	return sanitize(template)
}

//Take a list of items and return them with emphasis on the keyword
func emph_keyword(items []string, keywords []string) []string{
    if keywords == nil{ return items }

    res := make([]string, len(items))
    for i_i,item := range items{
        new_item := item
        for _,keyword := range keywords{
            new_item = strings.Replace(new_item, keyword,`\textbf{`+keyword+`}`,NB_REPEAT)
        }
        res[i_i]=new_item
    }
    return res
}

// Sanitize the template section with they version in the choosen language
func sanitize_template(template string,lang string) string {
        for english, translated := range translations[strings.ToUpper(lang)] {
            upper_translated := strings.ToUpper(translated)
            template = strings.ReplaceAll(template, strings.ToUpper(english), upper_translated)
            template = strings.ReplaceAll(template, english, upper_translated)
        }
	return template
}

// Sanitize the special charactere
func sanitize(template string) string {
	replacements := []struct {
		pattern     string
		replacement string
	}{
		{`([0-9])\%`, `$1\%`},
		{`\*\*(.*?)\*\*`, `\textbf{$1}`},
		{`\*(.*?)\*`, `\emph{$1}`},
		//{`\[(.*?)\]\((.*?)\)`, `\href{$2}{$1}`},
	}
	for _, r := range replacements {
		re := regexp.MustCompile(r.pattern)
		template = re.ReplaceAllString(template, r.replacement)
	}
	return template
}
