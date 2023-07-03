package account

import "scheduler/util"

type Instructor struct {
	Account `json:"account"`
	Contact `json:"contact"`
	InsID   util.UUID `json:"institutionID"`
	Name    string    `json:"name"`
}
