package controller

import (
	"scheduler/account"
	"scheduler/repository"
	"scheduler/util"
)

type Controller struct {
	repo repositories
}

func New() *Controller {
	return &Controller{
		repo: repositories{
			account: repository.NewAccountRepo(),
			user:    repository.NewUserRepo(),
		},
	}
}

type repositories struct {
	user    repository.Repository[account.User, util.UUID]
	account repository.Repository[account.Account, util.UUID]
	// academy repository.Repository[account.Academy]
}
