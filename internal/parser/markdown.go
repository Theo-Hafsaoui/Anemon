package parser

import (
	"errors"
	"regexp"
	"strings"
)

/*
Section represents a parsed Markdown section with up to four heading levels and a description.
with params_nb the max `depth` of the heading
*/
type Section struct {
    first string
    second  string
    third  string
    fourth  string
    description  string
    params_nb  int
}

/*
Parse parses a Markdown-like `paragraph` into a `Section`, extracting headings and description based on the number of leading hashtags. Returns an error if the format is invalid.
*/
func Parse(paragraph string) (Section,error){
    p_counter := map[int]bool{}
    r, _ := regexp.Compile("^#+")
    section := Section{}
    for _, line := range strings.Split(strings.TrimRight(paragraph, "\n"), "\n") {
        nb_hashtag := len(r.FindString(line))
        p_counter[nb_hashtag] = true
        switch{
        case nb_hashtag>0 && line[nb_hashtag] != ' ':
            return section, errors.New("Err: cannot parse this md line{"+line+"}  # should be followed by space")
        case nb_hashtag == 1:
            section.first=line[nb_hashtag+1:]
        case nb_hashtag == 2:
            section.second=line[nb_hashtag+1:]
        case nb_hashtag == 3:
            section.third=line[nb_hashtag+1:]
        case nb_hashtag == 4:
            section.fourth=line[nb_hashtag+1:]
        case nb_hashtag == 0:
            section.description += line
        case nb_hashtag > 4:
            return section, errors.New("Err: cannot parse this md line{"+line+"}")
    }
    }
    section.params_nb = len(p_counter)
    return section, nil
}
