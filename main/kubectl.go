package main

import (
	"os/exec"
)

func Kubectl(args ...string) ([]byte, error) {
	cmd := exec.Command("kubectl", args...)
	return cmd.CombinedOutput()
}
