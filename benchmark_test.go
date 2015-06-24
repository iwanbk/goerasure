package goerasure

import (
	"math/rand"
	"testing"

	"github.com/iwanbk/goerasure/jerasure2"
	"github.com/klauspost/reedsolomon"
)

func fillRandom(b []byte) {
	for i := range b {
		b[i] = byte(rand.Int() & 0xff)
	}
}

var byteFill = byte(1)

func fill(b []byte) {
	for i := range b {
		b[i] = byteFill
	}
}

func benchmarkReedSolEncode(b *testing.B, dataShards, parityShards, shardSize int) {
	r, err := reedsolomon.New(dataShards, parityShards)
	if err != nil {
		b.Fatal(err)
	}
	shards := make([][]byte, dataShards+parityShards)
	for s := range shards {
		shards[s] = make([]byte, shardSize)
	}

	//rand.Seed(0)
	for s := 0; s < dataShards; s++ {
		//fillRandom(shards[s])
		fill(shards[s])
	}

	b.SetBytes(int64(shardSize * dataShards))
	//b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = r.Encode(shards)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchmarkJerasure2BindingEncode(b *testing.B, k, m, blockSize int) {
	data := make([]byte, k*blockSize)
	fill(data)

	// i put ResetTimer here because the other benchmark lib already receive [][]byte not []byte
	b.ResetTimer()

	rsv := jerasure2.NewReedSolVand(k, m)

	b.SetBytes(int64(k * blockSize))
	for i := 0; i < b.N; i++ {
		_, _, _, err := rsv.Encode(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReedsolomonEncode10x2x10000(b *testing.B) {
	benchmarkReedSolEncode(b, 10, 2, 10000)
}
func BenchmarkJerasure2Encode10x2x10000(b *testing.B) {
	benchmarkJerasure2BindingEncode(b, 10, 20, 10000)
}

func BenchmarkReedSolomonEncode100x20x10000(b *testing.B) {
	benchmarkReedSolEncode(b, 100, 20, 10000)
}

func BenchmarkJerasure2100x20x10000(b *testing.B) {
	benchmarkJerasure2BindingEncode(b, 100, 20, 10000)
}

func BenchmarkReedSolomonEncode17x3x1M(b *testing.B) {
	benchmarkReedSolEncode(b, 17, 3, 1024*1024)
}

func BenchmarkJerasure217x3x1M(b *testing.B) {
	benchmarkJerasure2BindingEncode(b, 17, 3, 1024*1024)
}

// Benchmark 10 data shards and 4 parity shards with 16MB each.
func BenchmarkReedSolomonEncode10x4x16M(b *testing.B) {
	benchmarkReedSolEncode(b, 10, 4, 16*1024*1024)
}

// Benchmark 10 data shards and 4 parity shards with 16MB each.
func BenchmarkJerasure2Encode10x4x16M(b *testing.B) {
	benchmarkJerasure2BindingEncode(b, 10, 4, 16*1024*1024)
}

// Benchmark 5 data shards and 2 parity shards with 1MB each.
func BenchmarkReedSolomonEncode5x2x1M(b *testing.B) {
	benchmarkReedSolEncode(b, 5, 2, 1024*1024)
}

// Benchmark 5 data shards and 2 parity shards with 1MB each.
func BenchmarkJerasure2Encode5x2x1M(b *testing.B) {
	benchmarkJerasure2BindingEncode(b, 5, 2, 1024*1024)
}

// Benchmark 1 data shards and 2 parity shards with 1MB each.
func BenchmarkReedSolomonEncode10x2x1M(b *testing.B) {
	benchmarkReedSolEncode(b, 10, 2, 1024*1024)
}

// Benchmark 1 data shards and 2 parity shards with 1MB each.
func BenchmarkJerasure2Encode10x2x1M(b *testing.B) {
	benchmarkJerasure2BindingEncode(b, 10, 2, 1024*1024)
}

// Benchmark 10 data shards and 4 parity shards with 1MB each.
func BenchmarkReedSolomonEncode10x4x1M(b *testing.B) {
	benchmarkReedSolEncode(b, 10, 4, 1024*1024)
}

// Benchmark 10 data shards and 4 parity shards with 1MB each.
func BenchmarkJerasure2Encode10x4x1M(b *testing.B) {
	benchmarkJerasure2BindingEncode(b, 10, 4, 1024*1024)
}

// Benchmark 50 data shards and 20 parity shards with 1MB each.
func BenchmarkReedSolomonEncode50x20x1M(b *testing.B) {
	benchmarkReedSolEncode(b, 50, 20, 1024*1024)
}

// Benchmark 50 data shards and 20 parity shards with 1MB each.
func BenchmarkJerasure2Encode50x20x1M(b *testing.B) {
	benchmarkJerasure2BindingEncode(b, 50, 20, 1024*1024)
}

// Benchmark 17 data shards and 3 parity shards with 16MB each.
func BenchmarkReedSolomonEncode17x3x16M(b *testing.B) {
	benchmarkReedSolEncode(b, 17, 3, 16*1024*1024)
}

// Benchmark 17 data shards and 3 parity shards with 16MB each.
func BenchmarkJerasure2Encode17x3x16M(b *testing.B) {
	benchmarkJerasure2BindingEncode(b, 17, 3, 16*1024*1024)
}
