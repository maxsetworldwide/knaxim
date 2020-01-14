package handlers

import (
	"regexp"

	"github.com/badoux/checkmail"
)

var validUserName = func(s string) bool {
	if s == "$PUBLIC" {
		return false
	}
	return regexp.MustCompile(`^[[:print:]]{6,100}$`).MatchString(s)
}
var validGroupName = regexp.MustCompile(`^[[:print:]]{3,100}$`).MatchString
var validDirName = regexp.MustCompile(`^[[:print:]]{1,100}$`).MatchString

func validEmail(email string) bool {
	if err := checkmail.ValidateFormat(email); err != nil {
		verbose("invalid email form (%s): %s", email, err.Error())
		return false
	}
	// if err := checkmail.ValidateHost(email); err != nil {
	// 	verbose("email not found (%s): %s", email, err.Error())
	// 	return false
	// }
	return true
}
