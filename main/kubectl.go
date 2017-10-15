package main

import (
	"os/exec"
	"strings"
	"fmt"
)

func KubectlSh(namespace string, cmd string) (string, error) {
	shell := fmt.Sprintf("kubectl -n %s %s", namespace, cmd)
	fmt.Printf("+ %s\n", shell)
	process := exec.Command("/bin/sh", "-c", shell)
	bytes, err := process.CombinedOutput()
	return strings.TrimRight(string(bytes), "\n"), err
}
