package jerasure2

import (
	"log"
	"reflect"
	"testing"
)

func TestReedSolVand(t *testing.T) {
	origData := []byte("The quick brown fox jumps over the lazy dog 12345")
	rsv := NewReedSolVand(16, 4)
	encodedData, encodedParity, blockSize, _ := rsv.Encode(origData)
	log.Printf("blockSize = %v\n", blockSize)

	recoveredData := rsv.Decode(encodedData, encodedParity, blockSize)

	if !reflect.DeepEqual(origData, recoveredData) {
		t.Fatalf("failed to decode data")
	}
}
