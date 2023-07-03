package repository

import (
	"scheduler/account"
	"scheduler/router/errors"
	"scheduler/util"
)

type AccountRepository interface {
	Repository[account.Account, util.UUID]
	GetByEmail(email string) account.Account
	VerifyPassword(UUID util.UUID, password account.Password) bool
}

type memoryAccRepo struct {
	db map[util.UUID]account.Account
}

func NewAccountRepo() *memoryAccRepo {
	return &memoryAccRepo{
		db: make(map[util.UUID]account.Account, 10),
	}
}

func (a *memoryAccRepo) Get(UUID util.UUID) account.Account {
	if _, ok := a.db[UUID]; !ok {
		return account.Account{}
	}
	return a.db[UUID]
}

func (a *memoryAccRepo) Save(u *account.Account) error {
	u.HashPassword()
	a.db[u.UUID] = *u
	u.Password = ""
	return nil
}

func (a *memoryAccRepo) Update(u account.Account) error {
	if _, ok := a.db[u.UUID]; !ok {
		return errors.BadRequest("user does not exist")
	}
	a.db[u.UUID] = u
	return nil
}

func (a *memoryAccRepo) Delete(id util.UUID) error {
	if _, ok := a.db[id]; !ok {
		return errors.BadRequest("user does not exist")
	}
	delete(a.db, id)
	return nil
}

func (a *memoryAccRepo) GetByEmail(email string) account.Account {
	for _, v := range a.db {
		if v.Email == email {
			return v
		}
	}
	return account.EmptyAccount
}

func (a *memoryAccRepo) VerifyPassword(id util.UUID, password account.Password) bool {
	return password.Compare(a.db[id].Password.Hash())
}
