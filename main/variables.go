package main

import (
	"regexp"
)

var (
	variables map[string]string
	pattern   *regexp.Regexp
)

func substituteForVars(text string) string {
	if pattern == nil {
		pattern = regexp.MustCompile(`[$]\w+\b`)
	}

	return pattern.ReplaceAllStringFunc(text, func(match string) string {
		if val, ok := variables[match[1:]]; ok {
			return val
		}
		return match
	})
}
