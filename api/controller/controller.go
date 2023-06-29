package controller

import (
	"scheduler/account"
	"scheduler/repository"
)

type Controller struct {
	repo repositories
}

func New() *Controller {
	return &Controller{
		repo: repositories{
			account: &repository.AccountRepository{},
		},
	}
}

type repositories struct {
	// user    repository.Repository[account.User]
	account repository.Repository[account.Account]
	// academy repository.Repository[account.Academy]
}
