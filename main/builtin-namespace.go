package main

import (
	"fmt"
	"strings"
	"encoding/json"
	"os/exec"
	"bytes"
)

type builtinNamespace struct{}

func (b builtinNamespace) init() error {
	return nil
}

func (b builtinNamespace) filter(command string) bool {
	parts := strings.Split(command, " ")
	return parts[0] == "namespace" || parts[0] == "ns"
}

func (b builtinNamespace) run(command string) error {
	parts := strings.Split(command, " ")
	if len(parts) > 1 {
		return switchNamespace(parts[1])
	}
	return pickNamespace()
}

type KubernetesV1NamespaceList struct {
	Items []KubernetesV1Namespace `json:"items"`
}
type KubernetesV1Namespace struct {
	Metadata struct {
		Name string `json:"name"`
	} `json:"metadata"`
}

func namespaceSelector(selector func([]string) (string, error)) error {
	cmd := exec.Command("/bin/sh", "-c", kubectl("get namespaces --output=json"))
	var jsonOut bytes.Buffer
	var cmdErr bytes.Buffer
	cmd.Stdout = &jsonOut
	cmd.Stderr = &cmdErr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%v: %v", err, cmdErr.String())
	}

	var namespaces KubernetesV1NamespaceList
	err = json.Unmarshal(jsonOut.Bytes(), &namespaces)
	if err != nil {
		return err
	}

	targets := make([]string, len(namespaces.Items))
	for num, ns := range namespaces.Items {
		targets[num] = ns.Metadata.Name
	}

	response, err := selector(targets)
	if err != nil {
		return err
	}
	namespace = closestString(response, targets)
	return nil
}

func pickNamespace() error {
	return namespaceSelector(func(namespaces []string) (string, error) {
		for n, ns := range namespaces {
			key := fmt.Sprintf("%d", n + 1)
			variables[key] = []string{ns}
			printIndexedLine(key, ns)
		}
		return prompt()
	})
}

func switchNamespace(ns string) error {
	return namespaceSelector(func(namespaces []string) (string, error) {
		return ns, nil
	})
}
