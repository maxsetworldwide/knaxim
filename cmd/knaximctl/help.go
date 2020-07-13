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

var helpstr = `knaximctl is for admin controls of a knaxim server
command form: knaximctl [Options]... [Action] [Arguments]...

Actions

addUser
userInfo
addRole
removeRole
updateFileCount
updateFileSpace
initDB
addAcronyms
help
`

var helpstrs = map[string]string{
	"help":            "display help message of actions\nknaximctl help",
	"addrole":         "adds a role to a user\nknaximctl addRole [username] [role]",
	"removerole":      "removes a role to a user\nknaximctl removeRole [username] [role]",
	"updatefilecount": "change the file count limit for a user\nknaimxctl updateFileCount [username] [amount]",
	"updatefilespace": "change the file space limit for a user\nknaximctl updateFileSpace [username] [amount]",
	"initdb":          "initialize database\nknaximctl initDB",
	"addacronyms":     "add acyonyms to database\nknaximctl addAcronyms [filepath]",
	"adduser":         "add user to database\nknaximctl addUser [username] [email] [password,optional]",
	"userinfo":        "display information about a user\nknaximctl userInfo [username]",
}
