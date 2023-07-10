package services

import (
	"scheduler/account"
	"scheduler/repository"
	"scheduler/util"
)

type User struct {
	repo repository.User
}

func NewUserServ(u repository.User) User {
	return User{
		repo: u,
	}
}

func (s User) Update(u *account.User) error {
	if err := validateUserUpdate(u); err != nil {
		return nil
	}
	return s.repo.Update(u)
}

func (s User) Register(u *account.User) error {
	u.UUID = util.NewUUID()
	if err := u.Validate(); err != nil {
		return err
	} else if err := s.repo.Save(u); err != nil {
		return repository.ErrUserAlreadyExist
	} else {
		return nil
	}
}

func (s User) Get(id util.UUID) (account.User, error) {
	return s.repo.Get(id)
}
