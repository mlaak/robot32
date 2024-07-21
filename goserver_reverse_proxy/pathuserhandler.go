package main

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
)

func PathUserRateLimitersSelect(path string, usertype string) (*RateLimiter,*RateLimiter){
    p := path
    if strings.HasPrefix(p,"/experts/") { 
        return ExpertsLimitersSelect(path,usertype); 
        
    } else if strings.HasPrefix(p,"/recieved_images/") || strings.HasPrefix(p,"/welcome2.jpg") {
       return ImagesLimitersSelect(path,usertype); 

    } else if strings.HasPrefix(p,"/units/github_login") || strings.HasPrefix(p,"/units/google_login"){
        return loginLimiter,    unLimiter

    } else {
        return staticIPLimiter, unLimiter
    }
}

func ImagesLimitersSelect(path string, usertype string) (*RateLimiter,*RateLimiter){
    switch usertype{
        case UserTypeUnverified: return imageIPLimiter,  unLimiter
        case UserTypeCaptchad:   return imageIPLimiter,  unLimiter
        case UserTypeGithub:     return imageUSRLimiter, unLimiter  
        case UserTypeGoogle:     return imageUSRLimiter, unLimiter
        default:                 return imageIPLimiter,  unLimiter
    }
}
func ExpertsLimitersSelect(path string, usertype string) (*RateLimiter,*RateLimiter){
    switch usertype{
        case UserTypeUnverified: return code498Limiter,  code498Limiter
        case UserTypeCaptchad:   return expertIPLimiter, expertAIPLimiter
        case UserTypeGithub:     return expertUSRLimiter,unLimiter
        case UserTypeGoogle:     return expertUSRLimiter,unLimiter    
        default:                 return expertIPLimiter, expertAIPLimiter
    }
}
