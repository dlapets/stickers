package matcher

import (
	"strconv"
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
	for _, letter := range word {
		if _, ok := runeCounts[letter]; ok {
			runeCounts[letter]++
		} else {
			runeCounts[letter] = 1
		}
	}
	return runeCounts
}

func wordHash(word string) string {
	rc := wordRuneCounts(word)
	b := strings.Builder{}
	for r, c := range rc {
		b.WriteRune(r)
		b.WriteRune(':')
		b.WriteString(strconv.Itoa(c))
		b.WriteRune(';')
	}
	return b.String()
}
