package main

import (
	"fmt"
	"github.com/asarcar/go_test/misc"
)

var sqr misc.Fn = func(x int) int { return x * x }

var inc misc.Fn = func(x int) int { return x + 1 }

var dec misc.Fn = func(x int) int { return x - 1 }

func main() {
	misc.DumpFlags()
	misc.DumpStr()
	misc.DumpGreek()
	misc.DumpTypePrint(float32(3.200123432), 10, "hello raju", float64(3.200123432))

	misc.DumpValue()

	fmt.Println("---------------------------------------------")
	fmt.Printf("inc(sqr(dec(1)))=%d, dec(sqr(inc(1)))=%d\n",
		misc.Compose(inc, sqr)(dec)(1), misc.Compose(dec, sqr)(inc)(1))
	fmt.Println("---------------------------------------------")
}
