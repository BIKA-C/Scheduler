package controller

import "scheduler/account"

type accountUpdate struct {
	account.Account
	Verify account.Password `json:"verify,omitempty"`
}
