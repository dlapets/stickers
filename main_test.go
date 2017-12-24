package main_test

import (
	"sort"
	"testing"

	"github.com/dlapets/stickers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const TestDictionaryPath = "data/test_dictionary.txt"

func TestLoadDictionary(t *testing.T) {
	dictionary, err := main.LoadDictionary(TestDictionaryPath)
	require.NoError(t, err)
	assert.Equal(t, []string{"hell", "hello", "help", "well"}, dictionary)
}

func TestWordsMatching(t *testing.T) {
	dictionary := []string{"hi", "i", "you", "get", "height"}
	matching := main.WordsMatching("height", dictionary)
	sort.Strings(matching) // don't expect any particular order
	assert.Equal(t, []string{"get", "height", "hi", "i"}, matching)
}

func BenchmarkWordsMatching_Simple(b *testing.B) {
	dictionary := []string{"hi", "i", "you", "get", "height"}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		main.WordsMatching("height", dictionary)
	}
}

func BenchmarkWordsMatching_RealDictionary(b *testing.B) {
	dictionary, err := main.LoadDictionary(main.DictionaryPath)
	if err != nil {
		b.Fatalf("cannot load dictionary: %s", err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		main.WordsMatching("height", dictionary)
	}
}
