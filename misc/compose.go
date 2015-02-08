package misc

import (
	"fmt"
)

type Fn func(int) int

func DumpCompose() {
	var sqr Fn = func(x int) int { return x * x }
	var inc Fn = func(x int) int { return x + 1 }
	var dec Fn = func(x int) int { return x - 1 }
	fmt.Println("---------------------------------------------")
	fmt.Printf("inc(sqr(dec(1)))=%d, dec(sqr(inc(1)))=%d\n",
		Compose(inc, sqr)(dec)(1), Compose(dec, sqr)(inc)(1))
	fmt.Println("---------------------------------------------")
}

// Must implement f(x)
type ComposeFns interface {
	compose(Fn) Fn
}

func (g Fn) compose(f Fn) Fn {
	return func(x int) int {
		return g(f(x))
	}
}

func Compose(a, b ComposeFns) func(Fn) Fn {
	return func(f Fn) Fn {
		return a.compose(b.compose(f))
	}
}
