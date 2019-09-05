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

type Result struct {
	words     []string
	remainder string
}

func (t *WordTree) FindAll(word string) [][]string {
	wordsOnly := [][]string{}
	for _, result := range findAll(0, t, []rune{}, []rune(wordHash(word))) {
		wordsOnly = append(wordsOnly, result.words)
	}
	return wordsOnly
}

func findAll(level int, cur *WordTree, remainderRunes []rune, givenRunes []rune) []Result {
	lprintf(level, "findAll: called with level: %d, remainderRunes: %s, givenRunes: %s\n", level, string(remainderRunes), string(givenRunes))

	found := []Result{}

	// There are words at this node, that means we can take them.
	if len(cur.words) > 0 {
		newRemainderRunes := []rune{}
		newRemainderRunes = append(newRemainderRunes, remainderRunes...)
		if len(givenRunes) > 0 {
			newRemainderRunes = append(newRemainderRunes, givenRunes...)
		}

		result := Result{words: cur.words, remainder: normalizeRemainderRunes(newRemainderRunes)}
		lprintf(level, "findAll: adding result: %s\n", result)
		found = append(found, result)
	}

	var prev rune
	for i, r := range givenRunes {
		if r == prev {
			lprintf(level, "findAll: rune: %s: skipping duplicate\n", string(r))
			continue
		}

		// TODO see if we can avoid copy here...
		newRemainderRunes := []rune{}
		newRemainderRunes = append(newRemainderRunes, remainderRunes...)
		newRemainderRunes = append(newRemainderRunes, givenRunes[0:i]...)

		if next, ok := cur.children[r]; ok {
			// Add all words found in subtrees
			nextRunes := givenRunes[i+1:]
			lprintf(level, "findAll: rune: %s: recursing with remainderRunes: %s, nextRunes: %s\n", string(r), string(newRemainderRunes), string(nextRunes))
			newFound := findAll(level+1, next, newRemainderRunes, nextRunes)
			found = append(found, newFound...)
		} else {
			lprintf(level, "findAll: rune: %s: no children\n", string(r))
		}
		prev = r
	}

	lprintf(level, "findAll: return %s\n", found)
	return found
}

func normalizeRemainderRunes(runes []rune) string {
	return wordHash(string(runes))
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

func lprintf(level int, str string, stuff ...interface{}) (int, error) {
	return fmt.Printf(lsprintf(level, str, stuff...))
}
