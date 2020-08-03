// Copyright August 2020 Maxset Worldwide Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"git.maxset.io/web/knaxim/pkg/srverror"
)

// UserI interface tyoe for representing users
type UserI interface {
	Owner
	GetLock() UserCredentialI
	SetLock(UserCredentialI)
	NewCookies(time.Time, time.Time) []*http.Cookie
	CheckCookie(*http.Request) bool
	RefreshCookie(time.Time)
	GetCookieTimeouts() (time.Time, time.Time)
	GetRole(string) bool
	GetRoles() []string
	SetRole(string, bool)
	GetEmail() string
	ChangeEmail(string)
	GetTotalSpace() int64
}

// UserCredentialI interface representing types that can be used to validate user
type UserCredentialI interface {
	Valid(c map[string]interface{}) bool
}

// User basic user object
type User struct {
	ID               OwnerID         `json:"uid" bson:"id"`
	Name             string          `json:"name" bson:"name"`
	Pass             UserCredential  `json:"pass" bson:"pass"`
	Email            string          `json:"email" bson:"email"`
	CookieSig        []byte          `json:"cs" bson:"cs"`
	CookieInactivity time.Time       `json:"ci" bson:"ci"`
	CookieTimeout    time.Time       `json:"ct" bson:"ct"`
	Roles            map[string]bool `json:"roles,omitempty" bson:"roles,omitempty"`
	Space            int64           `json:"space,omitempty" bson:"space,omitempty"`
	Max              int64           `json:"max,omitempty" bson:"max,omitempty"`
}

// NewUser with username, email and password
func NewUser(name, password, email string) *User {
	n := new(User)

	n.ID.Type = 'u'
	n.ID.UserDefined = name2userdefined(name)

	n.Name = name
	n.Pass = NewUserCredential(password)
	n.Email = email

	n.CookieSig = make([]byte, cookiesize)
	n.Roles = make(map[string]bool)
	return n
}

// MaxFiles returns the maximum number of files that an owner can own
func (u *User) MaxFiles() int64 {
	return u.Max
}

// GetName implements UserI
func (u *User) GetName() string {
	return u.Name
}

// GetEmail implements UserI
func (u *User) GetEmail() string {
	return u.Email
}

// ChangeEmail implements UserI
func (u *User) ChangeEmail(path string) {
	u.Email = path
}

// GetID implements UserI
func (u *User) GetID() OwnerID {
	return u.ID
}

// GetLock implements UserI
func (u *User) GetLock() UserCredentialI {
	return u.Pass
}

// SetLock implements UserI
func (u *User) SetLock(c UserCredentialI) {
	if cred, ok := c.(UserCredential); !ok {
		panic("Invalid UserCredential, Not Sha256 User Credential")
	} else {
		u.Pass = cred
	}
}

var cookiekeys = []string{"KnaximSessionSig", "KnaximSessionUID"}
var cookiesize = 7

// NewCookies creates new cookies for a new session
func (u *User) NewCookies(inactivity time.Time, timeout time.Time) []*http.Cookie {
	if _, err := rand.Read(u.CookieSig); err != nil {
		panic(err)
	}
	u.CookieTimeout = timeout
	u.CookieInactivity = inactivity
	result := make([]*http.Cookie, 0, len(cookiekeys))
	ma := int(timeout.Sub(time.Now()) / time.Second)
	{
		n := new(http.Cookie)
		n.Name = cookiekeys[0]
		n.Value = hex.EncodeToString(u.CookieSig)
		n.Path = "/"
		n.Expires = timeout
		n.MaxAge = ma
		n.HttpOnly = true
		result = append(result, n)
	}
	{
		n := new(http.Cookie)
		n.Name = cookiekeys[1]
		n.Value = u.ID.String()
		n.Path = "/"
		n.Expires = timeout
		n.MaxAge = ma
		n.HttpOnly = true
		result = append(result, n)
	}
	return result
}

// CheckCookie returns true if request has cookied session cookies
func (u *User) CheckCookie(r *http.Request) bool {
	if u.CookieTimeout.Before(time.Now()) || u.CookieInactivity.Before(time.Now()) {
		//log.Printf("Cookie Timed Out: %s %s\n", u.CookieInactivity.String(), u.CookieTimeout.String())
		return false
	}
	cookie, err := r.Cookie(cookiekeys[0])
	if err != nil {
		//log.Printf("unable to get cookie sig: %v", err)
		return false
	}
	var sig []byte
	sig, err = hex.DecodeString(cookie.Value)
	if err != nil {
		//log.Printf("unable to decode cookie value: %s", err)
		return false
	}
	if len(sig) != len(u.CookieSig) {
		//log.Printf("miss match sig length")
		return false
	}
	for i, s := range u.CookieSig {
		if sig[i] != s {
			//log.Printf("cookie sig no match")
			return false
		}
	}
	return true
}

// RefreshCookie sets inactivity timeout
func (u *User) RefreshCookie(inactive time.Time) {
	u.CookieInactivity = inactive
}

// GetCookieTimeouts returns current session timeouts
func (u *User) GetCookieTimeouts() (time.Time, time.Time) {
	return u.CookieInactivity, u.CookieTimeout
}

// GetCookieUID gets OwnerID of current session
func GetCookieUID(r *http.Request) (OwnerID, error) {
	cookie, err := r.Cookie(cookiekeys[1])
	if err != nil {
		return OwnerID{}, srverror.New(err, 401, "Unable to identify user")
	}
	return DecodeOwnerIDString(cookie.Value)
}

// Match returns true if owner is equal to user
func (u *User) Match(o Owner) bool {
	return u.Equal(o)
}

// Equal if id is Equal
func (u *User) Equal(o Owner) bool {
	switch v := o.(type) {
	case *User:
		return u.ID.Equal(v.GetID())
	default:
		return false
	}
}

// Copy builds a new instance of the User
func (u *User) Copy() Owner {
	if u == nil {
		return nil
	}
	nu := new(User)
	*nu = *u
	return nu
}

// GetRole is true if user has that role
func (u *User) GetRole(k string) bool {
	return u.Roles[k]
}

// GetRoles returns all roles that user has
func (u *User) GetRoles() []string {
	var out []string
	for k, v := range u.Roles {
		if v {
			out = append(out, k)
		}
	}
	return out
}

// SetRole assigns user to role if v is true, removes role if v is false
func (u *User) SetRole(k string, v bool) {
	if u.Roles == nil {
		u.Roles = make(map[string]bool)
	}
	u.Roles[k] = v
}

// GetTotalSpace returns current Total Space value of user
func (u *User) GetTotalSpace() int64 {
	return u.Space
}

// UserCredential is a salt and hashed of a password
type UserCredential struct {
	Salt []byte `json:"a" bson:"a"`
	Hash []byte `json:"b" bson:"b"`
}

var staticPassSalt = []byte("6`9McFWZ7]{HLR`[D7'")

func hashpass(pass, salt []byte) []byte {
	h := sha256.New()
	h.Write(staticPassSalt)
	h.Write(pass)
	h.Write(salt)
	return h.Sum(nil)
}

// NewUserCredential build UserCredential from password
func NewUserCredential(pass string) UserCredential {
	var n UserCredential

	n.Salt = make([]byte, 32)
	rand.Read(n.Salt)

	n.Hash = hashpass([]byte(pass), n.Salt)

	return n
}

// Valid returns true if credential has the field "pass" has the correct password.
func (up UserCredential) Valid(credential map[string]interface{}) bool {
	temp := credential["pass"]
	var pass []byte
	switch v := temp.(type) {
	case string:
		pass = []byte(v)
	case []byte:
		pass = v
	default:
		return false
	}
	h := hashpass(pass, up.Salt)
	if len(h) != len(up.Hash) {
		return false
	}
	for i := range h {
		if h[i] != up.Hash[i] {
			return false
		}
	}
	return true
}
