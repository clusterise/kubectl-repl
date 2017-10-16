package main

import (
	"fmt"
	"github.com/fatih/color"
	"os/exec"
	"strings"
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
