// Copyright 2019 by David A. Golden. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package recycling

import (
	"math/bits"
	"sync"
)

// As this gets very large, all pools perform similarly because
// resizing isn't necessary.  (But consumes more memory than necessary.)
// We make it a var instead of const so we can tweak it in benchmarking.
var StartCap = 256

type BytePool interface {
	Get() []byte                      // Get a zero slice of some capacity
	Put([]byte)                       // Return a slice
	Resize(orig []byte, n int) []byte // Set slice length to new size
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func powerOfTwo(n int) int {
	return bits.Len(uint(n - 1))
}

//////////////////////////////////////////////

type NullPool struct{}

func NewNullPool() *NullPool {
	return &NullPool{}
}

func (np NullPool) Get() []byte {
	return make([]byte, 0, StartCap)
}

func (np NullPool) Put(bs []byte) {}

func (np NullPool) Resize(orig []byte, size int) []byte {
	if size < cap(orig) {
		return orig[0:size]
	}
	temp := make([]byte, size, max(size, cap(orig)*2))
	copy(temp, orig)
	return temp
}

//////////////////////////////////////////////

type SyncPool struct {
	pool *sync.Pool
}

// N.B. There is no 'New' field in the Pool construction because constructing
// new slices manually when the pool is empty lets us be smarter about when we
// have to zero memory.  If `sp.Pool.Get` always gave us a slice, we wouldn't
// know if it's newly created (and thus already zeroed) or if we need to zero.
func NewSyncPool() *SyncPool {
	return &SyncPool{pool: &sync.Pool{}}
}

func (sp SyncPool) Get() []byte {
	bp := sp.pool.Get()
	if bp == nil {
		return make([]byte, 0, StartCap)
	}
	buf := bp.([]byte)
	for i := range buf {
		buf[i] = 0
	}
	return buf[0:0]
}

func (sp SyncPool) Put(bs []byte) {
	sp.pool.Put(bs)
}

func (sp SyncPool) Resize(orig []byte, size int) []byte {
	if size < cap(orig) {
		return orig[0:size]
	}
	temp := make([]byte, size, max(size, cap(orig)*2))
	copy(temp, orig)
	sp.Put(orig)
	return temp
}

//////////////////////////////////////////////

type Power2Pool struct {
	pools []*sync.Pool
}

func NewPower2Pool() *Power2Pool {
	sp := &Power2Pool{pools: make([]*sync.Pool, 63)}
	for i := range sp.pools {
		sp.pools[i] = &sync.Pool{}
	}
	return sp
}

func (sp Power2Pool) Get() []byte {
	return sp.getn(StartCap)
}

// XXX precompute 1<<uint(p) in an array?
func (sp Power2Pool) getn(n int) []byte {
	p := powerOfTwo(n)
	bp := sp.pools[p].Get()
	if bp == nil {
		return make([]byte, 0, 1<<uint(p))
	}
	buf := bp.([]byte)
	for i := range buf {
		buf[i] = 0
	}
	return buf[0:0]
}

func (sp Power2Pool) Put(bs []byte) {
	sp.pools[powerOfTwo(cap(bs))].Put(bs)
}

func (sp Power2Pool) Resize(orig []byte, size int) []byte {
	if size < cap(orig) {
		return orig[0:size]
	}
	temp := sp.getn(size)
	temp = temp[0:size]
	copy(temp, orig)
	sp.Put(orig)
	return temp
}

//////////////////////////////////////////////

// N.B. 20 elements in the pool is very arbitrary -- it's tied to the number
// of goroutine workers in our benchmarking code.
type ReservedPool struct {
	sync.Mutex
	pool      [20][]byte
	size      int
	targetCap int
}

func NewReservedPool() *ReservedPool {
	return &ReservedPool{targetCap: StartCap}
}

func (rp *ReservedPool) Get() []byte {
	rp.Lock()
	if rp.size == 0 {
		newCap := rp.targetCap
		rp.Unlock()
		return make([]byte, 0, newCap)
	}
	buf := rp.pool[rp.size-1]
	rp.pool[rp.size-1] = nil
	rp.size--
	rp.Unlock()
	for i := range buf {
		buf[i] = 0
	}
	return buf[0:0]
}

func (rp *ReservedPool) Put(bs []byte) {
	rp.Lock()
	if cap(bs) < rp.targetCap || rp.size == len(rp.pool) {
		// drop it
		rp.Unlock()
		return
	}
	rp.pool[rp.size] = bs
	rp.size++
	rp.Unlock()
	return
}

func (rp *ReservedPool) Resize(orig []byte, size int) []byte {
	if size < cap(orig) {
		return orig[0:size]
	}
	rp.Lock()
	rp.targetCap = max(size, rp.targetCap*2)
	newCap := rp.targetCap
	rp.Unlock()
	temp := make([]byte, size, newCap)
	copy(temp, orig)
	return temp
}

//////////////////////////////////////////////

type LeakySyncPool struct {
	pool *sync.Pool
}

func NewLeakySyncPool() *LeakySyncPool {
	return &LeakySyncPool{pool: &sync.Pool{}}
}

func (sp LeakySyncPool) Get() []byte {
	bp := sp.pool.Get()
	if bp == nil {
		return make([]byte, 0, StartCap)
	}
	buf := bp.([]byte)
	for i := range buf {
		buf[i] = 0
	}
	return buf[0:0]
}

func (sp LeakySyncPool) Put(bs []byte) {
	sp.pool.Put(bs)
}

func (sp LeakySyncPool) Resize(orig []byte, size int) []byte {
	if size < cap(orig) {
		return orig[0:size]
	}
	temp := make([]byte, size, max(size, cap(orig)*2))
	copy(temp, orig)
	return temp
}
