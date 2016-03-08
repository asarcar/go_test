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
	kProtocol = "tcp"
	kHTTPAddr = "localhost" // address of local HTTP service
	kRPCAddr  = "localhost" // address of remote RPC service
	kRPCPort  = 4000        // port of remote RPC service
	kHTTPPort = 8000        // port of local HTTP service
)

// Server addr:port where server accepts RPC/HTTP requests
var (
	serverRPC  string
	serverHTTP string
)

// server: implements GoogleServer
type server struct{}

func main() {
	parseFlags()
	fmt.Println("Server Spawned: RPC-Addr=\"" + serverRPC + "\"" +
		": HTTP-Addr=\"" + serverHTTP + "\"")
	// allows one to retrieve RPC visibility at /debug/requests and /debug/events
	go spawnHTTPServer(serverHTTP)
	spawnRPCServer(serverRPC, &server{})
}

func spawnHTTPServer(httpAddr string) {
	err := http.ListenAndServe(httpAddr, nil)
	log.Fatal(err)
}

func spawnRPCServer(rpcAddr string, svr_p *server) {
	lis, err := net.Listen(kProtocol, rpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGoogleServer(s, svr_p)

	s.Serve(lis)
}

func (s *server) Search(ctx context.Context, req *pb.Request) (*pb.Results, error) {
	d := randomSleep(ctx)

	select {
	case <-time.After(d):
		return backend.Search(ctx, req.Query)
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s *server) Watch(req *pb.Request, stream pb.Google_WatchServer) error {
	ctx := stream.Context()
	for i := 0; i < 5; i++ {
		d := randomSleep(ctx)
		select {
		case <-time.After(d):
			err := stream.Send(&pb.Results{
				Res: []*pb.Result{
					&pb.Result{
						Title: fmt.Sprintf("[%d] for [%s] from backend %s",
							2*i, req.Query, serverRPC),
						Url:     "http://bakwaas.com",
						Content: "Na tera koi na mera koi",
					},
					&pb.Result{
						Title: fmt.Sprintf("[%d] for [%s] from backend %s",
							2*i+1, req.Query, serverRPC),
						Url:     "http://jhakaas.com",
						Content: "Sab mera sab tera",
					},
				},
			})
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

func randomSleep(ctx context.Context) time.Duration {
	d := randomDuration(1000 * time.Millisecond)
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

func parseFlags() {
	serverRPCPtr := flag.String("rpcserver",
		fmt.Sprintf("%s:%d", kRPCAddr, kRPCPort),
		"rpc server address \"addr:port\" to connect")
	serverHTTPPtr := flag.String("httpserver",
		fmt.Sprintf("%s:%d", kRPCAddr, kHTTPPort),
		"http server address \"addr:port\" to connect")
	flag.Parse()
	serverRPC = *serverRPCPtr
	serverHTTP = *serverHTTPPtr
}
