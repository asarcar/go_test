package main

import (
	"flag" // flag.Parse
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/asarcar/go_test/search/backend"
	pb "github.com/asarcar/go_test/search/protos"

	"golang.org/x/net/context"
	"golang.org/x/net/trace"
	"google.golang.org/grpc"
)

const (
	kProtocol       = "tcp"
	kServerAddr     = "localhost"
	kServerHTTPPort = 8000
	kServerRPCPort  = 4000
)

// Server addr:port where server accepts RPC/HTTP requests
var serverRPC string
var serverHTTP string

type result struct {
	res *pb.Result
	err error
}

func main() {
	parseFlags()
	fmt.Println("Server Spawned: RPC-Addr=\"" + serverRPC + "\"" +
		": HTTP-Addr=\"" + serverHTTP + "\"")
	lis, err := net.Listen(kProtocol, serverRPC)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGoogleServer(s, &server{})

	// allows one to retrieve RPC visibility at /debug/requests and /debug/events
	go func() {
		err2 := http.ListenAndServe(serverHTTP, nil)
		log.Fatal(err2)
	}()

	s.Serve(lis)
}

func parseFlags() {
	serverRPCPtr := flag.String("rpcserver",
		fmt.Sprintf("%s:%d", kServerAddr, kServerRPCPort),
		"rpc server address \"addr:port\" to connect")
	serverHTTPPtr := flag.String("httpserver",
		fmt.Sprintf("%s:%d", kServerAddr, kServerHTTPPort),
		"http server address \"addr:port\" to connect")
	flag.Parse()
	serverRPC = *serverRPCPtr
	serverHTTP = *serverHTTPPtr
}

// server: implements GoogleServer
type server struct{}

func (s *server) Search(ctx context.Context, req *pb.Request) (*pb.Results, error) {
	d := randomSleep(ctx)

	select {
	case <-time.After(d):
		return backend.Search(ctx, req.Query)
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func randomSleep(ctx context.Context) time.Duration {
	d := randomDuration(100 * time.Millisecond)
	if tr, ok := trace.FromContext(ctx); ok {
		tr.LazyPrintf("sleeping for " + d.String())
	}
	return d
}

func randomDuration(max time.Duration) time.Duration {
	src := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(src)
	return time.Duration(rand.Int63n(int64(max)))
}
