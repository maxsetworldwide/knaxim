package database

import (
	"context"
	"errors"

	"git.maxset.io/web/knaxim/internal/database/filehash"
	"git.maxset.io/web/knaxim/internal/database/tag"
	"git.maxset.io/web/knaxim/pkg/srverror"
)

type ContextKey byte

const (
	OWNER ContextKey = iota
	FILE
	STORE
	CONTENT
	TAG
	ACRONYM
)

var (
	ErrNotFound       = srverror.New(errors.New("Not Found in Database"), 404, "Not Found")
	ErrNoResults      = srverror.Basic(204, "Empty", "No results found")
	ErrNameTaken      = srverror.New(errors.New("Id is already in use"), 409, "Name Already Taken")
	ErrCorruptData    = srverror.New(errors.New("unable to decode data from the database"), 500, "Database Error 010")
	ErrPermission     = srverror.New(errors.New("User does not have appropriate permission"), 403, "Permission Denied")
	ErrIDNotReserved  = srverror.Basic(500, "Database Error 011", "ID has not been reserved for Insert")
	ErrIDUnrecognized = srverror.Basic(400, "Unrecognized ID")
)

type Database interface {
	Init(context.Context, bool) error
	Owner(context.Context) Ownerbase
	File(context.Context) Filebase
	Store(context.Context) Storebase
	Content(context.Context) Contentbase
	Tag(context.Context) Tagbase
	Acronym(context.Context) Acronymbase
	Close(context.Context) error
	GetContext() context.Context
}

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
}

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

type Storebase interface {
	Database
	Reserve(id filehash.StoreID) (filehash.StoreID, error)
	Insert(fs *FileStore) error
	Get(id filehash.StoreID) (*FileStore, error)
	MatchHash(h uint32) ([]*FileStore, error)
	UpdateMeta(fs *FileStore) error
	//Get Total size
}

type Contentbase interface {
	Database
	Insert(...ContentLine) error
	Len(id filehash.StoreID) (int64, error)
	Slice(id filehash.StoreID, start int, end int) ([]ContentLine, error)
	RegexSearchFile(regex string, file filehash.StoreID, start int, end int) ([]ContentLine, error)
}

type Tagbase interface {
	Database
	UpsertFile(filehash.FileID, ...tag.Tag) error
	UpsertStore(filehash.StoreID, ...tag.Tag) error
	FileTags(...filehash.FileID) (map[string][]tag.Tag, error)
	GetFiles([]tag.Tag, ...filehash.FileID) ([]filehash.FileID, []filehash.StoreID, error)
	SearchData(tag.Type, tag.Data) ([]tag.Tag, error)
}

type Acronymbase interface {
	Database
	Put(string, string) error
	Get(string) ([]string, error)
}
