package main

var commands = []builtin{
	builtinExit{},
	builtinNamespace{},
	builtinShell{},
	builtinGet{},
}

type builtin interface {
	filter(command string) bool
	run(command string) error
}
