package misc

import (
	"fmt"
	"sync"
)

// https://leetcode.com/problems/print-zero-even-odd/
//
// Suppose you are given the following code:
// class ZeroEvenOdd {
//	 public void zero(printNumber) { ... }  // only output 0's
//	 public void even(printNumber) { ... }  // only output even numbers
//	 public void odd(printNumber) { ... }   // only output odd numbers
// }
// The same instance of ZeroEvenOdd will be passed to three different threads:
// Thread A will call zero() which should only output 0's.
// Thread B will call even() which should only ouput even numbers.
// Thread C will call odd() which should only output odd numbers.
// Each of the threads is given a printNumber method to output an integer.
// Modify the given program to output the series 010203040506... where the length of the series must be 2n.
//

// Input: n = 2
// Output: "0102"
// Explanation: There are three threads being fired asynchronously.
// One of them calls zero(), the other calls even(), and the last one calls odd().
// "0102" is the correct output.

const (
	ZERO_FN int = iota
	ODD_FN
	EVEN_FN
	MAX_FN_TYPE
)

const (
	ZERO_SIGNAL int = iota
	ODD_SIGNAL
	EVEN_SIGNAL
	MAX_SIGNAL
)

type pNumFn func(val int)

type ZeroEvenOdd struct {
	n int
	// wg_odd_zero used by zero to wait on odd
	wg      [MAX_SIGNAL]sync.WaitGroup
	wg_over sync.WaitGroup
}

func NewZeroEvenOdd(val int) {
	if val <= 0 {
		return
	}

	p := &ZeroEvenOdd{n: val}
	for i := ZERO_SIGNAL; i < MAX_SIGNAL; i++ {
		p.wg[i].Add(1)
	}

	// wait for zero, odd, and even to complete
	p.wg_over.Add(3)

	printN := func(val int) {
		if val < 0 {
			fmt.Printf("\n")
			return
		}
		fmt.Print(val)
	}

	// seed initial state assuming -1 is printed and first zero can print
	p.wg[ZERO_SIGNAL].Done()

	go p.Zero(printN)
	go p.Odd(printN)
	go p.Even(printN)

	p.wg_over.Wait()

	printN(-1)

	return
}

func (p *ZeroEvenOdd) Zero(printN pNumFn) {
	for i := 0; i < p.n; i++ {
		p.wg[ZERO_SIGNAL].Wait()
		// set up for next wait call
		p.wg[ZERO_SIGNAL].Add(1)
		printN(0)
		// signal odd and even threads alternatively
		p.wg[i%2+ODD_SIGNAL].Done()
	}
	p.wg_over.Done()
}

func (p *ZeroEvenOdd) Odd(printN pNumFn) {
	for i := 0; i < (p.n+1)/2; i++ {
		p.wg[ODD_SIGNAL].Wait()
		// set up for next wait call
		p.wg[ODD_SIGNAL].Add(1)
		printN(2*i + 1)
		p.wg[ZERO_SIGNAL].Done()
	}
	p.wg_over.Done()
}

func (p *ZeroEvenOdd) Even(printN pNumFn) {
	for i := 0; i < p.n/2; i++ {
		p.wg[EVEN_SIGNAL].Wait()
		// set up for next wait call
		p.wg[EVEN_SIGNAL].Add(1)
		printN(2*i + 2)
		p.wg[ZERO_SIGNAL].Done()
	}
	p.wg_over.Done()
}
