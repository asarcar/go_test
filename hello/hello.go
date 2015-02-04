package main

import (
	"fmt"
	"github.com/asarcar/go_test/newmath"
	"log"
	"net"
	"time"
)

const listenAddr = "localhost:4000"
const waitTimeInSecs = 10

func main() {
	// Hello World
	fmt.Printf("Hello, World: From Go!\n")
	// Library example
	fmt.Printf("Square Root(2) = %v\n", newmath.Sqrt(2))

	// Network Service example
	addr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	l.SetDeadline(time.Now().Add(time.Second * waitTimeInSecs))

	c, err := l.Accept()
	if err != nil {
		if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
			fmt.Printf("Quit... Listen connection Timed out in %d secs.\n", waitTimeInSecs)
			return
		}
		log.Fatal(err)
	}
	defer c.Close()

	fmt.Fprintln(c, "Hello, Net: From Go!")
}
