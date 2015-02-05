package newmath

import "testing"

func TestFib(t *testing.T) {
	const in, out = 10, 55
	if x := Fib(in); x != out {
		t.Errorf("Fib(%v) = %v, want %v", in, x, out)
	}
}
