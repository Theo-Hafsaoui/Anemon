package core

import (
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"sync"
)

var TresholdPageOverFlow struct {
	value int
	mutex sync.Mutex
}

type CVService struct {
	root              string
	source            Source
	paramsSource      SourceParams
	templateReader    TemplateReader
	templateProcessor TemplateProcessor
	compiler          Compiler
}

type BuilderService struct {
	root              string
	source            Source
	paramsSource      SourceParams
	templateReader    TemplateReader
	templateProcessor TemplateProcessor
	compiler          Compiler
}

const DO_NOT_RM_EDUCATION = true

// generate the templates for the cvs defined in the assets directory
func (s *CVService) GenerateTemplates() error {
	slog.Info("--Generating CVs--")
	cvs, err := s.source.GetCVsFrom(s.root)
	if err != nil {
		return fmt.Errorf("failed to get CVs: %w", err)
	}
	params, err := s.paramsSource.GetParamsFrom(s.root)
	if err != nil {
		return fmt.Errorf("failed to get parameters: %w", err)
	}
	generiqueTemplate, err := s.templateReader.ReadCVTemplate(s.root, params)
	if err != nil {
		return fmt.Errorf("failed to read generic template: %w", err)
	}
	if err := s.generateAllCVs(cvs, params, generiqueTemplate, false); err != nil {
		return err
	}
	return s.compileWithOverflowHandling(cvs, params, generiqueTemplate)
}

// Compile the CV into PDF, and if they are too long regenrate theme with one less section
func (s *CVService) compileWithOverflowHandling(cvs []CV, params Params, template string) error {
	threshold := GetOverflowThreshold()
	maxNbPage, err := s.compiler.CompileTemplate(s.root)
	if err != nil {
		return fmt.Errorf("failed to compile template: %w", err)
	}
	for maxNbPage > threshold {
		slog.Warn("Page overflow detected; adjusting layout and regenerating CVs")
		if err := s.generateAllCVs(cvs, params, template, true); err != nil {
			return err
		}
		maxNbPage, err = s.compiler.CompileTemplate(s.root)
		if err != nil {
			return fmt.Errorf("failed to recompile template after adjustment: %w", err)
		}
	}
	return nil
}

// Use the the slice of CV to write in the output directory the new CV template
func (s *CVService) generateAllCVs(cvs []CV, params Params, template string, adjustLayout bool) error {
	for _, cv := range cvs {
		if err := generateCVFrom(cv, params, s.root, template, s.templateProcessor, adjustLayout); err != nil {
			return fmt.Errorf("failed to generate CV: %w", err)
		}
	}
	return nil
}

// generate a template for the cv with the given params
func generateCVFrom(cv CV, params Params, root string,
	template string, processor TemplateProcessor, shouldBeShorter bool) error {
	var err error
	if len(params.Variante) == 0 { //if no variante create simple CV
		params.Variante = map[string][]string{"simple": nil}
	}
	for vari, keywords := range params.Variante {
		cvName := "CV-" + cv.Lang + "-" + vari + ".tex"
		slog.Info("Generating for:" + cvName)
		cvTemplate := template
		if shouldBeShorter {
			removeLowestParagraphe(&cv, keywords)
		}
		for _, section := range cv.Sections {
			for _, paragraph := range section.Paragraphes {
				headers := []string{paragraph.H1, paragraph.H2,
                                                    paragraph.H3, paragraph.H4}
				items := paragraph.Items
				sortByScore(items, keywords)
				cvTemplate, err = processor.ApplySectionToTemplate(
                                cvTemplate, headers, items, section.Title, keywords)
				if err != nil {
					return err
				}
			}
			err = processor.MakeNewTemplate(root, cvTemplate, cvName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Remove the lowest paragraph in the project and experience sections
func removeLowestParagraphe(cv *CV, keywords []string) {
	if len(cv.Sections) == 0 {
		return
	}
	s_i,p_i := getLowestParagraphAcrossSections(*cv, keywords)
        if s_i<0 || p_i < 0{
            slog.Warn("No paragraphe to remove found")
            return
        }
        cv.Sections[s_i].Paragraphes = append(cv.Sections[s_i].Paragraphes[:p_i], cv.Sections[s_i].Paragraphes[p_i+1:]...)
}

// Return the index of the paragraph across all sections with the lowest score
func getLowestParagraphAcrossSections(cv CV, keywords []string) (s_i int, p_i int) {
    minScore := int(^uint(0) >> 1)
    s_i,p_i = -1,-1
    for secIdx, section := range cv.Sections {
        slog.Warn(section.Title)
        if section.Title == "education" && DO_NOT_RM_EDUCATION{
            continue
        }
        parIdx := getLowestParagraphe(section, keywords)
        if parIdx >= 0 {
            currentScore := getScoreParagraphe(section.Paragraphes[parIdx], keywords)
            if currentScore < minScore {
                minScore, s_i, p_i = currentScore, secIdx, parIdx
            }
        }
    }
    return s_i, p_i
}

// Return the index of the paragraph from the section with the lowest score
func getLowestParagraphe(section Section, keywords []string) int {
	return getLowestIndex(section.Paragraphes, keywords, getScoreParagraphe)
}

// generic function to get the index of the element with the lowest score
// TODO ? an interface to avoid the any ?
func getLowestIndex[T any](items []T, keywords []string, getScore func(T, []string) int) int {
        if len(items)==0{
            return -1
        }
	min_score := getScore(items[0], keywords)
	min_idx := 0
	for idx, item := range items[1:] {
		current_score := getScore(item, keywords)
		if current_score < min_score {
			min_idx = idx + 1
			min_score = current_score
		}
	}
	return min_idx
}

// Get items of a single paragraph and sum they score to get global score of the items
func getScoreParagraphe(paragraph Paragraphe, keywords []string) int {
	res := 0
	for _, item := range paragraph.Items {
		res += getScore(item, keywords)
	}
	return res
}

// Sorte a slice of items by the number of keyword
//
// The sort is done in ascending order as the section append work like a stack(Lifo)
func sortByScore(items []string, keywords []string) {
	sort.Slice(items, func(i, j int) bool {
		return getScore(items[i], keywords) > getScore(items[j], keywords)
	})
}

// take and item and a list of keyword and return the number of keyword inside the item
func getScore(item string, keywords []string) int {
	score := 0
	for _, keyword := range keywords {
		score += strings.Count(item, keyword)
	}
	return score
}

//Get all info of the cvs and prompt them in the stdout
func GetInfoAllCvs(root string, source Source){
    cvs, err := source.GetCVsFrom(root)
    if err != nil{
        slog.Warn(err.Error())
    }
    for _,cv := range cvs{
        cv.Print()
    }
}

// Set the threshold value dynamically
func SetOverflowThreshold(newThreshold int) {
	TresholdPageOverFlow.mutex.Lock()
	defer TresholdPageOverFlow.mutex.Unlock()
	TresholdPageOverFlow.value = newThreshold
}

// Get the current threshold value
func GetOverflowThreshold() int {
	TresholdPageOverFlow.mutex.Lock()
	defer TresholdPageOverFlow.mutex.Unlock()
	return TresholdPageOverFlow.value
}

func init() {
	TresholdPageOverFlow.value = 1
}

//--Builder--

func (cv *BuilderService) SetRoot(root string) {
	cv.root = root
}

func (cv *BuilderService) SetSource(source Source) {
	cv.source = source
}

func (cv *BuilderService) SetParamsSource(paramsSource SourceParams) {
	cv.paramsSource = paramsSource
}

func (cv *BuilderService) SetTemplateReader(templateReader TemplateReader) {
	cv.templateReader = templateReader
}

func (cv *BuilderService) SetTemplateProcessor(templateProcessor TemplateProcessor) {
	cv.templateProcessor = templateProcessor
}

func (cv *BuilderService) SetCompiler(compiler Compiler) {
	cv.compiler = compiler
}

func (s *BuilderService) GetService() CVService {
	return CVService{
		root:              s.root,
		source:            s.source,
		paramsSource:      s.paramsSource,
		templateReader:    s.templateReader,
		templateProcessor: s.templateProcessor,
		compiler:          s.compiler,
	}
}
