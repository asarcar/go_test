package main

import (
	"github.com/asarcar/go_test/misc"
)

func main() {
	misc.DumpFlags()
	misc.DumpStr()
	misc.DumpGreek()
	misc.DumpTypePrint(float32(3.200123432), 10, "hello raju", float64(3.200123432))
	misc.DumpValue()
	misc.DumpCompose()
	misc.DumpMarkovWords()
}
