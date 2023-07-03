package controller

import (
	"scheduler/account"
	"scheduler/repository"
	"scheduler/util"
)

type Controller struct {
	repo repositories
}

var accRepo = repository.NewAccountRepo()
var usrRepo = repository.NewUserRepo(accRepo)

func New() *Controller {
	return &Controller{
		repo: repositories{
			account: accRepo,
			user:    usrRepo,
		},
	}
}

type accountRepo repository.AccountRepository

type repositories struct {
	user    repository.Repository[account.User, util.UUID]
	account repository.AccountRepository
	// academy repository.Repository[account.Academy]
}
