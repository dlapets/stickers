package matcher_test

import (
	"fmt"
	"testing"

	"github.com/dlapets/stickers/matcher"
	"github.com/stretchr/testify/require"
)

func TestWordTree_AddAndFind(t *testing.T) {
	tree := matcher.NewWordTree()
	tree.Add("help")

	require.Equal(t, []string{"help"}, tree.Find("help"))

	tree.Add("pleh")
	require.Equal(t, []string{"help", "pleh"}, tree.Find("help"))

	require.Nil(t, tree.Find("helper"))
	require.Nil(t, tree.Find("doctor"))
}

func TestWordTree_FindMany(t *testing.T) {
	tree := matcher.NewWordTree()
	tree.Add("shit")
	tree.Add("hist")
	tree.Add("history")
	tree.Add("ass")
	tree.Add("asshat")

	fmt.Println(tree.String())

	expected := [][]string{
		{"ass"},
		{"shit", "hist"},
	}

	require.Equal(t, expected, tree.FindAll("shitass"))

	//tree.FindAllGroups("shitasszzz")

	tree.FindAllGroups2("shitassztoryas")
}

func TestWordTree_Whatever(t *testing.T) {
	tree := matcher.NewWordTree()
	dict, err := matcher.LoadDictionary("../data/celine.txt")
	require.NoError(t, err)

	for word := range dict.Words() {
		tree.Add(word)
	}

	tree.FindAllGroups2("iiiiiiii")
}
