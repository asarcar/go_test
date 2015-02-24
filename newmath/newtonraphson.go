package newmath

import (
	"fmt"
	"math"
)

type OptimizationType uint

const (
	ZeroPoint OptimizationType = iota
	MinMaxPoint
)

type Fn func(float64) float64

const (
	// chosen to avoid areas of curve where function converges
	StartValue       = 1.0
	ApproxEqualError = 0.0001
	Epsilon          = 0.0001
	MaxIterations    = 16
)

func ApproxEqual(newV, prevV float64) bool {
	return math.Abs(newV-prevV) < math.Abs(prevV*ApproxEqualError)
}

func SolveNewtonRaphson(f Fn, t OptimizationType) float64 {
	// Choose the Function whose derivate needs to be computed
	derF := func(fn Fn) Fn {
		return func(val float64) float64 {
			return (fn(val+Epsilon) - fn(val-Epsilon)) / (2 * Epsilon)
		}
	}

	// Derivative Function
	dF := derF(f)
	// Derivative of Derivative Function
	d2F := derF(dF)

	// computes next iteration of Netwon Raphson value

	// Linear Variant
	// x(t+1) = x(t) -  f(x(t))/f'(x(t))
	nrLinearFn := func(x float64) float64 {
		return x - (f(x) / dF(x))
	}

	// Quadratic Variant
	// x(t+1) = x(t) - f'(x(t))/f''(x(t))
	nrQuadraticFn := func(x float64) float64 {
		return x - (dF(x) / d2F(x))
	}

	// choose function based on Optimization Type:
	nrFn := func(opt OptimizationType) Fn {
		if t == ZeroPoint {
			return nrLinearFn
		}
		return nrQuadraticFn
	}(t)

	prevValue, newValue := StartValue, nrFn(StartValue)
	for i := 0; i < MaxIterations; i++ {
		if ApproxEqual(newValue, prevValue) {
			return prevValue
		}

		prevValue = newValue
		newValue = nrFn(prevValue)
	}

	var str string
	fmt.Sprintf(str,
		"Bug in Code! %d Iterations: prevValue %0.2f: newValue %0.2f\n",
		MaxIterations, prevValue, newValue)
	panic(str)

	return 0.0
}
