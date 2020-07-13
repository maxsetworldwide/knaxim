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

package mongo

import (
	"context"
	"sync"

	"git.maxset.io/web/knaxim/internal/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	//"log"
	"git.maxset.io/web/knaxim/pkg/srverror"
)

// Database implements database.Database for mongodb
type Database struct {
	trackOwners
	URI       string            `json:"uri"`
	DBName    string            `json:"db"`
	CollNames map[string]string `json:"coll"`
	client    *mongo.Client
	ctx       context.Context
	cancel    context.CancelFunc
}

// Init tests the connection to the database
// if reset is true, empties out the databases and sets up indexex
func (d *Database) Init(ctx context.Context, reset bool) error {
	d.CollNames = initcoll(d.CollNames)
	//try connecting
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	//testclient, err := mongo.Connect(ctx, d.URI)
	testclient, err := mongo.Connect(ctx, options.Client().ApplyURI(d.URI))
	if err != nil {
		return err
	}
	defer testclient.Disconnect(ctx)
	err = testclient.Ping(ctx, nil)
	if err != nil {
		return err
	}
	dbnames, err := testclient.ListDatabaseNames(ctx, bson.M{
		"name": d.DBName,
	})
	if err != nil {
		return err
	}
	switch {
	case reset:
		if err = testclient.Database(d.DBName).Drop(ctx); err != nil {
			return err
		}
		fallthrough
	case len(dbnames) == 0:
		initIndexes := []func(context.Context, *Database, *mongo.Client) error{
			initViewIndex,
			initFileIndex,
			initUserIndex,
			initResetIndex,
			initGroupIndex,
			initChunkIndex,
			initStoreIndex,
			initAcronymIndex,
			initContentIndex,
			initStoreTagIndex,
			initFileTagsIndex,
		}
		var wg sync.WaitGroup
		wg.Add(len(initIndexes))
		indexctx, cancel := context.WithCancel(ctx)
		defer cancel()
		cherr := make(chan error, len(initIndexes))
		for _, initI := range initIndexes {
			go func(i func(context.Context, *Database, *mongo.Client) error) {
				defer wg.Done()
				if err := i(indexctx, d, testclient); err != nil {
					select {
					case cherr <- err:
					case <-indexctx.Done():
					}
				}
			}(initI)
		}
		wg.Wait()
		cherr <- nil
		err := <-cherr
		if err != nil {
			return err
		}
	}
	return testclient.Disconnect(ctx)
}

func initcoll(c map[string]string) map[string]string {
	if c == nil {
		c = make(map[string]string)
	}
	if _, ok := c["user"]; !ok {
		c["user"] = "user"
	}
	if _, ok := c["group"]; !ok {
		c["group"] = "group"
	}
	if _, ok := c["file"]; !ok {
		c["file"] = "file"
	}
	if _, ok := c["store"]; !ok {
		c["store"] = "store"
	}
	if _, ok := c["lines"]; !ok {
		c["lines"] = "lines"
	}
	if _, ok := c["chunk"]; !ok {
		c["chunk"] = "chunk"
	}
	if _, ok := c["filetags"]; !ok {
		c["filetags"] = "filetags"
	}
	if _, ok := c["storetags"]; !ok {
		c["storetags"] = "storetags"
	}
	if _, ok := c["acronym"]; !ok {
		c["acronym"] = "acronym"
	}
	if _, ok := c["reset"]; !ok {
		c["reset"] = "reset"
	}
	if _, ok := c["view"]; !ok {
		c["view"] = "view"
	}
	return c
}

// Owner opens a new connection to the database if provided a context and returns Ownerbase type
// if provided context is nil the resulting Ownerbase will reuse the existing connection
func (d *Database) Owner() database.Ownerbase {
	n := new(Ownerbase)
	n.Database = *d
	return n
}

// File opens a new connection to the database if provided a context and returns Filebase type
// if provided context is nil it will reuse the existing connection
func (d *Database) File() database.Filebase {
	n := new(Filebase)
	n.Database = *d
	return n
}

// Store opens a new connection to the database if provided a context and returns Storebase type
// if provided context is nil it will reuse the existing connection
func (d *Database) Store() database.Storebase {
	n := new(Storebase)
	n.Database = *d
	return n
}

// Content opens a new connection to the database if provided a context and returns Contentbase type
// if provided context is nil it will reuse the existing connection
func (d *Database) Content() database.Contentbase {
	n := new(Contentbase)
	n.Database = *d
	return n
}

// Tag opens a new connection to the database if provided a context and returns Tagbase type
// if provided context is nil it will reuse the existing connection
func (d *Database) Tag() database.Tagbase {
	n := new(Tagbase)
	n.Database = *d
	return n
}

// Acronym opens a new connection to the database if provided a context and returns Acronymbase type
// if provided context is nil it will reuse the existing connection
func (d *Database) Acronym() database.Acronymbase {
	n := new(Acronymbase)
	n.Database = *d
	return n
}

// View opens a new connection to the database if provided a context and returns Viewbase type
// if provided context is nil it will reuse the existing connection
func (d *Database) View() database.Viewbase {
	n := new(Viewbase)
	n.Database = *d
	return n
}

// Connect establishes a new connection to the mongodb
func (d *Database) Connect(ctx context.Context) (database.Database, error) {
	nd := new(Database)
	*nd = *d
	nd.trackOwners = newTrackOwners()
	var err error
	nd.ctx, nd.cancel = context.WithCancel(ctx)
	nd.client, err = mongo.Connect(nd.ctx, options.Client().ApplyURI(d.URI))
	return nd, err
}

// Close stops the active connection necessary or else there can be memory leak from unclosed connections
func (d *Database) Close(ctx context.Context) error {
	defer func() {
		d.client = nil
		d.cancel = nil
		d.ctx = nil
	}()
	if ctx == nil {
		ctx = d.ctx
	}
	if d.client != nil {
		if err := d.client.Disconnect(ctx); err != nil {
			return srverror.New(err, 500, "Error 101", "unable to close connection to Database")
		}
	}
	if d.cancel != nil {
		d.cancel()
	}
	return nil
}

// GetContext returns context of the current open connection
func (d *Database) GetContext() context.Context {
	return d.ctx
}
