package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
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
			cmd.Process.Kill()
		}
	}()

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
