package memory

import (
	"crypto/rand"
	"encoding/base64"

	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/internal/database/types/errors"
	"git.maxset.io/web/knaxim/pkg/srverror"
)

// Ownerbase is a wrapper for the memory database for owner operations
type Ownerbase struct {
	Database
}

// Reserve is the first step to adding a new Owner, returns OwnerID
// that was reserved, it might be a mutated value of the input
func (ob *Ownerbase) Reserve(id types.OwnerID, name string) (types.OwnerID, error) {
	lock.Lock()
	defer lock.Unlock()
	if id.Type == 'u' {
		if _, ok := ob.Owners.UserName[name]; ok {
			return id, errors.ErrNameTaken
		}
	} else if id.Type == 'g' {
		if _, ok := ob.Owners.GroupName[name]; ok {
			return id, errors.ErrNameTaken
		}
	} else {
		return id, srverror.Basic(500, "Error MO1", "unrecognized id type")
	}

	for idstr := id.String(); true; idstr = id.String() {
		if _, ok := ob.Owners.ID[idstr]; !ok {
			break
		}
		id = id.Mutate()
	}
	ob.Owners.ID[id.String()] = nil
	switch id.Type {
	case 'u':
		ob.Owners.UserName[name] = nil
	case 'g':
		ob.Owners.GroupName[name] = nil
	}
	return id, nil
}

// Insert adds owner to database
func (ob *Ownerbase) Insert(u types.Owner) error {
	lock.Lock()
	defer lock.Unlock()
	idstr := u.GetID().String()
	if expectnil, ok := ob.Owners.ID[idstr]; !ok {
		return errors.ErrIDNotReserved
	} else if expectnil != nil {
		return errors.ErrNameTaken
	}
	switch v := u.(type) {
	case types.UserI:
		expectnil, ok := ob.Owners.UserName[v.GetName()]
		if !ok {
			return errors.ErrIDNotReserved
		}
		if expectnil != nil {
			return errors.ErrNameTaken
		}
		ob.Owners.UserName[v.GetName()] = v
	case types.GroupI:
		expectnil, ok := ob.Owners.GroupName[v.GetName()]
		if !ok {
			return errors.ErrIDNotReserved
		}
		if expectnil != nil {
			return errors.ErrNameTaken
		}
		ob.Owners.GroupName[v.GetName()] = v
	default:
		return srverror.Basic(500, "Error MO2", "Unrecognized Owner Type")
	}
	ob.Owners.ID[idstr] = u
	return nil
}

// Get pulls owner out of database
func (ob *Ownerbase) Get(id types.OwnerID) (types.Owner, error) {
	lock.RLock()
	defer lock.RUnlock()
	return ob.get(id)
}

func (ob *Ownerbase) get(id types.OwnerID) (types.Owner, error) {
	if ob.Owners.ID[id.String()] == nil {
		return nil, errors.ErrNotFound
	}
	return ob.Owners.ID[id.String()].Copy(), nil
}

// FindUserName returns user that has a particular username
func (ob *Ownerbase) FindUserName(name string) (types.UserI, error) {
	lock.RLock()
	defer lock.RUnlock()
	if ob.Owners.UserName[name] == nil {
		return nil, errors.ErrNotFound
	}
	return ob.Owners.UserName[name].Copy().(types.UserI), nil
}

// FindGroupName returns group that has a particular name
func (ob *Ownerbase) FindGroupName(name string) (types.GroupI, error) {
	lock.RLock()
	defer lock.RUnlock()
	if ob.Owners.GroupName[name] == nil {
		return nil, errors.ErrNotFound
	}
	return ob.Owners.GroupName[name].Copy().(types.GroupI), nil
}

// GetGroups returns groups that are owned by owner and groups the
// owner is a member of
func (ob *Ownerbase) GetGroups(id types.OwnerID) (owned []types.GroupI, member []types.GroupI, err error) {
	lock.RLock()
	defer lock.RUnlock()
LOOP:
	for _, grp := range ob.Owners.GroupName {
		if grp.GetOwner().GetID().Equal(id) {
			owned = append(owned, grp.Copy().(types.GroupI))
			continue
		}
		for _, mem := range grp.GetMembers() {
			if mem.GetID().Equal(id) {
				member = append(member, grp.Copy().(types.GroupI))
				continue LOOP
			}
		}
	}
	return
}

// Update owner
func (ob *Ownerbase) Update(o types.Owner) error {
	lock.Lock()
	defer lock.Unlock()
	if ob.Owners.ID[o.GetID().String()] == nil {
		return errors.ErrNotFound
	}
	switch v := o.(type) {
	case types.UserI:
		ob.Owners.UserName[v.GetName()] = v
	case types.GroupI:
		ob.Owners.GroupName[v.GetName()] = v
	default:
		return srverror.Basic(500, "Error MO3", "Unrecognized owner type")
	}
	ob.Owners.ID[o.GetID().String()] = o.Copy()
	return nil
}

// GetSpace returns the total amount of filesize owned by owner
func (ob *Ownerbase) GetSpace(o types.OwnerID) (int64, error) {
	lock.RLock()
	defer lock.RUnlock()
	if ob.Owners.ID[o.String()] == nil {
		return 0, errors.ErrNotFound
	}
	var total int64
	for _, file := range ob.Files {
		if file != nil && file.GetOwner().GetID().Equal(o) {
			total += ob.Stores[file.GetID().StoreID.String()].FileSize
		}
	}
	return total, nil
}

// GetTotalSpace returns the total file space available to owner
func (ob *Ownerbase) GetTotalSpace(o types.OwnerID) (int64, error) {
	lock.RLock()
	defer lock.RUnlock()
	if ob.Owners.ID[o.String()] == nil {
		return 0, errors.ErrNotFound
	}
	if o.Type == 'u' {
		return 50 << 20, nil
	}
	return 0, nil
}

// GetResetKey generates new password reset key
func (ob *Ownerbase) GetResetKey(id types.OwnerID) (key string, err error) {
	newkey := make([]byte, 32)
	_, err = rand.Read(newkey)
	if err != nil {
		return "", srverror.New(err, 500, "Error MO4", "Unable to generate new password reset key")
	}
	str := base64.RawURLEncoding.EncodeToString(newkey)
	ob.Owners.Reset[str] = id
	return str, nil
}

// CheckResetKey returns associated ownerid of reset key
func (ob *Ownerbase) CheckResetKey(keystr string) (id types.OwnerID, err error) {
	id, assigned := ob.Owners.Reset[keystr]
	if !assigned {
		return id, errors.ErrNotFound
	}
	return
}

// DeleteResetKey removes resetkey
func (ob *Ownerbase) DeleteResetKey(id types.OwnerID) error {
	for k, v := range ob.Owners.Reset {
		if v.Equal(id) {
			delete(ob.Owners.Reset, k)
			break
		}
	}
	return nil
}
