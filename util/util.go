package util

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/reiver/go-porterstemmer"
)

// Check checks for the existence of an error, and if any, displays the specified message in a standard stream
func Check(err error, format string) {
	if err != nil {
		fmt.Printf(format, err)
	}
}

// Abs takes a number x and returns its module
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// CleanUserInput takes an array of words and formats it, removing the word stop and highlighting tokens
func CleanUserInput(word string, consumer func(input string)) {
	word = strings.TrimFunc(word, func(r rune) bool {
		return !unicode.IsLetter(r)
	})
	if !EnglishStopWordChecker(word) {
		word = porterstemmer.StemString(word)
		if len(word) > 0 {
			consumer(word)
		}
	}
}
