package situation

import (
//	"fmt"
//	"net/http"
//	"sync"
//	"time"
//	"strconv"
)

type RequestContext struct {
	iporid string
}

func NewRequestContext(iporid string) *RequestContext{
	return &RequestContext{
		iporid: iporid,
	}
}