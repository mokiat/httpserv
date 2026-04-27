package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const shutdownTimeout = 5 * time.Second

var (
	host = flag.String("host", "127.0.0.1", "http server host")
	port = flag.String("port", "8080", "http server port")
	dir  = flag.String("dir", ".", "directory to serve")
)

func main() {
	flag.Parse()

	ctxSignal, ctxStop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer ctxStop()

	log.Println("Starting server...")
	if err := run(ctxSignal); err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Println("Server stopped.")
}

func run(ctx context.Context) error {
	rootDir, err := os.OpenRoot(*dir)
	if err != nil {
		return fmt.Errorf("failed to open %q directory: %w", *dir, err)
	}
	defer rootDir.Close()

	fs := http.FileServerFS(rootDir.FS())

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("[%s] %s", req.Method, req.URL.Path)
		fs.ServeHTTP(w, req)
	})

	hostport := net.JoinHostPort(*host, *port)
	log.Printf("Serving directory %q on %q ...", *dir, hostport)

	server := &http.Server{
		Addr:    hostport,
		Handler: mux,
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
		log.Printf("Shutting down server...")
		ctxShutdown, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := server.Shutdown(ctxShutdown); err != nil {
			return fmt.Errorf("server shutdown error: %w", err)
		}
		return nil
	}
}
