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

func ComposeFn(a ComposeFns, args ...Fn) Fn {
	numargs := len(args)
	if numargs == 0 {
		return a.(Fn)
	}

	// numargs >= 1
	res := args[numargs-1]
	for i := numargs - 2; i >= 0; i-- {
		res = args[i].compose(res)
	}
	return a.compose(res)
}
