# Go binding for C erasure coding library

## Install

### Install the Go part
```
go get github.com/iwanbk/goerasure
```

### Install jerasure2 C library (http://lab.jerasure.org/jerasure/jerasure).

Full install script can be found on : install_jerasure2_c_lib.sh  

In case you have installed it, you need to fix some include file
```
cd /usr/local/include/
sudo ln -s jerasure/galois.h
```


## Test & Benchmark

This lib has been benchmarked with other library (https://github.com/klauspost/reedsolomon).

Get the other library
```
go get github.com/klauspost/reedsolomon
```

Do the test

```
go test ./...
```

Do the benchmark

```
go test -bench=.
```


## Benchmark Results

TODO
