package repository

import "time"

type Meta struct {
	CreatedAt time.Time
	LastLogin time.Time
}

type Repository[T any] interface {
	Save(T) error
}
