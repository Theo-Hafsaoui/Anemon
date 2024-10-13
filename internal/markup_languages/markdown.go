package markuplanguages

import (
	"errors"
        "regexp"
	"fmt"
	"strings"
)

/*
Section represents a parsed Markdown section with up to four heading levels and a description.
*/
type Section struct {
    first string
    second  string
    third  string
    fourth  string
    description  []string
}

func (s Section) String() string {
    return fmt.Sprintf("1: %s\n2: %s\n3: %s\n4: %s\nitems: %v",
        s.first,s.second,s.third,s.fourth,s.description)
}

/*
Parse parses a Markdown-like `paragraph` into a `Section`,
extracting headings and description based on the number of leading hashtags or stars.
Returns an error if the format is invalid.
*/
func Parse(paragraph string) (Section,error){
    section := Section{}
    if len(strings.Split(paragraph, "\n\n")) > 1{
        return section, errors.New("Tried to parse mutiple paragraph into a single section")
    }
    hashtag_regex, _ := regexp.Compile("^#+")
    wasASkill := false
    for _, line := range strings.Split(strings.TrimRight(paragraph, "\n"), "\n") {
        nb_hashtag := len(hashtag_regex.FindString(line))

        if len(line) == 0{
            continue
        }

        if wasASkill {
            wasASkill = false
            section.second= strings.TrimLeft(line,"- ")
        }

        if nb_hashtag == 0 && string(line[0])=="*" && len(strings.Trim(line,"*")) == len(line) - 4 {//Trim should **tt** -> tt
            section.first= strings.Trim(line,"*")
            wasASkill = true
        }


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
        case nb_hashtag == 0 && len(line)>1:
            items := strings.Split(line, "\n")
            section.description = append(section.description, items...)
        case nb_hashtag > 4:
            return section, errors.New("Err: cannot parse this md line{"+line+"}")
    }
    }
    return section, nil
}
