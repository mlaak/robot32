package ratelimiter

import (
	. "grp/ttd"
	"strconv"
	"time"
	//	"fmt"
)

type IRate interface {
	ResetIfTime(c int64, iporid string, now time.Time) bool
	Addbytes(c int64, user string, count int64)
	AddRequest(c int64, user string)

	IsRequestLimitBroken(c int64, user string) bool
	IsBytesLimitBroken(c int64, user string) bool

	GetWaitTimeStr(c int64, user string, now time.Time) string
	GetMaxRequests() int
	GetMaxBytes() int64
}

type Rate struct {
	requestCounts  map[string]int
	bytesCounts    map[string]int64
	lastResetTimes map[string]time.Time

	maxRequests int
	maxBytes    int64
	period      time.Duration
}

func NewRate(period time.Duration, maxRequests int, maxBytes int64) IRate {
	return &Rate{
		requestCounts:  make(map[string]int),
		bytesCounts:    make(map[string]int64),
		lastResetTimes: make(map[string]time.Time),

		maxRequests: maxRequests,
		maxBytes:    maxBytes,
		period:      period,
	}
}

func (r *Rate) Addbytes(c int64, user string, count int64) {
	r.bytesCounts[user] += count
	TTD(c, "Added bytes", "user", user, "add", count, "new", r.bytesCounts[user])
}

func (r *Rate) AddRequest(c int64, user string) {
	r.requestCounts[user] += 1
	TTD(c, "Added request", "user", user, "new", r.requestCounts[user])
}

func (r *Rate) IsRequestLimitBroken(c int64, user string) bool {
	TTD(c, "Testing for Request limits", "user", user, "requestCounts", r.requestCounts[user], "maxRequests", r.maxRequests)

	if r.maxRequests != -1 && r.requestCounts[user] > r.maxRequests {
		TTD(c, "Failed Request limits", "user", user, "requestCounts", r.requestCounts[user], "maxRequests", r.maxRequests)
		return true
	} else {
		TTD(c, "Passed Request limits", "user", user, "requestCounts", r.requestCounts[user], "maxRequests", r.maxRequests)
		return false
	}
}

func (r *Rate) GetWaitTimeStr(c int64, user string, now time.Time) string {
	//TODO: make it also talk in minutes hours
	secs := r.period.Seconds() - now.Sub(r.lastResetTimes[user]).Seconds()
	return strconv.Itoa(int(secs)) + "s"
}

func (r *Rate) IsBytesLimitBroken(c int64, user string) bool {
	TTD(c, "Testing for Bytes limits", "user", user, "bytesCounts", r.bytesCounts[user], "maxBytes", r.maxBytes)

	if r.maxBytes != -1 && r.bytesCounts[user] > r.maxBytes {
		TTD(c, "Failed Bytes limits", "user", user, "bytesCounts", r.bytesCounts[user], "maxBytes", r.maxBytes)
		return true
	} else {
		TTD(c, "Passed Bytes limits", "user", user, "bytesCounts", r.bytesCounts[user], "maxBytes", r.maxBytes)
		return false
	}
}

func (r *Rate) ResetIfTime(c int64, iporid string, now time.Time) bool {
	lastReset, exists := r.lastResetTimes[iporid]
	TTD(c, "Testing if it is time to reset limits", "iporid", iporid, "lastReset", lastReset, " now.Sub(lastReset)", now.Sub(lastReset), "r.period", r.period)

	if !exists || now.Sub(lastReset) > r.period {
		TTD(c, "Reseting limits")
		r.requestCounts[iporid] = 0
		r.bytesCounts[iporid] = 0
		r.lastResetTimes[iporid] = now
		return true
	}
	TTD(c, "Not reseting limits")
	return false
}

func (r *Rate) GetMaxRequests() int {
	return r.maxRequests
}

func (r *Rate) GetMaxBytes() int64 {
	return r.maxBytes
}
