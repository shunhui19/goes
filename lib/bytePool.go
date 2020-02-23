package lib

import "sync"

// BytePool response the slice byte pool.
type BytePool struct {
	size int
	pool sync.Pool
}

// byteBuffer response a slice byte.
type byteBuffer struct {
	// B a []byte use in conn.Read.
	B []byte
}

// NewBytesPool return a bytes pool instance.
func NewBytesPool(size int) *BytePool {
	return &BytePool{pool: sync.Pool{}, size: size}
}

// Get get a byteBuffer from pool.
func (bp *BytePool) Get() *byteBuffer {
	v := bp.pool.Get()
	if v != nil {
		return v.(*byteBuffer)
	}
	return &byteBuffer{B: make([]byte, bp.size)}
}

// Put put byteBuffer back to pool.
func (bp *BytePool) Put(bf *byteBuffer) {
	bp.pool.Put(bf)
}
