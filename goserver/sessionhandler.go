package main

import (
	"fmt"
//	"io"
    "os"
	"log"
	"net/http"
//	"net/http/httputil"
//	"net/url"
//	"errors"
	"strings" 
//	"bytes"
	"io/ioutil"
    "regexp"
    "path/filepath"
)

var hexRegex = regexp.MustCompile("^[0-9a-fA-F]+$");

func IsHex(s string) bool {
	return hexRegex.MatchString(s)
}

func getCookieValue(r *http.Request, cookieName string) string {

    /*cookies := r.Cookies()
    for _, cookie := range cookies {
        fmt.Println("Cookie name:", cookie.Name)
        fmt.Println("Cookie value:", cookie.Value)
    }*/


	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return ""
	}
	return cookie.Value
}

const (
    UserTypeUnverified = "unverified_ip"
    UserTypeCaptchad = "captchad_ip"
    UserTypeGithub = "github"
    UserTypeGoogle = "google"   
)

func GetUser(req *http.Request)(string,string){
    ip := strings.Split(req.RemoteAddr, ":")[0]
    c := getCookieValue(req,"r_ression_id")
    fmt.Println("Cookie is",c)
    
    if !IsHex(c){
        c = "";
    }
    
    if c == ""{
        return UserTypeUnverified,ip
    }
    
    workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}
	
    fileName := c+".txt";
    filePath := filepath.Join(workingDir, "..", "working_data","sessions", fileName)   
    
    data, err := ioutil.ReadFile(filePath)
	if err == nil { //file exists, could read
	    fileContents := string(data)
	    
	    parts := strings.SplitN(fileContents, ",", 3)
	    if len(parts)!=3{
	        log.Fatalf("Expected 3 parts (sessionhandler.go)")
	    }
	    utype := strings.TrimSpace(parts[0])
	    uid := strings.TrimSpace(parts[1])
	    
	    if utype == "ip"{
	        if uid != "ipus"+ip {
	            return UserTypeUnverified,ip
	        } else {
	            return UserTypeCaptchad,uid
	        }
	    } else if utype == "github"{
	        return UserTypeGithub, uid
	    } else if utype == "google"{
	        return UserTypeGoogle, uid
	    } else {
	        return UserTypeUnverified,ip
	    }
	    //fmt.Println("!!!!!!!!!",fileContents)
	    //return "1","1"
		
	} else {
	    return UserTypeUnverified,ip
	}
       
}
