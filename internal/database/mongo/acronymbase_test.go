package mongo

import (
	"context"
	"testing"
	"time"
)

func TestAcronym(t *testing.T) {
	t.Parallel()
	var ab *Acronymbase
	{
		db := new(Database)
		*db = *configuration.DB
		db.DBName = "TestAcronym"
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		if err := db.Init(ctx, true); err != nil {
			t.Fatal("Unable to Init database", err)
		}
		methodtesting, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		ab = db.Acronym(methodtesting).(*Acronymbase)
	}
	t.Run("Put", func(t *testing.T) {
		err := ab.Put("ab", "Acronymbase")
		if err != nil {
			t.Fatal("Unable to add acronym, ", err)
		}
	})
	if t.Failed() {
		t.FailNow()
	}
	t.Run("Get", func(t *testing.T) {
		result, err := ab.Get("ab")
		if err != nil {
			t.Fatal("Unable to get acronym, ", err)
		}
		if len(result) != 1 || result[0] != "Acronymbase" {
			t.Fatal("incorrect result: ", result)
		}
	})
}
