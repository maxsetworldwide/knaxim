// Copyright August 2020 Maxset Worldwide Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
