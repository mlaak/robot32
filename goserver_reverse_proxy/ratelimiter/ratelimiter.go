/***********************************************************
   ______
  /      \     IS USED TO LIMIT USERS
 /  STOP  \			* NR O<F REQUESTS
 \        /			* PARALLEL REQUESTS
  \______/			* AMOUNT OF DATA
     ||		  BASED ON
     ||				* MINUTE, HOUR, DAY
     ||
***********************************************************/

package ratelimiter

import (
	. "grp/translator" //lint:ignore ST1001 Sometimes dot imports are needed
	. "grp/ttd"        //lint:ignore ST1001 Sometimes dot imports are needed
	"sync"
	"time"
)

type IRateLimiter interface {
	Allow(c int64, iporid string) (bool, int, string)
	CountUpOneConnection(c int64, iporid string) func()
	Addbytes(c int64, iporid string, bytesCount int64)
	SetResponseCode(int)
	GetNr() int
}

type RateLimiter struct {
	minuteLimit IRate
	hourLimit   IRate
	dayLimit    IRate

	activeConnections   map[string]int
	maxParallelRequests int

	mu           sync.Mutex
	nr           int
	ResponseCode int
}

func NewRateLimiter(nr, maxRequestsPerMinute, maxRequestsPerHour, maxRequestsPerDay, maxParallelRequests int, maxBytesPerMinute int64, maxBytesPerHour int64, maxBytesPerDay int64) *RateLimiter {
	return &RateLimiter{
		minuteLimit: NewRate(time.Minute, maxRequestsPerMinute, maxBytesPerMinute),
		hourLimit:   NewRate(time.Hour, maxRequestsPerHour, maxBytesPerHour),
		dayLimit:    NewRate(time.Hour*24, maxRequestsPerDay, maxBytesPerDay),

		activeConnections:   make(map[string]int),
		maxParallelRequests: maxParallelRequests,

		nr:           nr,
		ResponseCode: 429,
	}
}

func (rl *RateLimiter) Allow(c int64, iporid string) (bool, int, string) {
	// *********** PREPARATIONS ***********************************
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := time.Now()

	// *********** RESET LIMITS ***********************************
	if rl.minuteLimit.ResetIfTime(c, iporid, now) {
		//in case we have miscounted active connection (edge cases when we dont detect disconnect)
		//lets reset that also - otherwise we get a system stop for the user
		rl.activeConnections[iporid] = 0
	}
	rl.hourLimit.ResetIfTime(c, iporid, now)
	rl.dayLimit.ResetIfTime(c, iporid, now)

	// ********** CHECK IF ACCESS IF FULLY BLOCKED ****************
	if rl.minuteLimit.GetMaxRequests() == 0 || rl.hourLimit.GetMaxRequests() == 0 || rl.dayLimit.GetMaxRequests() == 0 {
		txt := TR(c, "Not allowed (maybe you need to login or prove you are not a robot or something).")
		return TTX3(c, false, rl.ResponseCode, txt)
	}

	// ********** CHECK REQUEST LIMITS  ***************************
	if rl.minuteLimit.IsRequestLimitBroken(c, iporid) {
		txt := TR(c, "Requests per minute exceeded, wait "+rl.minuteLimit.GetWaitTimeStr(c, iporid, now))
		return TTX3(c, false, rl.ResponseCode, txt)
	}
	if rl.hourLimit.IsRequestLimitBroken(c, iporid) {
		txt := TR(c, "Requests per hour exceeded, wait "+rl.minuteLimit.GetWaitTimeStr(c, iporid, now))
		return TTX3(c, false, rl.ResponseCode, txt)
	}
	if rl.dayLimit.IsRequestLimitBroken(c, iporid) {
		txt := TR(c, "Requests per day exceeded, wait "+rl.minuteLimit.GetWaitTimeStr(c, iporid, now))
		return TTX3(c, false, rl.ResponseCode, txt)
	}

	// ********** CHECK DATA FLOW LIMITS  *************************
	if rl.minuteLimit.IsBytesLimitBroken(c, iporid) {
		txt := TR(c, "Characters per minute exceeded, wait "+rl.minuteLimit.GetWaitTimeStr(c, iporid, now))
		return TTX3(c, false, rl.ResponseCode, txt)
	}
	if rl.hourLimit.IsBytesLimitBroken(c, iporid) {
		txt := TR(c, "Characters per hour exceeded, wait "+rl.minuteLimit.GetWaitTimeStr(c, iporid, now))
		return TTX3(c, false, rl.ResponseCode, txt)
	}
	if rl.dayLimit.IsBytesLimitBroken(c, iporid) {
		txt := TR(c, "Characters per day exceeded, wait "+rl.minuteLimit.GetWaitTimeStr(c, iporid, now))
		return TTX3(c, false, rl.ResponseCode, txt)
	}

	// ** INCREASE COUNTERS AND SET TRIGGER FOR CONNECTION CLOSE **
	rl.minuteLimit.AddRequest(c, iporid)
	rl.hourLimit.AddRequest(c, iporid)
	rl.dayLimit.AddRequest(c, iporid)

	// ********** RETURN SUCCESS **********************************
	return true, 200, ""
}

func (rl *RateLimiter) CountUpOneConnection(c int64, iporid string) func() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.activeConnections[iporid]++
	releasefunc := func() { //we need to track the connection close
		rl.mu.Lock()
		defer rl.mu.Unlock()
		if rl.activeConnections[iporid] > 0 {
			rl.activeConnections[iporid]--
		}
	}
	return releasefunc
}

func (rl *RateLimiter) Addbytes(c int64, iporid string, bytesCount int64) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.minuteLimit.Addbytes(c, iporid, bytesCount)
	rl.hourLimit.Addbytes(c, iporid, bytesCount)
	rl.dayLimit.Addbytes(c, iporid, bytesCount)
}

func (rl *RateLimiter) SetResponseCode(rc int) {
	rl.ResponseCode = rc
}

func (rl *RateLimiter) GetNr() int {
	return rl.nr
}
