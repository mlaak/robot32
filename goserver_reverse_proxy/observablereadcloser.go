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