package account

import (
	"bytes"
	"scheduler/internal"
	"scheduler/router/errors"
	"scheduler/util"
	"time"

	"crypto/sha512"
)

// Password can't have a length
// greater than 30 characters
type Password string

type Meta struct {
	UpdatedAt time.Time `json:"-"`
	CreatedAt time.Time `json:"-"`
	LastLogin time.Time `json:"-"`
}

type Account struct {
	Meta     Meta      `json:"-"`
	Password Password  `json:"password,omitempty"`
	Email    string    `json:"email"`
	UUID     util.UUID `json:"id"`
}

var EmptyAccount Account = Account{}

func (a *Account) Validate() error {
	err := errors.DefaultAccError
	if len(a.Password) < 8 {
		err.Set("password", "Password too short")
	}

	if !email.MatchString(a.Email) {
		err.Set("email", "Not a valid email")
	}

	return err.Build()
}

func (a *Account) HashPassword() {
	a.Password = Password(util.B2S(a.Password.Hash()))
}

func (a *Account) Commit(s internal.Saver[Account]) error {
	return s.Save(a)
}

func (p Password) MarshalJSON() ([]byte, error) {
	return nil, nil
}

var key = []byte(util.RandomString(16))

// Hash hashes the password.
// if the password is already hashed
// p will be returned
func (p Password) Hash() []byte {
	if len(p) == sha512.Size {
		return util.S2B(string(p))
	}
	var hash = sha512.New()
	hash.Write(key)
	hash.Write(util.S2B(string(p)))
	return hash.Sum(nil)
}

func (p Password) Compare(h []byte) bool {
	return bytes.Equal(h, p.Hash())
}
