/***********************************************************
  CLIENT            MIDLESITTER
  ____                                  APACHE WEB
 /  O \_            //////                SERVER
(     __\          / /  O \             ////
 \____/ PLA PLAPLA (     __\          //   0\___
   ||               \____/ PLA PLAPLA (   )   __\
                      ||               \ _____/
                                         | |
    Sits in the middle of the client and the apache server
    (in both directions)
    Uses ratelimiter to limit usage.
    Could use cache, loadbalancing etc... (just implement it)
***********************************************************/
package main

import (
	"fmt"
//	"io"
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

/**************************** ROUNDTRIP ***************************************
    Here we actually facilitate the middlesitting.
    We get request from client, forwards it to apache, get response and forward
    (backward?) it to the client. Note that we get the body -stream- and forward
    it back to client (meaning the body data might be still flowing).  
    The headers we have though.
*******************************************************************************/

func (t *MiddleSitterTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    usertype, iporid := GetUser(req) 
    context := NewRequestContext(iporid)
    orc := &ObservableReadCloser{ip:iporid}
    orc.request = req.URL.String()


    rateLimiter,coLimiter := PathUserRateLimitersSelect(req.URL.Path,usertype)

    if allowed, ecode, ertext := rateLimiter.Allow(iporid,orc,context);!allowed {       		
			orc.ReleaseAll()
			return MakeHttpErrorResponse(ecode,ertext)
    }
    if allowed, ecode, ertext := coLimiter.Allow("ALLTOGETHER",orc,context);!allowed{
			orc.ReleaseAll()
			return MakeHttpErrorResponse(ecode,TR("For all IP (non logged in) users combined (consider logging in):",context)+ertext);  
    }

    // forward request to apache
	resp, err := t.originalTransport.RoundTrip(req)
	if err != nil {
	    orc.ReleaseAll()
		return nil, err
	}


	_, shouldWeMeterBytes := resp.Header["Meter-Bytes"];
	if shouldWeMeterBytes {
	    rateLimiter.Addbytes(iporid,req.ContentLength) //add request bytes. Can this be tricked by the user?
	    meterfunc := func(data []byte, n int64){
	       rateLimiter.Addbytes(iporid,n) //add downloaded bytes
	    }
	    orc.AddStreamObserver(meterfunc)
	}
    

	orc.ReadCloser = resp.Body
	resp.Body = orc
	analyzeResponse(resp)
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




// *** Stuff useful for debugging ***
    func debugPrint(req *http.Request,resp *http.Response, orc *ObservableReadCloser){
	fmt.Println(req.URL.String());
    fmt.Println("Client headers:")
    for key, values := range req.Header {
        fmt.Println("  "+key, values)
    }
	
	fmt.Println("Server headers:")
    for key, values := range resp.Header {
        if(key=="Openrouter-Id"){ orc.openrouterId = values[0]; }
        fmt.Println("  "+key, values)
    }
}
