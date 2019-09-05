package matcher

import (
	"fmt"
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
	fmt.Println("add called:", word)
	if len(word) == 0 {
		return
	}

	curTree := t
	for _, r := range wordHash(word) {
		fmt.Println("add rune:", string(r))

		// Create child if it doesn't exist
		if _, ok := curTree.children[r]; !ok {
			curTree.children[r] = NewWordTree()
		}

		curTree = curTree.children[r]
	}
	fmt.Println("append", word)
	curTree.words = append(curTree.words, word)
}

func (t *WordTree) Find(word string) []string {
	curTree := t
	for _, r := range wordHash(word) {
		fmt.Println("find rune:", string(r))

		// If there's no child that means it was not found!
		if _, ok := curTree.children[r]; !ok {
			fmt.Println("not found:", word)
			return nil
		}

		curTree = curTree.children[r]
	}

	return curTree.words
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

func (t *WordTree) FindAll(word string) [][]string {
	return findAll(0, t, []rune(wordHash(word)))
}

func findAll(level int, cur *WordTree, givenRunes []rune) [][]string {
	lprintf(level, "findAll: called with level: %d, givenRunes: %s\n", level, string(givenRunes))

	found := [][]string{}

	// There are words at this node, that means we can take them.
	if len(cur.words) > 0 {
		lprintf(level, "findAll: adding words: %s\n", cur.words)
		found = append(found, cur.words)
	}

	var prev rune
	for i, r := range givenRunes {
		if r == prev {
			lprintf(level, "findAll: rune: %s: skipping duplicate\n", string(r))
			continue
		}
		if next, ok := cur.children[r]; ok {
			// Add all words found in subtrees
			nextRunes := givenRunes[i+1:]
			lprintf(level, "findAll: rune: %s: recursing with nextRunes: %s\n", string(r), string(nextRunes))
			newFound := findAll(level+1, next, nextRunes)
			found = append(found, newFound...)
		} else {
			lprintf(level, "findAll: rune: %s: no children\n", string(r))
		}
		prev = r
	}

	lprintf(level, "findAll: return %s\n", found)
	return found
}

func lsprintf(level int, str string, stuff ...interface{}) string {
	var padding string
	if level > 0 {
		padding += strings.Repeat("  ", level)
	}

	return fmt.Sprintf(padding+str, stuff...)
}

func lprintf(level int, str string, stuff ...interface{}) (int, error) {
	return fmt.Printf(lsprintf(level, str, stuff...))
}
