package main

import (
	"regexp"
)

var (
	variables map[string][]string
	pattern   *regexp.Regexp
)

func substituteForVars(text string) (string, error) {
	if pattern == nil {
		pattern = regexp.MustCompile(`[$]\w+\b`)
	}

	var err error
	output := pattern.ReplaceAllStringFunc(text, func(match string) string {
		if val, ok := variables[match[1:]]; ok {
			return val[0]
		}
		return match
	})
	return output, err
}
