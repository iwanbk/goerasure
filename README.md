# Go binding for C erasure coding library

It was planned to support various C erasure coding libraries.
But currently only support jerasure2 lib.

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

This lib has been benchmarked with other pure Go library (https://github.com/klauspost/reedsolomon).

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

### Encode benchmark

Each benchmark line will be something like : BenchmarkReedsolomonEncode10x2x10000      100000             21265 ns/op        4702.37 MB/s, it means:

- Reedsolomon -> it use reedsolomon library (https://github.com/klauspost/reedsolomon), jerasure2 binding has 'Jerasure2RS' string
- 10x2x10000 -> shards size = 10, parity size = 2, blocksize = 10000 bytes
- 100000 -> benchmarked by doing the operation 100000 times (decided by Go runtime)
- 21265 ns/op -> need 21265 ns per operation
- 4702.37 MB/s -> processing speed 

Benchmark on 1GB RAM 1 CPU Core Digital Ocean VM
```
BenchmarkReedsolomonEncode10x2x10000      100000             21265 ns/op        4702.37 MB/s
BenchmarkJerasure2RSEncode10x2x10000        5000            388215 ns/op         257.59 MB/s

BenchmarkReedSolomonEncode100x20x10000      1000           2173794 ns/op         460.03 MB/s
BenchmarkJerasure2RSEncode100x20x10000       500           2380024 ns/op         420.16 MB/s

BenchmarkReedSolomonEncode17x3x1M            200           7771865 ns/op        2293.63 MB/s
BenchmarkJerasure2RSEncode17x3x1M            200           9370720 ns/op        1902.29 MB/s

BenchmarkReedSolomonEncode10x4x16M            10         146259245 ns/op        1147.09 MB/s
BenchmarkJerasure2RSEncode10x4x16M            10         146811227 ns/op        1142.77 MB/s

BenchmarkReedSolomonEncode5x2x1M            1000           1394550 ns/op        3759.55 MB/s
BenchmarkJerasure2RSEncode5x2x1M            1000           1819554 ns/op        2881.41 MB/s

BenchmarkReedSolomonEncode10x2x1M            500           3100032 ns/op        3382.47 MB/s
BenchmarkJerasure2RSEncode10x2x1M            500           3603983 ns/op        2909.49 MB/s

BenchmarkReedSolomonEncode10x4x1M            200           5999499 ns/op        1747.77 MB/s
BenchmarkJerasure2RSEncode10x4x1M            200           7347378 ns/op        1427.14 MB/s

BenchmarkReedSolomonEncode50x20x1M            10         171563736 ns/op         305.59 MB/s
BenchmarkJerasure2RSEncode50x20x1M            10         173780297 ns/op         301.70 MB/s

BenchmarkReedSolomonEncode17x3x16M             5         234370903 ns/op        1216.93 MB/s
BenchmarkJerasure2RSEncode17x3x16M            10         182773941 ns/op        1560.47 MB/s
```

Benchmark on 2GB RAM 2 CPU Core Digital Ocean VM
```
BenchmarkReedsolomonEncode10x2x10000       50000             22397 ns/op        4464.78 MB/s
BenchmarkJerasure2RSEncode10x2x10000        3000            613834 ns/op         162.91 MB/s

BenchmarkReedSolomonEncode100x20x10000       500           2616154 ns/op         382.24 MB/s
BenchmarkJerasure2RSEncode100x20x10000       500           2973498 ns/op         336.30 MB/s

BenchmarkReedSolomonEncode17x3x1M            200           8855725 ns/op        2012.91 MB/s
BenchmarkJerasure2RSEncode17x3x1M            100          10384135 ns/op        1716.64 MB/s

BenchmarkReedSolomonEncode10x4x16M             5         241958575 ns/op         693.39 MB/s
BenchmarkJerasure2RSEncode10x4x16M            10         178249398 ns/op         941.22 MB/s

BenchmarkReedSolomonEncode5x2x1M            1000           1634996 ns/op        3206.66 MB/s
BenchmarkJerasure2RSEncode5x2x1M            1000           2053897 ns/op        2552.65 MB/s

BenchmarkReedSolomonEncode10x2x1M            300           3969426 ns/op        2641.63 MB/s
BenchmarkJerasure2RSEncode10x2x1M            300           4345243 ns/op        2413.16 MB/s

BenchmarkReedSolomonEncode10x4x1M            200           6608054 ns/op        1586.82 MB/s
BenchmarkJerasure2RSEncode10x4x1M            200           7898203 ns/op        1327.61 MB/s

BenchmarkReedSolomonEncode50x20x1M            10         183947510 ns/op         285.02 MB/s
BenchmarkJerasure2RSEncode50x20x1M            10         290146121 ns/op         180.70 MB/s

BenchmarkReedSolomonEncode17x3x16M             5         249068728 ns/op        1145.12 MB/s
BenchmarkJerasure2RSEncode17x3x16M            10         211232004 ns/op        1350.23 MB/s
```

From above results, pure Go library from klauspost has better result
