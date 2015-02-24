package main

import (
	"code.google.com/p/go.net/websocket"
	"flag"
	"fmt"
	"github.com/asarcar/go_test/misc"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func main() {
	setGlobals()
	setHandleFuncs()
}

var gChain *misc.Chain
var gChatTemplate *template.Template

const (
	kListenAddr = "localhost:4000"
	// Server allows onto to chat with a virtual bot
	// max time before one is paired with bot
	kMaxChatWaitTime = 10
	// markov chain prefix length used to generate sentences by bot
	kMaxPrefixLen = 2
	// max # of words spewed at a time by bot
	kMaxWordsGen = 10
	// Default String with which generator is seeded
	kStr = `By now, most high school seniors planning to attend college 
       in the fall have selected their chosen institute of higher 
       education. It’s an exciting time for you, Wildcats '13, 
       and you probably have some questions about your future. 
       Such as, who will I meet? What clubs will I join? 
       What if my roommate only wants to stay in the room 
       eating cold cuts and watching Moesha re-runs? 
       Will I decide to buy a body pillow from Bed Bath and 
       Beyond? (Yes, besides being extremely comfortable body 
       pillows are an excellent way to block you from other 
       people's booger walls). In an effort to get to know each 
       other a little better before the fall rolls around, several 
       members of Columbia University’s future class of 2017 
       uploaded their college application essays into a shared 
       Google doc. That Google doc, which contains 70 essays that 
       either answer the Columbia essay prompt or the Common app 
       prompt, was then shared with us. And now with you.`
)

func setGlobals() {
	htmlDir := parseFlags()
	gChatTemplate = getTemplate(htmlDir, "chat.html")

	// Seeds the first set of random text spewed by Bot
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator.
	gChain = misc.NewChain(kMaxPrefixLen)
	gChain.Build(strings.NewReader(kStr))
}

func setHandleFuncs() {
	http.HandleFunc("/", rootHandler)
	http.Handle("/chat", websocket.Handler(socketHandler))

	log.Println("Chat Server Started...")
	log.Fatal(http.ListenAndServe(kListenAddr, nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	gChatTemplate.Execute(w, kListenAddr)
}

type socket struct {
	// A pass through for read/write can also be implemented by
	// simply listing the interface name: io.ReadWriter
	// io.ReadWriter
	io.Reader
	io.Writer
	done chan bool
}

func (s socket) Close() error {
	s.done <- true
	return nil
}

// websocket.Conn is held open by its handler function
// The handler is kept running using the channel receive
// until an explicit Close is called on the socket
// When we are simulating chat with a virtual person
// randomly emitting text
func socketHandler(ws *websocket.Conn) {
	r, w := io.Pipe()

	go func() {
		_, err := io.Copy(io.MultiWriter(w, gChain), ws)
		w.CloseWithError(err)
	}()

	s := socket{r, ws, make(chan bool)}
	go matchChat(s)

	<-s.done
}

// Bot returns an io.ReadeWriteCloser that responds to
// each incoming write with a generated sentence
type bot struct {
	io.ReadCloser
	out io.Writer
}

// Read or Close of Bot delegated to Read/Close(Pipe)
func Bot() io.ReadWriteCloser {
	r, out := io.Pipe()
	return bot{r, out}
}

// Pipe links the Write(Bot) to the Write(Pipe)
// WR(bot) throws away any data
// Also triggers WR(Pipe) with generated data ->
// Data available on RD end of Pipe ->
// RD(bot) == RD(Pipe) will get generated data
func (b bot) Write(buf []byte) (int, error) {
	go b.speak()
	return len(buf), nil
}

func (b bot) speak() {
	time.Sleep(time.Second)
	b.out.Write([]byte(gChain.Generate(kMaxWordsGen)))
}

// Channel where pairing channel received/sent
var pair = make(chan io.ReadWriteCloser)

func matchChat(c io.ReadWriteCloser) {
	var p io.ReadWriteCloser

	// Channel where termination of chat is received
	errc := make(chan error)

	// For s:
	// (a) p channel is first read from partner channel or
	// (b) c channel is sent to partner.
	// If (a) happens first then by definition case (b) is
	// executed on c and if (b) happens then case (a) on p.
	select {
	case p = <-pair:
		log.Println("Chat pair established!")
	case pair <- c:
		// p's goroutine would handle the chat session: we are done
		return
	case <-time.After(kMaxChatWaitTime * time.Second):
		log.Println("Chat pair simulated via Bot!")
		p = Bot()
	}

	// stich the read and write ends to facilitate the chat
	go cp(c, p, errc)
	go cp(p, c, errc)
	defer c.Close()
	defer p.Close()

	// Wait for this chat session to complete
	err := <-errc
	if err != nil {
		log.Fatal(err)
	}

	return
}

func cp(c, p io.ReadWriteCloser, errc chan<- error) { // write only channel
	fmt.Fprint(c, "CHAT Session Established!\n")
	_, err := io.Copy(c, p)
	errc <- err
}

func parseFlags() string {
	dPtr := flag.String("d",
		"/home/asarcar/git/go_test/src/github.com/asarcar/go_test/chat/html/",
		"full path to directory where template html files exit\n")
	flag.Parse()

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
