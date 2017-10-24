package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClosestString(t *testing.T) {
	assert.Equal(t, "alpha", closestString("al", []string{"aX", "alpha"}),
		"should prefer insertion over replacement and deletion")

	assert.Equal(t, "alpha", closestString("aXpha", []string{"aX", "alpha"}),
		"should correct typos")

	assert.Equal(t, "alpha", closestString("alph", []string{"alpha", "alphafoo"}),
		"should prefer shorter prefix match")

	assert.Equal(t, "alpha", closestString("alpha", []string{"alpha", "alphafoo"}),
		"should prefer exact match to prefix")
}
