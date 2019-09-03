package matcher_test

import (
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
