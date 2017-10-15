package main

import (
	"bufio"
	"os"
	"fmt"
	"strings"
	"log"
)

var (
	Input *bufio.Reader
	Namespace string
	Variables map[string]string
)

func prompt(text string) (string, error) {
	fmt.Print(text + " ")
	response, err := Input.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.Trim(response, "\n"), nil
}

func pickNamespace() error {
	namespaces, err := GetNamespaces()
	if err != nil {
		return err
	}

	targets := make([]string, len(namespaces.Items))
	for num, ns := range namespaces.Items {
		fmt.Printf("$%d\t %s\n", num, ns.Name)
		targets[num] = ns.Name
	}

	response, err := prompt("Select namespace:")
	if err != nil {
		return err
	}
	Namespace = ClosestString(response, targets)
	return nil
}

func switchNamespace(ns string) error {
	namespaces, err := GetNamespaces()
	if err != nil {
		return err
	}

	targets := make([]string, len(namespaces.Items))
	for num, ns := range namespaces.Items {
		targets[num] = ns.Name
	}
	Namespace = ClosestString(ns, targets)
	return nil
}

func repl() error {
	command, err := prompt(Namespace)
	if err != nil {
		return err
	}

	for from, to := range Variables {
		command = strings.Replace(command, from, to, -1)
	}

	parts := strings.Split(command, " ")
	if parts[0] == "namespace" || parts[0] == "ns" {
		if len(parts) > 1 {
			switchNamespace(parts[1])
		} else {
			pickNamespace()
		}
		return nil
	}

	output, err := KubectlSh(command)

	if strings.HasPrefix(command, "get") {
		variableIndex := 0
		for _, line := range strings.Split(output, "\n") {
			if strings.HasPrefix(line, "NAME ") {
				fmt.Printf("   \t%s\n", line)
			} else {
				variableIndex++
				fmt.Printf("$%v \t%s\n", variableIndex, line)
			}
			key := fmt.Sprintf("$%d", variableIndex)
			Variables[key] = strings.Split(line, " ")[0]
		}
	} else {
		fmt.Println(output)
	}
	return err
}

func assert(v interface{}) {
	if v != nil {
		log.Fatal(v)
	}
}

func main() {
	Variables = make(map[string]string)
	Input = bufio.NewReader(os.Stdin)
	assert(KubernetesSetup())
	assert(pickNamespace())

	for {
		repl()
	}
}
