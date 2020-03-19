package util

import "fmt"

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
