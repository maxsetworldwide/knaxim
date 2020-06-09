package main

var helpstr = `knaximctl is for admin controls of a knaxim server
command form: knaximctl [Options]... [Action] [Arguments]...

Actions
help
addRole
removeRole
updateFileCount
updateFileSpace
`

var helpstrs = map[string]string{
	"help":            `display help message of actions`,
	"addRole":         `adds a role to a user`,
	"removeRole":      `removes a role to a user`,
	"updateFileCount": `change the file count limit for a user`,
	"updateFileSpace": `change the file space limit for a user`,
}
