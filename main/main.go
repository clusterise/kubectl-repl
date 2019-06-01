package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/chzyer/readline"
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
	context   string
	verbose   bool
	rl *readline.Instance
)

func prompt() (string, error) {
	color.New(color.Bold).Print("# ")
	prompt := bytes.NewBufferString("")

	if context != "" {
		color.New(color.FgBlack, color.Italic).Fprint(prompt, context)
		fmt.Fprint(prompt, " ")
	}

	if namespace != "" {
		color.New(color.Bold).Fprint(prompt, namespace)
	} else {
		color.New(color.Bold).Fprint(prompt, "namespace")
	}
	fmt.Fprint(prompt, " ")

	rl.SetPrompt(prompt.String())
	line, err := rl.Readline()
	if err != nil {
		return "", err
	}
	response := strings.TrimSpace(line)
	expandedCmd, err := substituteForVars(response)
	if err != nil {
		return "", err
	}
	err = rl.Operation.SaveHistory(expandedCmd)
	return expandedCmd, err
}

func printIndexedLine(index, line string) {
	coloredIndex := color.New(color.FgBlue).Sprintf("$%s", index)
	fmt.Printf("%s \t%s\n", coloredIndex, line)
}

func repl(commands Commands) error {
	command, err := prompt()
	if err != nil {
		return err
	}

	if command == "" {
		return nil
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
	flag.StringVar(&context, "context", "", "Override current context")
	flag.StringVar(&namespace, "namespace", "", "Override current context namespace")
	flag.Parse()

	if namespace == "" {
		namespace = os.Getenv("KUBECTL_NAMESPACE")
	}

	if version {
		fmt.Println(versionString)
		return
	}

	var err error
	rl, err = readline.New("_ ")
	rl.Config.DisableAutoSaveHistory = true
	if err != nil {
		log.Fatal(err)
	}
	defer rl.Close()

	commands := Commands{
		builtinExit{},
		builtinNamespace{},
		builtinShell{},
		builtinGet{},
	}
	err = commands.Init()
	if err != nil {
		log.Fatal(err)
	}

	variables = make(map[string][]string)
	input = bufio.NewReader(os.Stdin)

	if namespace == "" {
		err = pickNamespace()
		if err == io.EOF {
			return
		} else if err != nil {
			log.Fatal(err)
		}
	}

	for {
		err = repl(commands)
		if err == io.EOF {
			break
		} else if err != nil {
			color.New(color.FgRed).Println(err)
		}
	}
}
