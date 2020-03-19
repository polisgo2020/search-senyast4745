package util

import (
	"errors"
	"fmt"
	"github.com/bbalet/stopwords"
	"github.com/reiver/go-porterstemmer"
	"strings"
	"unicode"
)

func Check(err error, format string) {
	if err != nil {
		fmt.Printf(format, err)
	}
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

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
