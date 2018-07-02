package main

import (
	"fmt"
	"strings"
	"regexp"
)

var (
	outputRegexp *regexp.Regexp
)

type builtinGet struct{}

// Apply to all "get" commands, ignoring flags
func (b builtinGet) filter(command string) bool {
	if outputRegexp == nil {
		outputRegexp, _ = regexp.Compile(`^([^|]*)(-o|--output)(\s*=\s*|\s+)(json|yaml)`)
	}

	return strings.HasPrefix(command, "get") && !outputRegexp.MatchString(command)
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
		variables[key] = strings.Split(line, " ")[0]
	})
}
