package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	host = flag.String("host", "127.0.0.1", "http server host")
	port = flag.Uint("port", 8080, "http server port")
	dir  = flag.String("dir", ".", "directory to serve")
)

func main() {
	flag.Parse()

	fs := http.FileServer(http.Dir(*dir))
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Printf("[%s] %s\n", req.Method, req.URL.Path)
		fs.ServeHTTP(w, req)
	}))

	hostport := fmt.Sprintf("%s:%d", *host, *port)
	log.Printf("Serving directory %q on %q ...\n", *dir, hostport)

	if err := http.ListenAndServe(hostport, nil); err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}
