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

package middlesitter

import (
	"io"
)



type IObservableReadCloser interface {
    AddOnCloseFunc(onCloseFunc func())
    AddStreamObserver(observerFunc func([]byte,int64 ))
    SetReadCloser(ReadCloser io.ReadCloser)
    SetRequestStr(rstr string)

    CallAllOnCloseFuncs()
    Close() error
    Read(p []byte) (int, error)
}

// Lets us observe (Observable) data (Read) that is flowing through, and hadle the stream close (Closer) 
type ObservableReadCloser struct {
    ReadCloser io.ReadCloser // Embed the original ReadCloser.
    
    dataObserved   []byte

    rstr string
    //Functions that need to be executed upon stream close. Used for example by ratelimiter
    onCloseFuncs []func()
    streamObservers []func([]byte,int64)
}

func NewObservableReadCloser() (IObservableReadCloser){
    return &ObservableReadCloser{}
}

/*
    Registers a function that is called when stream closes
*/
func (w *ObservableReadCloser) AddOnCloseFunc(onCloseFunc func()) {
    w.onCloseFuncs = append(w.onCloseFuncs, onCloseFunc)
}


func (w *ObservableReadCloser) SetReadCloser(ReadCloser io.ReadCloser){
    w.ReadCloser = ReadCloser
}

func (w *ObservableReadCloser) SetRequestStr(rstr string) {
    w.rstr = rstr
}

/*
    Registers a function that is called when chunk of data moves in stream
*/
func (w *ObservableReadCloser) AddStreamObserver(observerFunc func([]byte,int64 )){
    w.streamObservers = append(w.streamObservers, observerFunc)
}

/*
 Call all releaser functions
*/
func (w *ObservableReadCloser) CallAllOnCloseFuncs() {
    for key, onCloseFunc := range w.onCloseFuncs {
        if(onCloseFunc != nil){
            onCloseFunc()
        }
        w.onCloseFuncs[key] = nil
    }
}

/*
    Gets automatically called when stream is closed
    NB. Might not be called if server error?
*/    
func (w *ObservableReadCloser) Close() error {
    //fmt.Println("CLOSING!")
    w.CallAllOnCloseFuncs()
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