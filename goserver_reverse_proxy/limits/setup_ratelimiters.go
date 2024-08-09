package limits

import (
	rl "grp/ratelimiter"
)

var UnLimiter rl.IRateLimiter
var StopLimiter rl.IRateLimiter

//var rateLimiter *IRateLimiter;
//var overallIPLimiter *IRateLimiter;

var ExpertIPLimiter rl.IRateLimiter  //non-logged-in users using experts
var ExpertAIPLimiter rl.IRateLimiter //non-logged-in users using experts total. Meant to switch service off for non-logged-in users in case of DDOS

var ExpertUSRLimiter rl.IRateLimiter // logged in users using experts
var ExpertKEYLimiter rl.IRateLimiter // users who have their own openrouter key (we dont pay for their usage)
var ExpertPAYLimiter rl.IRateLimiter // paying users - a man can dream!

var StaticIPLimiter rl.IRateLimiter // Limit for static assets (html, css, tiny images)
var ImageIPLimiter rl.IRateLimiter  // Limit for images for non-logged in users
var ImageUSRLimiter rl.IRateLimiter // Limit for images for logged in users
var ImagePAYLimiter rl.IRateLimiter // Limit for images for paying

var LoginLimiter rl.IRateLimiter   // Limit for google and login
var Code498Limiter rl.IRateLimiter //returns code 498

func SetupRateLimiters() {
	// Parse the URL of the Apache server
	//                                reqminute|reqhour|reqday|paralconn|   bytesmin| byteshour|  bytesday
	UnLimiter = rl.NewRateLimiter(-1, -1, -1, -1, -1, -1, -1, -1)
	StopLimiter = rl.NewRateLimiter(0, 0, 0, 0, 0, 0, 0, 0)

	Code498Limiter = rl.NewRateLimiter(1, 0, 0, 0, 0, 0, 0, 0)
	Code498Limiter.SetResponseCode(498) //expired or otherwise invalid token

	ExpertIPLimiter = rl.NewRateLimiter(2, 60, 300, 3000, 10, 1000, 500000, 5000000)
	ExpertAIPLimiter = rl.NewRateLimiter(3, 6000, 60000, 600000, 500, 5000000, 50000000, 500000000)

	ExpertUSRLimiter = rl.NewRateLimiter(4, 60, 300, 3000, 10, 5000, 500000, 5000000)
	ExpertPAYLimiter = rl.NewRateLimiter(5, 60, 300, 3000, 10, 5000, 500000, 5000000)
	ExpertKEYLimiter = rl.NewRateLimiter(6, 60, 300, 3000, 10, 5000, 500000, 5000000)

	StaticIPLimiter = rl.NewRateLimiter(7, 6000, 60000, 600000, 1000, -1, -1, -1)

	ImageIPLimiter = rl.NewRateLimiter(8, 6000, 60000, 600000, 1000, -1, -1, -1)
	ImageUSRLimiter = rl.NewRateLimiter(9, 6000, 60000, 600000, 1000, -1, -1, -1)
	ImagePAYLimiter = rl.NewRateLimiter(10, 6000, 60000, 600000, 1000, -1, -1, -1)

	LoginLimiter = rl.NewRateLimiter(11, 6000, 60000, 600000, 1000, -1, -1, -1)
}
