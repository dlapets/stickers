package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dlapets/stickers/internal"
)

const dictionaryPath = "data/celine.txt"

func main() {
	log.SetOutput(os.Stderr)

	target := strings.ToLower(strings.Join(os.Args[1:], ""))
	if target == "" {
		log.Panicf("no input")
	}

	words, err := internal.LoadWords(dictionaryPath)
	if err != nil {
		log.Panicf("failed to read dictionary: %s", err)
	}

	wordTree := internal.NewWordTree()
	for _, word := range words {
		wordTree.Add(word)
	}

	if results := wordTree.WordCombos(target); len(results) != 0 {
		for _, result := range results {
			fmt.Println(result)
		}
		os.Exit(0)
	}

	log.Println("Doesn't look like you can make anything with that sticker!")
	os.Exit(0)
}
