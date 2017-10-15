package main

import (
	"regexp"
)

var (
	Variables map[string]string
	pattern *regexp.Regexp
)

func SubstituteForVars(text string) string {
	if pattern == nil {
		pattern = regexp.MustCompile(`[$]\w+\b`)
	}

	return pattern.ReplaceAllStringFunc(text, func(match string) string {
		if val, ok := Variables[match[1:]]; ok {
			return val
		}
		return match
	})
}
