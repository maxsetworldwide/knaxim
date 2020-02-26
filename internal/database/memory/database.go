package memory

// This package provides an in memory implementation of the database interface.
// Primarily used for testing.

import (
	"context"
	"errors"
	"sync"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/tag"
)

var lock = new(sync.RWMutex)

// Database is an implementation of database.Database that operates within local machine memory
type Database struct {
	ctx context.Context

	Owners struct {
		ID        map[string]database.Owner // key Owner.ID.String()
		UserName  map[string]database.UserI
		GroupName map[string]database.GroupI
		Reset     map[string]database.OwnerID
	}
	Files     map[string]database.FileI         // key filehash.FileID.String()
	Stores    map[string]*database.FileStore    // key filehash.StoreID.String()
	Lines     map[string][]database.ContentLine // key filehash.StoreID.String()
	TagFiles  map[string]map[string]tag.Tag     // key filehash.FileID.String() => word string => tag
	TagStores map[string]map[string]tag.Tag     // key filehash.StoreID.String() => word string => tag
	Views     map[string]*database.ViewStore    // key filehash.StoreID.String()
	Acronyms  map[string][]string
}

// Init preps an instance of the Database for use. if reset is true, it will allocate new maps to store the
// data
func (db *Database) Init(_ context.Context, reset bool) error {
	if db == nil {
		return errors.New("Memory Database Unallocated")
	}
	if !reset {
		return nil
	}
	lock.Lock()
	defer lock.Unlock()
	db.Owners.ID = make(map[string]database.Owner)
	db.Owners.UserName = make(map[string]database.UserI)
	db.Owners.GroupName = make(map[string]database.GroupI)
	db.Owners.Reset = make(map[string]database.OwnerID)
	db.Files = make(map[string]database.FileI)
	db.Stores = make(map[string]*database.FileStore)
	db.Lines = make(map[string][]database.ContentLine)
	db.TagFiles = make(map[string]map[string]tag.Tag)
	db.TagStores = make(map[string]map[string]tag.Tag)
	db.Views = make(map[string]*database.ViewStore)
	db.Acronyms = make(map[string][]string)
	return nil
}

var connectionCount int
var countLock sync.Mutex

// CurrentOpenConnections returns the current number of open connections to the database.
func CurrentOpenConnections() int {
	countLock.Lock()
	defer countLock.Unlock()
	return connectionCount
}

// Owner opens a connection to the database and returns Ownerbase wrapping of the
// Database
func (db *Database) Owner(c context.Context) database.Ownerbase {
	lock.Lock()
	defer lock.Unlock()
	countLock.Lock()
	defer countLock.Unlock()
	out := &Ownerbase{
		Database: *db,
	}
	out.ctx = c
	connectionCount++
	return out
}

// File opens a connection to the database and returns Filebase wrapping of the
// Database
func (db *Database) File(c context.Context) database.Filebase {
	lock.Lock()
	defer lock.Unlock()
	countLock.Lock()
	defer countLock.Unlock()
	out := &Filebase{
		Database: *db,
	}
	out.ctx = c
	connectionCount++
	return out
}

// Store opens a connection to the database and returns Storebase wrapping of the
// Database
func (db *Database) Store(c context.Context) database.Storebase {
	lock.Lock()
	defer lock.Unlock()

	return db.store(c)
}

func (db *Database) store(c context.Context) database.Storebase {
	countLock.Lock()
	defer countLock.Unlock()
	out := &Storebase{
		Database: *db,
	}
	out.ctx = c
	connectionCount++
	return out
}

// Content opens a connection to the database and returns Contentbase wrapping of the
// Database
func (db *Database) Content(c context.Context) database.Contentbase {
	lock.Lock()
	defer lock.Unlock()
	countLock.Lock()
	defer countLock.Unlock()
	out := &Contentbase{
		Database: *db,
	}
	out.ctx = c
	connectionCount++
	return out
}

// Tag opens a connection to the database and returns Tagbase wrapping of the
// Database
func (db *Database) Tag(c context.Context) database.Tagbase {
	lock.Lock()
	defer lock.Unlock()
	countLock.Lock()
	defer countLock.Unlock()
	out := &Tagbase{
		Database: *db,
	}
	out.ctx = c
	connectionCount++
	return out
}

// Acronym opens a connection to the database and returns Acronymbase wrapping of the Database
func (db *Database) Acronym(c context.Context) database.Acronymbase {
	lock.Lock()
	defer lock.Unlock()
	countLock.Lock()
	defer countLock.Unlock()
	out := &Acronymbase{
		Database: *db,
	}
	out.ctx = c
	connectionCount++
	return out
}

// View opens a connection to the database and returns Viewbase wrapping of the Database
func (db *Database) View(c context.Context) database.Viewbase {
	lock.Lock()
	defer lock.Unlock()
	countLock.Lock()
	defer countLock.Unlock()
	out := &Viewbase{
		Database: *db,
	}
	out.ctx = c
	connectionCount++
	return out
}

// Close closes the open connection, meant to be called by wrapping objects
func (db *Database) Close(_ context.Context) error {
	lock.Lock()
	defer lock.Unlock()
	return db.close()
}

func (db *Database) close() error {
	countLock.Lock()
	defer countLock.Unlock()
	db.ctx = nil
	connectionCount += -1
	return nil
}

// GetContext returns the context of the active connection
func (db *Database) GetContext() context.Context {
	lock.RLock()
	defer lock.RUnlock()
	return db.ctx
}
