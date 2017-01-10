package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"go4.org/net/throttle"
)

var (
	throttleLatency = flag.Duration("latency", 10*time.Millisecond, "latency")
	throttleRate    = flag.Int("rate", 10, "rate in Kbps")
	listenAddr      = flag.String("port", ":3000", "listen address")
)

// Server is a http server.
type Server struct {
	mux      *http.ServeMux
	listener net.Listener
}

// New creates a Server type.
func New() *Server {
	return &Server{
		mux: http.NewServeMux(),
	}
}

// HandleFunc implements ServeMux.HandleFunc method.
func (s *Server) HandleFunc(pattern string, fn func(http.ResponseWriter, *http.Request)) {
	s.mux.HandleFunc(pattern, fn)
}

func (s *Server) Handle(pattern string, handler http.Handler) {
	s.mux.Handle(pattern, handler)
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	s.mux.ServeHTTP(rw, req)
}

// Listen starts listening on the given host:port addr.
func (s *Server) Listen(addr string) error {
	if s.listener != nil {
		return nil
	}
	var err error
	s.listener, err = net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("Failed to listen on %s: %v", addr, err)
	}
	return nil
}

func (s *Server) Serve() {
	if err := s.Listen(""); err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Handler: s,
	}

	err := srv.Serve(s.throttleListener())
	if err != nil {
		log.Printf("Error in http server: %v\n", err)
		os.Exit(1)
	}

}

func (s *Server) throttleListener() net.Listener {
	rate := throttle.Rate{
		KBps:    *throttleRate,
		Latency: *throttleLatency,
	}

	return &throttle.Listener{
		Listener: s.listener,
		Down:     rate,
		Up:       rate,
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!! the time is %v\n", time.Now().UTC())
}

func main() {
	flag.Parse()
	s := New()

	s.HandleFunc("/", indexHandler)

	if err := s.Listen(*listenAddr); err != nil {
		log.Fatalf("Listen: %v", err)
	}

	log.Printf("Server started on port %s", *listenAddr)

	s.Serve()
}
