// Package newmath is a trivial example package.
package newmath

// Sqrt returns an approximation to the square root of x.
func Sqrt(x float64) float64 {
	sqRtFn := func(kValue float64) Fn {
		if kValue < 0 {
			panic("Invalid Function")
		}
		return func(xValue float64) float64 {
			return xValue*xValue - kValue
		}
	}

	// Newton Raphson for Fixed Point
	return SolveNewtonRaphson(sqRtFn(x), ZeroPoint)
}

// Cubert returns an approximation to the cube root of x.
func Cubert(x float64) float64 {
	cubeRtFn := func(kValue float64) Fn {
		return func(xValue float64) float64 {
			return xValue*xValue*xValue - kValue
		}
	}

	// Newton Raphson for Fixed Point
	return SolveNewtonRaphson(cubeRtFn(x), ZeroPoint)
}
