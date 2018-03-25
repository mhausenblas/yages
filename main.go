package main

import (
	"flag"
	"log"
	"net"

	"github.com/mhausenblas/yages/yages"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	var ipnport string
	flag.StringVar(&ipnport, "ipnport", "0.0.0.0:9000", "IP and port to bind to in the format `ip:port` with a default of `0.0.0.0:9000`.")
	flag.Parse()
	err := New(ipnport).Listen()
	if err != nil {
		log.Fatalf("Failed to launch the YAGES due to %v", err)
	}
}

// Server wraps the gRPC server and implements yages.EchoServer
type Server struct {
	bind string
}

// New creates a new rpc server.
func New(bind string) *Server {
	return &Server{bind}
}

// Listen binds the server to the indicated interface:port.
func (s *Server) Listen() error {
	ln, err := net.Listen("tcp", s.bind)
	if err != nil {
		return err
	}
	srv := grpc.NewServer()
	yages.RegisterEchoServer(srv, s)
	return srv.Serve(ln)
}

// Send returns the same message it received.
func (s *Server) Send(ctx context.Context, msg *yages.Content) (*yages.Content, error) {
	return &yages.Content{Text: msg.Text}, nil
}
