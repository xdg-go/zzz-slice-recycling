// Copyright 2019 by David A. Golden. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package recycling

import (
	"testing"
)

type testCases struct {
	name string
	pool BytePool
}

var cases = []testCases{
	{name: "nullPool", pool: NewNullPool()},
	{name: "syncPool", pool: NewSyncPool()},
	{name: "pow2SyncPool", pool: NewPower2Pool()},
	{name: "reservedPool", pool: NewReservedPool()},
	{name: "leakySyncPool", pool: NewLeakySyncPool()},
}

func TestEncode(t *testing.T) {
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			enc := NewPoolEncoder(c.pool)
			enc.EncodeKV("A", 0)
			want := []byte{65, 0, 0, 0, 0, 0, 0, 0, 0, 0}
			if string(enc.Bytes()) != string(want) {
				t.Fatalf("encoding error; got '%v', want '%v'", enc.Bytes(), want)
			}
			enc.EncodeKV("B", 2)
			want = []byte{65, 0, 0, 0, 0, 0, 0, 0, 0, 0, 66, 0, 2, 0, 0, 0, 0, 0, 0, 0}
			if string(enc.Bytes()) != string(want) {
				t.Fatalf("encoding error; got '%v', want '%v'", enc.Bytes(), want)
			}
			enc.Release()
			enc = NewPoolEncoder(c.pool)
			enc.EncodeKV("A", 0)
			want = []byte{65, 0, 0, 0, 0, 0, 0, 0, 0, 0}
			if string(enc.Bytes()) != string(want) {
				t.Fatalf("encoding error; got '%v', want '%v'", enc.Bytes(), want)
			}
		})
	}
}

func Benchmark(b *testing.B) {
	InitBench()
	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			b.ReportAllocs()
			RunBench(b.N, Procs, c.pool)
		})
	}
}
