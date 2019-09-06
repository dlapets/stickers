package matcher_test

import (
	"testing"

	"github.com/dlapets/stickers/matcher"
	"github.com/stretchr/testify/require"
)

func TestWordTree_Add_Find(t *testing.T) {
	tree := matcher.NewWordTree()

	tree.Add("help")
	require.True(t, tree.Find("help"))
	require.False(t, tree.Find("pleh"))

	tree.Add("pleh")
	require.True(t, tree.Find("help"))
	require.True(t, tree.Find("pleh"))

	require.False(t, tree.Find("helper"))
	require.False(t, tree.Find("doctor"))
}

func TestWordTree_Words(t *testing.T) {
	tree := matcher.NewWordTree()
	tree.Add("shit")
	tree.Add("hist")
	tree.Add("history")
	tree.Add("ass")
	tree.Add("asshat")

	expected := []string{"ass", "shit", "hist"}
	got := tree.Words("shitass")

	require.Equal(t, len(expected), len(got))
	for _, word := range expected {
		require.Contains(t, got, word)
	}
}

func TestWordTree_WordCombos(t *testing.T) {
	tree := matcher.NewWordTree()
	tree.Add("shit")
	tree.Add("hist")
	tree.Add("history")
	tree.Add("ass")
	tree.Add("asshat")

	expected := [][]string{
		{"asshat"},
		{"ass", "ass"},
		{"ass", "history"},
		{"ass", "shit"},
		{"ass", "hist"},
	}

	got := tree.WordCombos("shitassztoryas")

	require.Equal(t, len(expected), len(got))
	for _, words := range expected {
		require.Contains(t, got, words)
	}
}
