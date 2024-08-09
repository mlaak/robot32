// TODO: This file is a bit of a mess
package usersession

import (
	"fmt"
	. "grp/ttd" //lint:ignore ST1001 Sometimes dot imports are needed
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	UserTypeUnverified = "unverified_ip"
	UserTypeCaptchad   = "captchad_ip"
	UserTypeGithub     = "github"
	UserTypeGoogle     = "google"
)

func GetUser(c int64, req *http.Request) (int64, string, string) {
	ip := strings.Split(req.RemoteAddr, ":")[0]
	cook := getCookieValue(req, "r_ression_id")
	fmt.Println("Cookie is", cook)

	/*workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}

	fileReader := func(filename string)(string,error){
		str,err := ioutil.ReadFile(filename)
		return string(str),err
	}*/
	TTD(c, "Getting user", "Cookie", cook, "IP", ip)
	return GetUserByCookieIP(c, cook, ip, os.Getwd, ioutil.ReadFile)
	//return GetUserByCookieIPWDFileFunc(c,ip,workingDir,fileReader)
}

func GetUserByCookieIP(c int64, cook string, ip string, Getwd func() (string, error), ReadFile func(string) ([]byte, error)) (int64, string, string) {

	//workingDir string, fileReader func(fnam string)(string,error)

	if Getwd == nil {
		Getwd = os.Getwd
	}
	if ReadFile == nil {
		ReadFile = ioutil.ReadFile
	}

	if cook == "" || !IsHex(cook) {
		return c, UserTypeUnverified, ip
	}

	workingDir, err := Getwd()
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}

	fileName := cook + ".txt"
	filePath := filepath.Join(workingDir, "..", "working_data", "sessions", fileName)

	data, err := ReadFile(filePath)
	if err == nil { //file exists, could read
		fileContents := string(data)

		TTD(c, "read user file", "File", fileContents)

		parts := strings.SplitN(fileContents, ",", 3)
		if len(parts) != 3 {
			// log.Fatalf("Expected 3 parts (sessionhandler.go)")
			TTD(c, "Error parsing user file", "parts", parts)
			return c, UserTypeUnverified, ip
		}
		utype := strings.TrimSpace(parts[0])
		uid := strings.TrimSpace(parts[1])

		if utype == "ip" {
			if uid != "ipus"+ip {
				return c, UserTypeUnverified, ip
			} else {
				return c, UserTypeCaptchad, uid
			}
		} else if utype == "github" {
			return c, UserTypeGithub, uid
		} else if utype == "google" {
			return c, UserTypeGoogle, uid
		} else {
			return c, UserTypeUnverified, ip
		}
	} else {
		return c, UserTypeUnverified, ip
	}
}

func getCookieValue(r *http.Request, cookieName string) string {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return ""
	}
	return cookie.Value
}

var hexRegex = regexp.MustCompile("^[0-9a-fA-F]+$")

func IsHex(s string) bool {
	return hexRegex.MatchString(s)
}
