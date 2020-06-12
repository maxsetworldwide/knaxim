package main

var helpstr = `knaximctl is for admin controls of a knaxim server
command form: knaximctl [Options]... [Action] [Arguments]...

Actions

addUser
addRole
removeRole
updateFileCount
updateFileSpace
initDB
addAcronyms
help
`

var helpstrs = map[string]string{
	"help":            `display help message of actions`,
	"addrole":         `adds a role to a user`,
	"removerole":      `removes a role to a user`,
	"updatefilecount": `change the file count limit for a user`,
	"updatefilespace": `change the file space limit for a user`,
	"initdb":          `initialize database`,
	"addacronyms":     `add acyonyms to database`,
	"adduser":         `add user to database knaximctl addUser [username] [email] [password,optional]`,
}
