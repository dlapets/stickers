package main

import (
	"io/ioutil"
	"log"
)

const DictionaryPath = "dictionary.txt"

func main() {
	_, err := LoadDictionary(DictionaryPath)
	if err != nil {
		log.Panicf("failed to read dictionary: %s", err)
	}

}

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

func WordsMatching(target string, dictionary []string) []string {
	matching := []string{}
	for _, word := range dictionary {
		if wordContains(word, target) {
			matching = append(matching, word)
		}
	}

	return matching
}

func wordContains(word, target string) bool {
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
