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
	"grp/limits"
	"grp/middlesitter"
	"os"

	//	"runtime"
	//. "grp/ttd"
	"strings"
)

const (
	apacheAddr = "http://localhost:8000" // Apache server address
)

var proxyAddr string

func main() {

	//TTD(1, "Hi", "you", "and", "you")

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
	proxy.Transport = middlesitter.NewMiddleSitterTransport(http.DefaultTransport)

	limits.SetupRateLimiters()

	// Start the proxy server
	proxyAddr = os.Args[1]
	fmt.Printf("Starting proxy server on %s\n", proxyAddr)

	if strings.HasSuffix(proxyAddr, "443") {
		certFile := "/etc/letsencrypt/live/003232.xyz/fullchain.pem"
		keyFile := "/etc/letsencrypt/live/003232.xyz/privkey.pem"
		http.ListenAndServeTLS(proxyAddr, certFile, keyFile, proxy)
	} else {
		log.Fatal(http.ListenAndServe(proxyAddr, proxy))
	}
}

func analyzeRequest(req *http.Request) {
	fmt.Printf("Received request: %s %s\n", req.Method, req.URL.Path)
}
