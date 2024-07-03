package main

import (
	"fmt"
//	"net/http"
	"sync"
	"time"
	"strconv"

)

type RateLimiter struct {
	mu              sync.Mutex
	requestMinuteCounts   map[string]int
	tokensMinuteCounts   map[string]int
	lastMinuteResetTimes   map[string]time.Time
	
	requestHourCounts   map[string]int
	tokensHourCounts   map[string]int
	lastHourResetTimes   map[string]time.Time
	
	requestDayCounts   map[string]int
	tokensDayCounts   map[string]int
	lastDayResetTimes   map[string]time.Time
	
	
	
	activeSessions  map[string]int
	
	maxRequestsPerMinute int
	maxRequestsPerHour   int
	maxRequestsPerDay   int
	maxParallelRequests  int
	
	maxTokensPerMinute int
	maxTokensPerHour   int
	maxTokensPerDay   int
	
}

func NewRateLimiter(maxRequestsPerMinute, maxRequestsPerHour, maxRequestsPerDay,  maxParallelRequests int, maxTokensPerMinute int, maxTokensPerHour int, maxTokensPerDay int) *RateLimiter {
	return &RateLimiter{
		requestMinuteCounts:   make(map[string]int),
		tokensMinuteCounts:   make(map[string]int),
		lastMinuteResetTimes:   make(map[string]time.Time),
		
		requestHourCounts:   make(map[string]int),
		tokensHourCounts:   make(map[string]int),
		lastHourResetTimes:   make(map[string]time.Time),
		
		requestDayCounts:   make(map[string]int),
		tokensDayCounts:   make(map[string]int),		
		lastDayResetTimes:   make(map[string]time.Time),
		
		
		
		activeSessions:  make(map[string]int),
		
		maxRequestsPerMinute: maxRequestsPerMinute,
		maxRequestsPerHour: maxRequestsPerHour,
		maxRequestsPerDay: maxRequestsPerDay,
	
		maxParallelRequests:  maxParallelRequests,
		
		maxTokensPerMinute: maxTokensPerMinute,
		maxTokensPerHour: maxTokensPerHour,
		maxTokensPerDay: maxTokensPerDay,
	
		
		
	}
}

func (rl *RateLimiter) AddTokens(ip string, tokenCount int){
    rl.mu.Lock()
    defer rl.mu.Unlock()
    
    rl.tokensMinuteCounts[ip]+=tokenCount
	rl.tokensHourCounts[ip]+=tokenCount
	rl.tokensDayCounts[ip]+=tokenCount
    
}

func (rl *RateLimiter) Allow(ip string, w interface {AddReleaser(func())} ) (bool,string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Reset counter if a minute has passed
	if lastReset, exists := rl.lastMinuteResetTimes[ip]; !exists || now.Sub(lastReset) > time.Minute {
		rl.activeSessions[ip] = 0
		rl.requestMinuteCounts[ip] = 0
		rl.tokensMinuteCounts[ip] = 0
		rl.lastMinuteResetTimes[ip] = now
	}
	
	if lastReset, exists := rl.lastHourResetTimes[ip]; !exists || now.Sub(lastReset) > time.Hour {
		rl.requestHourCounts[ip] = 0
		rl.tokensHourCounts[ip] = 0
		rl.lastHourResetTimes[ip] = now
	}
	
	if lastReset, exists := rl.lastDayResetTimes[ip]; !exists || now.Sub(lastReset) > (time.Hour*24) {
		rl.requestDayCounts[ip] = 0
		rl.tokensDayCounts[ip] = 0
		rl.lastDayResetTimes[ip] = now
	}
	
	
	// Check if the IP has exceeded the rate limit
	if rl.requestMinuteCounts[ip] >= rl.maxRequestsPerMinute {
		return false,"Requests per minute exceeded"
	}

	// Check if the IP has exceeded the rate limit
	if rl.requestHourCounts[ip] >= rl.maxRequestsPerHour {
		return false,"Requests per hour exceeded"
	}
	
	// Check if the IP has exceeded the rate limit
	if rl.requestDayCounts[ip] >= rl.maxRequestsPerDay {
		return false,"Requests per day exceeded"
	}

	// Check if the IP has exceeded the parallel request limit
	if rl.activeSessions[ip] >= rl.maxParallelRequests {
		return false,"Max parallel requests exceeded"
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
	
	
	fmt.Println("Active sessions for "+ip+" is now "+strconv.Itoa(rl.activeSessions[ip])) 
	

	return true,""
}

func (rl *RateLimiter) Release(ip string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.activeSessions[ip] > 0 {
	    rl.activeSessions[ip]--
	    fmt.Println(ip+" closed count is "+strconv.Itoa(rl.activeSessions[ip]))
	}
}

