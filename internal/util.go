package internal

import (
	"fmt"
	"log"
	"sort"
	"strings"
)

func wordHash(word string) string {
	return wordHashDiff(word, "")
}

// wordHashDiff gets a new word hash from lhs, subtracting any runes present in
// rhs; rhs expected to be subset of lhs.
func wordHashDiff(lhs, rhs string) string {
	runeCounts := wordRuneCounts(lhs)

	// Subtract rhs from the counts.
	for _, r := range rhs {
		runeCounts[r]--
	}

	// Put all the runes in order.
	runes := []rune{}
	for r, count := range runeCounts {
		if count > 0 {
			runes = append(runes, r)
		}
	}
	sort.Slice(runes, func(i, j int) bool { return runes[i] < runes[j] })

	// Write sorted results to a string.
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
