package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

const (
	listenAddr = "localhost:4000"
)

func main() {
	log.Println("Chat Daemon Started...")

	l := getListener()
	defer l.Close()

	// Keep looping wait for pair of interested chatters
	for {
		spawnPartners(l)
	}

	log.Println("Chat Daemon Stopped.")

	return
}

func getListener() *net.TCPListener {
	// Network Service example
	addr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	return l
}

func spawnPartners(l *net.TCPListener) {
	// Channel where pairing channel received/sent
	pair := make(chan io.ReadWriteCloser)

	for i := 0; i < 2; i++ {
		// Spawn pair of sessions to allow chat
		c, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go matchChat(c, pair)
	}
}

func matchChat(c io.ReadWriteCloser, pair chan io.ReadWriteCloser) {
	fmt.Fprint(c, "Waiting for pair... ")

	// Channel where termination of chat is received
	errc := make(chan error)

	// For s:
	// (a) p channel is first read from partner channel or
	// (b) c channel is sent to partner.
	// If (a) happens first then by definition case (b) is
	// executed on c and if (b) happens then case (a) on p.
	select {
	case p := <-pair:
		log.Println("Chat pair established!")
		// stich the read and write ends to facilitate the chat
		go cp(c, p, errc)
		go cp(p, c, errc)
		defer c.Close()
		defer p.Close()
	case pair <- c:
		return
	}

	// Wait for this chat session to complete
	err := <-errc
	if err != nil {
		log.Fatal(err)
	}

	return
}

func cp(c, p io.ReadWriteCloser, errc chan<- error) { // write only channel
	fmt.Fprint(c, "Matched!\n")
	_, err := io.Copy(c, p)
	errc <- err
}
