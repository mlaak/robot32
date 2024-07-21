package limits

import (
//	"fmt"
//	"io"
//	"log"
//	"net/http"
//	"net/http/httputil"
//	"net/url"
//	"errors"
	"strings" 
//	"bytes"
//	"io/ioutil"
    . "grp/ratelimiter"
    s "grp/session"
)


func GetLimitersForPathAndUserType(path string, usertype string) (*RateLimiter,*RateLimiter){
    p := path
    if strings.HasPrefix(p,"/experts/") { 
        return GetExpertsLimiters(path,usertype); 
        
    } else if strings.HasPrefix(p,"/recieved_images/") || strings.HasPrefix(p,"/welcome2.jpg") {
       return GetImagesLimiters(path,usertype); 

    } else if strings.HasPrefix(p,"/units/github_login") || strings.HasPrefix(p,"/units/google_login"){
        return loginLimiter,    unLimiter

    } else {
        return staticIPLimiter, unLimiter
    }
}

func GetExpertsLimiters(path string, usertype string) (*RateLimiter,*RateLimiter){
    switch usertype{
        case s.UserTypeUnverified: return code498Limiter,  code498Limiter
        case s.UserTypeCaptchad:   return expertIPLimiter, expertAIPLimiter
        case s.UserTypeGithub:     return expertUSRLimiter,unLimiter
        case s.UserTypeGoogle:     return expertUSRLimiter,unLimiter    
        default:                 return expertIPLimiter, expertAIPLimiter
    }
}

func GetImagesLimiters(path string, usertype string) (*RateLimiter,*RateLimiter){
    switch usertype{
        case s.UserTypeUnverified: return imageIPLimiter,  unLimiter
        case s.UserTypeCaptchad:   return imageIPLimiter,  unLimiter
        case s.UserTypeGithub:     return imageUSRLimiter, unLimiter  
        case s.UserTypeGoogle:     return imageUSRLimiter, unLimiter
        default:                 return imageIPLimiter,  unLimiter
    }
}
