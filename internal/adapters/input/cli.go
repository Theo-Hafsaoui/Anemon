package input

import (
	"anemon/internal/adapters/output"
	"anemon/internal/core"
)

// Use the implementation for markdown and latex to generate latex CV from a tree dir of mardown document
func GenerateCVFromMarkDownToLatex(root string) error {
	var builder core.BuilderService = core.BuilderService{}
	builder.SetRoot(root)
	builder.SetSource(&MarkdownSource{})
	builder.SetParamsSource(&YamlSource{})
	builder.SetTemplateReader(&output.LatexReader{})
	builder.SetTemplateProcessor(&output.LatexProccesor{})
	builder.SetCompiler(&output.LatexCompiler{})
	service := builder.GetService()
	return service.GenerateTemplates()
}

// Change the threshold for the regeration of the PDF
func ChangeOverflowThreshold(newThreshold int) {
	core.SetOverflowThreshold(newThreshold)
}
// Print info in the prompt for all the cvs
func PrintAllCvs(root string) {
        source := MarkdownSource{}
        core.GetInfoAllCvs(root,&source)
}

