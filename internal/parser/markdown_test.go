package parser

import (
    "testing"
)

func TestParseHappyPath(t *testing.T) {
    input := `
# Title
## Skill
### Date
#### Url
Description`
    result, err := Parse(input)
    if err != nil {
        t.Fatalf("Unexecpted eroor :%s", err.Error())
    }

    first := "Title"
    if result.first != first {
        t.Fatalf("want %s got %s", first, result.first)
    }


    second := "Skill"
    if result.second != second {
        t.Fatalf("want %s got %s", second, result.second)
    }

    third := "Date"
    if result.third != third {
        t.Fatalf("want %s got %s", third, result.third)
    }

    fourth := "Url"
    if result.fourth != fourth {
        t.Fatalf("want %s got %s", fourth, result.fourth)
    }


    description := "Description"
    if result.description != description {
        t.Fatalf("want %s got %s", description, result.description)
    }
}

func TestParseIfTitleIsBadlyFormated(t *testing.T) {
    input := `
#Title
## Skill
### Date
#### Url
Description`
    _, err := Parse(input)
    want := "Err: cannot parse this md line{#Title}  # should be followed by space"
    
    if err == nil {
        t.Fatalf("expected an error")
    }

    if err.Error() != want  {
        t.Fatalf("should return a formating error, got %s should %s",err.Error(),want)
    }

}
func TestParseIfIsBadlyFormated(t *testing.T) {
    input := `
# Title
## Skill
### Date
#### Url
##### oups
Description`
    _, err := Parse(input)
    want := "Err: cannot parse this md line{##### oups}"
    
    if err == nil {
        t.Fatalf("expected an error")
    }
    if err.Error() != want  {
        t.Fatalf("should return a formating error, got %s should %s",err.Error(),want)
    }
}
