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
//
//	"strings" 
	"bytes"
	"io/ioutil"
)

// see function RoundTrip
type MiddleSitterTransport struct {
	originalTransport http.RoundTripper
}

// Lets us observe (Observable) data (Read) that is flowing through, and hadle the stream close (Closer) 
type ObservableReadCloser struct {
    io.ReadCloser // Embed the original ReadCloser.
    dataObserved   []byte
    
    //Here we also keep our notes about the request and response
    request string 
    ip string
    openrouterId string
    
    //Functions that need to be executed upon stream close. Used for example by ratelimiter
    releasers []func()
    streamObservers []func([]byte,int64)
}


/*
    Here we actually facilitate the middlesitting.
    We get request from client, forwards it to apache, get response and forward (backward?) it to the client.
    Note that we get the body stream and forward it back to client (meaning the body data might be still flowing, the headers we have though).
*/
func (t *MiddleSitterTransport) RoundTrip(req *http.Request) (*http.Response, error) {

    usertype, iporid := GetUser(req) 
    
    //ip := strings.Split(req.RemoteAddr, ":")[0]
    
    orc := &ObservableReadCloser{ip:iporid}
    orc.request = req.URL.String()
    
    // **********  Select rate limiters based on url and user
    
    rateLimiter,coLimiter := PathUserRateLimitersSelect(req.URL.Path,usertype)
    fmt.Println("")
    fmt.Println("Path",       req.URL.Path)
    fmt.Println("Usertype",   usertype)
    fmt.Println("iporid",     iporid)    
    fmt.Println("Limiter nr", rateLimiter.nr)
       
       
    // **********   Aplly ratelimiter ********************************
    if allowed, ecode, ertext := rateLimiter.Allow(iporid,orc);!allowed {       		
			orc.ReleaseAll()
			return MakeHttpErrorResponse(ecode,ertext)
    }
    if allowed, ecode, ertext := coLimiter.Allow("ALLTOGETHER",orc);!allowed{
			orc.ReleaseAll()
			return MakeHttpErrorResponse(ecode,"For all IP (non logged in) users combined (consider logging in):"+ertext);  
    }
    
   
    
	// *************** Send the request to the Apache server *********
	resp, err := t.originalTransport.RoundTrip(req)
	if err != nil {
	    orc.ReleaseAll()
		return nil, err
	}

	// *************** Meter bytes? *********************	
	_, meterBytes := resp.Header["Meter-Bytes"];
	if meterBytes {
	    //fmt.Println("Metering bytes")
	    rateLimiter.Addbytes(iporid,req.ContentLength) //can this be tricked by the user?
	    meterfunc := func(data []byte, n int64){
	       rateLimiter.Addbytes(iporid,n)
	    }
	    orc.AddStreamObserver(meterfunc)
	}
	
	// *************** Send reply back to client *********************
	
	
	orc.ReadCloser = resp.Body
	resp.Body = orc
	analyzeResponse(resp)
    
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
	
    return resp, nil

} 

func MakeHttpErrorResponse(status int,err string) (*http.Response, error) {
    errorResponse :=  &http.Response{
        StatusCode: status,//429,
        Status:     err,
        Body:       http.NoBody,
        Header:     make(http.Header),
    }    
    
    body := []byte(err)
    errorResponse.Body = ioutil.NopCloser(bytes.NewBuffer(body))
    errorResponse.ContentLength = int64(len(body))
    errorResponse.Header.Set("Error-reason", err)    
    return errorResponse, nil
    
}


/*
    http body flows through it
*/
func (w *ObservableReadCloser) Read(p []byte) (int, error) {
    n, err := w.ReadCloser.Read(p) // Call the original Read method.
    
    for _, streamObserver := range w.streamObservers {
        streamObserver(p,int64(n))
    }
    
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

func (w *ObservableReadCloser) AddStreamObserver(observerfunc func([]byte,int64 )){
    w.streamObservers = append(w.streamObservers, observerfunc)
}




