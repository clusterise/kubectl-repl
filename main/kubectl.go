package main

import (
	"fmt"
	"github.com/fatih/color"
	"os/exec"
	"os"
	"bufio"
)

func sh(shell string) error {
	if verbose {
		color.Yellow(fmt.Sprintf("+ %s\n", shell))
	}
	cmd := exec.Command("/bin/sh", "-c", shell)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

  	err := cmd.Start()
  	if err != nil {
  		return err
	}
  	return cmd.Wait()
}

func shHandler(shell string, outputHandler func(string)) error {
	if verbose {
		color.Yellow(fmt.Sprintf("+ %s\n", shell))
	}
	cmd := exec.Command("/bin/sh", "-c", shell)

	cmd.Stdin = os.Stdin

	stdout, err := cmd.StdoutPipe()
  	if err != nil {
  		return err
  	}
    reader := bufio.NewReader(stdout)

  	cmd.Start()
  	for {
  		line, _, err := reader.ReadLine()
        if err != nil {
            return cmd.Wait()
        }
        outputHandler(string(line))
    }
  	return cmd.Wait()
}

func kubectl(cmd string) string {
	return fmt.Sprintf("kubectl -n %s %s", namespace, cmd)
}
