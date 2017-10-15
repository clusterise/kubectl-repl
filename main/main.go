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
)

func prompt(text string) (string, error) {
	fmt.Printf("%s: ", text)
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

	for num, ns := range namespaces.Items {
		fmt.Printf("$%d\t %s\n", num, ns.Name)
	}
	_, err = prompt("Select namespace")
	return err
}

func assert(v interface{}) {
	if v != nil {
		log.Fatal(v)
	}
}

func main() {
	Input = bufio.NewReader(os.Stdin)
	assert(KubernetesSetup())
	assert(pickNamespace())
}
