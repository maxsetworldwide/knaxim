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

package main

import (
	"flag"
	"fmt"

	"git.maxset.io/web/knaxim/internal/database/types"
)

func adjustRole(add bool) {
	setup(false)
	if flag.NArg() <= 2 {
		if add {
			fmt.Println(helpstrs["addrole"])
		} else {
			fmt.Println(helpstrs["removerole"])
		}
		return
	}
	var username = flag.Arg(1)
	adjustUser(username, func(user types.UserI) (types.UserI, error) {
		var role = flag.Arg(2)
		user.SetRole(role, add)
		vPrintf("added %s to user %s", role, username)
		return user, nil
	})
	vPrintf("complete\n")
}
