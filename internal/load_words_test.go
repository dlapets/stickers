package internal_test

import (
	"testing"

	"github.com/dlapets/stickers/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadWords_Success(t *testing.T) {
	// TODO consider moving test data somewhere else
	words, err := internal.LoadWords("../data/test_dictionary.txt")
	require.NoError(t, err)
	assert.Equal(t, []string{"hell", "hello", "help", "well"}, words)
}

func TestLoadWords_NotFound(t *testing.T) {
	words, err := internal.LoadWords("../not_a_file.txt")
	require.Error(t, err)
	assert.Empty(t, words)
}
