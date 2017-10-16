package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSubstitution(t *testing.T) {
	Variables = map[string]string{
		"a":   "alpha",
		"aaa": "triple",
	}
	assert.Equal(t, "pre alpha post", SubstituteForVars("pre $a post"))
	assert.Equal(t, "alphaalphaalpha", SubstituteForVars("$a$a$a"))
	assert.Equal(t, "triple", SubstituteForVars("$aaa"))
}
