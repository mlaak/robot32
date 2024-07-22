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
package middlesitter

import (
	"fmt"
	"net/http"
	"bytes"
	"io/ioutil"
    "grp/situation"
    "grp/limits"
    . "grp/ratelimiter"  
    "grp/usersession"
    . "grp/translator"
    
)

// see function RoundTrip
type MiddleSitterTransport struct {
	OriginalTransport http.RoundTripper
    GetUser func(*http.Request)(string,string) 
    GetLimitersForPathAndUserType func(string,string) (IRateLimiter,IRateLimiter)
    Observers IObservableReadCloser
    RequestContext situation.IRequestContext

}

func NewMiddleSitterTransport(Original http.RoundTripper)(*MiddleSitterTransport){
    t := &MiddleSitterTransport{
	    OriginalTransport:Original,
        GetUser: usersession.GetUser,
        GetLimitersForPathAndUserType: limits.GetLimitersForPathAndUserType,
        Observers: NewObservableReadCloser(),
        RequestContext: situation.NewRequestContext(),
    }
    return t;
}

/**************************** ROUNDTRIP ***************************************
    Here we actually facilitate the middlesitting.
    We get request from client, forwards it to apache, get response and forward
    (backward?) it to the client. Note that we get the body -stream- and forward
    it back to client (meaning the body data might be still flowing).  
    The headers we have though.
*******************************************************************************/

func (t *MiddleSitterTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    

    
    usertype, iporid := t.GetUser(req) 
    
    context := situation.NewRequestContext()
    //context := t.RequestContext
    context.SetIporid(iporid)

    //TODO: fix the bug here and for goodness sake, write some integration tests
    //observers := t.Observers;
    observers := NewObservableReadCloser()
    observers.SetRequestStr(req.URL.String())

    rateLimiter,coLimiter := t.GetLimitersForPathAndUserType(req.URL.Path,usertype)

    //fmt.Println(rateLimiter.GetNr());

    if allowed, ecode, ertext := rateLimiter.Allow(iporid,context); allowed {       		
        countDownOneConnectionFunc := rateLimiter.CountUpOneConnection(iporid);
        observers.AddOnCloseFunc(countDownOneConnectionFunc)
    } else {
        fmt.Println(rateLimiter.GetNr());
        observers.CallAllOnCloseFuncs()
		return MakeHttpErrorResponse(ecode,ertext)
    }

    if allowed, ecode, ertext := coLimiter.Allow("ALLTOGETHER",context); allowed{
        countDownOneConnectionFunc := coLimiter.CountUpOneConnection(iporid);
        observers.AddOnCloseFunc(countDownOneConnectionFunc)  
    } else {
        observers.CallAllOnCloseFuncs()
		return MakeHttpErrorResponse(ecode,TR("For all IP (non logged in) users combined (consider logging in):",context)+ertext);
    }


    // forward request to apache
	resp, err := t.OriginalTransport.RoundTrip(req)
	if err != nil {
	    observers.CallAllOnCloseFuncs()
		return nil, err
	}


	_, shouldWeMeterBytes := resp.Header["Meter-Bytes"];
	if shouldWeMeterBytes {
        fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!");
        fmt.Println(req.URL.String());
        
	    rateLimiter.Addbytes(iporid,req.ContentLength) //add request bytes. Can this be tricked by the user?
	    meterfunc := func(data []byte, n int64){
	       rateLimiter.Addbytes(iporid,n) //add downloaded bytes
	    }
	    observers.AddStreamObserver(meterfunc)
	}
    
    observers.SetReadCloser(resp.Body)
	//observers.ReadCloser = resp.Body
	resp.Body = observers
	
	return resp, nil
} 


func MakeHttpErrorResponse(status int,err string) (*http.Response, error) {
    errorResponse :=  &http.Response{
        StatusCode: status,
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
    /*for key, values := range resp.Header {
        if(key=="Openrouter-Id"){ orc.openrouterId = values[0]; }
        fmt.Println("  "+key, values)
    }*/
}
