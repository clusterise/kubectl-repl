package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"syscall"
)

var (
	cmdPattern   *regexp.Regexp
)

func sh(shell string) error {
	if verbose {
		color.Yellow(fmt.Sprintf("+ %s\n", shell))
	}
	cmd := exec.Command("/bin/sh", "-c", shell)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	trap := make(chan os.Signal, 1)
	signal.Notify(trap, syscall.SIGINT)
	defer close(trap)
	defer signal.Stop(trap)

	err := cmd.Start()
	if err != nil {
		return err
	}

	go func() {
		_, ok := <-trap
		if ok {
			fmt.Println("^C")
			cmd.Process.Kill()
		}
	}()

	return cmd.Wait()
}

func shHandler(shell string, outputHandler func(string)) error {
	if verbose {
		color.Yellow(fmt.Sprintf("+ %s\n", shell))
	}
	cmd := exec.Command("/bin/sh", "-c", shell)

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	reader := bufio.NewReader(stdout)

	trap := make(chan os.Signal, 1)
	signal.Notify(trap, syscall.SIGINT)
	defer close(trap)
	defer signal.Stop(trap)

	cmd.Start()
	go func() {
		_, ok := <-trap
		if ok {
			fmt.Println("^C")
			cmd.Process.Kill()
		}
	}()

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err != io.EOF {
				return err
			}
			return cmd.Wait()
		}
		outputHandler(string(line))
	}
}

func kubectl(cmd string) string {
	// Kubectl plugins are only invoked only if all arguments are specified
	// after the command name itself.
	//   kubectl foo --namespace=foo # works
	//   kubectl --namespace=foo foo # fails
	// https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/
	//
	// The `cmd` given can contain a very complex bash input. Instead of parsing it as a whole,
	// we will find the first control character and place the arguments just before it.
	//
	// This fails on some edge case inputs, such as `foo bar=z`, but currently this syntax always
	// takes an option such as `get pods -l foo=bar`, so the `-l` would be parsed first.

	if cmdPattern == nil {
		cmdPattern = regexp.MustCompile(`(\s-|--|[[^\w\s]&&[^.-]])`)
	}

	splitAt := cmdPattern.FindStringIndex(cmd)
	if splitAt == nil {
		splitAt = []int{len(cmd), len(cmd)}
	}
	cmdA := cmd[:splitAt[0]]
	cmdB := cmd[splitAt[0]:]

	buffer := bytes.NewBufferString("kubectl ")
	buffer.WriteString(cmdA)
	if context != "" {
		buffer.WriteString(" --context=")
		buffer.WriteString(context)
	}
	if namespace != "" {
		buffer.WriteString(" --namespace=")
		buffer.WriteString(namespace)
	}
	buffer.WriteString(" ")
	buffer.WriteString(cmdB)
	return buffer.String()
}
