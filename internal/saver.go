package internal

type Saver[T any] interface {
	Save(*T) error
}
