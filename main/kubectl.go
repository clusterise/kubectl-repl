package main

import (
	"fmt"
	"github.com/fatih/color"
	"os/exec"
	"os"
	"bufio"
)

func sh(shell string, outputHandler func(string)) error {
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

func kubectlSh(cmd string, outputHandler func(string)) error {
	shell := fmt.Sprintf("kubectl -n %s %s", namespace, cmd)
	return sh(shell, outputHandler)
}
