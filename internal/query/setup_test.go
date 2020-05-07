package query

import (
	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/memory"
	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/internal/database/types/tag"
)

var DB database.Database

func init() {
	DB = new(memory.Database)
	DB.Init(nil, true)
	initOwners(DB)
	initFiles(DB)
}

var owners = []types.Owner{
	&types.User{
		ID: types.OwnerID{
			Type:        'u',
			UserDefined: [3]byte{'t', 'e', 's'},
			Stamp:       []byte{'t', 'u', 's', 'e', 'r'},
		},
		Name: "testuser",
	},
	&types.Group{
		ID: types.OwnerID{
			Type:        'g',
			UserDefined: [3]byte{'t', 'e', 's'},
			Stamp:       []byte{'t', 'g', 'r', 'o'},
		},
		Name: "testgroup",
		// initialization should make testuser the owner of this group so that:
		// Permission: types.Permission{
		// 	Own: owners[0],
		// },
	},
}

func initOwners(db database.Database) {
	//make testuser own testgroup
	owners[1].(*types.Group).Permission.Own = owners[0]
	//add to database
	for _, o := range owners {
		db.Owner().Reserve(o.GetID(), o.GetName())
		db.Owner().Insert(o)
	}
}

type filedata struct {
	Name  string
	ID    types.FileID
	Owner types.Owner
	Text  string
	Tags  []tag.Tag
}

var fileinfo = []filedata{
	filedata{
		Name:  "first.txt",
		Owner: owners[0],
		ID: types.FileID{
			StoreID: types.StoreID{
				Hash:  1,
				Stamp: 1,
			},
			Stamp: []byte{'1'},
		},
		Text: "This is the first test file.",
		Tags: []tag.Tag{
			tag.Tag{
				Word: "first",
				Type: tag.CONTENT,
			},
			tag.Tag{
				Word: "Bobby",
				Type: tag.TOPIC,
			},
			tag.Tag{
				Word: "test",
				Type: tag.PROCESS,
			},
		},
	},
	filedata{
		Name:  "second.txt",
		Owner: owners[0],
		ID: types.FileID{
			StoreID: types.StoreID{
				Hash:  2,
				Stamp: 2,
			},
			Stamp: []byte{'2'},
		},
		Text: "This is the second test file.",
		Tags: []tag.Tag{
			tag.Tag{
				Word: "second",
				Type: tag.CONTENT,
			},
			tag.Tag{
				Word: "Peggy",
				Type: tag.RESOURCE,
			},
			tag.Tag{
				Word: "test",
				Type: tag.PROCESS,
			},
		},
	},
	filedata{
		Name:  "third.txt",
		Owner: owners[1],
		ID: types.FileID{
			StoreID: types.StoreID{
				Hash:  3,
				Stamp: 3,
			},
			Stamp: []byte{'3'},
		},
		Text: "This is the third test file.",
		Tags: []tag.Tag{
			tag.Tag{
				Word: "third",
				Type: tag.CONTENT,
			},
			tag.Tag{
				Word: "Hank",
				Type: tag.USER,
			},
			tag.Tag{
				Word: "test",
				Type: tag.PROCESS,
			},
		},
	},
}

func initFiles(db database.Database) {
	for _, fd := range fileinfo {
		fs := &types.FileStore{
			ID:          fd.ID.StoreID,
			ContentType: "testtext",
		}
		db.Store().Reserve(fs.ID)
		db.Store().Insert(fs)
		file := &types.File{
			ID: fd.ID,
			Permission: types.Permission{
				Own: fd.Owner,
			},
			Name: fd.Name,
		}
		db.File().Reserve(file.ID)
		db.File().Insert(file)
		for _, t := range fd.Tags {
			ft := tag.FileTag{
				File:  fd.ID,
				Owner: fd.Owner.GetID(),
				Tag:   t,
			}
			db.Tag().Upsert(ft)
		}
	}
}
