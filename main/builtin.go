package main

type Commands []builtin

func (cmds Commands) Init() error {
	for _, cmd := range cmds {
		err := cmd.init()
		if err != nil {
			return err
		}
	}
	return nil
}

type builtin interface {
	init() error
	filter(command string) bool
	run(command string) error
}
