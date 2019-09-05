package matcher

import (
	"sort"
	"strings"
)

func maskWords(words []string, combo []int) string {
	b := strings.Builder{}
	for _, i := range combo {
		b.WriteString(words[i])
	}
	return b.String()
}

func getCombinations(size int) <-chan []int {
	combinations := make(chan []int)

	go func() {
		for i := 1; i <= size; i++ {
			for _, combo := range nChooseK(size, i) {
				combinations <- combo
			}
		}
		close(combinations)
	}()

	return combinations
}

// TODO This is ugly, clean this up?
func nChooseK(n, k int) [][]int {
	combinations := [][]int{}

	s := make([]int, k)
	last := k - 1

	var rc func(int, int) [][]int

	rc = func(i, next int) [][]int {
		for j := next; j < n; j++ {
			s[i] = j
			if i == last {
				combinations = append(combinations, copySlice(s))
			} else {
				rc(i+1, j+1)
			}
		}
		return combinations
	}

	return rc(0, 0)
}

func copySlice(s []int) []int {
	a := make([]int, len(s))
	for i, v := range s {
		a[i] = v
	}
	return a
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
