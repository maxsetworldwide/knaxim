package types_test

import (
	"testing"

	"git.maxset.io/web/knaxim/internal/database/memory"
	. "git.maxset.io/web/knaxim/internal/database/types"
)

func TestPermission(t *testing.T) {
	owner := NewUser("ownertest", "testtest", "test@test.test")
	viewer := NewUser("viewertest", "testtest", "test@test.test")
	groupViewer := NewGroup("grouptest", viewer)

	perm := new(Permission)
	perm.Own = owner
	perm.SetPerm(viewer, "view", true)
	perm.SetPerm(groupViewer, "view", true)

	if !perm.CheckPerm(viewer, "view") {
		t.Fatal("viewer not matching as viewer")
	}
	if !perm.CheckPerm(groupViewer, "view") {
		t.Fatal("viewer group not matching viewer")
	}
	{
		typs := perm.PermTypes()
		if len(typs) != 1 || typs[0] != "view" {
			t.Fatalf("incorrect PermTypes: %v", typs)
		}
	}
	{
		cpy := perm.CopyPerm(nil)
		if !cpy.GetOwner().Equal(owner) {
			t.Fatal("failed to copy perm")
		}
	}

	var db memory.Database
	db.Init(nil, true)
	{
		ob := db.Owner(nil)
		savedid, err := ob.Reserve(owner.GetID(), owner.GetName())
		if err != nil {
			t.Fatalf("unable to reseve owner id: %s", err)
		}
		if !savedid.Equal(owner.GetID()) {
			t.Fatal("id changed for owner")
		}
		err = ob.Insert(owner)
		if err != nil {
			t.Fatalf("unable to insert owner: %s", err)
		}

		savedid, err = ob.Reserve(viewer.GetID(), viewer.GetName())
		if err != nil {
			t.Fatalf("unable to reseve viewer id: %s", err)
		}
		if !savedid.Equal(viewer.GetID()) {
			t.Fatal("id changed for viewer")
		}
		err = ob.Insert(viewer)
		if err != nil {
			t.Fatalf("unable to insert viewer: %s", err)
		}

		savedid, err = ob.Reserve(groupViewer.GetID(), groupViewer.GetName())
		if err != nil {
			t.Fatalf("unable to reseve groupViewer id: %s", err)
		}
		if !savedid.Equal(groupViewer.GetID()) {
			t.Fatal("id changed for groupViewer")
		}
		err = ob.Insert(groupViewer)
		if err != nil {
			t.Fatalf("unable to insert groupViewer: %s", err)
		}
	}

	t.Logf("owner: %+v", owner)
	t.Logf("Database: %+v", db)

	{
		pjson, err := perm.MarshalJSON()
		if err != nil {
			t.Fatalf("unable to MarshalJSON: %s", err)
		}
		t.Logf("json: %s", string(pjson))
		np := new(Permission)
		err = np.UnmarshalJSON(pjson)
		if err != nil {
			t.Fatalf("unable to UnmarshalJSON: %s", err)
		}
		t.Logf("np: %+v", np)
		err = np.Populate(db.Owner(nil))
		if err != nil {
			t.Fatalf("Populate threw err: %s", err)
		}
		if np.GetOwner() == nil || !np.GetOwner().Equal(owner) {
			t.Fatal("failed to populate permission")
		}
	}

	{
		pbson, err := perm.MarshalBSON()
		if err != nil {
			t.Fatalf("unable to MarshalBSON: %s", err)
		}
		np := new(Permission)
		err = np.UnmarshalBSON(pbson)
		if err != nil {
			t.Fatalf("unable to UnmarshalBSON: %s, %#v", err, pbson)
		}
		err = np.Populate(db.Owner(nil))
		if err != nil {
			t.Fatalf("Populate threw err: %s", err)
		}
		if np.GetOwner() == nil || !np.GetOwner().Equal(owner) {
			t.Fatal("failed to populate permission")
		}
	}
}
