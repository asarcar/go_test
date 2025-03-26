package main

import (
	"flag"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/websocket"
)

const (
	kListenAddr      = "localhost:5432"
	kEchoWebFileName = "echoweb.html"
	kRootUrl         = "/"
	kEchoUrl         = "/echo"
	kEchoHtmlDir     = "/src/github.com/asarcar/go_test/echoweb/html/"
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
	dPtr := flag.String("d", kEchoHtmlDir,
		"relative to GOPATH the directory where template html files exist\n")
	flag.Parse()

	return *dPtr
}

func main() {
	dir := parseFlags()
	goPathDir := os.Getenv("GOPATH")
	htmlDir := goPathDir + dir
	echoWebTemplate = getTemplate(htmlDir, kEchoWebFileName)

	log.Println("EchoWeb: started at " + kListenAddr)
	http.HandleFunc(kRootUrl, rootHandler)
	http.Handle(kEchoUrl, websocket.Handler(EchoServer))
	log.Fatal(http.ListenAndServe(kListenAddr, nil))
}
