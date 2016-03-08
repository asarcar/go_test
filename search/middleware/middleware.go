package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	pb "github.com/asarcar/go_test/search/protos"

	"golang.org/x/net/context"
	"golang.org/x/net/trace"

	"google.golang.org/grpc"
)

const (
	kHTTPAddr          = "localhost" // address of local HTTP service
	kProtocol          = "tcp"
	kLocalRPCAddr      = "localhost" // address of local RPC service
	kRemoteRPCAddr     = "localhost" // address of remote RPC service
	kLocalRPCPort      = 5000        // port of local RPC service
	kRemoteRPCPortBase = 4000        // port of remote RPC service
	kHTTPPort          = 8880        // port of local HTTP service
	kNumBackendServers = 1           // default number of backend servers
)

var (
	localRPCAddr      string
	httpAddr          string
	numBackendServers int      // number of backend servers
	remoteRPCAddrs    []string // RPC address of remote backend servers
)

type server struct {
	backends []pb.GoogleClient
}

func main() {
	parseFlags()

	fmt.Println("Middleware Spawned: ")
	fmt.Printf("    Local-RPC-Addr=\"%s\"\n", localRPCAddr)
	fmt.Printf("    HTTP-Addr=\"%s\"\n", httpAddr)
	fmt.Printf("    Backend Servers: [#: %d]\n", numBackendServers)

	go spawnHTTPServer()
	svr := server{
		backends: make([]pb.GoogleClient, numBackendServers),
	}

	for index, remRPCAddr := range remoteRPCAddrs {
		fmt.Printf("        RemoteRPCAddr[%d]=%s\n", index, remRPCAddr)
		var conn *grpc.ClientConn
		conn, svr.backends[index] = dialRPCServer(remRPCAddr)
		defer conn.Close()
	}
	spawnRPCServer(localRPCAddr, &svr)
}

func spawnHTTPServer() {
	err := http.ListenAndServe(httpAddr, nil)
	log.Fatal(err)
}

func dialRPCServer(rpcAddr string) (*grpc.ClientConn, pb.GoogleClient) {
	// Connect to Google Search RPC server:
	conn, err := grpc.Dial(rpcAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	return conn, pb.NewGoogleClient(conn)
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

type result struct {
	res *pb.Results
	err error
}

func (s *server) Search(ctx context.Context, req *pb.Request) (*pb.Results, error) {
	c := make(chan result, len(s.backends))
	for i, b := range s.backends {
		go func(index int, backend pb.GoogleClient) {
			if tr, ok := trace.FromContext(ctx); ok {
				tr.LazyPrintf("Request: Search-Query \"%s\": Backend[%d]\n",
					req.Query, index)
			}
			res, err := backend.Search(ctx, req)
			c <- result{res, err}
		}(i, b)
	}
	earliest_result := <-c
	tr, ok := trace.FromContext(ctx)
	if ok && earliest_result.err == nil && len(earliest_result.res.Res) > 0 {
		tr.LazyPrintf("Result: Query \"%s\": Res[0].Title \"%s\"\n",
			req.Query, earliest_result.res.Res[0].Title)
	}
	return earliest_result.res, earliest_result.err
}

func parseFlags() {
	localRPCAddrPtr := flag.String("rpcaddr",
		fmt.Sprintf("%s:%d", kLocalRPCAddr, kLocalRPCPort),
		"local RPC \"addr:port\" where clients connect")
	httpAddrPtr := flag.String("http",
		fmt.Sprintf("%s:%d", kHTTPAddr, kHTTPPort),
		"HTTP \"addr:port\" to connect")
	numBackendServersPtr := flag.Int("numbackends", kNumBackendServers,
		"# of server RPC endpoints where clients connect")
	flag.Parse()

	localRPCAddr = *localRPCAddrPtr
	httpAddr = *httpAddrPtr
	numBackendServers = *numBackendServersPtr
	remoteRPCAddrs = make([]string, numBackendServers, numBackendServers)
	for index, _ := range remoteRPCAddrs {
		remoteRPCAddrs[index] =
			fmt.Sprintf("%s:%d", kRemoteRPCAddr, kRemoteRPCPortBase+index)
	}
}
