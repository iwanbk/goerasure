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
	rsv := ReedSolVand{
		k: k,
		m: m,
		w: 8,
	}
	rsv.matrix = C.reed_sol_vandermonde_coding_matrix(C.int(k), C.int(m), C.int(rsv.w))
	return rsv
}

func ceill(f float64) int {
	return int(f + 0.9)
}
func getAlignedDataSize(k, w int, dataLen int) int {
	wordSize := w / 8
	alignmentMultiple := k * wordSize
	return ceill(float64(dataLen)/float64(alignmentMultiple)) * alignmentMultiple
}

func prepareDataForEncode(k, m, w int, data []byte) ([][]byte, int) {
	encodedData := make([][]byte, k)

	// Calculate data sizes, aligned_data_len guaranteed to be divisible by k
	dataLen := len(data)
	alignedDataLen := getAlignedDataSize(k, w, dataLen)

	blockSize := alignedDataLen / k
	//blockSize := (dataLen + (k - 1)) / k
	payloadSize := blockSize
	fmt.Printf("dataLen = %v, alignedDataLen=%v, blockSize=%v,payloadSize=%v\n", dataLen, alignedDataLen, blockSize, payloadSize)

	cursor := 0
	for i := 0; i < k; i++ {
		copySize := payloadSize
		if dataLen < payloadSize {
			copySize = dataLen
		}
		if dataLen > 0 {
			to := make([]byte, payloadSize)
			copy(to, data[cursor:cursor+copySize])
			//fmt.Printf("copy i = %v, cursor = %v, copySize=%v, len (data) = %v, copied=%v\n", i, cursor, copySize, len(data), copied)
			encodedData[i] = to
			//fmt.Printf("len to = %v\n", len(to))
		}
		cursor += copySize
		dataLen -= copySize
	}
	return encodedData, blockSize
}

// Encode encodes data using reed solomon
func (rsv ReedSolVand) Encode(data []byte) ([]*C.char, []*C.char, int, error) {
	encodedData, blockSize := prepareDataForEncode(rsv.k, rsv.m, rsv.w, data)

	ed := make([](*C.char), rsv.k)
	for i, v := range encodedData {
		ed[i] = (*C.char)(unsafe.Pointer(&v[0]))
	}

	ep := make([](*C.char), rsv.m)
	for k := 0; k < rsv.m; k++ {
		v := make([]byte, blockSize)
		ep[k] = (*C.char)(unsafe.Pointer(&v[0]))
	}

	C.jerasure_matrix_encode(C.int(rsv.k), C.int(rsv.m), C.int(rsv.w),
		rsv.matrix,
		(**C.char)(unsafe.Pointer(&ed[0])),
		(**C.char)(unsafe.Pointer(&ep[0])),
		C.int(blockSize))
	return ed, ep, blockSize, nil
}

// Decode decodes data
func (rsv ReedSolVand) Decode(encodedData, encodedParity [](*C.char), blockSize int, missingIDs []int) []byte {
	var data []byte

	fmt.Printf("k=%v,m=%v\n", rsv.k, rsv.m)
	/*for _, v := range encodedData {
		fmt.Printf("len = %v\n", len(C.GoString(v)))
	}
	fmt.Printf("hello")
	for _, v := range encodedParity {
		fmt.Printf("len = %v\n", len(C.GoString(v)))
	}*/

	C.jerasure_matrix_decode(C.int(rsv.k), C.int(rsv.m), C.int(rsv.w),
		rsv.matrix, 1,
		(*C.int)(unsafe.Pointer(&missingIDs[0])),
		(**C.char)(unsafe.Pointer(&encodedData[0])),
		(**C.char)(unsafe.Pointer(&encodedParity[0])),
		C.int(blockSize))

	for _, d := range encodedData {
		data = append(data, C.GoBytes(unsafe.Pointer(d), C.int(blockSize))...)
	}
	return data
}
