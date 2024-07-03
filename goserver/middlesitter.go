/*
    Sits in the middle of the client and the apache server
    Uses ratelimiter to limit usage.
    Can use cache.
*/


package main

import (
	"fmt"
	"io"
//	"log"
	"net/http"
//	"net/http/httputil"
//	"net/url"
//	"errors"
	"strings" 
	"bytes"
	"io/ioutil"
)

// see function RoundTrip
type MiddleSitterTransport struct {
	originalTransport http.RoundTripper
}

// Lets us observe (Observable) data (Read) flowing through, and hadle the stream close (Closer) 
type ObservableReadCloser struct {
    io.ReadCloser // Embed the original ReadCloser.
    dataObserved   []byte
    
    //Here we also keep our notes about the request and response
    request string 
    ip string
    openrouterId string
    
    //Functions that need to be executed upon stream close. Used for example by ratelimiter
    releasers []func()
}


/*
    http body flows through it
*/
func (w *ObservableReadCloser) Read(p []byte) (int, error) {
    n, err := w.ReadCloser.Read(p) // Call the original Read method.
    if n<1000 {
        //fmt.Println(p[:n])
    }
    w.dataObserved = append(w.dataObserved, p[:n]...)
    return n, err
}


/*
    When stream is closed
    
    NB. Might not be called if server error?
*/    
func (w *ObservableReadCloser) Close() error {
     fmt.Println("CLOSING!")
     w.ReleaseAll()
     //overallIPLimiter.Release("IPS")
     //rateLimiter.Release(w.ip)
    return w.ReadCloser.Close() // Call the original Close method.
}


func (w *ObservableReadCloser) ReleaseAll() {
    for _, releaser := range w.releasers {
        releaser()
    }
}

func (w *ObservableReadCloser) AddReleaser(releasefunc func()) {
    w.releasers = append(w.releasers, releasefunc)
}




/*
    Here we actually facilitate the middlesitting.
    We get request from client, forwards it to apache, get response and forward (backward?) it to the client.
    Note that we get the body stream and forward it back to client (meaning the body data might be still flowing, the headers we have though).
*/
func (t *MiddleSitterTransport) RoundTrip(req *http.Request) (*http.Response, error) {

    ip := strings.Split(req.RemoteAddr, ":")[0]
    orc := &ObservableReadCloser{ip:ip}
    
    // **********    Just incase, lets already prepare the error response *******
    
    errorResponse :=  &http.Response{
        StatusCode: 429,
        Status:     "Requests limit exeeded (minute or hourly or daily)",
        Body:       http.NoBody,
        Header:     make(http.Header),
    }
    
    body := []byte("Requests limit exeeded (minute or hourly or daily)")
    errorResponse.Body = ioutil.NopCloser(bytes.NewBuffer(body))
    errorResponse.ContentLength = int64(len(body))
  
    
    // **********   Aplly ratelimiter *******************************************
    
    if allowed, ertext := rateLimiter.Allow(ip,orc);!allowed {
			//http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
	        errorResponse.Header.Set("Error-reason", ertext)		
			return errorResponse, nil//errors.New("Rate limit exceeded") 
    }
        
    if allowed, ertext := overallIPLimiter.Allow("IPS",orc);!allowed{
			//http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			errorResponse.Header.Set("Error-reason", "For all IP users (consider logging in): "+ertext)
			//rateLimiter.Release(ip)
			orc.ReleaseAll()
			return errorResponse, nil//errors.New("Rate limit exceeded for non-logged in users") 
    }
    

	// *************** Send the request to the Apache server ***********************
	
	resp, err := t.originalTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// *************** Send reply back to client ***********************************
	
	orc.ReadCloser = resp.Body
	analyzeResponse(resp)
    orc.request = req.URL.String()
    orc.ip = ip
    resp.Body = orc
    
	return resp, nil
	
	// *** Stuff useful for debugging ***
	
	/*fmt.Println(req.URL.String());
    fmt.Println("Client headers:")
    for key, values := range req.Header {
        fmt.Println("  "+key, values)
    }
	
	fmt.Println("Server headers:")
    for key, values := range resp.Header {
        if(key=="Openrouter-Id"){ orc.openrouterId = values[0]; }
        fmt.Println("  "+key, values)
    }*/
	
	// Clone the response body
	/*bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	// Create a new response with the cloned body
	resp.Body = io.NopCloser(io.NewReader(bodyBytes))
	*/
	
    //return resp, nil

}
