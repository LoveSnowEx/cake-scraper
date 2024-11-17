package location

import (
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

const compareLength = 3
const ratio = 2.0 / 3

func evalScore(targetTokens, locationTokens []string) float64 {
	score := 0.0
	remainRatio := 1.0
	for i := 0; i < compareLength; i++ {
		var targetToken, locationToken string
		if i < len(targetTokens) {
			targetToken = targetTokens[len(targetTokens)-1-i]
		}
		if i < len(locationTokens) {
			locationToken = locationTokens[len(locationTokens)-1-i]
		}
		sum := len(targetToken) + len(locationToken)
		r := remainRatio * ratio
		if sum == 0 {
			score += r
			remainRatio -= r
			continue
		}
		dis := fuzzy.LevenshteinDistance(targetToken, locationToken)
		score += float64(sum-dis) / float64(sum) * r
		remainRatio -= r
	}
	return score
}

func FindBestMatch(target string) *Location {
	targetTokens := strings.Split(target, ", ")
	locations := LoadLocations()
	maxScore := 0.0
	var bestMatch *Location
	for _, location := range locations {
		locationTokens := strings.Split(location.Address(), ", ")
		score := evalScore(targetTokens, locationTokens)
		if score > maxScore {
			maxScore = score
			bestMatch = location
		}
	}
	if maxScore >= ratio {
		return bestMatch
	}
	return nil
}
