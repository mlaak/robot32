package main

import (
	"fmt"
//	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
//	"errors"
//	"strings" 
//	"bytes"
//	"io/ioutil"
    "os"
)

const (
	apacheAddr = "http://localhost:8000" // Apache server address
)

var proxyAddr string

func main() {

	apacheURL, err := url.Parse(apacheAddr)
	if err != nil {
		log.Fatal("Error parsing Apache URL:", err)
	}

	// Create a reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(apacheURL)

	// Customize the director function to modify the request if needed
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		analyzeRequest(req)
	}

	// Create a custom transport to intercept the response (see middlesitter.go !)
	proxy.Transport = &MiddleSitterTransport{
		originalTransport: http.DefaultTransport,
	}

	SetupRateLimiters();
    
	// Start the proxy server
	proxyAddr = os.Args[1]
	fmt.Printf("Starting proxy server on %s\n", proxyAddr)
	log.Fatal(http.ListenAndServe(proxyAddr, proxy))
}

func analyzeRequest(req *http.Request) {
	fmt.Printf("Received request: %s %s\n", req.Method, req.URL.Path)
	// Add more analysis as needed
}

func analyzeResponse(resp *http.Response) {
	fmt.Printf("Received response: Status %s, Content-Length %d\n", resp.Status, resp.ContentLength)
}


