package markuplanguages

import (
    "fmt"
    "reflect"
    "testing"
)

func TestParse(t *testing.T) {
    t.Run("Happy path should return the wanted struct", func (t *testing.T) {
        input := `
# Title
## Skill
### Date
#### Url
Item
Item2`
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


        description := []string{"Item","Item2"}
        if !reflect.DeepEqual(result.description, description){
            fmt.Printf("want: %q, got: %q\n", description, result.description)
            t.Fatalf("want %s got %s", description, result.description)
        }
    })

    t.Run("should return an error if title is badly formarted", func (t *testing.T) {
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

    })
    t.Run("should return an error if outside of allowed nb of #", func (t *testing.T) {
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
    })
t.Run("should return an error with mutiples section", func (t *testing.T) {
        input := `
# Title
## Skill
### Date
#### Url
Description

# Title2
## Skill2
### Date2
#### Url2
DescriptionTwo`
        got, err := Parse(input)
        want := "Tried to parse mutiple paragraph into a single section"
        fmt.Println(got)

        if err == nil {
            t.Fatalf("expected an error")
        }
        if err.Error() != want  {
            t.Fatalf("should return a formating error, got %s should %s",err.Error(),want)
        }
    })
}
