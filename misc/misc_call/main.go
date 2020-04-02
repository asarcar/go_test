package main

import (
	"fmt"
	"github.com/asarcar/go_test/misc"
)

func main() {
	misc.DumpFlags()
	misc.DumpStr()
	misc.DumpGreek()
	misc.DumpTypePrint(float32(3.200123432), 10, "hello raju", float64(3.200123432))
	misc.DumpMarkovWords()
	misc.Median([]int{3, 7, 20, 50}, []int{4, 10, 12, 15})
	fmt.Printf("Par{3}: %v\n", misc.NewPar(3).String())
}
