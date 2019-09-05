package matcher

import (
	"fmt"
	"log"
)

type result struct {
	valid    bool
	hash     string
	words    []string
	children []*result
}

var noResult = &result{} // not nil, not valid

func (r *result) String() string {
	return fmt.Sprintf(
		"[result valid:%5t, words:%s, hash:%s, children: %s]",
		r.valid,
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
		valid:    true,
		hash:     wordHash(words[0]), // TODO don't recalculate hash here
		words:    words,
		children: []*result{},
	}
}

func (res *result) shallowCopy() *result {
	return &result{
		valid:    res.valid,
		hash:     res.hash,
		words:    res.words,
		children: []*result{}, // doesn't get copied
	}
}
