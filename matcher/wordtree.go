package matcher

import "fmt"

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
