package util

import (
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func TestAbs(t *testing.T) {
	t.Parallel()
	i := -10
	require.Equal(t, -i, Abs(i), "Abs must be equal positive of i")
	i = 10
	require.Equal(t, i, Abs(i), "Abs must be equal positive of i")

	require.Equal(t, math.MaxInt64, Abs(math.MinInt64)-1)
	require.Equal(t, 0, Abs(0))
}

func TestCleanUserInput(t *testing.T) {
	t.Parallel()
	input := "a"
	called := false
	assertFunc := func(input string) {
		called = true
	}
	CleanUserInput(input, assertFunc)
	require.False(t, called, "function must not be called")

	input = "word"
	called = false
	assertFunc = func(i string) {
		called = true
		require.Equal(t, input, i)
	}
	CleanUserInput(input, assertFunc)
	require.True(t, called, "function must be called")

	input = "beautiful"
	called = false
	assertFunc = func(i string) {
		called = true
		require.NotEqual(t, input, i, "input must be stemmed")
	}
	CleanUserInput(input, assertFunc)
	require.True(t, called, "function must be called")

}

func TestEnglishStopWordChecker(t *testing.T) {
	t.Parallel()
	word := "you"
	require.True(t, EnglishStopWordChecker(word))
	word = "golang"
	require.False(t, EnglishStopWordChecker(word))
}
