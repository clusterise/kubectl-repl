package main

import (
	"os"
)

type builtinExit struct{}

func (b builtinExit) init() error {
	return nil
}

func (b builtinExit) filter(command string) bool {
	return command == "exit" || command == "quit"
}

func (b builtinExit) run(command string) error {
	os.Exit(0)
	return nil
}
