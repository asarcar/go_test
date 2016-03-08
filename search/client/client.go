package main

import (
	"flag" // flag.Parse
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	pb "github.com/asarcar/go_test/search/protos"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const (
	kHTTPAddr   = "localhost"
	kRPCAddr    = "localhost"
	kRPCPort    = 5000
	kHTTPPort   = 8800
	kSearchPath = "/search"
	kWatchPath  = "/watch"
)

// Server addr:port where server accepts RPC requests
var (
	serverRPC string
	httpAddr  string
	client    pb.GoogleClient
)

func main() {
	parseFlags()
	fmt.Println("Client Spawned: Connecting to Server-RPC-Addr=\"" + serverRPC + "\"" +
		": HTTPAddr=\"" + httpAddr + "\"")

	var conn *grpc.ClientConn
	conn, client = dialRPCServer(serverRPC)
	defer conn.Close()
	spawnHTTPServer(httpAddr)
}

func spawnHTTPServer(httpAddr string) {
	http.HandleFunc(kSearchPath, handleSearch)
	http.HandleFunc(kWatchPath, handleWatch)

	log.Fatal(http.ListenAndServe(httpAddr, nil))
}

func dialRPCServer(rpcAddr string) (*grpc.ClientConn, pb.GoogleClient) {
	// Connect to Google Search RPC server:
	conn, err := grpc.Dial(rpcAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	return conn, pb.NewGoogleClient(conn)
}

// handleSearch handles URLs like /search?q=golang&timeout=1s by forwarding the
// query to google.Search. If the query param includes timeout, the search is
// canceled after that duration elapses.
func handleSearch(w http.ResponseWriter, req *http.Request) {
	// QUERY
	query := req.FormValue("q")
	if query == "" {
		http.Error(w, "no query", http.StatusBadRequest)
		return
	}

	// ctx is the Context for this handler. Calling cancel closes the
	// ctx.Done channel, which is the cancellation signal for requests
	// started by this handler.
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	// TIMEOUT
	timeout, err := time.ParseDuration(req.FormValue("t"))
	// The request has a timeout, so create a context that is
	// canceled automatically when the timeout expires.
	if err == nil {
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel() // Cancel ctx as soon as handleSearch returns.

	// SEARCH
	// Run the Google search and print the results.
	start := time.Now()
	// Req {Query}
	rpcreq := &pb.Request{Query: query}
	// Res {Title, Url, Content}
	rpcres, err := client.Search(ctx, rpcreq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	elapsed := time.Since(start)

	// BROWSER DISPLAY
	if err := resultsTemplate.Execute(w, struct {
		Results          *pb.Results
		Timeout, Elapsed time.Duration
	}{
		Results: rpcres,
		Timeout: timeout,
		Elapsed: elapsed,
	}); err != nil {
		log.Print(err)
		return
	}
}

// handleWatch
func handleWatch(w http.ResponseWriter, req *http.Request) {
	// QUERY
	query := req.FormValue("q")
	if query == "" {
		http.Error(w, "no query", http.StatusBadRequest)
		return
	}

	// ctx is the Context for this handler. Calling cancel closes the
	// ctx.Done channel, which is the cancellation signal for requests
	// started by this handler.
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	// TIMEOUT
	timeout, err := time.ParseDuration(req.FormValue("t"))
	// The request has a timeout, so create a context that is
	// canceled automatically when the timeout expires.
	if err == nil {
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel() // Cancel ctx as soon as handleSearch returns.

	// WATCH
	// Run the Google watch query and print the results.
	start := time.Now()
	stream, err := client.Watch(ctx, &pb.Request{Query: query})
	for {
		rpcres, rpcerr := stream.Recv()
		// rpcerr == io.EOF fails: per rpc_util.go equivalent error identified
		if rpcerr != nil && grpc.Code(rpcerr) == codes.OutOfRange {
			w.Write([]byte(string("watch session ended")))
			return
		}
		if rpcerr != nil {
			// log.Print(rpcerr.Error())
			http.Error(w, "RpcError: "+rpcerr.Error(), http.StatusInternalServerError)
			return
		}
		elapsed := time.Since(start)

		// BROWSER DISPLAY
		if err := resultsTemplate.Execute(w, struct {
			Results          *pb.Results
			Timeout, Elapsed time.Duration
		}{
			Results: rpcres,
			Timeout: timeout,
			Elapsed: elapsed,
		}); err != nil {
			// log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func parseFlags() {
	serverRPCPtr := flag.String("rpcserver",
		fmt.Sprintf("%s:%d", kRPCAddr, kRPCPort),
		"server RPC \"addr:port\" to connect")
	httpAddrPtr := flag.String("httpclient",
		fmt.Sprintf("%s:%d", kHTTPAddr, kHTTPPort),
		"HTTP \"addr:port\" to connect")
	flag.Parse()
	serverRPC = *serverRPCPtr
	httpAddr = *httpAddrPtr
}

var resultsTemplate = template.Must(template.New("results").Parse(`
<html>
<head/>
<body>
  <ol>
  {{range .Results.Res}}
    <li>{{.Title}} - <a href="{{.Url}}">{{.Url}}</a></li>
    Snippet: {{.Content}}
  {{end}}
  </ol>
  <p>result obtained in {{.Elapsed}}; timeout {{.Timeout}}</p>
</body>
</html>
`))
