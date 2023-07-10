package services

import (
	"scheduler/account"
	"scheduler/repository"
	"scheduler/util"
)

type Account struct {
	repo repository.Account
}

func NewAccountServ(u repository.Account) Account {
	return Account{
		repo: u,
	}
}

func (s Account) Update(u *account.Account, password account.Password) error {
	if err := s.validateAccount(u, password); err != nil {
		return err
	}
	return s.repo.Update(u)
}

func (s Account) Get(id util.UUID) (account.Account, error) {
	return s.repo.Get(id)
}

func (s Account) VerifyPassword(id util.UUID, pass account.Password) bool {
	return s.repo.VerifyPassword(id, pass)
}
