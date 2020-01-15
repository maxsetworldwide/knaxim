package database

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"git.maxset.io/web/knaxim/pkg/srverror"
)

type UserI interface {
	Owner
	GetName() string
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
}

type UserCredentialI interface {
	Valid(c map[string]interface{}) bool
}

type User struct {
	ID               OwnerID         `json:"uid" bson:"id"`
	Name             string          `json:"name" bson:"name"`
	Pass             UserCredential  `json:"pass" bson:"pass"`
	Email            string          `json:"email" bson:"email"`
	CookieSig        []byte          `json:"cs" bson:"cs"`
	CookieInactivity time.Time       `json:"ci" bson:"ci"`
	CookieTimeout    time.Time       `json:"ct" bson:"ct"`
	Roles            map[string]bool `json:"roles,omitempty" bson:"roles,omitempty"`
}

func NewUser(name, password, email string) *User {
	n := new(User)

	n.ID.Type = 'u'
	n.ID.UserDefined = name2userdefined(name)
	n.ID.Stamp = newstamp()

	n.Name = name
	n.Pass = NewUserCredential(password)
	n.Email = email

	n.CookieSig = make([]byte, cookiesize)
	n.Roles = make(map[string]bool)
	return n
}

func (u *User) GetName() string {
	return u.Name
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) ChangeEmail(path string) {
	u.Email = path
}

func (u *User) GetID() OwnerID {
	return u.ID
}

func (u *User) GetLock() UserCredentialI {
	return u.Pass
}

func (u *User) SetLock(c UserCredentialI) {
	if cred, ok := c.(UserCredential); !ok {
		panic("Invalid UserCredential, Not Sha256 User Credential")
	} else {
		u.Pass = cred
	}
}

var cookiekeys = []string{"KnaximSessionSig", "KnaximSessionUID"}
var cookiesize = 7

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

func (u *User) RefreshCookie(inactive time.Time) {
	u.CookieInactivity = inactive
}

func (u *User) GetCookieTimeouts() (time.Time, time.Time) {
	return u.CookieInactivity, u.CookieTimeout
}

func GetCookieUid(r *http.Request) (OwnerID, error) {
	cookie, err := r.Cookie(cookiekeys[1])
	if err != nil {
		return OwnerID{}, srverror.New(err, 401, "Unable to identity user")
	}
	return DecodeObjectIDString(cookie.Value)
}

func (u *User) Match(o Owner) bool {
	return u.Equal(o)
}

func (u *User) Equal(o Owner) bool {
	switch v := o.(type) {
	case *User:
		return u.ID.Equal(v.GetID())
	default:
		return false
	}
}

func (u *User) GetRole(k string) bool {
	return u.Roles[k]
}

func (u *User) GetRoles() []string {
	var out []string
	for k, v := range u.Roles {
		if v {
			out = append(out, k)
		}
	}
	return out
}

func (u *User) SetRole(k string, v bool) {
	u.Roles[k] = v
}

type UserCredential struct {
	Salt []byte `json:"a" bson:"a"`
	Hash []byte `json:"b" bson:"b"`
}

var staticPassSalt = []byte("6`9McFWZ7]{HLR`[D7'")

func hash(pass, salt []byte) []byte {
	h := sha256.New()
	h.Write(staticPassSalt)
	h.Write(pass)
	h.Write(salt)
	return h.Sum(nil)
}

func NewUserCredential(pass string) UserCredential {
	var n UserCredential

	n.Salt = make([]byte, 32)
	rand.Read(n.Salt)

	n.Hash = hash([]byte(pass), n.Salt)

	return n
}

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
	h := hash(pass, up.Salt)
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

// type UserToken struct {
// 	UID     string    `json:"uid"`
// 	Timeout time.Time `json:"timeout"`
// 	Sig     string    `json:"sig"`
// }
//
// func (token *UserToken) GetUID() string {
// 	return token.UID
// }
//
// func (token *UserToken) GetTimeout() time.Time {
// 	return token.Timeout
// }
//
// func tokensig(t *UserToken, key []byte) string {
// 	b := new(bytes.Buffer)
// 	b.WriteString(t.UID)
// 	b.WriteString(t.Timeout.Format(time.RFC3339Nano))
// 	mssg := b.Bytes()
// 	h := sha256.New()
// 	if _, err := h.Write(mssg); err != nil {
// 		panic(err)
// 	}
// 	if _, err := h.Write(key); err != nil {
// 		panic(err)
// 	}
// 	sig := h.Sum(nil)
// 	sb := new(strings.Builder)
// 	encoder := base64.NewEncoder(base64.URLEncoding, sb)
// 	if _, err := encoder.Write(sig); err != nil {
// 		panic(err)
// 	}
// 	if err := encoder.Close(); err != nil {
// 		panic(err)
// 	}
// 	return sb.String()
// }
//
// func (token *UserToken) Sign(key []byte) error {
// 	token.Sig = tokensig(token, key)
// 	//log.Printf("generated token signature: %s", token.Sig)
// 	//log.Printf("token is: %v", token)
// 	return nil
// }
//
// func (token *UserToken) Check(key []byte) (bool, error) {
// 	expect := tokensig(token, key)
// 	// if expect != token.Sig {
// 	//   log.Printf("expected: %s\nDelivered: %s", expect, token.Sig);
// 	// }
// 	return token.Sig == expect, nil
// }
