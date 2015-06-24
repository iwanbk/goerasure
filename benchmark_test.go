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

func BenchmarkEncode100x20x10000(b *testing.B) {
	benchmarkReedSolEncode(b, 100, 20, 10000)
}

func BenchmarkJerasure2100x20x10000(b *testing.B) {
	benchmarkJerasure2BindingEncode(b, 100, 20, 10000)
}

func BenchmarkEncode17x3x1M(b *testing.B) {
	benchmarkReedSolEncode(b, 17, 3, 1024*1024)
}

func BenchmarkJerasure217x3x1M(b *testing.B) {
	benchmarkJerasure2BindingEncode(b, 17, 3, 1024*1024)
}

// Benchmark 10 data shards and 4 parity shards with 16MB each.
func BenchmarkEncode10x4x16M(b *testing.B) {
	benchmarkReedSolEncode(b, 10, 4, 16*1024*1024)
}

// Benchmark 10 data shards and 4 parity shards with 16MB each.
func BenchmarkJerasure2Encode10x4x16M(b *testing.B) {
	benchmarkJerasure2BindingEncode(b, 10, 4, 16*1024*1024)
}
