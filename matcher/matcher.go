package matcher

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

type Matcher interface {
	WordsMatching(string) []string
}

// reads a list of words split by newlines from a file
func LoadDictionary(filename string) ([]string, error) {
	dictionaryBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	dictionary := []string{}
	lastOffset := 0
	for currentOffset, b := range dictionaryBytes {
		if b != byte('\n') {
			continue
		}
		dictionary = append(
			dictionary,
			string(dictionaryBytes[lastOffset:currentOffset]),
		)
		lastOffset = currentOffset + 1
	}
	return dictionary, nil
}

type SimpleMatcher []string

func (m SimpleMatcher) WordsMatching(target string) []string {
	matching := []string{}
	for _, word := range m {
		if wordContains(word, target) {
			matching = append(matching, word)
		}
	}

	return matching
}

func (m SimpleMatcher) MultiWordsMatching(target string) [][]string {
	words := m.WordsMatching(target)

	matching := [][]string{}
	combinations := getCombinations(len(words))

	for _, combo := range combinations {
		word := maskWords(words, combo)
		if wordContains(word, target) {
			fmt.Printf("LOOKING FOR word = %s IN target = %s\n", word, target)

			// TODO clean this up; sort is here for tests only
			phrase := make([]string, len(combo))
			for i, j := range combo {
				phrase[i] = words[j]
			}
			sort.Strings(phrase)
			matching = append(matching, phrase)
		}
	}

	return matching
}

func maskWords(words []string, combo []int) string {
	b := strings.Builder{}
	for _, i := range combo {
		b.WriteString(words[i])
	}
	return b.String()
}

func getCombinations(size int) [][]int {
	combinations := [][]int{}
	for i := 1; i <= size; i++ {
		combinations = append(combinations, nChooseK(size, i)...)
	}
	return combinations
}

// TODO This is ugly, clean this up?
func nChooseK(n, k int) [][]int {
	combinations := [][]int{}

	s := make([]int, k)
	last := k - 1

	var rc func(int, int) [][]int

	rc = func(i, next int) [][]int {
		for j := next; j < n; j++ {
			s[i] = j
			if i == last {
				combinations = append(combinations, sliceCopy(s))
			} else {
				rc(i+1, j+1)
			}
		}
		return combinations
	}

	return rc(0, 0)
}

func sliceCopy(s []int) []int {
	a := make([]int, len(s))
	for i, v := range s {
		a[i] = v
	}
	return a
}

func wordContains(word, target string) bool {
	if len(word) > len(target) {
		fmt.Println("That's too long\n")
		return false
	}
	runeCounts := wordRuneCounts(target)
	for _, letter := range word {
		if remaining, ok := runeCounts[letter]; !ok || remaining <= 0 {
			return false
		}
		runeCounts[letter]--
	}
	return true
}

func wordRuneCounts(word string) map[rune]int {
	runeCounts := map[rune]int{}
	for _, letter := range word {
		if _, ok := runeCounts[letter]; ok {
			runeCounts[letter]++
		} else {
			runeCounts[letter] = 1
		}
	}
	return runeCounts
}