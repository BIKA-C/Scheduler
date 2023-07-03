package util

type UUID string
type ID uint32

func NewUUID() UUID {
	return UUID(RandomString(8))
}

func NewID() ID {
	return ID(RandomNumber())
}
