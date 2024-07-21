/***********************************************
 011001110001110011010011110011100111100111001
     |                                  |
    xxx                                 xx
   x   x                               x x
   x   x                                 x
    xxx  This one            This one    x
     |O/ looks alright           also  \O|      
      |                                 |
     / \                               / \
***********************************************

It lets us observe the bytes flying by (for ex. HTTP stream) 
and we also get notified if connection closes. 

*/

package main

import (
	"fmt"
	"io"
//	"log"
//	"net/http"
//	"net/http/httputil"
//	"net/url"
//	"errors"
//
//	"strings" 
//	"bytes"
//	"io/ioutil"
)



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
    Registers a function that is called when stream closes
*/
func (w *ObservableReadCloser) AddReleaser(releasefunc func()) {
    w.releasers = append(w.releasers, releasefunc)
}

/*
    Registers a function that is called when chunk of data moves in stream
*/
func (w *ObservableReadCloser) AddStreamObserver(observerfunc func([]byte,int64 )){
    w.streamObservers = append(w.streamObservers, observerfunc)
}

/*
 Call all releaser functions
*/
func (w *ObservableReadCloser) ReleaseAll() {
    for _, releaser := range w.releasers {
        releaser()
    }
}

/*
    Gets automatically called when stream is closed
    NB. Might not be called if server error?
*/    
func (w *ObservableReadCloser) Close() error {
    //fmt.Println("CLOSING!")
    w.ReleaseAll()
   return w.ReadCloser.Close() // Call the original Close method.
}

/*
    Gets automatically called when stream data flows.
   Stream data (eg http body) flows through it
*/
func (w *ObservableReadCloser) Read(p []byte) (int, error) {
    n, err := w.ReadCloser.Read(p) // Call the original Read method.
    
    for _, streamObserver := range w.streamObservers {
        streamObserver(p,int64(n))
    }
    w.dataObserved = append(w.dataObserved, p[:n]...)
    return n, err
}