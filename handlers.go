package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request, logRequestDetails bool)

type EchoResponse struct {
	Hostname    string      `json:"hostname,omitempty"`
	IP          []string    `json:"ip,omitempty"`
	Headers     http.Header `json:"headers,omitempty"`
	URL         string      `json:"url,omitempty"`
	Host        string      `json:"host,omitempty"`
	Method      string      `json:"method,omitempty"`
	Name        string      `json:"name,omitempty"`
	RemoteAddr  string      `json:"remoteAddr,omitempty"`
	RequestBody string      `json:"requestBody,omitempty"`
}

// buildEchoResponse creates an EchoResponse from the request
func buildEchoResponse(r *http.Request, logRequestDetails bool) EchoResponse {
	queryParams := r.URL.Query()

	wait := queryParams.Get("wait")
	if wait == "" {
		wait = queryParams.Get("sleep")
	}
	if wait != "" {
		waitSeconds, err := strconv.Atoi(wait)
		if err != nil {
			duration, err := time.ParseDuration(wait)
			if err != nil {
				log.Printf("Invalid wait duration: %s", wait)
			} else {
				time.Sleep(duration)
			}
		} else {
			time.Sleep(time.Duration(waitSeconds) * time.Second)
		}
	}

	hostname, _ := os.Hostname()

	// Read request body
	bodyBytes, _ := io.ReadAll(r.Body)
	r.Body.Close()                                              // Close the body after reading
	r.Body = io.NopCloser(strings.NewReader(string(bodyBytes))) // Restore the body for other handlers

	echoResponse := EchoResponse{
		Hostname:    hostname,
		IP:          getClientIP(),
		Headers:     r.Header,
		URL:         r.URL.RequestURI(),
		Host:        r.Host,
		Method:      r.Method,
		Name:        name,
		RemoteAddr:  r.RemoteAddr,
		RequestBody: string(bodyBytes),
	}

	if logRequestDetails {
		t := "---"
		t += fmt.Sprintf("\nHostname: %s", echoResponse.Hostname)
		t += fmt.Sprintf("\nName: %s", echoResponse.Name)
		t += fmt.Sprintf("\nRemoteAddr: %s", echoResponse.RemoteAddr)
		for i, ip := range echoResponse.IP {
			t += fmt.Sprintf("\nIP[%d]: %s", i, ip)
		}
		t += fmt.Sprintf("\nMethod: %s", echoResponse.Method)
		t += fmt.Sprintf("\nURL: %s", echoResponse.URL)
		t += fmt.Sprintf("\nHost: %s", echoResponse.Host)
		for key, values := range echoResponse.Headers {
			if len(values) > 1 {
				for i, value := range values {
					t += fmt.Sprintf("\n%s[%d]: %s", key, i, value)
				}
			} else {
				t += fmt.Sprintf("\n%s: %s", key, values[0])
			}
		}
		if echoResponse.RequestBody != "" {
			t += fmt.Sprintf("\nRequest Body: %s", echoResponse.RequestBody)
		}
		t += "\n---"
		log.Print(t)
	}

	return echoResponse
}

func apiHandler(w http.ResponseWriter, r *http.Request, logRequestDetails bool) {
	data := buildEchoResponse(r, logRequestDetails)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func echoResponseToText(w http.ResponseWriter, r *http.Request, statusCode int, logRequestDetails bool) {
	data := buildEchoResponse(r, logRequestDetails)

	w.WriteHeader(statusCode)

	// Print in text format
	fmt.Fprintf(w, "Status: %d\n", statusCode)
	fmt.Fprintf(w, "Hostname: %s\n", data.Hostname)
	fmt.Fprintf(w, "Name: %s\n", data.Name)
	fmt.Fprintf(w, "RemoteAddr: %s\n", data.RemoteAddr)

	for i, ip := range data.IP {
		fmt.Fprintf(w, "IP[%d]: %s\n", i, ip)
	}

	fmt.Fprintf(w, "Method: %s\n", data.Method)
	fmt.Fprintf(w, "URL: %s\n", data.URL)
	fmt.Fprintf(w, "Host: %s\n", data.Host)

	fmt.Fprintf(w, "\nHeaders:\n")
	for key, values := range data.Headers {
		if len(values) > 1 {
			for i, value := range values {
				fmt.Fprintf(w, "%s[%d]: %s\n", key, i, value)
			}
		} else {
			fmt.Fprintf(w, "%s: %s\n", key, values[0])
		}
	}

	if data.RequestBody != "" {
		fmt.Fprintf(w, "\nRequest Body:\n%s\n", data.RequestBody)
	}
}

func echoHandler(w http.ResponseWriter, r *http.Request, logRequestDetails bool) {
	echoResponseToText(w, r, http.StatusOK, logRequestDetails)
}

func makeStatusHandler(endpoint string) HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request, logRequestDetails bool) {
		code := req.URL.Path[len(endpoint):]
		statusCode, err := strconv.Atoi(code)
		if err != nil {
			http.Error(w, "Invalid status code", http.StatusBadRequest)
			return
		}
		echoResponseToText(w, req, statusCode, logRequestDetails)
	}
}
