package mongo

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math"
	"time"

	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/internal/database/types/errors"
	"git.maxset.io/web/knaxim/pkg/srverror"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type trackOwners struct {
	gotten     map[string]types.Owner
	usernames  map[string]types.UserI
	groupnames map[string]types.GroupI
	groupperms map[string][]types.GroupI
}

func appendUnique(list []types.GroupI, elements ...types.GroupI) []types.GroupI {
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
	out.gotten = make(map[string]types.Owner)
	out.usernames = make(map[string]types.UserI)
	out.groupnames = make(map[string]types.GroupI)
	out.groupperms = make(map[string][]types.GroupI)
	return out
}

func (to trackOwners) put(o types.Owner) {
	switch v := o.(type) {
	case types.UserI:
		to.usernames[v.GetName()] = v
	case types.GroupI:
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

func (to trackOwners) get(id string) types.Owner {
	return to.gotten[id]
}

func (to trackOwners) getUser(name string) types.UserI {
	return to.usernames[name]
}

func (to trackOwners) getGroup(name string) types.GroupI {
	return to.groupnames[name]
}

func (to trackOwners) getGroupByPermission(oname string) []types.GroupI {
	return to.groupperms[oname]
}

// Ownerbase is a connection to the databae with owner operations
type Ownerbase struct {
	Database
}

func mapIDtoCollection(id types.OwnerID, db Database) string {
	switch id.Type {
	case 'u':
		return db.CollNames["user"]
	case 'g':
		return db.CollNames["group"]
	default:
		panic(errors.ErrIDUnrecognized)
	}
}

// Reserve an owner id, will mutate if owner id not available, returns reserved owner id
func (ob *Ownerbase) Reserve(id types.OwnerID, name string) (oid types.OwnerID, err error) {
	defer func() {
		if r := recover(); r != nil {
			oid = types.OwnerID{}
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
	var out *types.OwnerID
	for _, cname := range []string{ob.CollNames["user"], ob.CollNames["group"]} {
		result := ob.client.Database(ob.DBName).Collection(cname).FindOne(
			ob.ctx,
			bson.M{
				"name": name,
			},
		)
		if err := result.Err(); err != nil {
			if err != mongo.ErrNoDocuments {
				return types.OwnerID{}, srverror.New(err, 500, "Database Error U1.2", "unable to confirm name not taken")
			}
		} else {
			return types.OwnerID{}, errors.ErrNameTaken
		}
	}
	cname := mapIDtoCollection(id, ob.Database)
	for out == nil {
		timeout := time.Now().Add(time.Hour * 24)
		result, err := ob.client.Database(ob.DBName).Collection(cname).UpdateOne(ob.ctx,
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
			result, err = ob.client.Database(ob.DBName).Collection(cname).UpdateOne(
				ob.ctx,
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

// Insert adds owner to the database, owner id must first be reserved, see Reserve
func (ob *Ownerbase) Insert(u types.Owner) (err error) {
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
	cname := mapIDtoCollection(u.GetID(), ob.Database)
	result, e := ob.client.Database(ob.DBName).Collection(cname).UpdateOne(
		ob.ctx,
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
		return errors.ErrIDNotReserved
	}
	return nil
}

// Get returns owner based on id
func (ob *Ownerbase) Get(id types.OwnerID) (types.Owner, error) {
	if result := ob.get(id.String()); result != nil {
		return result, nil
	}
	switch id.Type {
	case 'p':
		return types.Public, nil
	case 'u':
		result := ob.client.Database(ob.DBName).Collection(ob.CollNames["user"]).FindOne(ob.ctx, bson.M{
			"id": id,
		})
		u := new(types.User)
		if err := result.Decode(u); err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, errors.ErrNotFound.Extend("unable to find user", id.String())
			}
			return nil, srverror.New(err, 500, "Database Error U3", "unable to get user")
		}
		ob.put(u)
		return u, nil
	case 'g':
		result := ob.client.Database(ob.DBName).Collection(ob.CollNames["group"]).FindOne(ob.ctx, bson.M{
			"id": id,
		})
		g := new(types.Group)
		if err := result.Decode(g); err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, errors.ErrNotFound.Extend("unable to find group", id.String())
			}
			return nil, srverror.New(err, 500, "DatabaseError U3.1", "unable to get group")
		}
		ob.put(g)
		err := g.Populate(ob)
		ob.put(g)
		if err != nil {
			return nil, err
		}
		return g, nil
	default:
		return nil, errors.ErrIDUnrecognized
	}
}

// FindUserName returns user based on username
func (ob *Ownerbase) FindUserName(name string) (types.UserI, error) {
	if result := ob.getUser(name); result != nil {
		return result, nil
	}
	result := ob.client.Database(ob.DBName).Collection(ob.CollNames["user"]).FindOne(ob.ctx, bson.M{
		"name": name,
	})
	user := new(types.User)
	if err := result.Decode(user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound.Extend("User name", name)
		}
		return nil, srverror.New(err, 500, "Database Error O1", "error finding user name")
	}
	ob.put(user)
	return user, nil
}

// FindGroupName finds group based on name
func (ob *Ownerbase) FindGroupName(name string) (types.GroupI, error) {
	if result := ob.getGroup(name); result != nil {
		return result, nil
	}
	result := ob.client.Database(ob.DBName).Collection(ob.CollNames["group"]).FindOne(ob.ctx, bson.M{
		"name": name,
	})
	group := new(types.Group)
	if err := result.Decode(group); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound.Extend("Group name", name)
		}
		return nil, srverror.New(err, 500, "Database Error O2", "error finding group name")
	}
	ob.put(group)
	err := group.Populate(ob)
	ob.put(group)
	if err != nil {
		return nil, err
	}
	return group, nil
}

// GetGroups returns owned groups and groups an owner is a member of. based on the id of the owner
func (ob *Ownerbase) GetGroups(id types.OwnerID) (owned []types.GroupI, member []types.GroupI, err error) {
	grouplist := ob.getGroupByPermission(id.String())
	gids := make([]types.OwnerID, 0, len(grouplist))
	for _, g := range grouplist {
		gids = append(gids, g.GetID())
	}
	cursor, err := ob.client.Database(ob.DBName).Collection(ob.CollNames["group"]).Find(ob.ctx, bson.M{
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
			return nil, nil, errors.ErrNoResults.Extend("No associated groups", id.String())
		}
		return nil, nil, srverror.New(err, 500, "Database Error O3", "unable to find groups")
	}
	var groups []*types.Group
	if err = cursor.All(ob.ctx, &groups); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil, errors.ErrNoResults.Extend("No associated groups decoded", id.String())
		}
		return nil, nil, srverror.New(err, 500, "Database Error O4", "unable to decode groups")
	}
	for _, group := range groups {
		ob.put(group)
		group.Populate(ob)
		ob.put(group)
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

// Update owner
func (ob *Ownerbase) Update(u types.Owner) (err error) {
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
	cname := mapIDtoCollection(u.GetID(), ob.Database)
	result, err := ob.client.Database(ob.DBName).Collection(cname).UpdateOne(
		ob.ctx,
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
	if result.MatchedCount == 0 {
		return errors.ErrNotFound.Extend("no owner to update", u.GetID().String())
	}
	ob.gotten[u.GetID().String()] = nil
	return nil
}

// GetSpace returns amount of used space an owner has
func (ob *Ownerbase) GetSpace(id types.OwnerID) (int64, error) {
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

// GetTotalSpace returns the amount of total space available to an owner
func (ob *Ownerbase) GetTotalSpace(id types.OwnerID) (int64, error) {
	own, err := ob.Get(id)
	if err != nil {
		return 0, err
	}
	switch v := own.(type) {
	case types.GroupI:
		return 0, nil
	case types.UserI:
		customspace := v.GetTotalSpace()
		if customspace == 0 {
			if v.GetRole("guest") {
				return 0, nil
			} else if v.GetRole("admin") {
				return math.MaxInt64, nil
			} else {
				return 50 << 20, nil
			}
		}
		return customspace, nil
	default:
		return 0, errors.ErrNotFound.Extend("unrecognized user")
	}
}

// GetResetKey generates a password reset key
func (ob *Ownerbase) GetResetKey(id types.OwnerID) (key string, err error) {
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

// CheckResetKey looks up a reset key and which owner is associated with that key
func (ob *Ownerbase) CheckResetKey(keystr string) (id types.OwnerID, err error) {
	key, err := base64.RawURLEncoding.DecodeString(keystr)
	if err != nil {
		return types.OwnerID{}, srverror.New(err, 400, "Bad Reset", "malformed reset key string")
	}
	result := ob.client.Database(ob.DBName).Collection(ob.CollNames["reset"]).FindOne(ob.ctx, bson.M{
		"key": key,
	})
	if result.Err() != nil {
		return types.OwnerID{}, srverror.New(result.Err(), 404, "Not Found")
	}
	var resetDoc struct {
		User   types.OwnerID `bson:"user"`
		Key    []byte        `bson:"key"`
		Expire time.Time     `bson:"expire"`
	}
	err = result.Decode(&resetDoc)
	if err != nil {
		return types.OwnerID{}, errors.ErrNotFound.Extend("no key", err.Error())
	}
	if resetDoc.Expire.Before(time.Now()) {
		return types.OwnerID{}, srverror.Basic(404, "Not Found", "reset key expired")
	}
	return resetDoc.User, nil
}

// DeleteResetKey deletes resetkey pairing to owner id
func (ob *Ownerbase) DeleteResetKey(id types.OwnerID) error {
	_, err := ob.client.Database(ob.DBName).Collection(ob.CollNames["reset"]).DeleteOne(ob.ctx, bson.M{
		"user": id,
	})
	if err != nil {
		return srverror.New(err, 500, "Server Error", "unable to remove reset key")
	}
	return nil
}
