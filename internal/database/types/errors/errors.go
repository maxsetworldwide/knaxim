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

package errors

import (
	"errors"

	"git.maxset.io/web/knaxim/pkg/srverror"
)

// Error types for use across different database implementations
var (
	ErrNotFound       = srverror.New(errors.New("Not Found in Database"), 404, "Not Found")
	ErrNoResults      = srverror.Basic(204, "Empty", "No results found")
	ErrNameTaken      = srverror.New(errors.New("Id is already in use"), 409, "Name Already Taken")
	ErrCorruptData    = srverror.New(errors.New("unable to decode data from the database"), 500, "Error 010")
	ErrPermission     = srverror.New(errors.New("User does not have appropriate permission"), 403, "Permission Denied")
	ErrIDNotReserved  = srverror.Basic(500, "Error 011", "ID has not been reserved for Insert")
	ErrIDUnrecognized = srverror.Basic(400, "Unrecognized ID")

	FileLoadInProgress = &Processing{Status: 202, Message: "Processing File"}
)

// Processing is a record of an error that occured during processing of a file, and how to respond
type Processing struct {
	Status  int    `json:"status" bson:"s"`
	Message string `json:"msg" bson:"m"`
}

// Equal is true if the status and message are equal
func (pe *Processing) Equal(oth *Processing) bool {
	if pe.Status != oth.Status {
		return false
	}
	return pe.Message == oth.Message
}

// Error implements error, returns message
func (pe *Processing) Error() string {
	return pe.Message
}
