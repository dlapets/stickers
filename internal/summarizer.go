package internal

import (
	"fmt"
	"sort"
)

func summarize(result *result) [][]string {
	return (&summarizer{}).updateCombos(result)
}

type summarizer struct {
	wordCombos [][]string
}

func (s *summarizer) updateCombos(result *result) [][]string {
	s.wordCombos = [][]string{}
	for _, res := range result.children {
		s.traverseResult([]string{}, res)
	}

	// TODO this is ugly; would be better to find a way to prevent the solver
	// from returning duplicate results.
	deduper := map[string]struct{}{}
	dedupedCombos := [][]string{}
	for _, combo := range s.wordCombos {
		sort.Strings(combo)
		k := fmt.Sprintf("%s", combo)
		if _, ok := deduper[k]; !ok {
			deduper[k] = struct{}{}
			dedupedCombos = append(dedupedCombos, combo)
		}
	}
	return dedupedCombos
}

func (s *summarizer) traverseResult(previousWords []string, res *result) {
	for _, word := range res.words {
		newPreviousWords := []string{}
		newPreviousWords = append(newPreviousWords, previousWords...)
		newPreviousWords = append(newPreviousWords, word)

		if len(res.children) == 0 {
			s.wordCombos = append(s.wordCombos, newPreviousWords)
			continue
		} else {
			for _, child := range res.children {
				s.traverseResult(newPreviousWords, child)
			}
		}
	}
}
