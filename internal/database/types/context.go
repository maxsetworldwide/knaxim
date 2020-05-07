package types

// ContextKey is used to store connections to a database in the values of a context
type ContextKey byte

// Context Keys for each type of database connection
const (
	DATABASE ContextKey = iota
	OWNER
	FILE
	STORE
	CONTENT
	TAG
	ACRONYM
	VIEW
)
