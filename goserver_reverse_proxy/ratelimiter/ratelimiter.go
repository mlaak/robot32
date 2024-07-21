/***********************************************************
   ______
  /      \     IS USED TO LIMIT USERS 
 /  STOP  \			* NR OF REQUESTS
 \        /			* PARALLEL REQUESTS
  \______/			* AMOUNT OF DATA
     ||		  BASED ON
     ||				* MINUTE, HOUR, DAY
     ||
***********************************************************/

package ratelimiter

import (
//	"fmt"
//	"net/http"
	"sync"
	"time"
	"strconv"
	"grp/situation"
	. "grp/translator"
)

type IRateLimiter interface{
	Allow(iporid string, context situation.IRequestContext) (bool,int,string)
	CountUpOneConnection(iporid string)(func())
	Addbytes(iporid string, bytesCount int64)
	SetResponseCode(int)
	GetNr()(int)
}

type RateLimiter struct {
	minuteLimit *Rate
	hourLimit *Rate
	dayLimit *Rate

	activeConnections  map[string]int
	maxParallelRequests  int

	mu              sync.Mutex
	nr int
	ResponseCode int
}

func NewRateLimiter(nr,maxRequestsPerMinute, maxRequestsPerHour, maxRequestsPerDay,  maxParallelRequests int, maxBytesPerMinute int64, maxBytesPerHour int64, maxBytesPerDay int64) *RateLimiter {
	return &RateLimiter{
		minuteLimit:NewRate(time.Minute,  maxRequestsPerMinute, maxBytesPerMinute),
		hourLimit:  NewRate(time.Hour,    maxRequestsPerHour,   maxBytesPerHour),
		dayLimit:   NewRate(time.Hour*24, maxRequestsPerDay,    maxBytesPerDay),

		activeConnections:  make(map[string]int),
		maxParallelRequests:  maxParallelRequests,
		
		nr:nr,
		ResponseCode: 429,
	}
}



func (rl *RateLimiter) Allow(iporid string, context situation.IRequestContext) (bool,int,string) {
// *********** PREPARATIONS ***********************************
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := time.Now()

// *********** RESET LIMITS ***********************************
	if rl.minuteLimit.ResetIfTime(iporid,now){
		//in case we have miscounted active connection (edge cases when we dont detect disconnect)
		//lets reset that also - otherwise we get a system stop for the user
		rl.activeConnections[iporid] = 0 
	}
	rl.hourLimit.ResetIfTime(iporid,now);
	rl.dayLimit.ResetIfTime(iporid,now);
	
// ********** CHECK IF ACCESS IF FULLY BLOCKED ****************
	if(rl.minuteLimit.maxRequests == 0 || rl.hourLimit.maxRequests == 0 || rl.dayLimit.maxRequests == 0){
		txt:=TR("Not allowed (maybe you need to login or prove you are not a robot or something).",context)
		return false, rl.ResponseCode, txt
	}

// ********** CHECK REQUEST LIMITS  ***************************
	if(rl.minuteLimit.IsRequestLimitBroken(iporid)){
		txt:=TR("Requests per minute exceeded, wait "+rl.minuteLimit.GetWaitTimeStr(iporid,now),context);
		return false, rl.ResponseCode, txt;
	}
	if(rl.hourLimit.IsRequestLimitBroken(iporid)){
		txt:=TR("Requests per hour exceeded, wait "+rl.minuteLimit.GetWaitTimeStr(iporid,now),context);
		return false, rl.ResponseCode, txt;
	}
	if(rl.dayLimit.IsRequestLimitBroken(iporid)){
		txt:=TR("Requests per day exceeded, wait "+rl.minuteLimit.GetWaitTimeStr(iporid,now),context);
		return false, rl.ResponseCode, txt;
	}

// ********** CHECK DATA FLOW LIMITS  *************************
	if(rl.minuteLimit.IsBytesLimitBroken(iporid)){
		txt:=TR("Characters per minute exceeded, wait "+rl.minuteLimit.GetWaitTimeStr(iporid,now),context);
		return false, rl.ResponseCode, txt;
	}
	if(rl.hourLimit.IsBytesLimitBroken(iporid)){
		txt:=TR("Characters per hour exceeded, wait "+rl.minuteLimit.GetWaitTimeStr(iporid,now),context);
		return false, rl.ResponseCode, txt;
	}
	if(rl.dayLimit.IsBytesLimitBroken(iporid)){
		txt:=TR("Characters per day exceeded, wait "+rl.minuteLimit.GetWaitTimeStr(iporid,now),context);
		return false, rl.ResponseCode, txt;
	}

// ** INCREASE COUNTERS AND SET TRIGGER FOR CONNECTION CLOSE **
	rl.minuteLimit.AddRequest(iporid);
	rl.hourLimit.AddRequest(iporid);
	rl.dayLimit.AddRequest(iporid);
	

// ********** RETURN SUCCESS **********************************
	return true, 200, ""
}

func (rl *RateLimiter) CountUpOneConnection(iporid string)(func()){
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.activeConnections[iporid]++
	releasefunc := func(){ //we need to track the connection close
		rl.mu.Lock()
		defer rl.mu.Unlock()
		if rl.activeConnections[iporid] > 0 {
			rl.activeConnections[iporid]--
		}
	}
	return releasefunc
}


func (rl *RateLimiter) Addbytes(iporid string, bytesCount int64){
    rl.mu.Lock()
    defer rl.mu.Unlock()

	rl.minuteLimit.Addbytes(iporid,bytesCount);
	rl.hourLimit.Addbytes(iporid,bytesCount);
	rl.dayLimit.Addbytes(iporid,bytesCount);    
}

func (rl *RateLimiter) SetResponseCode(rc int){
	rl.ResponseCode = rc;
}

func (rl *RateLimiter) GetNr()int{
	return rl.nr;
}

type Rate struct {
	requestCounts  map[string]int
	bytesCounts   map[string]int64
	lastResetTimes   map[string]time.Time

	maxRequests int	
	maxBytes int64
	period time.Duration
}

func NewRate(period time.Duration, maxRequests int,maxBytes int64) *Rate {
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