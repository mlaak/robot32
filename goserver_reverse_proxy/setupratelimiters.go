package main

import (
	"fmt"
//	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
//	"errors"
//	"strings" 
//	"bytes"
//	"io/ioutil"
    "os"
)

var unLimiter *RateLimiter;
var stopLimiter *RateLimiter;
//var rateLimiter *RateLimiter;
//var overallIPLimiter *RateLimiter;

var expertIPLimiter *RateLimiter;  //non-logged-in users using experts 
var expertAIPLimiter *RateLimiter; //non-logged-in users using experts total. Meant to switch service off for non-logged-in users in case of DDOS

var expertUSRLimiter *RateLimiter; // logged in users using experts
var expertKEYLimiter *RateLimiter; // users who have their own openrouter key (we dont pay for their usage)
var expertPAYLimiter *RateLimiter; // paying users - a man can dream!

var staticIPLimiter *RateLimiter;  // Limit for static assets (html, css, tiny images)
var imageIPLimiter *RateLimiter;   // Limit for images for non-logged in users
var imageUSRLimiter *RateLimiter;  // Limit for images for logged in users
var imagePAYLimiter *RateLimiter;  // Limit for images for paying

var loginLimiter *RateLimiter;     // Limit for google and login 
var code498Limiter  *RateLimiter;  //returns code 498


fn SetupRateLimiters(){
		// Parse the URL of the Apache server
	//                                reqminute|reqhour|reqday|paralconn|   bytesmin| byteshour|  bytesday
    unLimiter =        NewRateLimiter(-1,    -1,     -1,    -1,       -1,         -1,        -1,        -1)
	stopLimiter =      NewRateLimiter(0,      0,      0,     0,        0,          0,         0,         0)
	
	code498Limiter =   NewRateLimiter(1,      0,      0,     0,        0,          0,         0,         0)
	code498Limiter.responseCode = 498; //expired or otherwise invalid token
	
	expertIPLimiter =  NewRateLimiter(2,     60,    300,  3000,       10,       1000,    500000,   5000000)
	expertAIPLimiter = NewRateLimiter(3,   6000,  60000,600000,      500,    5000000,  50000000, 500000000)
	
	expertUSRLimiter = NewRateLimiter(4,     60,    300,  3000,       10,        5000,    500000,   5000000)
	expertPAYLimiter = NewRateLimiter(5,     60,    300,  3000,       10,        5000,    500000,   5000000)
	expertKEYLimiter = NewRateLimiter(6,     60,    300,  3000,       10,        5000,    500000,   5000000)

    staticIPLimiter =  NewRateLimiter(7,   6000,  60000,600000,     1000,         -1,        -1,        -1)

    imageIPLimiter =   NewRateLimiter(8,   6000,  60000,600000,     1000,         -1,        -1,        -1)
    imageUSRLimiter =  NewRateLimiter(9,   6000,  60000,600000,     1000,         -1,        -1,        -1)
    imagePAYLimiter =  NewRateLimiter(10,  6000,  60000,600000,     1000,         -1,        -1,        -1)
      
	loginLimiter =     NewRateLimiter(11,  6000,  60000,600000,     1000,         -1,        -1,        -1)
}

