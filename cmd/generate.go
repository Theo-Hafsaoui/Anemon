package cmd

import (
	"anemon/internal/parser"
	"anemon/internal/walker"
	"errors"
	"os"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
    Use:   "generate",
    Short: "Generate a CV",
    Long:  `Generate a CV using the Markdown CV directory in the current work directory`,
    RunE: func(cmd *cobra.Command, args []string) error{
        dir, err := os.Getwd()  
        if err != nil{  
            return err
        }  
        CV,err := getSectionMapFrom(dir)
        if err != nil {
            return err
        }
        err = createLatexCVFrom(dir,CV)
        if err != nil {
            return err
        }
        return nil
    },
}

//Use a CV map created by `getSectionMapFrom` and write for each lang key a latex CV using the given information
func createLatexCVFrom(dir string, CV map[string]map[string]string )(error){
        for lang := range CV{
            err := parser.Init_output(lang+"-CV",dir)
            if err != nil{
                return err
            }
            for sec_name := range CV[lang]{
                sec, err := parser.Parse(CV[lang][sec_name])
                if err != nil {
                    return err
                }
                _,err = parser.ApplyToSection(sec,sec_name,dir+"/assets/latex/output/"+lang+"-CV.tex")
                if err != nil {
                    return err
                }
            }
        }
        return nil
}
//TODO consider struct for this map of map
//Walk throught the CV directory and return a map of lang within each their is a map of section
func getSectionMapFrom(dir string)(map[string]map[string]string,error){
       cv_path := dir+"/cv"
        _, err := os.Stat(cv_path)
        if err != nil{  
            if os.IsNotExist(err) {  return nil, errors.New("No `cv` directory found at:"+cv_path) }
            return nil,err
        }  
        result, err := walker.WalkCV(cv_path)
        if err != nil{  
            return nil,err
        }
        return result,nil
}

func init() {
    rootCmd.AddCommand(generateCmd)
}
