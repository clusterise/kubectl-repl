package main

import (
	"bufio"
	"os"
	"fmt"
	"strings"
	"log"
	"flag"
	"github.com/fatih/color"
)

var (
	Input *bufio.Reader
	Namespace string
	Variables map[string]string
	Verbose bool
)

func prompt(text string) (string, error) {
	fmt.Print(color.New(color.Bold).Sprintf(text + " "))
	line, err := Input.ReadString('\n')
	if err != nil {
		return "", err
	}
	response := strings.Trim(line, "\n")
	for from, to := range Variables {
		response = strings.Replace(response, from, to, -1)
	}
	return response, nil
}

func namespaceSelector(selector func([]string)(string, error)) error {
	namespaces, err := GetNamespaces()
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
	Namespace = ClosestString(response, targets)
	return nil
}

func printIndexedLine(index, line string) {
	coloredIndex := color.New(color.FgBlue).Sprintf("%s", index)
	fmt.Printf("%s \t%s\n", coloredIndex, line)
}

func pickNamespace() error {
	return namespaceSelector(func(namespaces []string) (string, error) {
		for n, ns := range namespaces {
			key := fmt.Sprintf("$%d", n)
			Variables[key] = ns
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

func repl() error {
	command, err := prompt("# " + Namespace)
	if err != nil {
		return err
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
				printIndexedLine(fmt.Sprintf("$%v", variableIndex), line)
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
	flag.BoolVar(&Verbose, "verbose", false, "Verbose")
	flag.Parse()

	Variables = make(map[string]string)
	Input = bufio.NewReader(os.Stdin)
	assert(KubernetesSetup())
	assert(pickNamespace())

	for {
		repl()
	}
}
