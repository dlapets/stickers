package internal

import (
	"fmt"
	"log"
	"sort"
	"strings"
)

// Some words are more equal than others!
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

func wordHash(word string) string {
	return wordHashDiff(word, "")
}

// gets a new word hash; rhs expected to be subset of lhs.
// TODO this can be more efficient
func wordHashDiff(lhs, rhs string) string {
	runeCounts := wordRuneCounts(lhs)

	// subtract rhs
	for _, r := range rhs {
		runeCounts[r]--
	}

	// put all the runes in order
	runes := []rune{}
	for r, count := range runeCounts {
		if count > 0 {
			runes = append(runes, r)
		}
	}
	sort.Slice(runes, func(i, j int) bool { return runes[i] < runes[j] })

	// write to string
	b := strings.Builder{}
	for _, r := range runes {
		for i := 0; i < runeCounts[r]; i++ {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func wordRuneCounts(word string) map[rune]int {
	runeCounts := map[rune]int{}
	for _, r := range word {
		runeCounts[r]++
	}
	return runeCounts
}

// Helper func for recursive printing: pads left depending on value of level.
func lsprintf(level int, str string, stuff ...interface{}) string {
	var padding string
	if level > 0 {
		padding += strings.Repeat("  ", level)
	}

	return fmt.Sprintf(padding+str, stuff...)
}

func lprintf(level int, str string, stuff ...interface{}) {
	log.Printf(lsprintf(level, str, stuff...))
}
