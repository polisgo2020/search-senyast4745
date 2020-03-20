package util

import (
	"errors"
	"fmt"
	"github.com/bbalet/stopwords"
	"github.com/reiver/go-porterstemmer"
	"strings"
	"unicode"
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
func CleanUserInput(words []string) ([]string, error) {
	var data []string
	for _, v := range words {
		word := strings.TrimFunc(v, func(r rune) bool {
			return !unicode.IsLetter(r)
		})
		word = stopwords.CleanString(word, "en", true)
		if len(word) > 0 {
			data = append(data, porterstemmer.StemString(word))
		}
	}
	if len(data) == 0 {
		return nil, errors.New(fmt.Sprintf("bad search words %+v", words))
	}
	return data, nil
}
