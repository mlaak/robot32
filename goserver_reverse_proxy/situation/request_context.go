package situation

import (
	"math/rand"
	"strconv"
	"sync"
	"time"
)

const hexChars = "0123456789abcdef"

func GenerateRandomHex() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 16)
	for i := range b {
		b[i] = hexChars[rand.Intn(len(hexChars))]
	}
	return string(b)
}

type IRequestContext interface {
	SetIporid(iporid string)
	GetIporid() string
	GetReqCode() int64
	GetReqCodeStr() string
}
type RequestContext struct {
	iporid  string
	reqCode int64
}

var mu sync.Mutex

var identNum int64 = 0

func getIdent() int64 {
	mu.Lock()
	defer mu.Unlock()
	if identNum == 0 {
		identNum = time.Now().Unix() * 100000000
	} else {
		identNum++
	}
	return identNum * 8
}

func NewRequestContext() IRequestContext {
	rc := &RequestContext{}
	//rc.reqCode = GenerateRandomHex()
	//rc.reqNum
	rc.reqCode = getIdent()
	return rc
}

func (rc *RequestContext) SetIporid(iporid string) {
	rc.iporid = iporid
}

func (rc *RequestContext) GetIporid() string {
	return rc.iporid
}

func (rc *RequestContext) GetReqCode() int64 {
	return rc.reqCode
}

func (rc *RequestContext) GetReqCodeStr() string {
	return strconv.FormatInt(rc.reqCode, 10)
}
