package main

import (
	"fmt"
)

const (
	MAXI int = 5
)

func main() {
	var ch chan int = make(chan int, MAXI)
	var v int
	for i := 0; i < MAXI; i++ {
		v = i * 2 // closure refers to a variable that keeps changing
		go func() {
			ch <- v
		}()
	}
	for i := 0; i < MAXI; i++ {
		select {
		case v, ok := <-ch:
			if ok == false {
				ch = nil
			} else {
				fmt.Printf("v=%d\n", v)
			}
		}
	}
}
