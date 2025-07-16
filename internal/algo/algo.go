package algo

import (
	"sort"
	"strings"
)

type Match struct {
	str   string
	score int
}

// fuzzyFind returns matched strings sorted by relevance
func FuzzyFind(needle string, haystack []string) []string {
	if needle == "" {
		return haystack
	}

	// Convert needle to lowercase for case-insensitive matching
	needle = strings.ToLower(needle)
	matches := make([]Match, 0)

	// Score and filter matches
	for _, str := range haystack {
		if score := fuzzyMatch(needle, strings.ToLower(str)); score > 0 {
			matches = append(matches, Match{str: str, score: score})
		}
	}

	// Sort matches by score (highest first)
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].score > matches[j].score
	})

	// Extract just the strings
	result := make([]string, len(matches))
	for i, match := range matches {
		result[i] = match.str
	}

	return result
}

// fuzzyMatch calculates match score. Higher score means better match
func fuzzyMatch(needle, haystack string) int {
	score := 0
	needleIdx := 0
	lastMatchIdx := -1

	for i, char := range haystack {
		if needleIdx >= len(needle) {
			break
		}

		if byte(char) == needle[needleIdx] {
			score++

			// Bonus for consecutive matches
			if lastMatchIdx != -1 && i == lastMatchIdx+1 {
				score += 3
			}

			// Bonus for matching start of words
			if i == 0 || haystack[i-1] == ' ' {
				score += 2
			}

			lastMatchIdx = i
			needleIdx++
		}
	}

	// Return 0 if not all characters were found
	if needleIdx != len(needle) {
		return 0
	}

	return score
}
