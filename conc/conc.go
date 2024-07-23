package main

import (
	"flag"
	"fmt"
)

const (
	MAXI int = 5
)

func main() {
	sharedPtr := flag.Bool("shared", false, "boolean - closure uses shared variable when true")
	flag.Parse()

	var ch chan int = make(chan int, MAXI)
	var v int
	for i := 0; i < MAXI; i++ {
		v = i * 2 // closure refers to a variable that keeps changing
		go func() {
			if *sharedPtr == true {
				ch <- v
			} else {
				ch <- i * 2
			}
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
	close(ch)
}
