package main

import (
	"flag"
	"fmt"
	"math"
)

type OptimizationType uint

const (
	Zero OptimizationType = iota
	MinMax
)

type Fn func(float64) float64

// Zero Point Functions
func SqRootFn(kValue float64, t OptimizationType) Fn {
	if kValue < 0 {
		panic("Invalid Function")
	}

	// Optimization is Zero (FixedPoint)
	if t == Zero {
		return func(xValue float64) float64 {
			return xValue*xValue - kValue
		}
	}

	// t == MinMax
	return func(xValue float64) float64 {
		return (xValue*xValue - kValue) * (xValue*xValue - kValue)
	}
}

func CubeRootFn(kValue float64, t OptimizationType) Fn {
	// Optimization is Zero (FixedPoint)
	if t == Zero {
		return func(xValue float64) float64 {
			return xValue*xValue*xValue - kValue
		}
	}

	// t == MinMax
	return func(xValue float64) float64 {
		return (xValue*xValue*xValue - kValue) * (xValue*xValue*xValue - kValue)
	}
}

const (
	// chosen to avoid areas of curve where function converges
	StartValue          = 32.0
	CloseEnoughFraction = 0.0001
	Epsilon             = 0.0001
	MaxIterations       = 16
)

func SolveNewtonRaphson(f Fn, t OptimizationType) float64 {
	// Terminates loop when value does not change measurably
	closeFn := func(newV, prevV float64) bool {
		return math.Abs(newV-prevV) < math.Abs(prevV*CloseEnoughFraction)
	}

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
		if t == Zero {
			return nrLinearFn
		}
		return nrQuadraticFn
	}(t)

	prevValue, newValue := StartValue, nrFn(StartValue)
	for i := 0; i < MaxIterations; i++ {
		if closeFn(newValue, prevValue) {
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

func main() {
	flag.Parse()

	fPtr := flag.Float64("v", 64.0, "value to square root\n")
	val := *fPtr

	t := Zero
	fmt.Println("Linear Method - Computes Zero Point: ")
	fmt.Printf("  Square Root of %0.2f is %0.4f\n",
		*fPtr, SolveNewtonRaphson(SqRootFn(val, t), t))
	fmt.Printf("  Cube Root of %0.2f is %0.4f\n",
		*fPtr, SolveNewtonRaphson(CubeRootFn(val, t), t))

	t = MinMax
	fmt.Println("Quadratic Method - Computes Minimum or Maximum Point: ")
	fmt.Printf("  Square Root of %0.2f is %0.4f\n",
		*fPtr, SolveNewtonRaphson(SqRootFn(val, t), t))
	fmt.Printf("  Cube Root of %0.2f is %0.4f\n",
		*fPtr, SolveNewtonRaphson(CubeRootFn(val, t), t))
}
