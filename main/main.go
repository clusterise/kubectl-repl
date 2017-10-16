package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"io"
	"log"
	"os"
	"strings"
)

var (
	input     *bufio.Reader
	namespace string
	verbose   bool
)

func prompt(text string) (string, error) {
	fmt.Print(color.New(color.Bold).Sprintf(text + " "))
	line, err := input.ReadString('\n')
	if err != nil {
		return "", err
	}
	response := strings.Trim(line, "\n")
	return substituteForVars(response), nil
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

func printIndexedLine(index, line string) {
	coloredIndex := color.New(color.FgBlue).Sprintf("$%s", index)
	fmt.Printf("%s \t%s\n", coloredIndex, line)
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

func repl() error {
	command, err := prompt("# " + namespace)
	if err != nil {
		return err
	}

	if strings.HasPrefix(command, ";") {
		output, err := sh(strings.TrimPrefix(command, ";"))
		if output != "" {
			fmt.Println(output)
		}
		return err
	}

	parts := strings.Split(command, " ")
	if parts[0] == "exit" || parts[0] == "quit" {
		os.Exit(0)
	}
	if parts[0] == "namespace" || parts[0] == "ns" {
		if len(parts) > 1 {
			switchNamespace(parts[1])
		} else {
			pickNamespace()
		}
		return nil
	}

	output, err := kubectlSh(command)
	if output == "" {
		return err
	}

	if err == nil && strings.HasPrefix(command, "get") {
		variableIndex := 0
		for _, line := range strings.Split(output, "\n") {
			if strings.HasPrefix(line, "NAME ") {
				fmt.Printf("   \t%s\n", line)
			} else {
				variableIndex++
				key := fmt.Sprintf("%d", variableIndex)
				printIndexedLine(key, line)
			}
			key := fmt.Sprintf("%d", variableIndex)
			variables[key] = strings.Split(line, " ")[0]
		}
	} else {
		fmt.Println(output)
	}
	return err
}

func main() {
	flag.BoolVar(&verbose, "verbose", false, "Verbose")
	flag.Parse()

	variables = make(map[string]string)
	input = bufio.NewReader(os.Stdin)

	err := kubernetesSetup()
	if err != nil {
		log.Fatal(err)
	}

	err = pickNamespace()
	if err == io.EOF {
		return
	} else if err != nil {
		log.Fatal(err)
	}

	for {
		err = repl()
		if err == io.EOF {
			break
		}
	}
}
