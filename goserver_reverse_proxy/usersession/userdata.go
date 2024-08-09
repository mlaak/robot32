package usersession

import (
	. "grp/ttd" //lint:ignore ST1001 Sometimes dot imports are needed
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

var isAlphanumericRegex *regexp.Regexp = regexp.MustCompile("^[a-zA-Z0-9_]+$")

func isAlphanumeric_(s string) bool {
	return isAlphanumericRegex.MatchString(s)
}

func getUserProperty(c int64, username string, property string) string {
	TTD(c, "getUserProperty", "username", username, "property", property)

	if !isAlphanumeric_(username) || !isAlphanumeric_(property) {
		return TTX(c, "") // Return empty string for invalid input
	}

	cwd, err := os.Getwd()
	if err != nil {
		// Handle error getting current working directory
		TTD(c, "Error getting current directory???", "err", err)
		return TTX(c, "")
	}

	filePath := filepath.Join(cwd, "..", "working_data", "users", username, property+".txt")
	/**/ TTD(c, "", "cwd", cwd, "filePath", filePath)

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return TTX(c, "") // File doesn't exist, return an empty string
		}
		TTD(c, "General error reading file??? We should do something, like handle it?", "err", err)
		return TTX(c, "")
	}

	// Return the content of the file as a string
	return TTX(c, string(content))
}

func GetUserData(c int64, usertype string, userid string, data string) string {
	TTD(c, "getUserData", "usertype", usertype, "userid", userid, "data", data)
	if usertype != UserTypeGithub && usertype != UserTypeGoogle {
		return TTX(c, "")
	}
	if !isAlphanumeric_(userid) || !isAlphanumeric_(data) {
		return TTX(c, "")
	}
	return TTX(c, getUserProperty(c, userid, data))
}
