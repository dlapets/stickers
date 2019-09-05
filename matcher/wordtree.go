package matcher

import (
	"fmt"
	"log"
	"sort"
	"strings"
)

type WordTree struct {
	children map[rune]*WordTree
	words    []string
}

func NewWordTree() *WordTree {
	return &WordTree{
		children: map[rune]*WordTree{},
	}
}

func (t *WordTree) Add(word string) {
	//log.Println("add called:", word)
	if len(word) == 0 {
		return
	}

	curTree := t
	for _, r := range wordHash(word) {
		//log.Println("add rune:", string(r))

		// Create child if it doesn't exist
		if _, ok := curTree.children[r]; !ok {
			curTree.children[r] = NewWordTree()
		}

		curTree = curTree.children[r]
	}
	//log.Println("append", word)
	curTree.words = append(curTree.words, word)
}

func (t *WordTree) Find(word string) []string {
	curTree := t
	for _, r := range wordHash(word) {
		//log.Println("find rune:", string(r))

		// If there's no child that means it was not found!
		if _, ok := curTree.children[r]; !ok {
			//log.Println("not found:", word)
			return nil
		}

		curTree = curTree.children[r]
	}

	return curTree.words
}

func (t *WordTree) FindAllGroups2(givenWord string) {
	knownSolutions := map[string]*result{}
	givenWordHash := wordHash(givenWord)

	var fillSubresults func(int, string, *result)
	fillSubresults = func(level int, originalHash string, currentResult *result) {
		if currentResult.valid == false {
			return
		}

		hashDiff := wordHashDiff(originalHash, currentResult.hash)

		if knownSolution, ok := knownSolutions[hashDiff]; ok {
			currentResult.children = []*result{knownSolution.shallowCopy()}
		} else {
			currentResult.children = findAll(level+1, t, []rune(hashDiff))
		}

		for _, child := range currentResult.children {
			if _, ok := knownSolutions[child.hash]; !ok {
				knownSolutions[child.hash] = child.shallowCopy()
			}
			fillSubresults(level+1, hashDiff, child)
		}
	}

	for _, newResult := range findAll(0, t, []rune(givenWordHash)) {
		knownSolutions[newResult.hash] = newResult
		fillSubresults(0, givenWordHash, newResult)
	}

	log.Println("KNOWN SOLUTIONS:")
	resultsToPrint := []*result{}
	for _, res := range knownSolutions {
		resultsToPrint = append(resultsToPrint, res)
	}
	rprintSolutions(0, resultsToPrint)
	log.Println("SUMMARIZED RESULTS:")
	(&solutionSummarizer{}).updateCombos(resultsToPrint)
}

func rprintSolutions(level int, results []*result) {
	for _, res := range results {
		if level == 0 {
			log.Println("SOLUTION:")
		}
		for _, word := range res.words {
			lprintf(level+1, "%s\n", word)
			rprintSolutions(level+2, res.children)
		}
	}
}

type solutionSummarizer struct {
	wordCombos [][]string
}

func (s *solutionSummarizer) updateCombos(results []*result) {
	s.wordCombos = [][]string{}
	for _, res := range results {
		s.traverseSolution([]string{}, res)
	}

	// TODO this is ugly; would be better to find a way to prevent the solver
	// from returning duplicate results.
	deduper := map[string]struct{}{}

	for _, combo := range s.wordCombos {
		sort.Strings(combo)
		k := fmt.Sprintf("%s", combo)
		if _, ok := deduper[k]; !ok {
			deduper[k] = struct{}{}
			log.Println(k)
		}
	}
}

func (s *solutionSummarizer) traverseSolution(previousWords []string, res *result) {
	//log.Println("traverseSolution", previousWords, res)
	for _, word := range res.words {
		//log.Println("looking at word:", word)
		newPreviousWords := []string{}
		newPreviousWords = append(newPreviousWords, previousWords...)
		newPreviousWords = append(newPreviousWords, word)

		if len(res.children) == 0 {
			s.wordCombos = append(s.wordCombos, newPreviousWords)
			continue
		} else {
			for _, child := range res.children {
				s.traverseSolution(newPreviousWords, child)
			}
		}
	}
}

// FindAll finds all single word groups which can be derived from givenWord
func (t *WordTree) FindAll(givenWord string) [][]string {
	wordsOnly := [][]string{}
	for _, result := range findAll(0, t, []rune(wordHash(givenWord))) {
		wordsOnly = append(wordsOnly, result.words)
	}
	return wordsOnly
}

func findAll(level int, cur *WordTree, givenRunes []rune) []*result {
	//lprintf(level, "findAll: called with level: %d, remainderRunes: %s, givenRunes: %s\n", level, string(remainderRunes), string(givenRunes))

	found := []*result{}

	// There are words at this node, that means we can take them.
	if len(cur.words) > 0 {
		result := newResult(cur.words)
		//lprintf(level, "findAll: adding result: %s\n", result)
		found = append(found, result)
	}

	var prev rune
	for i, r := range givenRunes {
		if r == prev {
			//lprintf(level, "findAll: rune: %s: skipping duplicate\n", string(r))
			continue
		}

		if next, ok := cur.children[r]; ok {
			// Add all words found in subtrees
			nextRunes := givenRunes[i+1:]
			//lprintf(level, "findAll: rune: %s: recursing with remainderRunes: %s, nextRunes: %s\n", string(r), string(nextRunes))
			newFound := findAll(level+1, next, nextRunes)
			found = append(found, newFound...)
		} else {
			//lprintf(level, "findAll: rune: %s: no children\n", string(r))
		}
		prev = r
	}

	//lprintf(level, "findAll: return %s\n", found)
	return found
}

func (t *WordTree) String() string {
	return t.draw("root", 0)
}

func (t *WordTree) draw(label string, level int) string {
	s := lsprintf(level, "%s words: %s\n", label, t.words)
	for r, child := range t.children {
		s += child.draw(string(r), level+1)
	}
	return s
}

func lsprintf(level int, str string, stuff ...interface{}) string {
	var padding string
	if level > 0 {
		padding += strings.Repeat("  ", level)
	}

	return fmt.Sprintf(padding+str, stuff...)
}

func lprintf(level int, str string, stuff ...interface{}) {
	//return 0, nil
	log.Printf(lsprintf(level, str, stuff...))
}
