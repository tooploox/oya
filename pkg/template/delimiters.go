package template

import(
	"fmt"
	"strings"
)

type Delimiters struct {
	start string
	end   string
}

type invalidDelimitersFormat struct {
	Delimiters string
}

func (e *invalidDelimitersFormat) Error() string {
	return fmt.Sprintf("Invalid template delimiters \"%v\". Use 2 chars, seperate with 3 dots. for ex: \"{{...}}\"", e.Delimiters)
}

func ParseDelimiters(s string) (Delimiters, error) {
	arr := strings.Split(s, "...")
	if len(arr) != 2 {
		return Delimiters{}, &invalidDelimitersFormat{s}
	}
	return Delimiters{arr[0], arr[1]}, nil
}
