package memory

import (
	"context"
	"errors"
	"sync"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/tag"
)

type Database struct {
	ctx    context.Context
	lock   *sync.RWMutex
	Owners struct {
		ID        map[string]database.Owner // key Owner.ID.String()
		UserName  map[string]database.UserI
		GroupName map[string]database.GroupI
	}
	Files     map[string]database.FileI         // key filehash.FileID.String()
	Stores    map[string]*database.FileStore    // key filehash.StoreID.String()
	Lines     map[string][]database.ContentLine // key filehash.StoreID.String()
	TagFiles  map[string]map[string]tag.Tag     // key filehash.FileID.String() => word string => tag
	TagStores map[string]map[string]tag.Tag     // key filehash.StoreID.String() => word string => tag
	Acronyms  map[string][]string
}

func (db *Database) Init(_ context.Context, reset bool) error {
	if db == nil {
		return errors.New("Memory Database Unallocated")
	}
	if !reset {
		return nil
	}
	db.lock = new(sync.RWMutex)
	db.Owners.ID = make(map[string]database.Owner)
	db.Owners.UserName = make(map[string]database.UserI)
	db.Owners.GroupName = make(map[string]database.GroupI)
	db.Files = make(map[string]database.FileI)
	db.Stores = make(map[string]*database.FileStore)
	db.Lines = make(map[string][]database.ContentLine)
	db.FileTags = make(map[string]map[string]tag.Tag)
	db.StoreTags = make(map[string]map[string]tag.Tag)
	db.Acronyms = make(map[string][]string)
	return nil
}

func (db *Database) Owner(c context.Context) database.Ownerbase {
	out := &Ownerbase{
		Database: *db,
	}
	out.ctx = c
	return out
}

func (db *Database) File(c context.Context) database.Filebase {
	out := &Filebase{
		Database: *db,
	}
	out.ctx = c
	return out
}

func (db *Database) Store(c context.Context) database.Storebase {
	out := &Storebase{
		Database: *db,
	}
	out.ctx = c
	return out
}

func (db *Database) Content(c context.Context) database.Contentbase {
	out := &Contentbase{
		Database: *db,
	}
	out.ctx = c
	return out
}

func (db *Database) Tag(c context.Context) database.Tagbase {
	out := &Tagbase{
		Database: *db,
	}
	out.ctx = c
	return out
}

func (db *Database) Acronym(c context.Context) database.Acronymbase {
	out := &Acronymbase{
		Database: *db,
	}
	out.ctx = c
	return out
}

func (db *Database) Close(_ context.Context) error {
	db.ctx = nil
}

func (db *Database) GetContext() context.Context {
	return db.ctx
}
