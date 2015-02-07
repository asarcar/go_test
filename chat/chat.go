package main

import (
	"code.google.com/p/go.net/websocket"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
)

func main() {
	setGlobals()
	setHandleFuncs()
}

func setGlobals() {
	htmlDir := parseFlags()
	gChatTemplate = getTemplate(htmlDir, "chat.html")
}

const (
	listenAddr = "localhost:4000"
)

func setHandleFuncs() {
	http.HandleFunc("/", rootHandler)
	http.Handle("/chat", websocket.Handler(socketHandler))

	log.Println("Chat Server Started...")
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

var gChatTemplate *template.Template

func rootHandler(w http.ResponseWriter, r *http.Request) {
	gChatTemplate.Execute(w, listenAddr)
}

type socket struct {
	// A pass through for read/write can also be implemented by
	// simply listing the interface name: io.ReadWriter
	// io.ReadWriter
	conn *websocket.Conn
	done chan bool
}

func (s socket) Read(msg []byte) (n int, err error) {
	return s.conn.Read(msg)
}

func (s socket) Write(msg []byte) (n int, err error) {
	return s.conn.Write(msg)
}

func (s socket) Close() error {
	s.done <- true
	s.conn.Close()
	return nil
}

// websocket.Conn is held open by its handler function
// The handler is kept running using the channel receive
// until an explicit Close is called on the socket
func socketHandler(ws *websocket.Conn) {
	s := socket{conn: ws, done: make(chan bool)}
	go matchChat(s)
	<-s.done
}

// Channel where pairing channel received/sent
var pair = make(chan io.ReadWriteCloser)

func matchChat(c io.ReadWriteCloser) {
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

func parseFlags() string {
	flag.Parse()
	dPtr := flag.String("d",
		"/home/asarcar/git/go_test/src/github.com/asarcar/go_test/chat/html/",
		"full path to directory where template html files exit\n")

	return *dPtr
}

func getTemplate(dirPath, fileName string) *template.Template {
	f := dirPath + fileName
	t, err := template.ParseFiles(f)
	if err != nil {
		log.Fatal(err)
	}
	return t
}
