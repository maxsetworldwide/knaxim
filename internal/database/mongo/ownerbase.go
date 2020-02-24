package mongo

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math"
	"time"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/pkg/srverror"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type trackOwners struct {
	gotten     map[string]database.Owner
	usernames  map[string]database.UserI
	groupnames map[string]database.GroupI
	groupperms map[string][]database.GroupI
}

func appendUnique(list []database.GroupI, elements ...database.GroupI) []database.GroupI {
	for _, ele := range elements {
		found := false
		for _, l := range list {
			if ele.Equal(l) {
				found = true
				break
			}
		}
		if !found {
			list = append(list, ele)
		}
	}
	return list
}

func newTrackOwners() trackOwners {
	var out trackOwners
	out.gotten = make(map[string]database.Owner)
	out.usernames = make(map[string]database.UserI)
	out.groupnames = make(map[string]database.GroupI)
	out.groupperms = make(map[string][]database.GroupI)
	return out
}

func (to trackOwners) put(o database.Owner) {
	switch v := o.(type) {
	case database.UserI:
		to.usernames[v.GetName()] = v
	case database.GroupI:
		to.groupnames[v.GetName()] = v
		if o := v.GetOwner(); o != nil {
			oid := o.GetID().String()
			to.groupperms[oid] = appendUnique(to.groupperms[oid], v)
			for _, member := range v.GetMembers() {
				mid := member.GetID().String()
				to.groupperms[mid] = appendUnique(to.groupperms[mid], v)
			}
		}
	default:
		panic(srverror.Basic(500, "Database Error O0", "unrecognized Owner type"))
	}
	to.gotten[o.GetID().String()] = o
}

func (to trackOwners) get(id string) database.Owner {
	return to.gotten[id]
}

func (to trackOwners) getUser(name string) database.UserI {
	return to.usernames[name]
}

func (to trackOwners) getGroup(name string) database.GroupI {
	return to.groupnames[name]
}

func (to trackOwners) getGroupByPermission(oname string) []database.GroupI {
	return to.groupperms[oname]
}

type Ownerbase struct {
	Database
}

func mapIDtoCollection(id database.OwnerID, db Database) string {
	switch id.Type {
	case 'u':
		return db.CollNames["user"]
	case 'g':
		return db.CollNames["group"]
	default:
		panic(database.ErrIDUnrecognized)
	}
}

func (ub *Ownerbase) Reserve(id database.OwnerID, name string) (oid database.OwnerID, err error) {
	defer func() {
		if r := recover(); r != nil {
			oid = database.OwnerID{}
			switch v := r.(type) {
			case srverror.Error:
				err = v
			case error:
				err = srverror.New(v, 500, "Database Error U1.0")
			default:
				err = fmt.Errorf("Reserve Panic: %v", v)
			}
		}
	}()
	var out *database.OwnerID
	cname := mapIDtoCollection(id, ub.Database)
	result := ub.client.Database(ub.DBName).Collection(cname).FindOne(
		ub.ctx,
		bson.M{
			"name": name,
		},
	)
	if err := result.Err(); err != nil {
		if err != mongo.ErrNoDocuments {
			return database.OwnerID{}, srverror.New(err, 500, "Database Error U1.2", "unable to confirm name not taken")
		}
	} else {
		return database.OwnerID{}, database.ErrNameTaken
	}
	for out == nil {
		timeout := time.Now().Add(time.Hour * 24)
		result, err := ub.client.Database(ub.DBName).Collection(cname).UpdateOne(ub.ctx,
			bson.M{
				"id":      id,
				"reserve": bson.M{"$lte": time.Now()},
			},
			bson.M{
				"$set": bson.M{"reserve": timeout, "name": name},
			},
		)
		if err != nil {
			return id, srverror.New(err, 500, "Database Error U1", "Unable to update id reserve")
		}
		if result.ModifiedCount > 0 {
			out = &id
		} else {
			result, err = ub.client.Database(ub.DBName).Collection(cname).UpdateOne(
				ub.ctx,
				bson.M{"id": id},
				bson.M{"$setOnInsert": bson.M{
					"id":      id,
					"name":    name,
					"reserve": timeout,
				}},
				options.Update().SetUpsert(true),
			)
			if err != nil {
				return id, srverror.New(err, 500, "Database Error U1.1", "unable to upsert id")
			}
			if result.UpsertedCount > 0 {
				out = &id
			} else {
				id = id.Mutate()
			}
		}
	}
	return *out, nil
}

func (ub *Ownerbase) Insert(u database.Owner) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case srverror.Error:
				err = v
			case error:
				err = srverror.New(v, 500, "Database Error U2.0")
			default:
				err = fmt.Errorf("Insert Panic: %v", v)
			}
		}
	}()
	cname := mapIDtoCollection(u.GetID(), ub.Database)
	result, e := ub.client.Database(ub.DBName).Collection(cname).UpdateOne(
		ub.ctx,
		bson.M{
			"id":      u.GetID(),
			"reserve": bson.M{"$gt": time.Now()},
		},
		bson.M{
			"$unset": bson.M{"reserve": ""},
			"$set":   u,
		},
	)
	if e != nil {
		return srverror.New(e, 500, "Database Error U2", "unable to insert")
	}
	if result.ModifiedCount == 0 {
		return database.ErrIDNotReserved
	}
	return nil
}

func (ub *Ownerbase) Get(id database.OwnerID) (database.Owner, error) {
	if result := ub.get(id.String()); result != nil {
		return result, nil
	}
	switch id.Type {
	case 'p':
		return database.Public, nil
	case 'u':
		result := ub.client.Database(ub.DBName).Collection(ub.CollNames["user"]).FindOne(ub.ctx, bson.M{
			"id": id,
		})
		u := new(database.User)
		if err := result.Decode(u); err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, database.ErrNotFound.Extend("unable to find user", id.String())
			}
			return nil, srverror.New(err, 500, "Database Error U3", "unable to get user")
		}
		ub.put(u)
		return u, nil
	case 'g':
		result := ub.client.Database(ub.DBName).Collection(ub.CollNames["group"]).FindOne(ub.ctx, bson.M{
			"id": id,
		})
		g := new(database.Group)
		if err := result.Decode(g); err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, database.ErrNotFound.Extend("unable to find group", id.String())
			}
			return nil, srverror.New(err, 500, "DatabaseError U3.1", "unable to get group")
		}
		ub.put(g)
		err := g.Populate(ub)
		ub.put(g)
		if err != nil {
			return nil, err
		}
		return g, nil
	default:
		return nil, database.ErrIDUnrecognized
	}
}

func (ub *Ownerbase) FindUserName(name string) (database.UserI, error) {
	if result := ub.getUser(name); result != nil {
		return result, nil
	}
	result := ub.client.Database(ub.DBName).Collection(ub.CollNames["user"]).FindOne(ub.ctx, bson.M{
		"name": name,
	})
	user := new(database.User)
	if err := result.Decode(user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, database.ErrNotFound.Extend("User name", name)
		}
		return nil, srverror.New(err, 500, "Database Error O1", "error finding user name")
	}
	ub.put(user)
	return user, nil
}

func (ub *Ownerbase) FindGroupName(name string) (database.GroupI, error) {
	if result := ub.getGroup(name); result != nil {
		return result, nil
	}
	result := ub.client.Database(ub.DBName).Collection(ub.CollNames["group"]).FindOne(ub.ctx, bson.M{
		"name": name,
	})
	group := new(database.Group)
	if err := result.Decode(group); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, database.ErrNotFound.Extend("Group name", name)
		}
		return nil, srverror.New(err, 500, "Database Error O2", "error finding group name")
	}
	ub.put(group)
	err := group.Populate(ub)
	ub.put(group)
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (ub *Ownerbase) GetGroups(id database.OwnerID) (owned []database.GroupI, member []database.GroupI, err error) {
	grouplist := ub.getGroupByPermission(id.String())
	gids := make([]database.OwnerID, 0, len(grouplist))
	for _, g := range grouplist {
		gids = append(gids, g.GetID())
	}
	cursor, err := ub.client.Database(ub.DBName).Collection(ub.CollNames["group"]).Find(ub.ctx, bson.M{
		"$or": bson.A{
			bson.M{
				"own": id,
			},
			bson.M{
				"perm.%member%": id,
			},
		},
		"id": bson.M{
			"$not": bson.M{"$in": gids},
		},
	})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil, database.ErrNoResults.Extend("No associated groups", id.String())
		}
		return nil, nil, srverror.New(err, 500, "Database Error O3", "unable to find groups")
	}
	var groups []*database.Group
	if err = cursor.All(ub.ctx, &groups); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil, database.ErrNoResults.Extend("No associated groups decoded", id.String())
		}
		return nil, nil, srverror.New(err, 500, "Database Error O4", "unable to decode groups")
	}
	for _, group := range groups {
		ub.put(group)
		group.Populate(ub)
		ub.put(group)
		grouplist = append(grouplist, group)
	}
	for _, group := range grouplist {
		if id.Equal(group.GetOwner().GetID()) {
			owned = append(owned, group)
		}
		for _, mem := range group.GetMembers() {
			if id.Equal(mem.GetID()) {
				member = append(member, group)
				break
			}
		}
	}
	return
}

func (ub *Ownerbase) Update(u database.Owner) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case srverror.Error:
				err = v
			case error:
				err = srverror.New(v, 500, "Database Error O5.0")
			default:
				err = fmt.Errorf("Update Panic: %v", v)
			}
		}
	}()
	cname := mapIDtoCollection(u.GetID(), ub.Database)
	result, err := ub.client.Database(ub.DBName).Collection(cname).UpdateOne(
		ub.ctx,
		bson.M{
			"id": u.GetID(),
		},
		bson.M{
			"$set": u,
		},
	)
	if err != nil {
		return srverror.New(err, 500, "Database Error O5", "error updating owner")
	}
	if result.ModifiedCount == 0 {
		return database.ErrNotFound.Extend(u.GetID().String())
	}
	ub.gotten[u.GetID().String()] = nil
	return nil
}

func (ob *Ownerbase) GetSpace(id database.OwnerID) (int64, error) {
	cursor, err := ob.client.Database(ob.DBName).Collection(ob.CollNames["file"]).Aggregate(
		ob.ctx,
		bson.A{
			bson.M{"$match": bson.M{"own": id}},
			bson.M{"$project": bson.M{"_id": 0, "store": "$id.storeid"}},
			bson.M{"$lookup": bson.M{
				"from":         ob.CollNames["store"],
				"localField":   "store",
				"foreignField": "id",
				"as":           "data",
			}},
			bson.M{"$unwind": "$data"},
			bson.M{"$group": bson.M{
				"_id":  nil,
				"size": bson.M{"$sum": "$data.fsize"},
			}},
		},
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, nil
		}
		return 0, srverror.New(err, 500, "Database Error O6", "unable to send aggregation")
	}
	var result []struct {
		Size int64 `bson:"size"`
	}
	if err := cursor.All(ob.ctx, &result); err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, nil
		}
		return 0, srverror.New(err, 500, "Database Error O6.1", "unable to Decode aggregation")
	}
	if len(result) == 0 {
		return 0, nil
	}
	return result[0].Size, nil
}

func (ob *Ownerbase) GetTotalSpace(id database.OwnerID) (int64, error) {
	own, err := ob.Get(id)
	if err != nil {
		return 0, err
	}
	switch v := own.(type) {
	case database.GroupI:
		return 0, nil
	case database.UserI:
		if v.GetRole("guest") {
			return 0, nil
		} else if v.GetRole("admin") {
			return math.MaxInt64, nil
		} else {
			return 50 << 20, nil
		}
	default:
		return 0, database.ErrNotFound.Extend("unrecognized user")
	}
}

func (ob *Ownerbase) GetResetKey(id database.OwnerID) (key string, err error) {
	newkey := make([]byte, 32)
	_, err = rand.Read(newkey)
	if err != nil {
		return "", srverror.New(err, 500, "Server Error", "Unable to generate new password reset key")
	}
	_, err = ob.client.Database(ob.DBName).Collection(ob.CollNames["reset"]).UpdateOne(ob.ctx, bson.M{
		"user": id,
	}, bson.M{
		"$set": bson.M{
			"key":    newkey,
			"expire": time.Now().Add(time.Hour * 6),
		},
		"$setOnInsert": bson.M{
			"user": id,
		},
	}, options.Update().SetUpsert(true))
	if err != nil {
		return "", srverror.New(err, 500, "Server Error", "unable to insert key")
	}
	return base64.RawURLEncoding.EncodeToString(newkey), nil
}

func (ob *Ownerbase) CheckResetKey(keystr string) (id database.OwnerID, err error) {
	key, err := base64.RawURLEncoding.DecodeString(keystr)
	if err != nil {
		return database.OwnerID{}, srverror.New(err, 400, "Bad Reset", "malformed reset key string")
	}
	result := ob.client.Database(ob.DBName).Collection(ob.CollNames["reset"]).FindOne(ob.ctx, bson.M{
		"key": key,
	})
	if result.Err() != nil {
		return database.OwnerID{}, srverror.New(result.Err(), 404, "Not Found")
	}
	var resetDoc struct {
		User   database.OwnerID `bson:"user"`
		Key    []byte           `bson:"key"`
		Expire time.Time        `bson:"expire"`
	}
	err = result.Decode(&resetDoc)
	if err != nil {
		return database.OwnerID{}, database.ErrNotFound.("no key", err.Error())
	}
	if resetDoc.Expire.Before(time.Now()) {
		return database.OwnerID{}, srverror.Basic(404, "Not Found", "reset key expired")
	}
	return resetDoc.User, nil
}

func (ob *Ownerbase) DeleteResetKey(id database.OwnerID) error {
	_, err := ob.client.Database(ob.DBName).Collection(ob.CollNames["reset"]).DeleteOne(ob.ctx, bson.M{
		"user": id,
	})
	if err != nil {
		return srverror.New(err, 500, "Server Error", "unable to remove reset key")
	}
	return nil
}
