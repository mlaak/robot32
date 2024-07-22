package situation

import (

)

type IRequestContext interface{
	SetIporid(iporid string)
	GetIporid()(string)
}

type RequestContext struct {
	iporid string
}

func NewRequestContext() IRequestContext{
	return &RequestContext{}
}

func (rc *RequestContext) SetIporid(iporid string){
	rc.iporid = iporid
}

func (rc *RequestContext) GetIporid()(string){
	return rc.iporid
}
