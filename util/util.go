package util

import (
	"strings"
	"unicode"

	"github.com/reiver/go-porterstemmer"
)

// Abs takes a number x and returns its module
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// CleanUserInput takes an array of words and formats it, removing the word stop and highlighting tokens
func CleanUserInput(rawWord string, consumer func(input string)) {
	word := strings.TrimFunc(rawWord, func(r rune) bool {
		return !unicode.IsLetter(r)
	})
	if !EnglishStopWordChecker(word) {
		word = porterstemmer.StemString(word)
		if len(word) > 0 {
			consumer(word)
		}
	}
}
