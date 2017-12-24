package main_test

import (
	"testing"

	"github.com/dlapets/stickers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const TestDictionaryPath = "test_fixtures/dictionary.txt"

func TestLoadDictionary(t *testing.T) {
	dictionary, err := main.LoadDictionary(TestDictionaryPath)
	require.NoError(t, err)

	assert.Equal(t, []string{"hell", "hello", "help", "well"}, dictionary)
}
