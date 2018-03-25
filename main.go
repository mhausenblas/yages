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

var (
	release string
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
	if release == "" {
		release = "dev"
	}
	log.Printf("YAGES in version %v serving on %v is ready for gRPC clients …", release, s.bind)
	return gs.Serve(ln)
}

// Ping returns a "pong" (constant message).
func (s *Server) Ping(ctx context.Context, _ *yages.Empty) (*yages.Content, error) {
	return &yages.Content{Text: "pong"}, nil
}

// Reverse returns the message it received in reverse order.
func (s *Server) Reverse(ctx context.Context, msg *yages.Content) (*yages.Content, error) {
	revstr := func(s string) string {
		r := []rune(s)
		for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
			r[i], r[j] = r[j], r[i]
		}
		return string(r)
	}
	return &yages.Content{Text: revstr(msg.Text)}, nil
}
