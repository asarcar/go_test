package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

const (
	listenAddr          = "localhost:4000"
	chatProgramDuration = 100
)

func main() {
	log.Printf("Chat Daemon Started...\n")

	l := getListener()
	defer l.Close()

	// Channel where pairing channel received/sent
	pair := make(chan io.ReadWriteCloser)
	// Channel where termination of chat is received
	errc := make(chan error)

	// Spawn pair of sessions to allow chat
	for i := 0; i < 2; i++ {
		c, err := l.Accept()
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				log.Printf("Quit... Listen connection Timed out in %d secs\n", chatProgramDuration)
				return
			}
			log.Fatal(err)
		}
		defer c.Close()

		// Create a new pair for every odd session creation
		go matchChat(c, pair, errc)
	}

	if err := <-errc; err != nil {
		log.Fatal(err)
	}

	log.Println("Chat session over.")

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
	l.SetDeadline(time.Now().Add(time.Second * chatProgramDuration))
	return l
}

func matchChat(c io.ReadWriteCloser,
	pair chan io.ReadWriteCloser, // bidir channel
	errc chan<- error) { // write only channel
	fmt.Fprint(c, "Waiting for pair... ")
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
	case pair <- c:
	}
}

func cp(c, p io.ReadWriteCloser, errc chan<- error) { // write only channel
	fmt.Fprint(c, "Matched!\n")
	_, err := io.Copy(c, p)
	errc <- err
}
