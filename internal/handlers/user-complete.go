package handlers

import (
	"net/http"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/tag"
	"git.maxset.io/web/knaxim/internal/util"

	"git.maxset.io/web/knaxim/pkg/srvjson"
)

type userProfile struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Groups struct {
		Own    []string `json:"own"`
		Member []string `json:"member"`
	} `json:"groups"`
	Dirs  []string `json:"folders"`
	Files struct {
		Own  []string `json:"own"`
		View []string `json:"view"`
	} `json:"files"`
	Roles []string `json:"roles"`
	Data  struct {
		Current int64 `json:"current"`
		Total   int64 `json:"total"`
	} `json:"data"`
}

type groupProfile struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Owner   string   `json:"owner"`
	IsOwned bool     `json:"isOwned"`
	Members []string `json:"members"`
	Groups  struct {
		Own    []string `json:"own"`
		Member []string `json:"member"`
	} `json:"groups"`
	Dirs  []string `json:"folders"`
	Files struct {
		Own  []string `json:"own"`
		View []string `json:"view"`
	} `json:"files"`
}

func buildGP(g database.GroupI, isOwned bool, gown, gm, d, fo, fv []string) groupProfile {
	var out groupProfile
	out.ID = g.GetID().String()
	out.Name = g.GetName()
	out.Owner = g.GetOwner().GetID().String()
	out.IsOwned = isOwned
	for _, member := range g.GetMembers() {
		out.Members = append(out.Members, member.GetID().String())
	}
	out.Groups.Own = gown
	out.Groups.Member = gm
	out.Dirs = d
	out.Files.Own = fo
	out.Files.View = fv
	return out
}

type fileProfile struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Type    string            `json:"type"`
	Owner   string            `json:"owner"`
	IsOwned bool              `json:"isOwned"`
	Date    database.FileTime `json:"date"`
	Size    int64             `json:"size"`
	URL     string            `json:"url,omitempty"`
	Viewers []string          `json:"viewers"`
}

func buildFP(r database.FileI, isOwned bool, size int64) fileProfile {
	var out fileProfile
	out.ID = r.GetID().String()
	out.Name = r.GetName()
	out.Owner = r.GetOwner().GetID().String()
	out.IsOwned = isOwned
	for _, viewer := range r.GetPerm("view") {
		out.Viewers = append(out.Viewers, viewer.GetID().String())
	}
	if _, ok := r.(*database.WebFile); ok {
		out.Type = "webpage"
	} else {
		out.Type = "file"
	}
	out.Date = r.GetDate()
	out.Size = size
	if wf, ok := r.(*database.WebFile); ok {
		out.URL = wf.URL
	}
	return out
}

type publicProfile []string

// CompletePackage is a full set of metadata relating to a particular user
type CompletePackage struct {
	User    userProfile             `json:"user"`
	Public  publicProfile           `json:"public"`
	Groups  map[string]groupProfile `json:"groups"`
	Records map[string]fileProfile  `json:"files"`
}

func (cp *CompletePackage) addGroup(g database.GroupI, currentUser database.UserI, ownerbase database.Ownerbase, filebase database.Filebase, tagbase database.Tagbase) error {
	if _, ok := cp.Groups[g.GetID().String()]; !ok {
		var gown, gm, d, fo, fv []string
		if owned, member, err := ownerbase.GetGroups(g.GetID()); err == nil {
			for _, ele := range owned {
				cp.addGroup(ele, currentUser, ownerbase, filebase, tagbase)
				gown = append(gown, ele.GetID().String())
			}
			for _, ele := range member {
				cp.addGroup(ele, currentUser, ownerbase, filebase, tagbase)
				gm = append(gm, ele.GetID().String())
			}
		} else {
			return err
		}
		if tags, err := tagbase.SearchData(tag.USER, tag.Data{tag.USER: map[string]string{g.GetID().String(): dirflag}}); err == nil {
			for _, t := range tags {
				d = append(d, t.Word)
			}
		} else {
			return err
		}
		if owned, err := filebase.GetOwned(g.GetID()); err == nil {
			for _, o := range owned {
				cp.addRecord(currentUser, o, filebase)
				fo = append(fo, o.GetID().String())
			}
		} else {
			return err
		}
		if viewable, err := filebase.GetPermKey(g.GetID(), "view"); err == nil {
			for _, v := range viewable {
				cp.addRecord(currentUser, v, filebase)
				fv = append(fv, v.GetID().String())
			}
		}
		cp.Groups[g.GetID().String()] = buildGP(g, g.GetOwner().Match(currentUser), gown, gm, d, fo, fv)
	}
	return nil
}

func (cp *CompletePackage) addRecord(u database.UserI, r database.FileI, db database.Database) error {
	if _, ok := cp.Records[r.GetID().String()]; !ok {
		sb := db.Store(nil)
		fs, err := sb.Get(r.GetID().StoreID)
		if err != nil {
			return err
		}
		cp.Records[r.GetID().String()] = buildFP(r, r.GetOwner().Match(u), fs.FileSize)
	}
	return nil
}

func completeUserInfo(out http.ResponseWriter, r *http.Request) {
	w := out.(*srvjson.ResponseWriter)
	user := r.Context().Value(USER).(database.UserI)

	var info CompletePackage
	info.Groups = make(map[string]groupProfile)
	info.Records = make(map[string]fileProfile)
	filebase := r.Context().Value(database.FILE).(database.Filebase)
	ownerbase := r.Context().Value(database.OWNER).(database.Ownerbase)

	var err error
	info.User.ID = user.GetID().String()
	info.User.Name = user.GetName()
	info.User.Roles = user.GetRoles()
	info.User.Data.Total, err = ownerbase.GetTotalSpace(user.GetID())
	if err != nil {
		util.VerboseRequest(r, "error getting total space.")
		panic(err)
	}
	if info.User.Data.Current, err = ownerbase.GetSpace(user.GetID()); err != nil {
		util.VerboseRequest(r, "error getting current space.")
		panic(err)
	}
	if owned, members, err := ownerbase.GetGroups(user.GetID()); err == nil {
		for _, o := range owned {
			if err = info.addGroup(o, user, ownerbase, filebase, r.Context().Value(database.TAG).(database.Tagbase)); err != nil {
				util.VerboseRequest(r, "error adding group")
				panic(err)
			}
			info.User.Groups.Own = append(info.User.Groups.Own, o.GetID().String())
		}
		for _, m := range members {
			if err = info.addGroup(m, user, ownerbase, filebase, r.Context().Value(database.TAG).(database.Tagbase)); err != nil {
				util.VerboseRequest(r, "error adding member group")
				panic(err)
			}
			info.User.Groups.Member = append(info.User.Groups.Member, m.GetID().String())
		}
	} else {
		util.VerboseRequest(r, "error getting groups")
		panic(err)
	}
	if tags, err := r.Context().Value(database.TAG).(database.Tagbase).SearchData(tag.USER, tag.Data{tag.USER: map[string]string{user.GetID().String(): dirflag}}); err == nil {
		for _, t := range tags {
			info.User.Dirs = append(info.User.Dirs, t.Word)
		}
	} else {
		util.VerboseRequest(r, "error searching tag data")
		panic(err)
	}
	if owned, err := filebase.GetOwned(user.GetID()); err == nil {
		for _, o := range owned {
			info.addRecord(user, o, filebase)
			info.User.Files.Own = append(info.User.Files.Own, o.GetID().String())
		}
	} else {
		util.VerboseRequest(r, "unable to find owned files")
		panic(err)
	}
	if viewable, err := filebase.GetPermKey(user.GetID(), "view"); err == nil {
		for _, v := range viewable {
			info.addRecord(user, v, filebase)
			info.User.Files.View = append(info.User.Files.View, v.GetID().String())
		}
	} else {
		util.VerboseRequest(r, "unable to find Perm Key")
		panic(err)
	}
	if public, err := filebase.GetPermKey(database.Public.GetID(), "view"); err == nil {
		for _, p := range public {
			info.addRecord(user, p, filebase)
			info.Public = append(info.Public, p.GetID().String())
		}
	}

	w.Set("user", info.User)
	w.Set("public", info.Public)
	w.Set("groups", info.Groups)
	w.Set("files", info.Records)
}
