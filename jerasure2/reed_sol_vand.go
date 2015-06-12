package jerasure2

// #include <jerasure.h>
// #include <jerasure/reed_sol.h>
// #cgo CFLAGS: -I /usr/include/jerasure/
// #cgo LDFLAGS: -lJerasure
import "C"

import (
	"unsafe"
)

type ReedSolVand struct {
	matrix (*C.int)
	k      int
	m      int
	w      int
}

func getAlignedDataSize(k, w int, dataLen int) {
	wordSize := w / 8
}
func prepareFragmentsForEncode(k, m, w int, data []byte, blockSize int) ([][]byte, [][]byte, error) {
	// Calculate data sizes, aligned_data_len guaranteed to be divisible by k
	dataLen := len(data)
	alignedDataLen := getAlignedDataSize(dataLen)
}

func (rsv ReedSolVand) Encode(data []byte, blockSize int) ([][]byte, [][]byte, error) {
	encodedData, parity := prepareFragmentsForEncode(rsv.k, rsv.m, data, blockSize)

	C.jerasure_matrix_encode(C.int(rsv.k), C.int(rsv.m), C.int(rsv.w), rsv.matrix, (**C.char)(unsafe.Pointer(&data[0])),
		(**C.char)(unsafe.Pointer(&parity[0])), C.int(blockSize))
	return nil
}
