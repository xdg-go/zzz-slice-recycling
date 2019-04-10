// Copyright 2019 by David A. Golden. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package recycling

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	randSeed   int64 = 42
	maxKeys          = 1000
	maxKeyLen        = 200
	Iterations       = 400000

	// Procs has to be high enough to get enough churn in the pools
	Procs = 20
)

var randSizes [10000]int
var randKeys [10000]string

func InitBench() {
	// `go test -bench` will call this more than once so only execute
	// the initialization once
	var once sync.Once
	once.Do(func() {
		startCapEnv := os.Getenv("STARTCAP")
		if startCapEnv != "" {
			startCap, err := strconv.Atoi(startCapEnv)
			if err != nil {
				log.Fatal("Invalid value for STARTCAP: ", startCapEnv)
			}
			StartCap = startCap
		}

		rand.Seed(randSeed)
		for i := range randSizes {
			randSizes[i] = rand.Intn(maxKeys)
		}
		for i := range randKeys {
			randKeys[i] = strings.Repeat("a", rand.Intn(maxKeyLen))
		}
	})
}

func RunBench(N int, maxProcs int, pool BytePool) {
	var wg sync.WaitGroup
	for i := 0; i < maxProcs; i++ {
		wg.Add(1)
		go func() {
			// Pick random starting points in test data arrays for sizes
			// and keys so that the benchmarking loop doesn't have
			// to repeatedly call rand.Intn().
			sizeIdx := rand.Intn(len(randSizes))
			keyIdx := rand.Intn(len(randKeys))
			defer wg.Done()
			for j := 0; j < N/maxProcs; j++ {
				enc := NewPoolEncoder(pool)
				for i := 0; i < randSizes[sizeIdx]; i++ {
					k := randKeys[keyIdx]
					enc.EncodeKV(k, len(k))
					keyIdx = (keyIdx + 1) % len(randKeys)
				}
				enc.Release()
				sizeIdx = (sizeIdx + 1) % len(randSizes)
			}
		}()
	}
	wg.Wait()
}
