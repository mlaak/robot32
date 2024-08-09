/*
**********************************************************

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

**********************************************************
*/

package middlesitter

import (
	"bytes"
	"fmt"
	"grp/limits"

	rl "grp/ratelimiter"
	"grp/situation"
	. "grp/translator" //lint:ignore ST1001 Sometimes dot imports are needed

	. "grp/ttd" //lint:ignore ST1001 Sometimes dot imports are needed
	"grp/usersession"
	"io/ioutil"
	"net/http"
)

// see function RoundTrip
type MiddleSitterTransport struct {
	OriginalTransport             http.RoundTripper
	GetUser                       func(int64, *http.Request) (int64, string, string)
	GetLimitersForPathAndUserType func(int64, string, string) (rl.IRateLimiter, rl.IRateLimiter)
	Observers                     IObservableReadCloser
	RequestContext                situation.IRequestContext
}

func NewMiddleSitterTransport(Original http.RoundTripper) *MiddleSitterTransport {
	t := &MiddleSitterTransport{
		OriginalTransport:             Original,
		GetUser:                       usersession.GetUser,
		GetLimitersForPathAndUserType: limits.GetLimitersForPathAndUserType,
		Observers:                     NewObservableReadCloser(),
		RequestContext:                situation.NewRequestContext(),
	}
	return t
}

/**************************** ROUNDTRIP ***************************************
    Here we actually facilitate the middlesitting.
    We get request from client, forwards it to apache, get response and forward
    (backward?) it to the client. Note that we get the body -stream- and forward
    it back to client (meaning the body data might be still flowing).
    The headers we have though.
*******************************************************************************/

func (t *MiddleSitterTransport) RoundTrip(req *http.Request) (*http.Response, error) {

	nrc := situation.NewRequestContext()
	c := nrc.GetReqCode()
	c = TTDLEV(c, 1)
	/**/ TTD(c, "Starting Request")

	usertype, iporid := TTX2(t.GetUser(c, req))
	nrc.SetIporid(iporid)
	fmt.Println(" - " + nrc.GetReqCodeStr() + " " + req.URL.String())
	/**/ TTD(c, "Got request", "URL", req.URL.String(), "usertype", usertype, "iporid", iporid)

	//TODO: fix the bug here and for goodness sake, write some integration tests
	//observers := t.Observers;
	observers := NewObservableReadCloser()
	observers.SetRequestStr(TTX(c, req.URL.String()))

	rateLimiter, coLimiter := t.GetLimitersForPathAndUserType(c, req.URL.Path, usertype)
	/**/ TTD(c, "Selected ratelimiters", "rateLimiter", rateLimiter.GetNr(), "coLimiter", coLimiter.GetNr())

	if allowed, ecode, ertext := rateLimiter.Allow(c, iporid); allowed {
		countDownOneConnectionFunc := rateLimiter.CountUpOneConnection(c, iporid)
		observers.AddOnCloseFunc(c, countDownOneConnectionFunc)
		/**/ TTD(c, "Allowed for ratelimiter, continueing", "Limiter", rateLimiter.GetNr(), "IporID", iporid, "url", req.URL.String())
	} else {
		observers.CallAllOnCloseFuncs(c)
		/**/ TTD(c, "Failed ratelimiter, sending error response to user", "Limiter", rateLimiter.GetNr(), "Ecode", ecode, "Ertext", ertext)
		return MakeHttpErrorResponse(c, ecode, ertext)
	}

	if allowed, ecode, ertext := coLimiter.Allow(c, "ALLTOGETHER"); allowed {
		countDownOneConnectionFunc := coLimiter.CountUpOneConnection(c, iporid)
		observers.AddOnCloseFunc(c, countDownOneConnectionFunc)
		/**/ TTD(c, "Allowed for colimiter, continueing", "Limiter", coLimiter.GetNr(), "IporID", iporid, "url", req.URL.String())
	} else {
		observers.CallAllOnCloseFuncs(c)
		/**/ TTD(c, "Failed colimiter, sending error response to user", "Limiter", coLimiter.GetNr(), "Ecode", ecode, "Ertext", ertext)
		return MakeHttpErrorResponse(c, ecode, TR(c, "For all IP (non logged in) users combined (consider logging in):")+ertext)
	}

	observers.AddOnCloseFunc(c, func() {
		/**/ TTD(c, "CLOSED CONNECTION")
	})

	req.Header.Add("Ttd", nrc.GetReqCodeStr())

	/**/
	TTD(c, "Forward request to apache") //gofmt:ignore
	resp, err := t.OriginalTransport.RoundTrip(req)
	if err != nil {
		observers.CallAllOnCloseFuncs(c)
		/**/ TTD(c, "Failed reaching Apache", "Err", err)
		return nil, err
	}
	/**/ TTD(c, "Successfully got response from Apache", "Headers", resp.Header)

	_, shouldWeMeterBytes := resp.Header["Meter-Bytes"]
	resp.Header.Add("Req-Code", nrc.GetReqCodeStr())
	if shouldWeMeterBytes {

		rateLimiter.Addbytes(c, iporid, req.ContentLength) //add request bytes. Can this be tricked by the user?
		meterfunc := func(data []byte, n int64) {
			rateLimiter.Addbytes(c, iporid, n) //add downloaded bytes
			/**/ TTD(c, "Bytes added to limiter", "Limiter", rateLimiter.GetNr(), "Bytes", n, "IporID", iporid)
		}
		observers.AddStreamObserver(c, meterfunc)
	}

	observers.SetReadCloser(resp.Body)
	resp.Body = observers
	/**/ TTD(c, "Done proccessing, now send headers and body should be flowing soon")

	return resp, nil
}

func MakeHttpErrorResponse(c int64, status int, err string) (*http.Response, error) {
	errorResponse := &http.Response{
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
/*func debugPrint(req *http.Request, resp *http.Response, orc *ObservableReadCloser) {
	fmt.Println(req.URL.String())
	fmt.Println("Client headers:")
	_ = resp
	_ = orc
	for key, values := range req.Header {
		fmt.Println("  "+key, values)
	}

	fmt.Println("Server headers:")
	for key, values := range resp.Header {
	    if(key=="Openrouter-Id"){ orc.openrouterId = values[0]; }
	    fmt.Println("  "+key, values)
	}
}*/
