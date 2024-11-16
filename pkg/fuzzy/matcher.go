package fuzzy

import (
	"slices"
	"sort"

	"github.com/agnivade/levenshtein"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

func FindBestMatch(source string, targets []string) string {
	// Compute fuzzy matches
	fuzzyMatches := fuzzy.RankFindFold(source, targets)
	sort.Sort(fuzzyMatches)

	// Compute levenshtein distance
	bestMatch := ""
	bestDistance := -1
	for _, target := range targets {
		distance := levenshtein.ComputeDistance(source, target)
		// If the distance is greater than half the length of the target string, skip it
		if distance > slices.Min([]int{len(source) / 2, len(target) / 2}) {
			continue
		}
		if bestDistance == -1 || distance < bestDistance {
			bestMatch = target
			bestDistance = distance
		}
	}

	// Return the best match
	if bestDistance == -1 {
		return ""
	}
	if len(fuzzyMatches) != 0 {
		return fuzzyMatches[0].Target
	}
	return bestMatch
}
