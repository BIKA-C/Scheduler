package repository

import (
	"scheduler/account"
	"scheduler/router/errors"
	"scheduler/util"
)

type AccountRepository struct {
	db map[util.UUID]account.Account
}

func NewAccountRepo() *AccountRepository {
	return &AccountRepository{
		db: make(map[util.UUID]account.Account, 10),
	}
}

func (a *AccountRepository) Get(UUID util.UUID) account.Account {
	if _, ok := a.db[UUID]; !ok {
		return account.Account{}
	}
	return a.db[UUID]
}

func (a *AccountRepository) Save(u *account.Account) error {
	a.db[u.UUID] = *u
	return nil
}

func (a *AccountRepository) Update(u account.Account) error {
	if _, ok := a.db[u.UUID]; !ok {
		return errors.BadRequest("user does not exist")
	}
	a.db[u.UUID] = u
	return nil
}

func (a *AccountRepository) Delete(id util.UUID) error {
	if _, ok := a.db[id]; !ok {
		return errors.BadRequest("user does not exist")
	}
	delete(a.db, id)
	return nil
}
