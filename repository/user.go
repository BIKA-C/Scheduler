package repository

import (
	"scheduler/account"
	"scheduler/router/errors"
	"scheduler/util"
)

type UserRepository struct {
	db  map[util.UUID]account.User
	acc AccountRepository
}

func NewUserRepo(a AccountRepository) *UserRepository {
	return &UserRepository{
		db:  make(map[util.UUID]account.User, 10),
		acc: a,
	}
}

func (a *UserRepository) Get(UUID util.UUID) account.User {
	if _, ok := a.db[UUID]; !ok {
		return account.User{}
	}
	return a.db[UUID]
}

func (a *UserRepository) Save(u *account.User) error {
	u.UUID = util.NewUUID()
	a.acc.Save(&u.Account)
	a.db[u.UUID] = *u
	return nil
}

func (a *UserRepository) Update(u account.User) error {
	if _, ok := a.db[u.UUID]; !ok {
		return errors.BadRequest("user does not exist")
	}
	a.db[u.UUID] = u
	return nil
}

func (a *UserRepository) Delete(id util.UUID) error {
	if _, ok := a.db[id]; !ok {
		return errors.BadRequest("user does not exist")
	}
	a.acc.Delete(id)
	delete(a.db, id)
	return nil
}
