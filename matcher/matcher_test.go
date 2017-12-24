package matcher_test

import (
	"fmt"
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

func TestConcurrentMatcher_WordsMatching(t *testing.T) {
	for goroutines := 0; goroutines < 3; goroutines++ {
		t.Run(fmt.Sprintf("goroutines %d", goroutines), func(t *testing.T) {
			m := matcher.NewConcurrentMatcher(
				[]string{"hi", "i", "you", "get", "height"},
				goroutines,
			)
			matching := m.WordsMatching("height")
			sort.Strings(matching) // don't expect any particular order
			assert.Equal(t, []string{"get", "height", "hi", "i"}, matching)
		})
	}
}

func BenchmarkConcurrentMatcher_WordsMatching_SimpleDictionary(b *testing.B) {
	dictionary := []string{"hi", "i", "you", "get", "height"}
	for i := 1; i <= 4; i++ {
		b.Run(fmt.Sprintf("goroutines %d", i), func(b *testing.B) {
			m := matcher.NewConcurrentMatcher(dictionary, i)
			benchmarkMatcher(m, "height", b)
		})
	}
}

func BenchmarkConcurrentMatcher_WordsMatching_RealDictionary(b *testing.B) {
	dictionary, err := matcher.LoadDictionary(RealDictionaryPath)
	if err != nil {
		b.Fatalf("cannot load dictionary: %s", err)
	}
	for i := 1; i <= 4; i++ {
		b.Run(fmt.Sprintf("goroutines %d", i), func(b *testing.B) {
			m := matcher.NewConcurrentMatcher(dictionary, i)
			benchmarkMatcher(m, "height", b)
		})
	}
}

func benchmarkMatcher(m matcher.Matcher, w string, b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		m.WordsMatching(w)
	}
}
