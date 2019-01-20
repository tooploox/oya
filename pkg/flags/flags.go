package flags

import (
	"fmt"
	"regexp"
)

var switchFlagRx = regexp.MustCompile("^--([a-z][-a-zA-z_0-9]*)$")
var valueFlagRx = regexp.MustCompile("^--([a-z][-a-zA-z_0-9]*)=\"?([^\"]*)\"?$")
var maybeFlagRx = regexp.MustCompile("^--.*")

const TrueValue = "true"

type ErrRepeatedFlag struct {
	Arg string
}

func (err ErrRepeatedFlag) Error() string {
	return fmt.Sprintf("Flag %q specified more than once", err.Arg)
}

type ErrInvalidFlag struct {
	Arg string
}

func (err ErrInvalidFlag) Error() string {
	return fmt.Sprintf("Flag syntax invalid: %q", err.Arg)
}

func Parse(args []string) ([]string, map[string]string, error) {
	positionalArgs := make([]string, 0, len(args))
	flags := make(map[string]string)
	for _, arg := range args {
		flag, value, ok, err := parseFlag(arg)
		if err != nil {
			return nil, nil, err
		}
		if ok {
			if _, exists := flags[flag]; exists {
				return nil, nil, ErrRepeatedFlag{
					Arg: arg,
				}
			}
			flags[flag] = value
		} else {
			positionalArgs = append(positionalArgs, arg)
		}
	}
	return positionalArgs, flags, nil
}

func parseFlag(arg string) (string, string, bool, error) {
	if maybeFlagRx.MatchString(arg) {
		if matches := switchFlagRx.FindStringSubmatch(arg); len(matches) == 2 {
			return matches[1], TrueValue, true, nil

		} else if matches := valueFlagRx.FindStringSubmatch(arg); len(matches) == 3 {
			return matches[1], matches[2], true, nil
		}
		return "", "", false, ErrInvalidFlag{Arg: arg}
	}
	return arg, arg, false, nil
}
