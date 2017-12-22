package main

import (
	"fmt"
	"strings"
	"regexp"
	"github.com/k0kubun/go-ansi"
	"io/ioutil"
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

	lines := make(map[string]int, 0)

	return shHandler(kubectl(command), func(line string) {
		firstColumn := strings.Split(line, " ")[0]
		lineOffset, lineDefined := lines[firstColumn]

		if lineDefined {
			ansi.CursorPreviousLine(lineOffset)
			ansi.EraseInLine(2) // clear entire line
			fmt.Printf("   \t%s\n", line)

		} else {
			if strings.HasPrefix(line, "NAME ") {
				fmt.Printf("   \t%s\n", line)

			} else {
				variableIndex++
				key := fmt.Sprintf("%d", variableIndex)

				// shift existing lines
				for k, _ := range lines {
					lines[k] += 1
				}

				lines[firstColumn] = 1
				printIndexedLine(key, line)
			}
		}

		ioutil.WriteFile("/dev/ttys012", []byte(fmt.Sprintf("%#v\n", lines)), 777)
		ansi.CursorNextLine(999) // reset cursor

		// save variable
		key := fmt.Sprintf("%d", variableIndex)
		variables[key] = firstColumn
	})
}
