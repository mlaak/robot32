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

func NewRateLimiter(maxRequestsPerMinute, maxRequestsPerHour, maxRequestsPerDay,  maxParallelRequests int, maxBytesPerMinute int64, maxBytesPerHour int64, maxBytesPerDay int64) *RateLimiter {
	return &RateLimiter{
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

func (rl *RateLimiter) Addbytes(ip string, tokenCount int64){
    rl.mu.Lock()
    defer rl.mu.Unlock()
    
    rl.bytesMinuteCounts[ip]+=tokenCount
	rl.bytesHourCounts[ip]+=tokenCount
	rl.bytesDayCounts[ip]+=tokenCount
    
    //fmt.Println("!!!!!!!!!Metering bytes",ip,rl.bytesMinuteCounts[ip])
    
    
}

func (rl *RateLimiter) Allow(ip string, w interface {AddReleaser(func())} ) (bool,int,string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Reset counter if a minute has passed
	if lastReset, exists := rl.lastMinuteResetTimes[ip]; !exists || now.Sub(lastReset) > time.Minute {
		rl.activeSessions[ip] = 0
		rl.requestMinuteCounts[ip] = 0
		rl.bytesMinuteCounts[ip] = 0
		rl.lastMinuteResetTimes[ip] = now
	}
	
	if lastReset, exists := rl.lastHourResetTimes[ip]; !exists || now.Sub(lastReset) > time.Hour {
		rl.requestHourCounts[ip] = 0
		rl.bytesHourCounts[ip] = 0
		rl.lastHourResetTimes[ip] = now
	}
	
	if lastReset, exists := rl.lastDayResetTimes[ip]; !exists || now.Sub(lastReset) > (time.Hour*24) {
		rl.requestDayCounts[ip] = 0
		rl.bytesDayCounts[ip] = 0
		rl.lastDayResetTimes[ip] = now
	}
	
	if rl.maxRequestsPerMinute == 0 || rl.maxRequestsPerHour == 0 || rl.maxRequestsPerDay == 0{
	    return false, rl.responseCode, "Not allowed (maybe you need to login or prove you are not a robot)."
	}
	
	// Check if the IP has exceeded the rate limit for minute
	if rl.maxRequestsPerMinute!=-1 && rl.requestMinuteCounts[ip] >= rl.maxRequestsPerMinute {
		return false, rl.responseCode, "Requests per minute exceeded, wait  "+strconv.Itoa(60-int((now.Sub(rl.lastMinuteResetTimes[ip]).Seconds())))+" seconds!"
	}

	// Check if the IP has exceeded the rate limit for hour
	if rl.maxRequestsPerHour!=-1 && rl.requestHourCounts[ip] >= rl.maxRequestsPerHour {
		return false, rl.responseCode, "Requests per hour exceeded!"
	}
	
	// Check if the IP has exceeded the rate limit for day
	if rl.maxRequestsPerDay!=-1 && rl.requestDayCounts[ip] >= rl.maxRequestsPerDay {
		return false, rl.responseCode, "Requests per day exceeded!"
	}

	// Check if the IP has exceeded the parallel request limit
	if rl.maxParallelRequests!=-1 && rl.activeSessions[ip] >= rl.maxParallelRequests {
		return false, rl.responseCode, "Max parallel requests exceeded!"
	}
	
	
	
	// Check if the IP has exceeded the byte limit for minute
	if rl.maxBytesPerMinute!=-1 && rl.bytesMinuteCounts[ip] >= rl.maxBytesPerMinute {
		return false, rl.responseCode, "Characters per minute exceeded, wait "+strconv.Itoa(60-int((now.Sub(rl.lastMinuteResetTimes[ip]).Seconds())))+" seconds!" 
	}

	// Check if the IP has exceeded the byte limit for hour
	if rl.maxBytesPerHour!=-1 && rl.bytesHourCounts[ip] >= rl.maxBytesPerHour {
		return false, rl.responseCode, "Characters per hour exceeded!"
	}
	
	// Check if the IP has exceeded the byte limit for day
	if rl.maxBytesPerDay!=-1 && rl.bytesDayCounts[ip] >= rl.maxBytesPerDay {
		return false, rl.responseCode, "Characters per day exceeded!"
	}
	
	

	// Increment counters
	rl.requestMinuteCounts[ip]++
	rl.requestHourCounts[ip]++
	rl.requestDayCounts[ip]++
	
	rl.activeSessions[ip]++
	
	releasefunc := func(){
	    rl.Release(ip)
	}
	w.AddReleaser(releasefunc)
	//w.releasers = append(w.releasers, releasefunc)
	
	
	//fmt.Println("Active sessions for "+ip+" is now "+strconv.Itoa(rl.activeSessions[ip])) 
	

	return true, 200, ""
}

func (rl *RateLimiter) Release(ip string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.activeSessions[ip] > 0 {
	    rl.activeSessions[ip]--
	    //fmt.Println(ip+" closed count is "+strconv.Itoa(rl.activeSessions[ip]))
	}
}

