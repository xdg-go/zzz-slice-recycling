# Slice Recycling Performance and Pitfalls

This repository contains demonstration code for a talk on slice recycling
by David Golden.  Slides will be posted to https://xdg.me/talks/

## Running benchmarks

To run benchmarks run `go test -bench` from the root directory of the
repository:

```
$ go test -bench . ./recycling
```

## Finding maximum heap usage

To find the maximum heap usage of each memory strategy, compile the `main.go`
file into `./main`, then run the `find-max-heap.sh` program:

```
$ go build -o main .
$ ./find-max-heap.sh
```

## Generating profiling flame graphs and trace data

To generate profile/trace data, compile the `main.go` file into `./main`, then
execute `./main` with the desired pooltype and/or profile/trace output.  The
default profiling output is a CPU profile.  Results are in a subdirectry by
pool type.

```
$ go build -o main .
$ ./main -pooltype null      -profile cpu
$ ./main -pooltype sync      -profile cpu
$ ./main -pooltype power2    -profile cpu
$ ./main -pooltype reserved  -profile cpu
$ ./main -pooltype leakysync -profile cpu

$ ./main -pooltype null      -profile trace
$ ./main -pooltype sync      -profile trace
$ ./main -pooltype power2    -profile trace
$ ./main -pooltype reserved  -profile trace
$ ./main -pooltype leakysync -profile trace
```

To view the profile output, call `go tool pprof` with the web browser option
and paths for the binary and profile data.

```
$ go tool pprof -http :8080 ./main ./prof.null/cpu.pprof
```

To view the trace output, call `go tool trace` similarly.

```
$ go tool trace -http :8080 ./main ./prof.null/trace.out
```

## Changing the starting container capacity

If the STARTCAP environment variable is set to a number, it override the
default starting container capacity of 256 bytes.  This works for any of the
analyses above.

## License

Copyright 2019 by David A. Golden. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License"); you may
not use this file except in compliance with the License. You may obtain
a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
