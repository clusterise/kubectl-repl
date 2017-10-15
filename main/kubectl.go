package main

import (
	"os/exec"
)

func KubectlSh(cmd string) (string, error) {
	process := exec.Command("/bin/sh", "-c", "kubectl " + cmd)
	bytes, err := process.CombinedOutput()
	return string(bytes), err
}
