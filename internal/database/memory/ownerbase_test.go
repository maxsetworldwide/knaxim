package memory

import (
	"testing"

	"git.maxset.io/web/knaxim/internal/database"
)

func fillowners(db *Database) {
	//Users
	test1 := database.NewUser("testuser1", "testuserpass1", "test1@test.test")
	db.Owners.ID[test1.GetID().String()] = test1
	db.Owners.UserName[test1.GetName()] = test1

	test2 := database.NewUser("testuser2", "testuserpass2", "test2@test.test")
	if test1.GetID().Equal(test2.GetID()) {
		test2.ID = test2.ID.Mutate()
	}
	db.Owners.ID[test2.GetID().String()] = test2
	db.Owners.UserName[test2.GetName()] = test2

	group1 := database.NewGroup("group1", test1)
	group2 := database.NewGroup("group2", group1)
	group2.AddMember(test2)
	db.Owners.ID[group1.GetID().String()] = group1
	db.Owners.GroupName[group1.GetName()] = group1
	if group1.GetID().Equal(group2.GetID()) {
		group2.ID = group2.ID.Mutate()
	}
	db.Owners.ID[group2.GetID().String()] = group2
	db.Owners.GroupName[group2.GetName()] = group2

}

func TestOwners(t *testing.T) {
	t.Parallel()
	ob := DB.Owner(nil)
	defer ob.Close(nil)

	newUser := database.NewUser("testuser3", "testuserpass3", "test3@test.test")
	newGroup := database.NewGroup("group3", newUser)
	t.Run("Reserve", func(t *testing.T) {
		var err error
		newUser.ID, err = ob.Reserve(newUser.ID, newUser.Name)
		if err != nil {
			t.Fatalf("Unable to Reserve User: %+v", newUser)
		}
		newGroup.ID, err = ob.Reserve(newGroup.ID, newGroup.Name)
		if err != nil {
			t.Fatalf("Unable to Reserve Group: %+v", newGroup)
		}
	})
	if t.Failed() {
		t.FailNow()
	}
	t.Run("Insert", func(t *testing.T) {
		var err error
		err = ob.Insert(newUser)
		if err != nil {
			t.Fatalf("Unable to Insert User: %+v", newUser)
		}
		err = ob.Insert(newGroup)
		if err != nil {
			t.Fatalf("Unable to Inser Group: %+v", newGroup)
		}
	})
	if t.Failed() {
		t.FailNow()
	}
	t.Run("Get", func(t *testing.T) {
		gottenU, err := ob.Get(newUser.GetID())
		if err != nil {
			t.Fatalf("Unable to get User: %+v", err)
		}
		if gottenU.GetName() != newUser.GetName() {
			t.Fatalf("Gotten User does not match: %+v", gottenU)
		}
		gottenG, err := ob.Get(newGroup.GetID())
		if err != nil {
			t.Fatalf("Unable to get Group: %+v", err)
		}
		if gottenG.GetName() != newGroup.GetName() {
			t.Fatalf("Gotten Group does not match: %+v", gottenG)
		}
	})
	t.Run("Name", func(t *testing.T) {
		gottenU, err := ob.FindUserName(newUser.GetName())
		if err != nil {
			t.Fatalf("unable to find user: %+v", err)
		}
		if !gottenU.GetID().Equal(newUser.GetID()) {
			t.Fatalf("Gotten User does not match: %s, %s", gottenU.GetID(), newUser.GetID())
		}
		gottenG, err := ob.FindGroupName(newGroup.GetName())
		if err != nil {
			t.Fatalf("unable to find group: %+v", err)
		}
		if !gottenG.GetID().Equal(newGroup.GetID()) {
			t.Fatalf("Gotten Group does not match: %+v", gottenG)
		}
	})
	t.Run("GetGroups", func(t *testing.T) {
		owned, members, err := ob.GetGroups(newUser.GetID())
		if err != nil {
			t.Fatalf("unable to get groups: %+v", err)
		}
		if len(owned) != 1 {
			t.Fatalf("incorrect returned owned groups")
		}
		if len(members) != 0 {
			t.Fatalf("incorrect returned member groups")
		}
	})
	t.Run("Update", func(t *testing.T) {
		newUser.ChangeEmail("test4@test.test")
		err := ob.Update(newUser)
		if err != nil {
			t.Fatalf("Update Failed: %s", err)
		}
	})
	t.Run("Space", func(t *testing.T) {
		total, err := ob.GetTotalSpace(newUser.GetID())
		if err != nil {
			t.Fatalf("Unable get total space: %s", err)
		}
		if total != 50<<20 {
			t.Fatalf("incorrect total, %d", total)
		}
		space, err := ob.GetSpace(newUser.GetID())
		if err != nil {
			t.Fatalf("Unable to get space: %s", err)
		}
		if space != 0 {
			t.Fatalf("incorrect space. %d", space)
		}
	})
}
