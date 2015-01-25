package main

import (
	"fmt"
	"log"
	"net/http"
)

type String string

func (s String) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, s+"\n")
}

type Struct struct {
	Greeting string
	Punct    string
	Who      string
}

func (s Struct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s%s %s: %s\n", s.Greeting, s.Punct, s.Who, r.URL.Path[1:])
}

func MyFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "HandleFunc Called Custom Function\n")
}

func main() {
	http.Handle("/string", String("I'm a frayed knot."))
	http.Handle("/struct", &Struct{"Hello", ":", "Gophers!"})
	http.HandleFunc("/myfunc", MyFunc)
	log.Fatal(http.ListenAndServe("localhost:4000", nil))
}
