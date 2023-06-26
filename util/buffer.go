package util

import (
	"bytes"
	"unsafe"
)

func S2B(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func B2S(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

type BufferPool struct {
	list chan *bytes.Buffer
}

func NewBufferPool(poolSize int) *BufferPool {
	b := &BufferPool{
		list: make(chan *bytes.Buffer, poolSize),
	}

	return b
}

func (p *BufferPool) Put(b *bytes.Buffer) {
	b.Reset()
	select {
	case p.list <- b:
	default:
	}
}

func (p *BufferPool) Get() *bytes.Buffer {
	select {
	case b := <-p.list:
		return b
	default:
		return &bytes.Buffer{}
	}

}
