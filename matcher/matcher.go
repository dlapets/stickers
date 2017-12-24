package matcher

import "io/ioutil"

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

type ConcurrentMatcher struct {
	dictionary []string
	goroutines int
}

func NewConcurrentMatcher(
	dictionary []string,
	goroutines int,
) *ConcurrentMatcher {
	if goroutines < 1 {
		goroutines = 1
	}
	return &ConcurrentMatcher{
		dictionary: dictionary,
		goroutines: goroutines,
	}
}

func (m *ConcurrentMatcher) WordsMatching(target string) []string {
	wordsChan := make(chan string)
	go func() {
		for _, word := range m.dictionary {
			wordsChan <- word
		}
		close(wordsChan)
	}()

	resultsChan := make(chan string)
	doneChan := make(chan struct{})
	for i := 0; i < m.goroutines; i++ {
		go func() {
			for word := range wordsChan {
				if wordContains(word, target) {
					resultsChan <- word
				}
			}
			doneChan <- struct{}{}
		}()
	}

	matching := []string{}
	running := m.goroutines
	for {
		select {
		case word := <-resultsChan:
			matching = append(matching, word)
		case <-doneChan:
			running--
		}
		if running == 0 {
			break
		}
	}

	return matching
}
