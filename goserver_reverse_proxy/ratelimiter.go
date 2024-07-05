package main

import (
//	"fmt"
//	"net/http"
	"sync"
	"time"
	"strconv"

)

type RateLimiter struct {
	mu              sync.Mutex
	nr int
	requestMinuteCounts   map[string]int
	bytesMinuteCounts   map[string]int64
	lastMinuteResetTimes   map[string]time.Time
	
	requestHourCounts   map[string]int
	bytesHourCounts   map[string]int64
	lastHourResetTimes   map[string]time.Time
	
	requestDayCounts   map[string]int
	bytesDayCounts   map[string]int64
	lastDayResetTimes   map[string]time.Time
	
	
	
	activeSessions  map[string]int
	
	maxRequestsPerMinute int
	maxRequestsPerHour   int
	maxRequestsPerDay   int
	maxParallelRequests  int
	
	maxBytesPerMinute int64
	maxBytesPerHour   int64
	maxBytesPerDay   int64
	
	responseCode int
	
	
}

func NewRateLimiter(nr,maxRequestsPerMinute, maxRequestsPerHour, maxRequestsPerDay,  maxParallelRequests int, maxBytesPerMinute int64, maxBytesPerHour int64, maxBytesPerDay int64) *RateLimiter {
	return &RateLimiter{
	    nr:nr,
		requestMinuteCounts:   make(map[string]int),
		bytesMinuteCounts:   make(map[string]int64),
		lastMinuteResetTimes:   make(map[string]time.Time),
		
		requestHourCounts:   make(map[string]int),
		bytesHourCounts:   make(map[string]int64),
		lastHourResetTimes:   make(map[string]time.Time),
		
		requestDayCounts:   make(map[string]int),
		bytesDayCounts:   make(map[string]int64),		
		lastDayResetTimes:   make(map[string]time.Time),
		
		
		
		activeSessions:  make(map[string]int),
		
		maxRequestsPerMinute: maxRequestsPerMinute,
		maxRequestsPerHour: maxRequestsPerHour,
		maxRequestsPerDay: maxRequestsPerDay,
	
		maxParallelRequests:  maxParallelRequests,
		
		maxBytesPerMinute: maxBytesPerMinute,
		maxBytesPerHour: maxBytesPerHour,
		maxBytesPerDay: maxBytesPerDay,
	
		responseCode: 429,
		
	}
}

func (rl *RateLimiter) Addbytes(iporid string, tokenCount int64){
    rl.mu.Lock()
    defer rl.mu.Unlock()
    
    rl.bytesMinuteCounts[iporid]+=tokenCount
	rl.bytesHourCounts[iporid]+=tokenCount
	rl.bytesDayCounts[iporid]+=tokenCount
    
    //fmt.Println("!!!!!!!!!Metering bytes",iporid,rl.bytesMinuteCounts[iporid])
    
    
}

func (rl *RateLimiter) Allow(iporid string, w interface {AddReleaser(func())} ) (bool,int,string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Reset counter if a minute has passed
	if lastReset, exists := rl.lastMinuteResetTimes[iporid]; !exists || now.Sub(lastReset) > time.Minute {
		rl.activeSessions[iporid] = 0
		rl.requestMinuteCounts[iporid] = 0
		rl.bytesMinuteCounts[iporid] = 0
		rl.lastMinuteResetTimes[iporid] = now
	}
	
	if lastReset, exists := rl.lastHourResetTimes[iporid]; !exists || now.Sub(lastReset) > time.Hour {
		rl.requestHourCounts[iporid] = 0
		rl.bytesHourCounts[iporid] = 0
		rl.lastHourResetTimes[iporid] = now
	}
	
	
	
	if lastReset, exists := rl.lastDayResetTimes[iporid]; !exists || now.Sub(lastReset) > (time.Hour*24) {
		rl.requestDayCounts[iporid] = 0
		rl.bytesDayCounts[iporid] = 0
		rl.lastDayResetTimes[iporid] = now
	}
	
	
	//we have a zero limiter, the request is clearlu blocked
	if rl.maxRequestsPerMinute == 0 || rl.maxRequestsPerHour == 0 || rl.maxRequestsPerDay == 0{
	    return false, rl.responseCode, "Not allowed (maybe you need to login or prove you are not a robot or something)."
	}
	
	// Check if the iporid has exceeded the rate limit for minute
	if rl.maxRequestsPerMinute!=-1 && rl.requestMinuteCounts[iporid] >= rl.maxRequestsPerMinute {
		return false, rl.responseCode, "Requests per minute exceeded, wait  "+strconv.Itoa(60-int((now.Sub(rl.lastMinuteResetTimes[iporid]).Seconds())))+" seconds!"
	}

	// Check if the iporid has exceeded the rate limit for hour
	if rl.maxRequestsPerHour!=-1 && rl.requestHourCounts[iporid] >= rl.maxRequestsPerHour {
		return false, rl.responseCode, "Requests per hour exceeded!"
	}
	
	// Check if the iporid has exceeded the rate limit for day
	if rl.maxRequestsPerDay!=-1 && rl.requestDayCounts[iporid] >= rl.maxRequestsPerDay {
		return false, rl.responseCode, "Requests per day exceeded!"
	}

	// Check if the iporid has exceeded the parallel request limit
	if rl.maxParallelRequests!=-1 && rl.activeSessions[iporid] >= rl.maxParallelRequests {
		return false, rl.responseCode, "Max parallel requests exceeded!"
	}
	
	
	
	// Check if the iporid has exceeded the byte limit for minute
	if rl.maxBytesPerMinute!=-1 && rl.bytesMinuteCounts[iporid] >= rl.maxBytesPerMinute {
		return false, rl.responseCode, "Characters per minute exceeded, wait "+strconv.Itoa(60-int((now.Sub(rl.lastMinuteResetTimes[iporid]).Seconds())))+" seconds!" 
	}

	// Check if the iporid has exceeded the byte limit for hour
	if rl.maxBytesPerHour!=-1 && rl.bytesHourCounts[iporid] >= rl.maxBytesPerHour {
		return false, rl.responseCode, "Characters per hour exceeded!"
	}
	
	// Check if the iporid has exceeded the byte limit for day
	if rl.maxBytesPerDay!=-1 && rl.bytesDayCounts[iporid] >= rl.maxBytesPerDay {
		return false, rl.responseCode, "Characters per day exceeded!"
	}
	
	

	// Increment counters
	rl.requestMinuteCounts[iporid]++
	rl.requestHourCounts[iporid]++
	rl.requestDayCounts[iporid]++
	
	rl.activeSessions[iporid]++
	
	releasefunc := func(){
	    rl.Release(iporid)
	}
	w.AddReleaser(releasefunc)
	//w.releasers = append(w.releasers, releasefunc)
	
	
	//fmt.Println("Active sessions for "+iporid+" is now "+strconv.Itoa(rl.activeSessions[iporid])) 
	

	return true, 200, ""
}

func (rl *RateLimiter) Release(iporid string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.activeSessions[iporid] > 0 {
	    rl.activeSessions[iporid]--
	    //fmt.Println(iporid+" closed count is "+strconv.Itoa(rl.activeSessions[iporid]))
	}
}

