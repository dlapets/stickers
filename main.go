package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dlapets/stickers/matcher"
)

const DictionaryPath = "data/dictionary.txt"

func main() {
	target := strings.ToLower(strings.Join(os.Args[1:], ""))
	if target == "" {
		log.Panicf("no input")
	}

	dictionary, err := matcher.LoadDictionary(DictionaryPath)
	if err != nil {
		log.Panicf("failed to read dictionary: %s", err)
	}

	m := matcher.NewSimpleMatcher(dictionary)

	if results := m.MultiWordsMatching(target); len(results) != 0 {
		fmt.Println("YOU CAN MAKE THE FOLLOWING WITH YOUR STICKER!!")
		for _, result := range results {
			fmt.Println(result)
		}
		os.Exit(0)
	}

	fmt.Println("SORRY NOT TODAY!!")
	os.Exit(0)
}
