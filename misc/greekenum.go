package misc

import (
	"fmt"
)

type GreekLetter int

// Go Version >= 1.4 supports go generate used to
// auto-generate function.
// Example: below command auto-generates a
// String() member function for GreekLetter

//go:generate stringer -type=GreekLetter

const (
	alpha GreekLetter = iota
	beta
	gamma
)

func DumpGreek() {
	fmt.Println("DumpGreek\n-----------")
	for letter := alpha; letter <= gamma; letter++ {
		fmt.Println(letter)
	}
	fmt.Println("-----------")
}
