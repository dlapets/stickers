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
