// Copyright 2019 by David A. Golden. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package recycling

import (
	"encoding/binary"
)

// KVEncoder provides the API for a key-value encoder.
type KVEncoder interface {
	EncodeKV(key string, value int)
	Bytes() []byte
	Release()
}

// poolEncoder is private to require users to initialize it via a constructor
// to ensure the 'pool' field is non-nil and the 'buf' field is initialized.
type poolEncoder struct {
	pool BytePool
	buf  []byte
}

// NewPoolEncoder constructs a new encoder attached to a given pool.  If the
// pool argument is nil, the "null" pool will be used.
func NewPoolEncoder(pool BytePool) KVEncoder {
	if pool == nil {
		pool = NullPool{}
	}
	return &poolEncoder{pool: pool, buf: pool.Get()}
}

// EncodeKV appends a string and integer to the encoder's buffer.
func (pe *poolEncoder) EncodeKV(key string, value int) {
	writePos := len(pe.buf)
	pe.buf = pe.pool.Resize(pe.buf, len(pe.buf)+len(key)+1+8)
	writePos = writeCString(pe.buf, writePos, key)
	writeInt64(pe.buf, writePos, int64(value))
}

// Bytes provides access to the encoder's buffer.  Users must not
// not modify the buffer.
func (pe *poolEncoder) Bytes() []byte {
	return pe.buf
}

// Release returns the encoder's memory buffer to the pool.  After calling
// this function, the encoder must not be used again.
func (pe *poolEncoder) Release() {
	pe.pool.Put(pe.buf)
}

func writeInt64(dst []byte, offset int, value int64) int {
	binary.LittleEndian.PutUint64(dst[offset:offset+8], uint64(value))
	return offset + 8
}

func writeCString(dst []byte, offset int, value string) int {
	strlen := len(value)
	copy(dst[offset:offset+strlen], value)
	dst[offset+strlen] = 0
	return offset + strlen + 1
}
