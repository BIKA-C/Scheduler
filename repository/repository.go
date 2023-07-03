package repository

import (
	"scheduler/util"
	"time"
)

type Meta struct {
	CreatedAt time.Time
	LastLogin time.Time
}

type Repository[T any, ID interface{ util.UUID | util.ID }] interface {
	Save(*T) error
	Get(ID) T
	Delete(ID) error
}
