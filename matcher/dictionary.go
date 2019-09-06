package matcher

import (
	"fmt"
	"io/ioutil"
)

type Dictionary struct {
	words map[string][]string // word hash -> words matching it
}

func (d *Dictionary) Lookup(givenWord string) ([]string, bool) {
	w, ok := d.words[wordHash(givenWord)]
	return w, ok
}

func (d *Dictionary) Words() chan string {
	c := make(chan string)
	// TODO go/chan not the best thing here, but provides a simple interface.
	go func() {
		for _, words := range d.words {
			for _, word := range words {
				c <- word
			}
		}
		close(c)
	}()
	return c
}

// reads a list of words split by newlines from a file
func LoadDictionary(filename string) (*Dictionary, error) {
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
			fmt.Printf("WHAT THE HELL IS A %s??\n", word)
		} else {
			words = append(words, word)
		}
		lastOffset = currentOffset + 1
	}
	return NewDictionary(words), nil
}

func NewDictionary(words []string) *Dictionary {
	dictionary := &Dictionary{
		words: map[string][]string{},
	}
	for _, word := range words {
		h := wordHash(word)
		dictionary.words[h] = append(dictionary.words[h], word)
	}
	return dictionary
}
