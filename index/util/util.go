package util

import "fmt"

func Check(err error, format string) {
	if err != nil {
		fmt.Printf(format, err)
	}
}
