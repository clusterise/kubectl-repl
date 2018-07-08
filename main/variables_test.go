package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func assertSubstituion(t *testing.T, expected, input string) {
	out, err := substituteForVars(input)
	assert.Nil(t, err)
	assert.Equal(t, expected, out)
}

func TestSubstitution(t *testing.T) {
	variables = map[string][]string{
		"a":   {"alpha"},
		"aaa": {"triple"},
	}
	assertSubstituion(t, "pre alpha post", "pre $a post")
	assertSubstituion(t, "alphaalphaalpha", "$a$a$a")
	assertSubstituion(t, "triple", "$aaa")
}

func TestSubstitutionIndex(t *testing.T) {
	variables = map[string][]string{
		"a":   {"alpha", "alpha2", "alpha3"},
		"aaa": {"triple"},
	}
	assertSubstituion(t, "pre alpha post", "pre $a:1 post")
	assertSubstituion(t, "alphaalpha2alpha3", "$a:1$a:2$a:3")
	assertSubstituion(t, "triple", "$aaa")
}
