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

/*const (
    UserTypeUnverified = "unverified_ip"
    UserTypeCaptchad = "captchad_ip"
    UserTypeGithub = "github"
    UserTypeGoogle = "google"   
)*/

func PathUserRateLimitersSelect(path string, usertype string) (*RateLimiter,*RateLimiter){
    //p := req.URL.Path
     p := path
 
    //fmt.Println("!!!req URL",p)
    
    
    // /experts/general
    
    if strings.HasPrefix(p,"/experts/") {
        if usertype == UserTypeUnverified {
            return stopLimiter, stopLimiter
        }
        if usertype == UserTypeCaptchad {
            return expertIPLimiter, expertAIPLimiter    
        }
        if usertype == UserTypeGithub || usertype == UserTypeGoogle {
            return expertUSRLimiter, unLimiter    
        }
        return expertIPLimiter, expertAIPLimiter
        
    } else if strings.HasPrefix(p,"/recieved_images/") || strings.HasPrefix(p,"/welcome2.jpg") {
        if usertype == UserTypeUnverified {
            return imageIPLimiter, unLimiter
        }
        if usertype == UserTypeCaptchad {
            return imageIPLimiter, unLimiter    
        }
        if usertype == UserTypeGithub || usertype == UserTypeGoogle {
            return imageUSRLimiter, unLimiter    
        }
        return imageIPLimiter, unLimiter
        
    } else if strings.HasPrefix(p,"/units/github_login") || strings.HasPrefix(p,"/units/google_login"){
        return loginLimiter, unLimiter
    } else {
        return staticIPLimiter, unLimiter
    }
 
}


