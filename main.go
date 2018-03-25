package main

import (
	"log"
	"net"
	"os"

	"github.com/mhausenblas/yages/yages"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ipnport := "0.0.0.0:9000"
	if ie := os.Getenv("YAGES_BIND"); ie != "" {
		ipnport = ie
	}
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
	gs := grpc.NewServer()
	yages.RegisterEchoServer(gs, s)
	reflection.Register(gs)
	return gs.Serve(ln)
}

// Send returns the same message it received.
func (s *Server) Send(ctx context.Context, msg *yages.Content) (*yages.Content, error) {
	return &yages.Content{Text: msg.Text}, nil
}
