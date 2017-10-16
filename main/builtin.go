package main

var commands = []builtin{
	builtinExit{},
	builtinNamespace{},
}

type builtin interface {
	filter(command string) bool
	run(command string) error
}
