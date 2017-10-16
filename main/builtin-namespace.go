package main

import (
	"fmt"
	"strings"
)

type builtinNamespace struct{}

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

func namespaceSelector(selector func([]string) (string, error)) error {
	namespaces, err := getNamespaces()
	if err != nil {
		return err
	}

	targets := make([]string, len(namespaces.Items))
	for num, ns := range namespaces.Items {
		targets[num] = ns.Name
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
			key := fmt.Sprintf("%d", n)
			variables[key] = ns
			printIndexedLine(key, ns)
		}
		return prompt("# namespace")
	})
}

func switchNamespace(ns string) error {
	return namespaceSelector(func(namespaces []string) (string, error) {
		return ns, nil
	})
}
