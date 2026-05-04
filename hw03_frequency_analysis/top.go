package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(text string) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return nil
	}

	wordCounts := make(map[string]int, len(words))
	for _, word := range words {
		wordCounts[word]++
	}

	uniqueWords := make([]string, 0, len(wordCounts))
	for word := range wordCounts {
		uniqueWords = append(uniqueWords, word)
	}

	sort.Slice(uniqueWords, func(i, j int) bool {
		if wordCounts[uniqueWords[i]] != wordCounts[uniqueWords[j]] {
			return wordCounts[uniqueWords[i]] > wordCounts[uniqueWords[j]]
		}
		return uniqueWords[i] < uniqueWords[j]
	})

	resultLimit := 10
	if len(uniqueWords) > resultLimit {
		uniqueWords = uniqueWords[:resultLimit]
	}

	return uniqueWords
}
