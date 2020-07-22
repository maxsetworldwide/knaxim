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
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"git.maxset.io/web/knaxim/internal/config"
	"git.maxset.io/web/knaxim/internal/database/types"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "specify config file path")
	var configPathAlt string
	flag.StringVar(&configPathAlt, "c", "", "config alias")

	flag.Parse()

	if len(configPath) == 0 {
		configPath = configPathAlt
	}
	if len(configPath) == 0 {
		econfp := os.Getenv("KNAXIM_SERVER_CONFIG")
		if len(econfp) == 0 {
			configPath = "/etc/knaxim/conf.json"
		} else {
			configPath = econfp
		}
	}
}

func setup(initdb bool) {
	err := config.ParseConfig(configPath)
	if err != nil {
		log.Fatalf("unable to parse config: %s\n", err)
	}
	setupctx, cancel := context.WithTimeout(context.Background(), config.V.SetupTimeout.Duration)
	defer cancel()
	if err := config.DB.Init(setupctx, initdb); err != nil {
		log.Fatalf("database init error: %v\n", err)
	}
}

func main() {
	if flag.NArg() == 0 {
		fmt.Println(helpstr)
		return
	}
	switch strings.ToLower(flag.Arg(0)) {
	case "help":
		if flag.NArg() > 1 {
			message := helpstrs[strings.ToLower(flag.Arg(1))]
			if len(message) == 0 {
				message = helpstr
			}
			fmt.Println(message)
		} else {
			fmt.Println(helpstr)
		}
	case "addrole":
		adjustRole(true)
	case "removerole":
		adjustRole(false)
	case "updatefilecount":
		setup(false)
		if flag.NArg() < 3 {
			fmt.Println(helpstrs["updatefilecount"])
			return
		}
		var username = flag.Arg(1)
		var count, err = strconv.ParseInt(flag.Arg(2), 10, 64)
		if err != nil {
			log.Printf("count must be a number: %s\n%s\n", err, helpstrs["updatefilecount"])
			return
		}
		err = adjustUser(username, func(ui types.UserI) (types.UserI, error) {
			user, ok := ui.(*types.User)
			if !ok {
				return nil, fmt.Errorf("unrecognized user type")
			}
			user.Max = count
			return user, nil
		})
		if err != nil {
			log.Printf("adjust user error: %s\n", err)
			return
		}
		vPrintf("complete\n")
	case "updatefilespace":
		setup(false)
		if flag.NArg() < 3 {
			fmt.Println(helpstrs["updatefilespace"])
			return
		}
		var username = flag.Arg(1)
		var space, err = strconv.ParseInt(flag.Arg(2), 10, 64)
		if err != nil {
			log.Printf("count must be a number: %s\n%s\n", err, helpstrs["updatefilecount"])
			return
		}
		err = adjustUser(username, func(ui types.UserI) (types.UserI, error) {
			user, ok := ui.(*types.User)
			if !ok {
				return nil, fmt.Errorf("unrecognized user type")
			}
			user.Space = space
			return user, nil
		})
		if err != nil {
			log.Printf("adjust user error: %s\n", err)
			return
		}
		vPrintf("complete\n")
	case "cleardb":
		fallthrough
	case "initdb":
		vPrintf("reseting the database...\n")
		initArgs := flag.NewFlagSet("knaximctl/initDB", flag.ExitOnError)
		skipConfirmation := initArgs.Bool("y", false, "do not wait for confirmation before resetting the database. Useful for scripts. Warning using this flag will cause the tool to delete the database without confirmation from the user.")
		initArgs.Parse(flag.Args()[1:])
		if !*skipConfirmation {
			fmt.Print("Confirm that you wish to initialize the database, warning this will delete any preexisting data (y or n): ")
			var answer byte
			fmt.Scanf("%c", &answer)
			if answer != 'y' && answer != 'n' {
				fmt.Print("\n Please respond with y or n: ")
				fmt.Scanf("%c", &answer)
			}
			if answer != 'y' {
				return
			}
			fmt.Println("\nInitializing database...")
		}
		setup(true)
		fmt.Println("Done.")
	case "addacronyms":
		vPrintf("Adding Acronyms to database\n")
		setup(false)
		if flag.NArg() < 2 {
			fmt.Println(helpstrs["addacronyms"])
			return
		}
		filepath := flag.Arg(1)
		filein, err := os.Open(filepath)
		if err != nil {
			log.Printf("unable to open %s: %s", filepath, err)
			return
		}
		defer filein.Close()
		err = loadAcronyms(filein)
		if err != nil {
			log.Printf("unable to load acronyms: %s", err)
			return
		}
	case "adduser":
		setup(false)
		if flag.NArg() < 3 {
			fmt.Println(helpstrs["adduser"])
			return
		}
		username := flag.Arg(1)
		email := flag.Arg(2)
		var pass string
		if flag.NArg() == 3 {
			pass = generatePass()
		} else {
			pass = flag.Arg(3)
		}
		u, err := newUser(username, email, pass)
		if err != nil {
			log.Printf("unable to add user: %s", err)
			return
		}
		wrtr := json.NewEncoder(os.Stdout)
		wrtr.SetIndent("", "\t")
		if err = wrtr.Encode(u); err != nil {
			log.Printf("unable to output user data: %s", err)
			return
		}
		fmt.Printf("password: %s\n", pass)
	case "userinfo":
		setup(false)
		if flag.NArg() < 2 {
			fmt.Println(helpstrs["userinfo"])
			return
		}
		name := flag.Arg(1)
		err := userInfo(name)
		if err != nil {
			log.Printf("unable to output user data: %s", err)
		}
	default:
		fmt.Println("unrecognized command word.")
		fmt.Println(helpstr)
	}
}
