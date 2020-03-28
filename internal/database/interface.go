package database

import (
	"context"

	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/internal/database/types/tag"
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
	Reserve(id types.OwnerID, name string) (types.OwnerID, error)
	Insert(u types.Owner) error
	Get(id types.OwnerID) (types.Owner, error)
	FindUserName(name string) (types.UserI, error)
	FindGroupName(name string) (types.GroupI, error)
	GetGroups(id types.OwnerID) (owned []types.GroupI, member []types.GroupI, err error)
	Update(u types.Owner) error
	GetSpace(o types.OwnerID) (int64, error)
	GetTotalSpace(o types.OwnerID) (int64, error)
	GetResetKey(id types.OwnerID) (key string, err error)
	CheckResetKey(key string) (id types.OwnerID, err error)
	DeleteResetKey(id types.OwnerID) error
}

// Filebase is a database connection for file operations
type Filebase interface {
	Database
	Reserve(id types.FileID) (types.FileID, error)
	Insert(r types.FileI) error
	Get(fid types.FileID) (types.FileI, error)
	GetAll(fids ...types.FileID) ([]types.FileI, error)
	Update(r types.FileI) error
	Remove(r types.FileID) error
	GetOwned(uid types.OwnerID) ([]types.FileI, error)
	GetPermKey(uid types.OwnerID, pkey string) ([]types.FileI, error) // does not include owned records
	MatchStore(types.OwnerID, []types.StoreID, ...string) ([]types.FileI, error)
}

// Storebase is a database connection for file store operations
type Storebase interface {
	Database
	Reserve(id types.StoreID) (types.StoreID, error)
	Insert(fs *types.FileStore) error
	Get(id types.StoreID) (*types.FileStore, error)
	MatchHash(h uint32) ([]*types.FileStore, error)
	UpdateMeta(fs *types.FileStore) error
}

// Contentbase is a database connection for the content operations
type Contentbase interface {
	Database
	Insert(...types.ContentLine) error
	Len(id types.StoreID) (int64, error)
	Slice(id types.StoreID, start int, end int) ([]types.ContentLine, error)
	RegexSearchFile(regex string, file types.StoreID, start int, end int) ([]types.ContentLine, error)
}

// Tagbase is a database connection for the tag operations
type Tagbase interface {
	Database
	UpsertFile(types.FileID, ...tag.Tag) error
	UpsertStore(types.StoreID, ...tag.Tag) error
	FileTags(...types.FileID) (map[string][]tag.Tag, error)
	GetFiles([]tag.Tag, ...types.FileID) ([]types.FileID, []types.StoreID, error)
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
	Insert(*types.ViewStore) error
	Get(types.StoreID) (*types.ViewStore, error)
}
