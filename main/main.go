package main

import (
	"bufio"
	"os"
	"fmt"
	"strings"
	"log"
	"github.com/texttheater/golang-levenshtein/levenshtein"
	"math"
)

var (
	Input *bufio.Reader
	Namespace string
	Variables map[string]string
)

func prompt(text string) (string, error) {
	fmt.Printf("%s: ", text)
	response, err := Input.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.Trim(response, "\n"), nil
}

func mostSimilar(value string, targets []string) string {
	valueRunes := []rune(value)
	ops := levenshtein.Options{
		InsCost: 0,
		SubCost: 5,
		DelCost: 10,
		Matches: levenshtein.DefaultOptions.Matches,
	}

	distances := make(map[string]int, len(targets))
	for _, target := range targets {
		distances[target] = levenshtein.DistanceForStrings(valueRunes, []rune(target), ops)
	}

	best := struct {
		Distance int
		Value string
	}{math.MaxInt64,""}
	for target, distance := range distances {
		if distance < best.Distance {
			best.Distance = distance
			best.Value = target
		}
	}
	return best.Value
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

	response, err := prompt("Select namespace")
	if err != nil {
		return err
	}
	Namespace = mostSimilar(response, targets)
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

	output, err := KubectlSh(command)

	if strings.HasPrefix(command, "get") {
		for index, line := range strings.Split(output, "\n") {
			key := fmt.Sprintf("$%d", index)
			Variables[key] = strings.Split(line, " ")[0]

			if strings.HasPrefix(line, "NAME ") {
				fmt.Printf("   \t%s\n", line)
			} else {
				fmt.Printf("$%v \t%s\n", index, line)
			}
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
		assert(repl())
	}
}
