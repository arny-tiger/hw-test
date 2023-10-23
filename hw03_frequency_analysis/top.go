package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(str string) []string {
	wordsCount := getWordsCount(str)
	return sortByFrequency(wordsCount)
}

func getWordsCount(str string) map[string]int {
	wordsCount := make(map[string]int)
	words := strings.Fields(str)
	for _, word := range words {
		wordsCount[word]++
	}
	return wordsCount
}

func sortByFrequency(wordsCount map[string]int) []string {
	keys := make([]string, 0, len(wordsCount))
	for key := range wordsCount {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		if wordsCount[keys[i]] > wordsCount[keys[j]] {
			return true
		}
		if wordsCount[keys[i]] == wordsCount[keys[j]] {
			return strings.Compare(keys[i], keys[j]) < 0
		}
		return false
	})
	if len(keys) > 10 {
		return keys[:10]
	}
	return keys
}
