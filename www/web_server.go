package main

import (
	"fmt"
	"golang.org/x/net/trace"
	"log"
	"net/http"
)

const (
	kDomain  string = "localhost:4000"
	kPackage string = "github.com/asarcar/go_test.www"
)

var evlog trace.EventLog

type String string

func (s String) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, s+"\n")
	tr := trace.New(kDomain, r.URL.Path)
	defer tr.Finish()
	tr.LazyLog(s, true)
	evlog.Printf("StringEvLog: %s", s)
}

func (s String) String() string {
	return "StringLog: " + string(s)
}

type Struct struct {
	Greeting string
	Punct    string
	Who      string
}

func (s Struct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s%s %s: %s\n", s.Greeting, s.Punct, s.Who, r.URL.Path[1:])
	tr := trace.New(kDomain, r.URL.Path)
	defer tr.Finish()
	tr.LazyPrintf("Struct called with [Greeting-%s Punct-%s Who-%s]",
		s.Greeting, s.Punct, s.Who)
	evlog.Printf("StructEvLog: [Greeting-%s Punct-%s Who-%s]",
		s.Greeting, s.Punct, s.Who)
}

func MyFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "HandleFunc Called Custom Function\n")
	evlog.Printf("MyFunc: Terminating Eventlog")
	evlog.Finish()
}

func main() {
	evlog = trace.NewEventLog(kPackage, kDomain)
	http.Handle("/string", String("I'm a frayed knot."))
	http.Handle("/struct", &Struct{"Hello", ":", "Gophers!"})
	http.HandleFunc("/myfunc", MyFunc)
	log.Fatal(http.ListenAndServe(kDomain, nil))
}
