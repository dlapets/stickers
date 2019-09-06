package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dlapets/stickers/matcher"
)

const DictionaryPath = "data/celine.txt"

func main() {
	log.SetOutput(os.Stderr)

	target := strings.ToLower(strings.Join(os.Args[1:], ""))
	if target == "" {
		log.Panicf("no input")
	}

	dictionary, err := matcher.LoadDictionary(DictionaryPath)
	if err != nil {
		log.Panicf("failed to read dictionary: %s", err)
	}

	wordTree := matcher.NewWordTree()
	for word := range dictionary.Words() {
		wordTree.Add(word)
	}

	if results := wordTree.WordCombos(target); len(results) != 0 {
		fmt.Println("You can make the following with your sticker!")
		for _, result := range results {
			fmt.Println(result)
		}
		fmt.Println("Have fun!")
		os.Exit(0)
	}

	fmt.Println("Doesn't look like you can make anything with that sticker!")
	os.Exit(0)
}
