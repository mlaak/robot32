package limits

import (
	"strings" 
    . "grp/ratelimiter"
    s "grp/usersession"
)


func GetLimitersForPathAndUserType(path string, usertype string) (IRateLimiter,IRateLimiter){
    p := path
    if strings.HasPrefix(p,"/experts/") { 
        return GetExpertsLimiters(path,usertype); 
        
    } else if strings.HasPrefix(p,"/recieved_images/") || strings.HasPrefix(p,"/welcome2.jpg") {
       return GetImagesLimiters(path,usertype); 

    } else if strings.HasPrefix(p,"/units/github_login") || strings.HasPrefix(p,"/units/google_login"){
        return LoginLimiter,    UnLimiter

    } else {
        return StaticIPLimiter, UnLimiter
    }
}

func GetExpertsLimiters(path string, usertype string) (IRateLimiter,IRateLimiter){
    switch usertype{
        case s.UserTypeUnverified: return Code498Limiter,  Code498Limiter
        case s.UserTypeCaptchad:   return ExpertIPLimiter, ExpertAIPLimiter
        case s.UserTypeGithub:     return ExpertUSRLimiter,UnLimiter
        case s.UserTypeGoogle:     return ExpertUSRLimiter,UnLimiter    
        default:                   return ExpertIPLimiter, ExpertAIPLimiter
    }
}

func GetImagesLimiters(path string, usertype string) (IRateLimiter,IRateLimiter){
    switch usertype{
        case s.UserTypeUnverified: return ImageIPLimiter,  UnLimiter
        case s.UserTypeCaptchad:   return ImageIPLimiter,  UnLimiter
        case s.UserTypeGithub:     return ImageUSRLimiter, UnLimiter  
        case s.UserTypeGoogle:     return ImageUSRLimiter, UnLimiter
        default:                   return ImageIPLimiter,  UnLimiter
    }
}
