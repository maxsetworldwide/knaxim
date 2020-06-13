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
