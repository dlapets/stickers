package matcher

import (
	"io/ioutil"
	"log"
)

// LoadWords reads words separated by newlines from a file.
func LoadWords(filename string) ([]string, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	words := []string{}
	lastOffset := 0
	for currentOffset, b := range bytes {
		if b != byte('\n') {
			continue
		}
		word := string(bytes[lastOffset:currentOffset])
		if bullshit(word) {
			log.Printf("skipping bs word: %s", word)
		} else {
			words = append(words, word)
		}
		lastOffset = currentOffset + 1
	}
	return words, nil
}
