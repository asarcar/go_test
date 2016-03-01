package main

import (
	"flag"
	"golang.org/x/net/websocket"
	"html/template"
	"io"
	"log"
	"net/http"
)

const (
	kListenAddr      = "localhost:4000"
	kEchoWebFileName = "echoweb.html"
	kRootUrl         = "/"
	kEchoUrl         = "/echo"
)

var (
	echoWebTemplate *template.Template
)

func getTemplate(dirPath, fileName string) *template.Template {
	t, err := template.ParseFiles(dirPath + fileName)
	if err != nil {
		panic(err)
	}
	return t
}

func EchoServer(ws *websocket.Conn) {
	io.WriteString(ws, "Hello EchoWeb Begins...!")
	log.Println("EchoWeb Begins...")
	io.Copy(ws, ws)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	echoWebTemplate.Execute(w, kListenAddr+kEchoUrl)
}

func parseFlags() string {
	dPtr := flag.String("d",
		"/home/asarcar/git/go_test/src/github.com/asarcar/go_test/echoweb/html/",
		"full path to directory where template html files exist\n")
	flag.Parse()

	return *dPtr
}

func main() {
	htmlDir := parseFlags()
	echoWebTemplate = getTemplate(htmlDir, kEchoWebFileName)

	log.Println("EchoWeb: started at " + kListenAddr)
	http.HandleFunc(kRootUrl, rootHandler)
	http.Handle(kEchoUrl, websocket.Handler(EchoServer))
	log.Fatal(http.ListenAndServe(kListenAddr, nil))
}
