package main

import (
	"fmt"
	"os"

	"net/http"

	"github.com/asarcar/go_test/newmath"
	"github.com/nicholasjackson/env"
)

const (
	DEF_HNAME string = "Euler"
	DEF_ADDR  string = "localhost"
	DEF_PORT  int    = 9090
)

var hName = env.String("HNAME", false, DEF_HNAME, "name to greet hello")
var bAddr = env.String("BADDR", false, DEF_ADDR, "bind address for server")
var bPort = env.Int("BPORT", false, DEF_PORT, "bind port for server")

func main() {
	err := env.Parse()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	bAddrPortStr := fmt.Sprintf("%s:%d", *bAddr, *bPort)
	fmt.Printf("Spawning Web Server: hello name %s bind address/port %s\n",
		*hName, bAddrPortStr)
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(rw, "Hello %s - Square Root(2) = %v\n", *hName, newmath.Sqrt(2))
	})

	http.ListenAndServe(bAddrPortStr, nil)
}
