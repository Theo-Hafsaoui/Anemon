package parser

import (
	"fmt"
	"os"
	"strings"
)

type SectionName string


//Read the template file in the assets directory
func read_template()(string,error) {
    file, err := os.ReadFile("../../assets/latex/template/template.tex")
    if err != nil {
        return "", err
    }
    return string(file), nil
}

//Write the template file in the assets directory
func writeTemplate(template string, name string)error{
    err := os.WriteFile("../../assets/latex/output/"+name+".tex",
        []byte(template), 0644)
    return err
}

//Apply a section to a section type on a latex template
func applyToSection(section Section, section_type string)(string,error){
    replacements := []string{section.first, section.second, section.third, section.fourth}
    template := ""
    switch{

    case section_type == "Professional":
        template = replace_param(prof_template,NB_P_PROF,replacements)
        template = replace_items(template, section.description)

    case section_type == "Project":
        template = replace_param(proj_template,NB_P_PROJ,replacements)
        template = replace_items(template, section.description)

    case section_type == "Education":
        template = replace_param(edu_template,NB_P_EDU,replacements)

    case section_type == "Skill"://TODO https://github.com/Theo-Hafsaoui/Anemon/issues/1
        template = replace_param(sk_template,NB_P_SK,replacements)
    }
    return template,nil
}

//Search and replace the number in range of `nb_params` by their replacement
func replace_param(template string, nb_params int, replacements []string)string{
        for  i := 0; i < nb_params; i++ {
            position := fmt.Sprintf("%d", i+1)
            template = strings.Replace(template,
                position, replacements[i], 1)
        }
    return template
}

func replace_items(template string, section_items []string)string{
    items := ""
    for _,item := range section_items{
        items += strings.Replace(pro_item,"%ITEM%",item,1)
    }
    template = strings.Replace(template,
        "%ITEMS%", items, 1)
    return template
}
