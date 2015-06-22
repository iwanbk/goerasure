package jerasure2

// #include <jerasure.h>
// #include <jerasure/reed_sol.h>
// #cgo CFLAGS: -I /usr/include/jerasure/
// #cgo LDFLAGS: -lJerasure
import "C"

import (
	//"fmt"
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

func prepareDataForEncode(k, m, w int, data []byte) ([](*C.char), [](*C.char), int) {
	// Calculate data sizes, aligned_data_len guaranteed to be divisible by k
	dataLen := len(data)
	alignedDataLen := getAlignedDataSize(k, w, dataLen)

	blockSize := alignedDataLen / k
	payloadSize := blockSize
	//fmt.Printf("dataLen = %v, alignedDataLen=%v, blockSize=%v,payloadSize=%v\n", dataLen, alignedDataLen, blockSize, payloadSize)

	// prepare encoded data
	encodedData := make([](*C.char), k)
	cursor := 0
	for i := 0; i < k; i++ {
		copySize := payloadSize
		if dataLen < payloadSize {
			copySize = dataLen
		}
		if dataLen > 0 {
			to := make([]byte, payloadSize)
			copy(to, data[cursor:cursor+copySize])
			//fmt.Printf("copy i = %v, cursor = %v, copySize=%v, len (data) = %v\n", i, cursor, copySize, len(data))
			encodedData[i] = (*C.char)(unsafe.Pointer(&to[0]))
			//fmt.Printf("len to = %v, dataLen=%v\n", len(to), dataLen)
		}
		cursor += copySize
		dataLen -= copySize
	}

	// prepare encoded parity
	ep := make([](*C.char), m)
	for k := 0; k < m; k++ {
		v := make([]byte, blockSize)
		ep[k] = (*C.char)(unsafe.Pointer(&v[0]))
	}

	return encodedData, ep, blockSize
}

// Encode encodes data using reed solomon
func (rsv ReedSolVand) Encode(data []byte) ([][]byte, [][]byte, int, error) {
	ed, ep, blockSize := prepareDataForEncode(rsv.k, rsv.m, rsv.w, data)

	C.jerasure_matrix_encode(C.int(rsv.k), C.int(rsv.m), C.int(rsv.w),
		rsv.matrix,
		(**C.char)(unsafe.Pointer(&ed[0])),
		(**C.char)(unsafe.Pointer(&ep[0])),
		C.int(blockSize))

	// convert back to  [][]byte
	edBytes := make([][]byte, rsv.k)
	for i, v := range ed {
		edBytes[i] = C.GoBytes(unsafe.Pointer(v), C.int(blockSize))
	}
	epBytes := make([][]byte, rsv.m)
	for i, v := range ep {
		epBytes[i] = C.GoBytes(unsafe.Pointer(v), C.int(blockSize))
	}
	return edBytes, epBytes, blockSize, nil
}

// Decode decodes data
func (rsv ReedSolVand) Decode(encodedData, encodedParity [][]byte, blockSize int, missingIDs []int) []byte {
	var data []byte

	ed := make([](*C.char), rsv.k)
	for i, v := range encodedData {
		ed[i] = (*C.char)(unsafe.Pointer(&v[0]))
	}
	ep := make([](*C.char), rsv.m)
	for i, v := range encodedParity {
		ep[i] = (*C.char)(unsafe.Pointer(&v[0]))
	}

	C.jerasure_matrix_decode(C.int(rsv.k), C.int(rsv.m), C.int(rsv.w),
		rsv.matrix, 1,
		(*C.int)(unsafe.Pointer(&missingIDs[0])),
		(**C.char)(unsafe.Pointer(&ed[0])),
		(**C.char)(unsafe.Pointer(&ep[0])),
		C.int(blockSize))

	for _, d := range ed {
		data = append(data, C.GoBytes(unsafe.Pointer(d), C.int(blockSize))...)
	}
	return data
}
