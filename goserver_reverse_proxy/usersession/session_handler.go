//TODO: This file is a bit of a mess
package usersession

import (
	"fmt"
    "os"
	"log"
	"net/http"
	"strings" 
	"io/ioutil"
    "regexp"
    "path/filepath"
)


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
    
    
    /*workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}

	fileReader := func(filename string)(string,error){
		str,err := ioutil.ReadFile(filename)
		return string(str),err
	}*/
	return GetUserByCookieIP(c,ip,os.Getwd,ioutil.ReadFile)  
	//return GetUserByCookieIPWDFileFunc(c,ip,workingDir,fileReader)   
}

func GetUserByCookieIP(c string, ip string, Getwd func()(string,error), ReadFile func(string)([]byte,error))(string,string){

	//workingDir string, fileReader func(fnam string)(string,error)

	if(Getwd == nil){
		Getwd = os.Getwd
	}
	if(ReadFile == nil){
		ReadFile = ioutil.ReadFile
	}


	if c == "" || !IsHex(c){
        return UserTypeUnverified,ip
    }

	workingDir, err := Getwd()
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}
	
    fileName := c+".txt";
    filePath := filepath.Join(workingDir, "..", "working_data","sessions", fileName)   
    
    data, err := ReadFile(filePath)
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
	} else {
	    return UserTypeUnverified,ip
	}
}






func getCookieValue(r *http.Request, cookieName string) string {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return ""
	}
	return cookie.Value
}

var hexRegex = regexp.MustCompile("^[0-9a-fA-F]+$");
func IsHex(s string) bool {
	return hexRegex.MatchString(s)
}