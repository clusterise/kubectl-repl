package main

import (
	"os/exec"
	"strings"
	"fmt"
	"github.com/fatih/color"
)

func KubectlSh(cmd string) (string, error) {
	shell := fmt.Sprintf("kubectl -n %s %s", Namespace, cmd)

	if Verbose {
		color.Yellow(fmt.Sprintf("+ %s\n", shell))
	}
	process := exec.Command("/bin/sh", "-c", shell)
	bytes, err := process.CombinedOutput()
	return strings.TrimRight(string(bytes), "\n"), err
}
