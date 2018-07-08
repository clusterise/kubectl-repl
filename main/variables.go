package main

import (
	"regexp"
	"strconv"
	"fmt"
)

var (
	variables map[string][]string
	pattern   *regexp.Regexp
)

const (
	groupVar = 1
	groupColumn = 3
)

func substituteForVars(text string) (string, error) {
	if pattern == nil {
		pattern = regexp.MustCompile(`[$](\w+)(:(\d+))?\b`)
	}

	var outerError error
	output := pattern.ReplaceAllStringFunc(text, func(match string) string {
		parts := pattern.FindStringSubmatch(match)
		varName := parts[groupVar]
		if val, ok := variables[varName]; ok {
			col := parts[groupColumn]
			if col == "" {
				return val[0]
			}
			colInt, err := strconv.ParseInt(col, 10, 64)
			if err != nil {
				outerError = fmt.Errorf("could not parse %v as column index", match)
				return "(ERROR)"
			}
			index := int(colInt - 1)
			if index < 0 || index >= len(val) {
				outerError = fmt.Errorf("variable $%v has no value in column %v", varName, col)
				return "(ERROR)"
			}
			return val[index]
		}
		return match
	})
	return output, outerError
}
