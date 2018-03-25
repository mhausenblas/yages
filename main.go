package main

import (
	"log"
	"net"
	"os"
	"sort"
	"strings"

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

// Send returns the same message it received.
func (s *Server) Simple(ctx context.Context, msg *yages.Content) (*yages.Content, error) {
	return &yages.Content{Text: msg.Text}, nil
}

// Reverse returns the message it received in reverse order.
func (s *Server) Reverse(ctx context.Context, msg *yages.Content) (*yages.Content, error) {
	r := strings.Split(msg.Text, " ")
	sort.Sort(sort.Reverse(sort.StringSlice(r)))
	return &yages.Content{Text: strings.Join(r, " ")}, nil
}
