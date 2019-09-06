package matcher

type WordTree struct {
	children map[rune]*WordTree
	words    []string
}

func NewWordTree() *WordTree {
	return &WordTree{
		children: map[rune]*WordTree{},
	}
}

// String returns a string representation of the tree. Note that this may span
// multiple lines.
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

// Add inserts word into the tree.
// TODO handle duplicates!!
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

	// Add the word to the last subnode.
	curTree.words = append(curTree.words, word)
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
	for _, word := range curTree.words {
		if word == givenWord {
			return true
		}
	}
	return false
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

// WordCombos returns all combinations of words in the tree which can be
// created from the letters of givenWord. {word1, word2} is assumed to be
// equivalent to {word2, word1}.
func (t *WordTree) WordCombos(givenWord string) [][]string {
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

	return summarize(knownSolutions)
}
