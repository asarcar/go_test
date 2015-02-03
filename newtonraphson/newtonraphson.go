package main

import (
	"flag"
	"fmt"
	"math"
)

type Fn func(float64) float64

func SqRootFn(kValue float64) Fn {
	if kValue < 0 {
		panic("Invalid Function")
	}

	return func(xValue float64) float64 {
		return xValue*xValue - kValue
	}
}

func CubeRootFn(kValue float64) Fn {
	return func(xValue float64) float64 {
		return xValue*xValue*xValue - kValue
	}
}

const (
	StartValue          = 1.0
	CloseEnoughFraction = 0.0001
	Epsilon             = 0.0001
	MaxIterations       = 10
)

func SolveNewtonRaphson(f Fn) float64 {
	// Terminates loop when value does not change measurably
	closeFn := func(newV, prevV float64) bool {
		return math.Abs(newV-prevV) < math.Abs(prevV*CloseEnoughFraction)
	}

	// Derivative Function Computation
	derF := func(val float64) float64 {
		return (f(val+Epsilon) - f(val-Epsilon)) / (2 * Epsilon)
	}

	// computes next iteration of Netwon Raphson value
	// x(t+1) = x(t) -  f(x(t))/f'(x(t))
	nrFn := func(x float64) float64 {
		return x - (f(x) / derF(x))
	}

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

	fmt.Printf("Square Root of %0.2f is %0.4f\n",
		*fPtr, SolveNewtonRaphson(SqRootFn(val)))

	fmt.Printf("Cube Root of %0.2f is %0.4f\n",
		*fPtr, SolveNewtonRaphson(CubeRootFn(val)))
}
