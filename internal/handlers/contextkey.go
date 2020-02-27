package handlers

// ContextKey type used to map in values into request context
type ContextKey byte

const (
	// USER is for the currect user
	USER ContextKey = iota
	// GROUP is for the group the handler is meant to operate on or with
	GROUP
)
