package main

import (
	"os/exec"
	"strings"
	"fmt"
)

func Kubectl(args ...string) ([]byte, error) {
	cmd := exec.Command("kubectl", args...)
	return cmd.CombinedOutput()
}

func GetNamespaces() ([]string, error) {
	response, err := Kubectl("get", "namespaces", "-o", "jsonpath='{.items[*].metadata.name}'")
	if err != nil {
		return nil, err
	}
		fmt.Printf("%#v", response)

	return strings.Split(string(response), " "), nil
}
