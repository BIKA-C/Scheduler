package course

import (
	"scheduler/util"
	"time"
)

type Meta struct {
	ID               util.ID   `json:"id"`
	LastModification time.Time `json:"-"`
}
