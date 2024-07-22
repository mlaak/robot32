package limits

import (
	. "grp/ratelimiter"
)

var UnLimiter IRateLimiter;
var StopLimiter IRateLimiter;
//var rateLimiter *IRateLimiter;
//var overallIPLimiter *IRateLimiter;

var ExpertIPLimiter IRateLimiter;  //non-logged-in users using experts 
var ExpertAIPLimiter IRateLimiter; //non-logged-in users using experts total. Meant to switch service off for non-logged-in users in case of DDOS

var ExpertUSRLimiter IRateLimiter; // logged in users using experts
var ExpertKEYLimiter IRateLimiter; // users who have their own openrouter key (we dont pay for their usage)
var ExpertPAYLimiter IRateLimiter; // paying users - a man can dream!

var StaticIPLimiter IRateLimiter;  // Limit for static assets (html, css, tiny images)
var ImageIPLimiter IRateLimiter;   // Limit for images for non-logged in users
var ImageUSRLimiter IRateLimiter;  // Limit for images for logged in users
var ImagePAYLimiter IRateLimiter;  // Limit for images for paying

var LoginLimiter IRateLimiter;     // Limit for google and login 
var Code498Limiter  IRateLimiter;  //returns code 498




func SetupRateLimiters(){
		// Parse the URL of the Apache server
	//                                reqminute|reqhour|reqday|paralconn|   bytesmin| byteshour|  bytesday
    UnLimiter =        NewRateLimiter(-1,    -1,     -1,    -1,       -1,         -1,        -1,        -1)
	StopLimiter =      NewRateLimiter(0,      0,      0,     0,        0,          0,         0,         0)
	
	Code498Limiter =   NewRateLimiter(1,      0,      0,     0,        0,          0,         0,         0)
	Code498Limiter.SetResponseCode(498); //expired or otherwise invalid token
	
	ExpertIPLimiter =  NewRateLimiter(2,     60,    300,  3000,       10,       1000,    500000,   5000000)
	ExpertAIPLimiter = NewRateLimiter(3,   6000,  60000,600000,      500,    5000000,  50000000, 500000000)
	
	ExpertUSRLimiter = NewRateLimiter(4,     60,    300,  3000,       10,        5000,    500000,   5000000)
	ExpertPAYLimiter = NewRateLimiter(5,     60,    300,  3000,       10,        5000,    500000,   5000000)
	ExpertKEYLimiter = NewRateLimiter(6,     60,    300,  3000,       10,        5000,    500000,   5000000)

    StaticIPLimiter =  NewRateLimiter(7,   6000,  60000,600000,     1000,         -1,        -1,        -1)

    ImageIPLimiter =   NewRateLimiter(8,   6000,  60000,600000,     1000,         -1,        -1,        -1)
    ImageUSRLimiter =  NewRateLimiter(9,   6000,  60000,600000,     1000,         -1,        -1,        -1)
    ImagePAYLimiter =  NewRateLimiter(10,  6000,  60000,600000,     1000,         -1,        -1,        -1)
      
	LoginLimiter =     NewRateLimiter(11,  6000,  60000,600000,     1000,         -1,        -1,        -1)
}

