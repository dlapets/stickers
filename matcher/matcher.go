package matcher

import (
	"fmt"
	"io/ioutil"
	"sort"
)

type Matcher interface {
	WordsMatching(string) []string
}

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
	dict  *Dictionary
	cache map[string]map[string]bool
}

func NewSimpleMatcher(dict *Dictionary) *SimpleMatcher {
	return &SimpleMatcher{
		dict:  dict,
		cache: map[string]map[string]bool{},
	}
}

func (m *SimpleMatcher) WordsMatching(givenWord string) []string {
	matching := []string{}
	for word := range m.dict.Words() {
		if m.wordContains(word, givenWord) {
			matching = append(matching, word)
		}
	}
	return matching
}

func (m *SimpleMatcher) wordContains(word, givenWord string) bool {
	if _, ok := m.cache[word]; !ok {
		m.cache[word] = map[string]bool{}
	}
	if _, ok := m.cache[word][givenWord]; !ok {
		m.cache[word][givenWord] = m.wordContainsUncached(word, givenWord)
	}
	return m.cache[word][givenWord]

}

func (m *SimpleMatcher) wordContainsUncached(word, givenWord string) bool {
	if len(word) > len(givenWord) {
		//fmt.Println("That's too long\n")
		return false
	}
	runeCounts := wordRuneCounts(givenWord)
	for _, letter := range word {
		if remaining, ok := runeCounts[letter]; !ok || remaining <= 0 {
			return false
		}
		runeCounts[letter]--
	}
	return true
}

func (m *SimpleMatcher) MultiWordsMatching(givenWord string) [][]string {
	words := m.WordsMatching(givenWord)
	fmt.Printf("WORDS MATCHING ARE %v\n", words)

	matching := [][]string{}
	combinations := getCombinations(len(words))
	fmt.Printf("GOT COMBINATIONS\n")

	for combo := range combinations {
		word := maskWords(words, combo)
		if m.wordContains(word, givenWord) {
			//fmt.Printf("LOOKING FOR word = %s IN givenWord = %s\n", word, givenWord)

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
