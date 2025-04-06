package utils

import (
	"formi/data"
	"strings"
)

type FuzzyMatcher struct {
	locationMap map[string]data.IndianLocation
}

func NewFuzzyMatcher(locations []data.IndianLocation) *FuzzyMatcher {
	m := make(map[string]data.IndianLocation)
	for _, loc := range locations {
		normalized := normalizeString(loc.Name)
		m[normalized] = loc
	}
	return &FuzzyMatcher{locationMap: m}
}

func (f *FuzzyMatcher) FindBestMatch(input string) (string, string) {
	normalizedInput := normalizeString(input)
	
	if loc, exists := f.locationMap[normalizedInput]; exists {
		return loc.Name, loc.Type
	}
	
	bestScore := 2
	var bestMatch data.IndianLocation
	
	for normName, loc := range f.locationMap {
		if score := levenshtein(normalizedInput, normName); score <= bestScore {
			bestScore = score
			bestMatch = loc
		}
	}
	
	if bestScore <= 2 {
		return bestMatch.Name, bestMatch.Type
	}
	
	return "", ""
}

func normalizeString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func levenshtein(a, b string) int {
	d := make([][]int, len(a)+1)
	for i := range d {
		d[i] = make([]int, len(b)+1)
	}
	
	for i := range d {
		d[i][0] = i
	}
	for j := range d[0] {
		d[0][j] = j
	}
	
	for i := 1; i <= len(a); i++ {
		for j := 1; j <= len(b); j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			d[i][j] = min(
				d[i-1][j]+1,
				d[i][j-1]+1,
				d[i-1][j-1]+cost,
			)
			
			if i > 1 && j > 1 && a[i-1] == b[j-2] && a[i-2] == b[j-1] {
				d[i][j] = min(d[i][j], d[i-2][j-2]+1)
			}
		}
	}
	return d[len(a)][len(b)]
}

func min(nums ...int) int {
	m := nums[0]
	for _, num := range nums {
		if num < m {
			m = num
		}
	}
	return m
}