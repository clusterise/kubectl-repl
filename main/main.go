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

const (
	versionString = "kubectl-repl {{{VERSION}}}"
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

func printIndexedLine(index, line string) {
	coloredIndex := color.New(color.FgBlue).Sprintf("$%s", index)
	fmt.Printf("%s \t%s\n", coloredIndex, line)
}

func repl() error {
	command, err := prompt("# " + namespace)
	if err != nil {
		return err
	}

	for _, builtin := range commands {
		if builtin.filter(command) {
			return builtin.run(command)
		}
	}

	return sh(kubectl(command))
}

func main() {
	var version bool
	flag.BoolVar(&verbose, "verbose", false, "Verbose")
	flag.BoolVar(&version, "version", false, "Print current version")
	flag.Parse()

	if version {
		fmt.Println(versionString)
		return
	}

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
