package matcher_test

import (
	"sort"
	"testing"

	"github.com/dlapets/stickers/matcher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const TestDictionaryPath = "../data/test_dictionary.txt"
const RealDictionaryPath = "../data/dictionary.txt"

func TestLoadDictionary(t *testing.T) {
	dictionary, err := matcher.LoadDictionary(TestDictionaryPath)
	require.NoError(t, err)
	assert.Equal(t, []string{"hell", "hello", "help", "well"}, dictionary)
}

func TestSimpleMatcher_MultiWordsMatching(t *testing.T) {
	m := matcher.SimpleMatcher([]string{"hi", "there", "you"})
	matching := m.MultiWordsMatching("a hi youthere zzz")

	expectedPhrases := [][]string{
		{"hi", "there", "you"},
		{"hi", "there"},
		{"hi", "you"},
		{"hi"},
		{"there", "you"},
		{"there"},
		{"you"},
	}
	assert.Len(t, matching, len(expectedPhrases))
	for _, phrase := range expectedPhrases {
		assert.Contains(t, matching, phrase)
	}
}

func TestSimpleMatcher_WordsMatching(t *testing.T) {
	m := matcher.SimpleMatcher([]string{"hi", "i", "you", "get", "height"})
	matching := m.WordsMatching("height")
	sort.Strings(matching) // don't expect any particular order
	assert.Equal(t, []string{"get", "height", "hi", "i"}, matching)
}

func BenchmarkSimpleMatcher_WordsMatching_SimpleDictionary(b *testing.B) {
	m := matcher.SimpleMatcher([]string{"hi", "i", "you", "get", "height"})
	benchmarkMatcher(m, "height", b)
}

func BenchmarkSimpleMatcher_WordsMatching_RealDictionary(b *testing.B) {
	dictionary, err := matcher.LoadDictionary(RealDictionaryPath)
	if err != nil {
		b.Fatalf("cannot load dictionary: %s", err)
	}
	m := matcher.SimpleMatcher(dictionary)
	benchmarkMatcher(m, "height", b)
}

func benchmarkMatcher(m matcher.Matcher, w string, b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		m.WordsMatching(w)
	}
}
