package main

import (
	"fmt"
//	"io"
//	"log"
	"net/http"
//	"net/http/httputil"
//	"net/url"
//	"errors"
	"strings" 
//	"bytes"
//	"io/ioutil"
)


func RateLimitersSelect(req *http.Request) (*RateLimiter,*RateLimiter){
    p := req.URL.Path
    fmt.Println("!!!req URL",p)
    
    _, _ = GetUser(req)
    // /experts/general
    
    if strings.HasPrefix(p,"/experts/") {
        return expertIPLimiter, expertAIPLimiter
    } else if strings.HasPrefix(p,"/recieved_images/") || strings.HasPrefix(p,"/welcome2.jpg") {
        return imageIPLimiter, unLimiter
    } else if strings.HasPrefix(p,"/units/github_login") || strings.HasPrefix(p,"/units/google_login"){
        return loginLimiter, unLimiter
    } else {
        return staticIPLimiter, unLimiter
    }
 
}





