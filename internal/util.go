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

func wordRuneCounts(word string) map[rune]int {
	runeCounts := map[rune]int{}
	for _, r := range word {
		if _, ok := runeCounts[r]; ok {
			runeCounts[r]++
		} else {
			runeCounts[r] = 1
		}
	}

	return runeCounts
}

func wordHash(word string) string {
	rc := wordRuneCounts(word)

	runes := []rune{} //TODO use runes here ...
	for r := range rc {
		runes = append(runes, r)
	}

	sort.Slice(runes, func(i, j int) bool { return runes[i] < runes[j] })

	b := strings.Builder{}
	for _, r := range runes {
		for i := 0; i < rc[r]; i++ {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// gets a new word hash; rhs expected to be subset of lhs.
// TODO this can be more efficient
func wordHashDiff(lhs, rhs string) string {
	lhsCounts := wordRuneCounts(lhs)
	//fmt.Println(lhsCounts)
	for _, r := range rhs {
		lhsCounts[r]--
	}
	//fmt.Println(lhsCounts)

	b := strings.Builder{}
	for r, count := range lhsCounts {
		for i := 0; i < count; i++ {
			b.WriteRune(r)
		}
	}
	res := wordHash(b.String())
	//fmt.Printf("wordHashDiff: %s - %s = %s\n", lhs, rhs, res)
	return res
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
