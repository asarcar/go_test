package newmath

import (
	"log"
)

// Provides the kth fibonnaci sequece
func Fib(k int) int {
	if k < 0 {
		log.Fatal("Fib function only defined for natural numbers (integers > 0)")
	}

	j, i := 0, 1

	for m := 1; m < k; m++ {
		tmp := j
		j, i = i, tmp+i
	}

	return i
}
