package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

var (
	aliases []AliasSpec
)

func substituteForAliases(input string) string {
	for _, alias := range aliases {
		input = alias.Alias.ReplaceAllString(input, alias.Replacement)
	}
	return input
}

type AliasSpec struct {
	Alias       *regexp.Regexp
	Replacement string
}

func loadAliasesFromFile(file string) error {
	content, err  := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	delimiter := regexp.MustCompile(`\s+`)

	lines := strings.Split(string(content), "\n")
	for n, line := range lines {
		if line == "" {
			continue
		}

		parts := delimiter.Split(line, 2)
		if len(parts) != 2 {
			continue
		}

		alias, err := regexp.Compile(parts[0])
		if err != nil {
			return fmt.Errorf("failed to compile alias line %v '%v': %v", n, line, err)
		}

		aliases = append(aliases, AliasSpec{
			Alias: alias,
			Replacement: parts[1],
		})
	}
	return nil
}
