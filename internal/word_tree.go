package internal

import "log"

// WordTree is a structure for efficiently (?) identifying what words can be
// formed from the runes of a given string.
type WordTree struct {
	children map[rune]*WordTree
	words    []string
}

// NewWordTree creates an empty WordTree.
func NewWordTree() *WordTree {
	return &WordTree{
		children: map[rune]*WordTree{},
	}
}

// String returns a string representation of the tree. Note that this may span
// multiple lines.
func (t *WordTree) String() string {
	return t.draw(0, "root")
}

func (t *WordTree) draw(level int, label string) string {
	s := lsprintf(level, "%s words: %s\n", label, t.words)
	for r, child := range t.children {
		s += child.draw(level+1, string(r))
	}
	return s
}

// Add inserts word into the tree.
func (t *WordTree) Add(word string) {
	if len(word) == 0 {
		return
	}

	// Create subnodes and until every rune is accounted for.
	curTree := t
	for _, r := range wordHash(word) {
		if _, ok := curTree.children[r]; !ok {
			curTree.children[r] = NewWordTree()
		}
		curTree = curTree.children[r]
	}

	// Check that the word has not already been added.
	if sliceContainsWord(curTree.words, word) {
		log.Printf("attempt to add duplicate word: %s", word)
		return
	}

	// Add the word to the last subnode.
	curTree.words = append(curTree.words, word)
}

func sliceContainsWord(words []string, word string) bool {
	for _, oldWord := range words {
		if oldWord == word {
			return true
		}
	}
	return false
}

// Find returns true if word is in the tree
func (t *WordTree) Find(givenWord string) bool {
	curTree := t
	for _, r := range wordHash(givenWord) {
		// If there's no child that means it was not found!
		if _, ok := curTree.children[r]; !ok {
			return false
		}
		curTree = curTree.children[r]
	}
	return sliceContainsWord(curTree.words, givenWord)
}

// Words finds all single words in the tree which can be made from the letters
// of givenWord.
func (t *WordTree) Words(givenWord string) []string {
	words := []string{}
	for _, result := range findAll(0, t, []rune(wordHash(givenWord))) {
		words = append(words, result.words...)
	}
	return words
}

func findAll(level int, cur *WordTree, givenRunes []rune) []*result {
	// found contains all matching words we found so far.
	found := []*result{}

	// There are words at this node, we can use these words.
	if len(cur.words) > 0 {
		result := newResult(cur.words)
		found = append(found, result)
	}

	var prev rune
	for i, r := range givenRunes {
		if r == prev {
			continue
		}
		if next, ok := cur.children[r]; ok {
			// Add all words found in subtrees
			nextRunes := givenRunes[i+1:]
			newFound := findAll(level+1, next, nextRunes)
			found = append(found, newFound...)
		}
		prev = r
	}
	return found
}

// WordCombos returns all combinations of words in the tree which can be
// created from the letters of givenWord. {word1, word2} is assumed to be
// equivalent to {word2, word1}.
func (t *WordTree) WordCombos(givenWord string) [][]string {
	knownSolutions := map[string]*result{}
	hash := wordHash(givenWord)

	knownSolutions[hash] = newResult(nil)
	knownSolutions[hash].hash = hash

	t.fillSubresults(0, hash, knownSolutions[hash], knownSolutions)

	return summarize(knownSolutions)
}

func (t *WordTree) fillSubresults(
	level int,
	hash string,
	currentResult *result,
	knownSolutions map[string]*result,
) {
	if knownSolution, ok := knownSolutions[hash]; ok {
		// This has already been solved..., nothing to do:
		if len(knownSolution.children) > 0 {
			currentResult.children = knownSolution.children
			return
		}
	}

	currentResult.children = findAll(level+1, t, []rune(hash))

	for _, child := range currentResult.children {
		if _, ok := knownSolutions[child.hash]; !ok {
			knownSolutions[child.hash] = child
		}
		t.fillSubresults(level+1, wordHashDiff(hash, child.hash), child, knownSolutions)
	}
}
