package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	port              string
	name              string
	host              string
	logRequestDetails bool
)

func init() {
	flag.BoolVar(&logRequestDetails, "log-request-details", true, "Log request details including body")
	flag.StringVar(&port, "port", getEnv("ECKO_PORT", "8080"), "give me a port number")
	flag.StringVar(&name, "name", getEnv("ECKO_NAME", "ecko"), "give me a name")
	flag.StringVar(&host, "host", "0.0.0.0", "host to listen on")
}

func main() {
	flag.Parse()

	mux := http.NewServeMux()
	mux.Handle("/", handle(func(w http.ResponseWriter, r *http.Request, _ bool) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, "Available endpoints:")
		_, _ = fmt.Fprintln(w, "/api    - API endpoint, returns data as JSON")
		_, _ = fmt.Fprintln(w, "/echo   - Echo endpoint, returns request info")
		_, _ = fmt.Fprintln(w, "/status/{code} - Returns the specified HTTP status code")
	}, logRequestDetails))
	mux.Handle("/api", handle(apiHandler, logRequestDetails))
	mux.Handle("/echo", handle(echoHandler, logRequestDetails))
	mux.Handle("/status/", handle(statusHandler, logRequestDetails))

	log.Printf("Starting up on %s:%s", host, port)

	log.Fatal(http.ListenAndServe(host+":"+port, mux))
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func handle(next HandlerFunc, logRequestDetails bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := &statusWriter{ResponseWriter: w}

		next(sw, r, logRequestDetails)

		log.Printf("[%s] %s: %d", r.Method, r.URL.Path, sw.status)
	})
}
