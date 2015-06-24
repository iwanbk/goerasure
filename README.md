# Go binding for C erasure coding library

## Install

Install the Go part
```
go get github.com/iwanbk/goerasure
```

We also need to install jerasure2 C library (http://lab.jerasure.org/jerasure/jerasure).

TODO : add detailed steps


## Test

```
go test ./...
```

## Benchmark

This lib has been benchmarked with other library (https://github.com/klauspost/reedsolomon).

Get the other library
```
go get github.com/klauspost/reedsolomon
```

Do the benchmark

```
go test -bench=.
```


## Benchmark Results

TODO
