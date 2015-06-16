package jerasure2

import (
	"log"
	"reflect"
	"testing"
)

func TestReedSolVand(t *testing.T) {
	var origData []byte
	for i := 0; i < 128; i++ {
		origData = append(origData, []byte("aa")...)
	}
	rsv := NewReedSolVand(16, 4)
	encodedData, encodedParity, blockSize, _ := rsv.Encode(origData)
	log.Printf("blockSize = %v\n", blockSize)
	// fill mising IDs
	missingIDs := []int{0, rsv.k, -1}

	recoveredData := rsv.Decode(encodedData, encodedParity, blockSize, missingIDs)

	if !reflect.DeepEqual(origData, recoveredData) {
		t.Fatalf("failed to decode data")
	}
}
