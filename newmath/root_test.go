package newmath

import "testing"

func TestSqrt(t *testing.T) {
	const in, out = 4, 2
	if x := Sqrt(in); !ApproxEqual(x, out) {
		t.Errorf("Sqrt(%v) = %v, want %v", in, x, out)
	}
}

func TestCubert(t *testing.T) {
	const in, out = 8, 2
	if x := Cubert(in); !ApproxEqual(x, out) {
		t.Errorf("Cubert(%v) = %v, want %v", in, x, out)
	}
}

func TestMin(t *testing.T) {
	quadFn := func(k float64) float64 {
		return 3*k*k - 18*k - 73
	}
	const out = 3
	if x := SolveNewtonRaphson(quadFn, MinMaxPoint); !ApproxEqual(x, out) {
		t.Errorf("Min(3k^2 - 18k -73) = %v, k = %v, want %v", quadFn(x), x, out)
	}
}

func TestMax(t *testing.T) {
	quadFn := func(k float64) float64 {
		return -3*k*k + 18*k + 73
	}
	const out = 3
	if x := SolveNewtonRaphson(quadFn, MinMaxPoint); !ApproxEqual(x, out) {
		t.Errorf("Min(3k^2 - 18k -73) = %v, k = %v, want %v", quadFn(x), x, out)
	}
}
