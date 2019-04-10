// Copyright 2019 by David A. Golden. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package main

import (
	"flag"
	"log"

	"github.com/pkg/profile"
	"github.com/xdg-go/zzz-slice-recycling/recycling"
)

// XXX Really ought to parameterize starting capacity size and some of the bench parameters
var poolType = flag.String("pooltype", "null", "one of 'null', 'sync', 'power2', 'reserved', 'leakysync'")
var profileType = flag.String("profile", "cpu", "one of 'cpu', 'mem', 'trace'")

func main() {
	flag.Parse()

	var pool recycling.BytePool
	switch *poolType {
	case "null":
		pool = recycling.NewNullPool()
	case "sync":
		pool = recycling.NewSyncPool()
	case "power2":
		pool = recycling.NewPower2Pool()
	case "reserved":
		pool = recycling.NewReservedPool()
	case "leakysync":
		pool = recycling.NewLeakySyncPool()
	default:
		log.Fatalf("Unrecognized 'pooltype' argument '%s'", *poolType)
	}
	log.Printf("Benchmarking with pooltype '%s'", *poolType)

	recycling.InitBench()
	profPath := "./prof." + *poolType

	switch *profileType {
	case "cpu":
		defer profile.Start(profile.CPUProfile, profile.ProfilePath(profPath)).Stop()
	case "mem":
		defer profile.Start(profile.MemProfile, profile.ProfilePath(profPath)).Stop()
	case "trace":
		defer profile.Start(profile.TraceProfile, profile.ProfilePath(profPath)).Stop()
	default:
		log.Fatalf("Unrecognized 'profiletype' argument '%s'", *profileType)
	}

	recycling.RunBench(recycling.Iterations, recycling.Procs, pool)

}
