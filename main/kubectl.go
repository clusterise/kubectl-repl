package main

import (
	"os/exec"
	"strings"
)

func KubectlSh(cmd string) (string, error) {
	process := exec.Command("/bin/sh", "-c", "kubectl " + cmd)
	bytes, err := process.CombinedOutput()
	return strings.TrimRight(string(bytes), "\n"), err
}
