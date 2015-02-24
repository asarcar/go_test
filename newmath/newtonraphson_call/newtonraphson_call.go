package main

import (
	"flag"
	"fmt"
	"github.com/asarcar/go_test/newmath"
)

func SqrtFn(kValue float64) newmath.Fn {
	if kValue < 0 {
		panic("Invalid Function")
	}
	return func(xValue float64) float64 {
		return xValue*xValue - kValue
	}
}

func CubertFn(kValue float64) newmath.Fn {
	return func(xValue float64) float64 {
		return xValue*xValue*xValue - kValue
	}
}

func QuadFn(kValue float64) newmath.Fn {
	return func(xValue float64) float64 {
		return 8.0*xValue*xValue - kValue*xValue - 100.0
	}
}

func main() {
	fPtr := flag.Float64("root", 64.0, "value to square/cube/quadratic seed\n")
	flag.Parse()

	val := *fPtr

	t := newmath.ZeroPoint
	fmt.Println("Linear Method - Computes Zero Point: ")
	fmt.Printf("  Square Root of %0.2f is %0.4f\n",
		*fPtr, newmath.SolveNewtonRaphson(SqrtFn(val), t))
	fmt.Printf("  Cube Root of %0.2f is %0.4f\n",
		*fPtr, newmath.SolveNewtonRaphson(CubertFn(val), t))

	t = newmath.MinMaxPoint
	fmt.Println("Quadratic Method - Computes Minimum or Maximum Point: ")
	fmt.Printf("  Min of 8.0*x^2 - %0.2f*x - 100.0 is at x=%0.2f\n",
		val, newmath.SolveNewtonRaphson(QuadFn(val), t))
}
