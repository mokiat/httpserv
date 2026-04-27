package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
)

var (
	host = flag.String("host", "127.0.0.1", "http server host")
	port = flag.Uint("port", 8080, "http server port")
	dir  = flag.String("dir", ".", "directory to serve")
)

func main() {
	flag.Parse()

	ctxSignal, ctxStop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer ctxStop()

	log.Println("Starting server...")
	if err := run(ctxSignal); err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	log.Println("Server stopped.")
}

func run(ctx context.Context) error {
	fs := http.FileServer(http.Dir(*dir))
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("[%s] %s\n", req.Method, req.URL.Path)
		fs.ServeHTTP(w, req)
	})

	hostport := fmt.Sprintf("%s:%d", *host, *port)
	log.Printf("Serving directory %q on %q ...\n", *dir, hostport)

	server := &http.Server{
		Addr: hostport,
	}
	serverErr := make(chan error, 1)
	go func() {
		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			serverErr <- fmt.Errorf("server error: %w", err)
		}
	}()

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		log.Printf("Shutting down server ...\n")
		if err := server.Shutdown(context.Background()); err != nil {
			return fmt.Errorf("server shutdown error: %w", err)
		}
		return nil
	}
}
