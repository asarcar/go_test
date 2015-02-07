package misc

type Fn func(int) int

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
