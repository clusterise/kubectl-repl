package main

import (
	"strings"
)

type builtinShell struct{}

func (b builtinShell) init() error {
	return nil
}

func (b builtinShell) filter(command string) bool {
	return strings.HasPrefix(command, ";") || strings.HasPrefix(command, "#")
}

func (b builtinShell) run(command string) error {
	return sh(strings.TrimLeft(command, "; "))
}
