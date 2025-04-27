package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var (
	port              string
	name              string
	host              string
	logRequestDetails bool
	customEndpoints   string
)

func init() {
	flag.BoolVar(&logRequestDetails, "log-request-details", true, "log request details including body")
	flag.StringVar(&name, "name", getEnv("ECKO_SERVER_NAME", "ecko"), "name of the server, will appear in the response")
	flag.StringVar(&port, "port", getEnv("ECKO_LISTEN_PORT", "8080"), "port number to listen on")
	flag.StringVar(&host, "host", getEnv("ECKO_LISTEN_HOST", "0.0.0.0"), "host to listen")
	flag.StringVar(&customEndpoints, "custom-endpoints", "/echo:echo,/api:api,/status:status", "custom endpoints to listen on, format: /my-hook:echo,/my-api-hook:api,/my-status:status")
}

func main() {
	flag.Parse()

	endpoints := strings.Split(customEndpoints, ",")
	echoEndpoints := []string{}
	apiEndpoints := []string{}
	statusEndpoints := []string{}
	for _, endpoint := range endpoints {
		parts := strings.Split(endpoint, ":")
		if len(parts) != 2 {
			log.Fatalf("Invalid custom endpoint: %s", endpoint)
		}

		if parts[1] == "echo" {
			echoEndpoints = append(echoEndpoints, parts[0])
		} else if parts[1] == "api" {
			apiEndpoints = append(apiEndpoints, parts[0])
		} else if parts[1] == "status" {
			if strings.HasSuffix(parts[0], "/") {
				statusEndpoints = append(statusEndpoints, parts[0])
			} else {
				statusEndpoints = append(statusEndpoints, parts[0]+"/")
			}
		} else {
			log.Fatalf("Invalid custom endpoint: %s", endpoint)
		}
	}

	mux := http.NewServeMux()
	mux.Handle("/", handle(func(w http.ResponseWriter, r *http.Request, _ bool) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, "Available endpoints:")
		if len(echoEndpoints) > 0 {
			_, _ = fmt.Fprintln(w, strings.Join(echoEndpoints, ",")+": echo endpoint, returns request info")
		}
		if len(apiEndpoints) > 0 {
			_, _ = fmt.Fprintln(w, strings.Join(apiEndpoints, ",")+": api endpoint, returns data as JSON")
		}
		if len(statusEndpoints) > 0 {
			statusEndpointHelp := []string{}
			for _, endpoint := range statusEndpoints {
				statusEndpointHelp = append(statusEndpointHelp, endpoint+"{code}")
			}
			_, _ = fmt.Fprintln(w, strings.Join(statusEndpointHelp, ",")+": status endpoint, returns the specified HTTP status code")
		}
	}, logRequestDetails))
	for _, endpoint := range apiEndpoints {
		mux.Handle(endpoint, handle(apiHandler, logRequestDetails))
	}
	for _, endpoint := range echoEndpoints {
		mux.Handle(endpoint, handle(echoHandler, logRequestDetails))
	}
	for _, endpoint := range statusEndpoints {
		mux.Handle(endpoint, handle(makeStatusHandler(endpoint), logRequestDetails))
	}

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
