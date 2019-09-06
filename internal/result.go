package internal

import (
	"fmt"
	"log"
)

type result struct {
	hash     string
	words    []string
	children []*result
}

func (r *result) String() string {
	return fmt.Sprintf(
		"[result words:%s, hash:%s, children: %s]",
		r.words,
		r.hash,
		r.children,
	)
}

func newResult(words []string) *result {
	var hash string
	if len(words) > 0 {
		// TODO try to find a way to avoid recalculating the hash here.
		hash = wordHash(words[0])
	} else {
		log.Printf("creating result from empty word list")
	}
	return &result{
		hash:     hash,
		words:    words,
		children: []*result{},
	}
}

func (r *result) shallowCopy() *result {
	return &result{
		hash:     r.hash,
		words:    r.words,
		children: []*result{}, // doesn't get copied
	}
}
