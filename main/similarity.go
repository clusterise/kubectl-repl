package main

import (
	"github.com/texttheater/golang-levenshtein/levenshtein"
	"math"
)

func ClosestString(value string, targets []string) string {
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
		Value    string
	}{math.MaxInt16, ""}
	for target, distance := range distances {
		if distance < best.Distance {
			best.Distance = distance
			best.Value = target
		}
	}
	return best.Value
}
