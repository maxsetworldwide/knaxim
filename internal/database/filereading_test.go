package database_test

import (
	"context"
	"strings"
	"testing"
	"time"

	. "git.maxset.io/web/knaxim/internal/database"
	mongodb "git.maxset.io/web/knaxim/internal/database/mongo"
)

var db = &mongodb.Database{
	URI:    "mongodb://localhost:27017",
	DBName: "TestInjest",
}

var content = `This is an example text.
How do you know this is correct. Well you got this text back.
Good Bye.`

var testOwner = &User{
	ID: OwnerID{
		Type:        'u',
		UserDefined: [3]byte{'a', 'b', 'c'},
		Stamp:       []byte("test"),
	},
	Name: "testuser",
}

func TestInjustFile(t *testing.T) {
	//t.Parallel()
	initctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	err := db.Init(initctx, true)
	if err != nil {
		t.Fatal("unable to init database", err)
	}
	injestctx, cancel2 := context.WithTimeout(context.Background(), time.Minute)
	defer cancel2()
	file := &File{
		Name: "test.txt",
	}
	_, err = db.Owner(injestctx).Reserve(testOwner.GetID(), testOwner.GetName())
	if err != nil {
		t.Fatal("unable to reserve testOwner:", err)
	}
	err = db.Owner(injestctx).Insert(testOwner)
	file.Own = testOwner
	_, err = InjestFile(injestctx, file, "content/txt", strings.NewReader(content), db)
	if err != nil {
		t.Fatal("injest failed", err)
	}
}
