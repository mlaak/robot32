package ratelimiter

import (
	"time"
	"strconv"
)

type IRate interface {
	ResetIfTime(iporid string, now time.Time) bool
	Addbytes(user string,count int64)
	AddRequest(user string)
	
	IsRequestLimitBroken(user string) bool
	IsBytesLimitBroken(user string) bool

	GetWaitTimeStr(user string,now time.Time) (string)
	GetMaxRequests() int
	GetMaxBytes() int64
}

type Rate struct {
	requestCounts  map[string]int
	bytesCounts   map[string]int64
	lastResetTimes   map[string]time.Time

	maxRequests int	
	maxBytes int64
	period time.Duration
}

func NewRate(period time.Duration, maxRequests int,maxBytes int64) IRate {
	return &Rate{
		requestCounts:   make(map[string]int),
		bytesCounts:   make(map[string]int64),
		lastResetTimes:   make(map[string]time.Time),

		maxRequests: maxRequests,
		maxBytes: maxBytes,
		period: period,
	}
}

func (r *Rate) Addbytes(user string,count int64){
	r.bytesCounts[user]+=count;
}

func (r *Rate) AddRequest(user string){
	r.requestCounts[user]+=1;
}

func (r *Rate) IsRequestLimitBroken(user string)bool{
	if(r.maxRequests!=-1 && r.requestCounts[user]>r.maxRequests){
		return true
	} else {
		return false
	}
}

func (r *Rate) GetWaitTimeStr(user string,now time.Time) string{
	//TODO: make it also talk in minutes hours
	secs := r.period.Seconds()-now.Sub(r.lastResetTimes[user]).Seconds();
	return strconv.Itoa(int(secs))+"s";
}


func (r *Rate) IsBytesLimitBroken(user string) bool{
	if(r.maxBytes!=-1 && r.bytesCounts[user]>r.maxBytes){
		return true
	} else {
		return false
	}
}


func (r *Rate) ResetIfTime(iporid string, now time.Time) bool{
	if lastReset, exists := r.lastResetTimes[iporid]; !exists || now.Sub(lastReset) > r.period {
		r.requestCounts[iporid] = 0
		r.bytesCounts[iporid] = 0
		r.lastResetTimes[iporid] = now
		return true;
	}
	return false;
}

func (r *Rate) GetMaxRequests() int{
	return r.maxRequests
}

func (r *Rate) GetMaxBytes() int64 {
	return r.maxBytes
}