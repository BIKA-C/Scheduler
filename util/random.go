package util

import (
	"math/rand"
	"sync"
	"time"
)

var p sync.Pool = sync.Pool{
	New: func() any {
		return rand.NewSource(time.Now().UnixNano())
	},
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandomString(n int) string {
	b := make([]byte, n)
	src := p.Get().(rand.Source)
	defer p.Put(src)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return B2S(b)
}

const digitBytes = "0123456789"
const (
	digitIdxBits = 4                   // 4 bits to represent a digit index
	digitIdxMask = 1<<digitIdxBits - 1 // All 1-bits, as many as digitIdxBits
	digitIdxMax  = 63 / digitIdxBits   // # of digit indices fitting in 63 bits
)

func RandomNumberString(n int) string {
	b := make([]byte, n)
	src := p.Get().(rand.Source)
	defer p.Put(src)
	// A src.Int63() generates 63 random bits, enough for digitIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), digitIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), digitIdxMax
		}
		if idx := int(cache & digitIdxMask); idx < len(digitBytes) {
			b[i] = digitBytes[idx]
			i--
		}
		cache >>= digitIdxBits
		remain--
	}

	return B2S(b)
}

func RandomNumber() uint64 {
	src := p.Get().(rand.Source)
	defer p.Put(src)
	return uint64(src.Int63())
}
