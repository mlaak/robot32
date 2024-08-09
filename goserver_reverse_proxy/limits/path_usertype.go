package limits

import (
	rl "grp/ratelimiter"
	. "grp/ttd"
	s "grp/usersession"
	"strings"
)

func GetLimitersForPathAndUserType(c int64, path string, usertype string) (rl.IRateLimiter, rl.IRateLimiter) {
	p := path
	if strings.HasPrefix(p, "/experts/") {
		return GetExpertsLimiters(c, path, usertype)

	} else if strings.HasPrefix(p, "/recieved_images/") || strings.HasPrefix(p, "/welcome2.jpg") {
		return GetImagesLimiters(c, path, usertype)

	} else if strings.HasPrefix(p, "/units/github_login") || strings.HasPrefix(p, "/units/google_login") {
		return LoginLimiter, UnLimiter

	} else {
		return StaticIPLimiter, UnLimiter
	}
}

func GetExpertsLimiters(c int64, path string, usertype string) (rl.IRateLimiter, rl.IRateLimiter) {
	switch usertype {
	case s.UserTypeUnverified:
		return TTX2(c, Code498Limiter, Code498Limiter)
	case s.UserTypeCaptchad:
		return TTX2(c, ExpertIPLimiter, ExpertAIPLimiter)
	case s.UserTypeGithub:
		return TTX2(c, ExpertUSRLimiter, UnLimiter)
	case s.UserTypeGoogle:
		return TTX2(c, ExpertUSRLimiter, UnLimiter)
	default:
		return TTX2(c, ExpertIPLimiter, ExpertAIPLimiter)
	}
}

func GetImagesLimiters(c int64, path string, usertype string) (rl.IRateLimiter, rl.IRateLimiter) {
	switch usertype {
	case s.UserTypeUnverified:
		return TTX2(c, ImageIPLimiter, UnLimiter)
	case s.UserTypeCaptchad:
		return TTX2(c, ImageIPLimiter, UnLimiter)
	case s.UserTypeGithub:
		return TTX2(c, ImageUSRLimiter, UnLimiter)
	case s.UserTypeGoogle:
		return TTX2(c, ImageUSRLimiter, UnLimiter)
	default:
		return TTX2(c, ImageIPLimiter, UnLimiter)
	}
}
