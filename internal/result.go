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
	if len(words) == 0 { // TODO get rid of panic
		log.Panicf("there's no words on it!")
	}
	return &result{
		hash:     wordHash(words[0]), // TODO don't recalculate hash here
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
