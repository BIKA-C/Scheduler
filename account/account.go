package account

import (
	"bytes"
	"crypto/sha512"
	e "errors"
	"scheduler/router/errors"
	"scheduler/util"
	"time"
)

// Password can't have a length
// greater than 30 characters
type Password string

const minPassword = 8
const maxPassword = 50

type Meta struct {
	ID        int       `json:"-"`
	UpdatedAt time.Time `json:"-"`
	CreatedAt time.Time `json:"-"`
	LastLogin time.Time `json:"-"`
}

var EmptyMeta Meta = Meta{}

type Account struct {
	Meta     Meta      `json:"-"`
	Password Password  `json:"password,omitempty"`
	Email    string    `json:"email"`
	UUID     util.UUID `json:"id"`
}

var EmptyAccount Account = Account{}

func (a *Account) Validate() error {
	err := errors.DefaultAccError
	if passwordErr := a.Password.Validate(); passwordErr != nil {
		err.Set("password", passwordErr.Error())
	}
	if !util.EmailRegex.MatchString(a.Email) {
		err.Set("email", "Not a valid email")
	}

	return err.Build()
}

func (a *Account) HashPassword() {
	if a.Password == "" {
		return
	}
	a.Password = Password(util.B2S(a.Password.Hash()))
}

func (p Password) MarshalJSON() ([]byte, error) {
	return nil, nil
}

var key = []byte("4rfGvbnJuytfcvbnmKoip8gvbAdjfa-P")

// Hash hashes the password.
// if the password is already hashed
// p will be returned
func (p Password) Hash() []byte {
	if p == "" {
		return nil
	}
	if p.hashed() {
		return util.S2B(string(p))
	}
	var hash = sha512.New()
	hash.Write(key)
	hash.Write(util.S2B(string(p)))
	return hash.Sum(nil)
}

func (p Password) hashed() bool {
	return len(p) == sha512.Size
}

func (p Password) String() string {
	return util.B2S(p.Hash())
}

func (p Password) Compare(o Password) bool {
	return bytes.Equal(p.Hash(), o.Hash())
}

func (p Password) Validate() error {
	if len(p) < minPassword {
		return e.New("Password too short")
	} else if len(p) > maxPassword {
		return e.New("Password too long")
	}

	return nil
}
