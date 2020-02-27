package database

import (
	"context"
	"errors"

	"git.maxset.io/web/knaxim/internal/database/filehash"
	"git.maxset.io/web/knaxim/internal/database/tag"
	"git.maxset.io/web/knaxim/pkg/srverror"
)

// ContextKey is used to store connections to a database in the values of a context
type ContextKey byte

// Context Keys for each type of database connection
const (
	OWNER ContextKey = iota
	FILE
	STORE
	CONTENT
	TAG
	ACRONYM
	VIEW
)

// Error types for use across different database implementations
var (
	ErrNotFound       = srverror.New(errors.New("Not Found in Database"), 404, "Not Found")
	ErrNoResults      = srverror.Basic(204, "Empty", "No results found")
	ErrNameTaken      = srverror.New(errors.New("Id is already in use"), 409, "Name Already Taken")
	ErrCorruptData    = srverror.New(errors.New("unable to decode data from the database"), 500, "Database Error 010")
	ErrPermission     = srverror.New(errors.New("User does not have appropriate permission"), 403, "Permission Denied")
	ErrIDNotReserved  = srverror.Basic(500, "Database Error 011", "ID has not been reserved for Insert")
	ErrIDUnrecognized = srverror.Basic(400, "Unrecognized ID")

	FileLoadInProgress = &ProcessingError{Status: 202, Message: "Processing File"}
)

// Database is the root Database interface
type Database interface {
	Init(context.Context, bool) error
	Owner(context.Context) Ownerbase
	File(context.Context) Filebase
	Store(context.Context) Storebase
	Content(context.Context) Contentbase
	Tag(context.Context) Tagbase
	Acronym(context.Context) Acronymbase
	View(context.Context) Viewbase
	Close(context.Context) error
	GetContext() context.Context
}

// Ownerbase is a database connection for owner related actions
type Ownerbase interface {
	Database
	Reserve(id OwnerID, name string) (OwnerID, error)
	Insert(u Owner) error
	Get(id OwnerID) (Owner, error)
	FindUserName(name string) (UserI, error)
	FindGroupName(name string) (GroupI, error)
	GetGroups(id OwnerID) (owned []GroupI, member []GroupI, err error)
	Update(u Owner) error
	GetSpace(o OwnerID) (int64, error)
	GetTotalSpace(o OwnerID) (int64, error)
	GetResetKey(id OwnerID) (key string, err error)
	CheckResetKey(key string) (id OwnerID, err error)
	DeleteResetKey(id OwnerID) error
}

// Filebase is a database connection for file operations
type Filebase interface {
	Database
	Reserve(id filehash.FileID) (filehash.FileID, error)
	Insert(r FileI) error
	Get(fid filehash.FileID) (FileI, error)
	GetAll(fids ...filehash.FileID) ([]FileI, error)
	Update(r FileI) error
	Remove(r filehash.FileID) error
	GetOwned(uid OwnerID) ([]FileI, error)
	GetPermKey(uid OwnerID, pkey string) ([]FileI, error) // does not include owned records
	MatchStore(OwnerID, []filehash.StoreID, ...string) ([]FileI, error)
}

// Storebase is a database connection for file store operations
type Storebase interface {
	Database
	Reserve(id filehash.StoreID) (filehash.StoreID, error)
	Insert(fs *FileStore) error
	Get(id filehash.StoreID) (*FileStore, error)
	MatchHash(h uint32) ([]*FileStore, error)
	UpdateMeta(fs *FileStore) error
}

// Contentbase is a database connection for the content operations
type Contentbase interface {
	Database
	Insert(...ContentLine) error
	Len(id filehash.StoreID) (int64, error)
	Slice(id filehash.StoreID, start int, end int) ([]ContentLine, error)
	RegexSearchFile(regex string, file filehash.StoreID, start int, end int) ([]ContentLine, error)
}

// Tagbase is a database connection for the tag operations
type Tagbase interface {
	Database
	UpsertFile(filehash.FileID, ...tag.Tag) error
	UpsertStore(filehash.StoreID, ...tag.Tag) error
	FileTags(...filehash.FileID) (map[string][]tag.Tag, error)
	GetFiles([]tag.Tag, ...filehash.FileID) ([]filehash.FileID, []filehash.StoreID, error)
	SearchData(tag.Type, tag.Data) ([]tag.Tag, error)
}

// Acronymbase is a database connection for the acronym operations
type Acronymbase interface {
	Database
	Put(string, string) error
	Get(string) ([]string, error)
}

// Viewbase is a database connection for the view operations
type Viewbase interface {
	Database
	Insert(*ViewStore) error
	Get(filehash.StoreID) (*ViewStore, error)
}
