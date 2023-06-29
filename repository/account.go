package repository

import (
	"scheduler/account"
)

type AccountRepository struct {
	db map[string]account.Account
}

func (a *AccountRepository) Save(acc account.Account) error {
	a.db[acc.UUID] = acc
	return nil
}

