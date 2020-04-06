// Package index for building search
// csv of the index file by the file in the given directory,
// as well as the further search of the search phrase by the built index.
//
// The index is built taking into account the position of the token in the file,
// clearing each file of stop words and using stemming for each token.
// The search is also made taking into account the proximity of the search tokens.
//
// For information about inverted index see https://habr.com/ru/post/53987/
package index
