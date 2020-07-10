/*************************************************************************
 *
 * MAXSET CONFIDENTIAL
 * __________________
 *
 *  [2019] - [2020] Maxset WorldWide Inc.
 *  All Rights Reserved.
 *
 * NOTICE:  All information contained herein is, and remains
 * the property of Maxset WorldWide Inc. and its suppliers,
 * if any.  The intellectual and technical concepts contained
 * herein are proprietary to Maxset WorldWide Inc.
 * and its suppliers and may be covered by U.S. and Foreign Patents,
 * patents in process, and are protected by trade secret or copyright law.
 * Dissemination of this information or reproduction of this material
 * is strictly forbidden unless prior written permission is obtained
 * from Maxset WorldWide Inc.
 */

package handlers

import (
	"regexp"

	"git.maxset.io/web/knaxim/internal/util"
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
		util.Verbose("invalid email form (%s): %s", email, err.Error())
		return false
	}
	// if err := checkmail.ValidateHost(email); err != nil {
	// 	verbose("email not found (%s): %s", email, err.Error())
	// 	return false
	// }
	return true
}
