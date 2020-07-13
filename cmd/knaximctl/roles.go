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
