package jerasure2

// #include <jerasure.h>
// #include <jerasure/reed_sol.h>
// #cgo CFLAGS: -I /usr/include/jerasure/
// #cgo LDFLAGS: -lJerasure
import "C"

import (
	"fmt"
	"unsafe"
)

type ReedSolVand struct {
	matrix (*C.int)
	k      int
	m      int
	w      int
}

func NewReedSolVand(k, m int) ReedSolVand {
	return ReedSolVand{
		matrix: C.reed_sol_vandermonde_coding_matrix(C.int(k), C.int(m), C.int(8)),
		k:      k,
		m:      m,
		w:      8,
	}
}

func getAlignedDataSize(k, w int, dataLen int) int {
	wordSize := w / 8
	alignmentMultiple := k * wordSize
	return int(((float64(dataLen) / float64(alignmentMultiple)) * float64(alignmentMultiple)) + 0.9)
}

func prepareFragmentsForEncode(k, m, w int, data []byte) ([][]byte, [][]byte, int) {
	encodedData := make([][]byte, k)
	encodedParity := make([][]byte, m)

	// Calculate data sizes, aligned_data_len guaranteed to be divisible by k
	dataLen := len(data)
	alignedDataLen := getAlignedDataSize(k, w, dataLen)

	blockSize := int(alignedDataLen / k)
	payloadSize := alignedDataLen / k

	cursor := 0
	for i := 0; i < k; i++ {
		copySize := payloadSize
		if dataLen < payloadSize {
			copySize = dataLen
		}
		if dataLen > 0 {
			to := make([]byte, copySize)
			copied := copy(to, data[cursor:cursor+copySize])
			fmt.Printf("copy i = %v, cursor = %v, copySize=%v, len (data) = %v, copied=%v\n", i, cursor, copySize, len(data), copied)
			encodedData[i] = to
		}
		cursor += copySize
		dataLen -= copySize
	}
	return encodedData, encodedParity, blockSize
}

func prepareFragmentsForDecode(k, m int, encodedData, encodedParity [][]byte, missingIDs []int) {

}

// Decode decodes data
func (rsv ReedSolVand) Decode(encodedData, encodedParity [][]byte, blockSize int) []byte {
	var data []byte
	missingIDs := []int{}

	// fill mising IDs

	C.jerasure_matrix_decode(C.int(rsv.k), C.int(rsv.m), C.int(rsv.w), rsv.matrix, 1, (*C.int)(unsafe.Pointer(&missingIDs[0])), (**C.char)(unsafe.Pointer(&data[0])), (**C.char)(unsafe.Pointer(&encodedParity[0])), C.int(blockSize))
	return nil
	return data
}

// Encode encodes data using reed solomon
func (rsv ReedSolVand) Encode(data []byte) ([][]byte, [][]byte, int, error) {
	encodedData, encodedParity, blockSize := prepareFragmentsForEncode(rsv.k, rsv.m, rsv.w, data)

	ed := make([](*C.char), rsv.k)
	fmt.Printf("k = %v, len ed = %v\n", rsv.k, len(ed))
	for i, v := range encodedData {
		fmt.Printf("idx = %v v=%v, len v = %d\n", i, string(v), len(v))
		//v := []byte("hai")
		ed[i] = (*C.char)(unsafe.Pointer(&v[0]))
	}

	ep := make([](*C.char), rsv.m)
	for k, v := range encodedParity {
		v = make([]byte, blockSize)
		ep[k] = (*C.char)(unsafe.Pointer(&v[0]))
	}

	fmt.Printf("blockSize = %v\n", blockSize)
	C.jerasure_matrix_encode(C.int(rsv.k), C.int(rsv.m), C.int(rsv.w),
		rsv.matrix,
		(**C.char)(unsafe.Pointer(&ed[0])),
		(**C.char)(unsafe.Pointer(&ep[0])),
		C.int(blockSize))
	return encodedData, encodedParity, blockSize, nil
}
