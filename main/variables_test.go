package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSubstitution(t *testing.T) {
	variables = map[string][]string{
		"a":   {"alpha"},
		"aaa": {"triple"},
	}
	assert.Equal(t, "pre alpha post", substituteForVars("pre $a post"))
	assert.Equal(t, "alphaalphaalpha", substituteForVars("$a$a$a"))
	assert.Equal(t, "triple", substituteForVars("$aaa"))
}
