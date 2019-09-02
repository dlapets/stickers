package matcher

import (
	"fmt"
	"io/ioutil"
	"sort"
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
		word := string(dictionaryBytes[lastOffset:currentOffset])
		if bullshit(word) {
			fmt.Printf("WHAT THE HELL IS A %s??\n", word)
		} else {
			dictionary = append(dictionary, word)
		}
		lastOffset = currentOffset + 1
	}
	return dictionary, nil
}

func bullshit(word string) bool {
	if len(word) == 1 {
		switch word {
		case "a", "i":
			return false
		default:
			return true
		}
	}
	return false
}

type SimpleMatcher struct {
	dict []string
}

func NewSimpleMatcher(dict []string) *SimpleMatcher {
	return &SimpleMatcher{
		dict: dict,
	}
}

func (m SimpleMatcher) WordsMatching(target string) []string {
	matching := []string{}
	for _, word := range m.dict {
		if wordContains(word, target) {
			matching = append(matching, word)
		}
	}

	return matching
}

func (m SimpleMatcher) MultiWordsMatching(target string) [][]string {
	words := m.WordsMatching(target)
	fmt.Printf("WORDS MATCHING ARE %v\n", words)

	matching := [][]string{}
	combinations := getCombinations(len(words))
	fmt.Printf("GOT COMBINATIONS\n")

	for combo := range combinations {
		word := maskWords(words, combo)
		if wordContains(word, target) {
			//fmt.Printf("LOOKING FOR word = %s IN target = %s\n", word, target)

			// TODO clean this up; sort is here for tests only
			phrase := make([]string, len(combo))
			for i, j := range combo {
				phrase[i] = words[j]
			}
			sort.Strings(phrase)
			fmt.Printf("HOW ABOUT %v\n", phrase)
			matching = append(matching, phrase)
		}
	}

	return matching
}
