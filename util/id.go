package util

import "strconv"

type UUID string
type ID uint32

const UUID_LENGTH = 8

func NewUUID() UUID {
	return UUID(RandomString(UUID_LENGTH))
}

func (u UUID) IsUUID() bool {
	return len(u) == UUID_LENGTH
}

func NewID() ID {
	return ID(RandomNumber())
}

func ParseID(c string) ID {
	if i, err := strconv.ParseUint(c, 10, 32); err != nil {
		panic(err)
	} else {
		return ID(i)
	}
}

func (u UUID) String() string {
	return string(u)
}

func (u UUID) Str() string {
	return string(u)
}
