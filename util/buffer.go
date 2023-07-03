package util

import (
	"bytes"
	"sync"
	"unsafe"
)

func S2B(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func B2S(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

type BufferPool sync.Pool

func NewBufferPool() *BufferPool {
	b := &sync.Pool{
		New: func() any {
			return &bytes.Buffer{}
		},
	}
	return (*BufferPool)(b)
}

func (p *BufferPool) Put(b *bytes.Buffer) {
	(*sync.Pool)(p).Put(b)
}

func (p *BufferPool) Get() *bytes.Buffer {
	b := (*sync.Pool)(p).Get().(*bytes.Buffer)
	b.Reset()
	return b
}
