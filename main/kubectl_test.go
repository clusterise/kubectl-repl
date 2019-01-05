package main

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

const SPLIT_MARK = "\x11"

func assertInsertion(t *testing.T, input string) {
	cmd := strings.Replace(input, SPLIT_MARK, "", 1)
	expected := "kubectl " + strings.Replace(input, SPLIT_MARK, " --context=c --namespace=ns", 1)

	assert.Equal(t, expected, kubectl(cmd))
}

func TestKubectlInsertion(t *testing.T) {
	context = "c"
	namespace = "ns"

	assertInsertion(t, "get pods\x11")
	assertInsertion(t, "get pods\x11 -l foo=bar")
	assertInsertion(t, "\x11; echo foo")
	assertInsertion(t, "get pods\x11 -l foo | grep bar")
	assertInsertion(t, "\x11 -l foo get pods")
	assertInsertion(t, "get pods\x11|grep foo")
	assertInsertion(t, "describe node ip-172-23-76-75.eu-central-1.compute.internal\x11")
}
