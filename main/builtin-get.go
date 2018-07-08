package main

import (
	"fmt"
	"strings"
	"regexp"
)

var (
	splitRegexp *regexp.Regexp
	outputRegexp *regexp.Regexp
)

type builtinGet struct{}

func (b builtinGet) init() error {
	var err error
	splitRegexp, err = regexp.Compile(`\s+`)
	if err != nil {
		return err
	}
	outputRegexp, err = regexp.Compile(`^([^|]*)(-o|--output)(\s*=\s*|\s+)(json|yaml)`)
	return err
}

// Apply to all "get" commands, ignoring flags
func (b builtinGet) filter(command string) bool {
	return strings.HasPrefix(command, "get") && !outputRegexp.MatchString(command) &&
		!strings.Contains(command, " --help")
}

func (b builtinGet) run(command string) error {
	variableIndex := 0
	return shHandler(kubectl(command), func(line string) {
		if strings.HasPrefix(line, "NAME ") {
			fmt.Printf("   \t%s\n", line)
		} else {
			variableIndex++
			key := fmt.Sprintf("%d", variableIndex)
			printIndexedLine(key, line)
		}
		key := fmt.Sprintf("%d", variableIndex)
		variables[key] = splitRegexp.Split(line, -1)
	})
}
