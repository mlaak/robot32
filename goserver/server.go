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
)

const (
	proxyAddr  = ":8080"
	apacheAddr = "http://localhost:8000" // Apache server address
)


var rateLimiter *RateLimiter;
var overallIPLimiter *RateLimiter;

func main() {
	// Parse the URL of the Apache server
	
	rateLimiter =      NewRateLimiter(  60,  300,  3000,  5,     5000,   50000,   500000)
	overallIPLimiter = NewRateLimiter(6000,60000,600000,500,   500000, 5000000, 50000000 )
	
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

	// Create a custom transport to intercept the response
	proxy.Transport = &MiddleSitterTransport{
		originalTransport: http.DefaultTransport,
	}

	// Start the proxy server
	fmt.Printf("Starting proxy server on %s\n", proxyAddr)
	log.Fatal(http.ListenAndServe(proxyAddr, proxy))
}

func analyzeRequest(req *http.Request) {
	// Analyze the request here
	fmt.Printf("Received request: %s %s\n", req.Method, req.URL.Path)
	// Add more analysis as needed
}

func analyzeResponse(resp *http.Response) {
	// Analyze the response here
	fmt.Printf("Received response: Status %s, Content-Length %d\n", resp.Status, resp.ContentLength)
	// Add more analysis as needed
}


