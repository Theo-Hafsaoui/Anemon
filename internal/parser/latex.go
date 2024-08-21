package parser

import "os"

//Apply a section to a section type on a latex template
//func (section section, type section_type, template string)

//to write this template
//func write(template string, name string)
//Todo

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
