package markuplanguages

import (
	"errors"
	"regexp"
	"fmt"
	"os"
	"strings"
)

type SectionName string

//Read the template file in the assets directory
func read_template(path string)(string,error) {
    file, err := os.ReadFile(path+"/assets/latex/template/template.tex")
    if err != nil {
        return "", err
    }
    return string(file), nil
}

//Write the template file in the assets directory
func writeTemplate(path string, template string, name string)error{
    err := os.WriteFile(path+"/assets/latex/output/"+name,
        []byte(template), 0644)
    return err
}

//Create a new empty template in the output dir
func Init_output(name string, root_path string)error{
    template,err := read_template(root_path)
    if err != nil {
        return err
    }
    err = writeTemplate(root_path, template, name+".tex")
    return err
}

//Apply a section to a section type on a latex template
func ApplyToSection(section Section, section_type string, output_path string)(string,error){
    replacements := []string{section.first, section.second, section.third, section.fourth}
    b_template, err := os.ReadFile(output_path)
    //nolint:all
    section_type = strings.Title(section_type)
    if err != nil {
        return "", err
    }
    template := string(b_template)
    switch{

    case section_type == "Professional":
        template = strings.Replace(template,"%EXPERIENCE_SECTIONS%",
            "%EXPERIENCE_SECTIONS%\n"+replace_param(prof_template,NB_P_PROF,replacements),1)
        template = replace_items(template, section.description)

    case section_type == "Project":
        template = strings.Replace(template,"%PROJECTS_SECTIONS%",
            "%PROJECTS_SECTIONS%\n"+replace_param(proj_template,NB_P_PROJ,replacements),1)
        template = replace_items(template, section.description)

    case section_type == "Education":
        template = strings.Replace(template,"%EDUCATION_SECTIONS%",
            "%EDUCATION_SECTIONS%\n"+replace_param(edu_template,NB_P_EDU,replacements),1)

    case section_type == "Skill"://TODO https://github.com/Theo-Hafsaoui/Anemon/issues/1
        template = strings.Replace(template,"%SKILL_SECTIONS%",
            "%SKILL_SECTIONS%\n"+replace_param(sk_template,NB_P_SK,replacements),1)
    default:
        return "",errors.New("Don't know type "+section_type)
    }
    path_name := strings.Split(output_path, "/assets/latex/output/")
    if len(path_name) == 1 {
        return template,errors.New("Trying to save outside of output file, at "+output_path)
    }
    if err != nil{
        return "",err
    }
    err = writeTemplate(path_name[0],template,path_name[1])
    return template,err
}

//Search and replace the number in range of `nb_params` by their replacement
func replace_param(template string, nb_params int, replacements []string)string{
        for  i := 0; i < nb_params; i++ {
            position := fmt.Sprintf("$%d", i+1)
            template = strings.Replace(template,
                position, replacements[i], 1)
        }
    return sanitize(template)
}

func replace_items(template string, section_items []string)string{
    items := ""
    for _,item := range section_items{
        items += strings.Replace(pro_item,"%ITEM%",item,1)
    }
    template = strings.Replace(template,
        "%ITEMS%", items, 1)
    return sanitize(template)
}

//Sanitize the special charactere
func sanitize(template string)(string){
	replacements := []struct {
		pattern     string
		replacement string
	}{
		{`[0-9]\%`, `\\%`},
		{`\*\*(.*?)\*\*`, `\textbf{$1}`},
		{`\*(.*?)\*`, `\emph{$1}`},
		//{`\[(.*?)\]\((.*?)\)`, `\href{$2}{$1}`},
	}
	for _, r := range replacements {
		re := regexp.MustCompile(r.pattern)
		template = re.ReplaceAllString(template, r.replacement)
	}
	return  template
}
