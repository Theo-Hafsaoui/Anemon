package output

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

const COMPILER = "pdflatex"

var REGEX = regexp.MustCompile(`(\d+) page`)

type LatexCompiler struct{}

// Compile the templates into PDFs and return the number of page of the longest one
func (*LatexCompiler) CompileTemplate(root string) (int, error) {
	templates, err := getListOfTemplate(root)
	if err != nil {
		return 0, err
	}
	max_nb_of_page := 0
	page_nb := make(chan int, len(templates))
	var wg sync.WaitGroup
	for _, template := range templates {
		wg.Add(1)
		go compile(template, root, &wg, page_nb)
		pdf_pages_nb := <-page_nb
		slog.Info("Number of pages", "pages", pdf_pages_nb)
		max_nb_of_page = max(pdf_pages_nb, max_nb_of_page)
	}
	close(page_nb)
	wg.Wait()
	return max_nb_of_page, nil
}

// Compile the template into a pdf
func compile(template string, root string, wg *sync.WaitGroup, c_page_nb chan int) {
	defer wg.Done()
	cmd := exec.Command(COMPILER, "-interaction=nonstopmode",
		"-output-directory="+root+"/assets/latex/output", template)
	log, err := cmd.Output()
	if err != nil {
		slog.Warn("error(s) to compile file:" + template)
	}
	page_nb := -1
	log_page := REGEX.FindStringSubmatch(string(log))
	if len(log_page) < 1 {
		slog.Error("failed to compile file didnt get the number of pages:" + template)
		fmt.Println(log_page)
	} else {
		page_nb, err = strconv.Atoi(log_page[1])
		if err != nil {
			slog.Error(err.Error())
		}
	}
	c_page_nb <- page_nb
}

// Return the path of latex file inside the template directory
func getListOfTemplate(root string) ([]string, error) {
	var res []string
	templates, err := os.ReadDir(root + "/assets/latex/output")
	if err != nil {
		slog.Error("failed to read directory because: " + err.Error())
		return res, err
	}
	for _, template := range templates {
		if strings.HasSuffix(template.Name(), ".tex") {
			res = append(res, root+"/assets/latex/output/"+template.Name())
		}
	}
	fmt.Println(res)
	return res, nil
}
