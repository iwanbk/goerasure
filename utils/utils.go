package utils

import "C"

import (
	"unsafe"
)

func ceill(f float64) int {
	return int(f + 0.9)
}

func getAlignedDataSize(k, w int, dataLen int) int {
	wordSize := w / 8
	alignmentMultiple := k * wordSize
	return ceill(float64(dataLen)/float64(alignmentMultiple)) * alignmentMultiple
}

// PrepareDataForEncode prepare needed data structure to do encoding
// it returns three data
// - encoded data array
// - encoded parity array
// - blocksize of the data
func PrepareDataForEncode(k, m, w int, data []byte) ([][]byte, [][]byte, int) {
	// Calculate data sizes, aligned_data_len guaranteed to be divisible by k
	dataLen := len(data)
	alignedDataLen := getAlignedDataSize(k, w, dataLen)

	blockSize := alignedDataLen / k
	payloadSize := blockSize
	//fmt.Printf("k=%v, m=%v,w=%v, dataLen = %v, alignedDataLen=%v, blockSize=%v,payloadSize=%v\n", k, m, w, dataLen, alignedDataLen, blockSize, payloadSize)

	// prepare encoded data
	//encodedData := make([](*C.char), k)
	encodedData := make([][]byte, k)
	cursor := 0
	for i := 0; i < k; i++ {
		copySize := payloadSize
		if dataLen < payloadSize {
			copySize = dataLen
		}
		if dataLen > 0 {
			//to := data[cursor : cursor+copySize]
			//encodedData[i] = (*C.char)(unsafe.Pointer(&to[0]))
			encodedData[i] = data[cursor : cursor+copySize]
		}
		cursor += copySize
		dataLen -= copySize
	}

	// prepare encoded parity
	ep := make([][]byte, m)
	for k := 0; k < m; k++ {
		ep[k] = make([]byte, blockSize)
	}

	return encodedData, ep, blockSize
}

// PrepareDataForDecode prepare all data needed to do decoding
// it convert encoded data and encoded parity to data type that ready to be used by cgo
func PrepareDataForDecode(k, m int, encodedData, encodedParity [][]byte) ([]*C.char, []*C.char) {
	ed := make([](*C.char), k)
	for i, v := range encodedData {
		ed[i] = (*C.char)(unsafe.Pointer(&v[0]))
	}
	ep := make([](*C.char), m)
	for i, v := range encodedParity {
		ep[i] = (*C.char)(unsafe.Pointer(&v[0]))
	}
	return ed, ep
}

// ConvertResultData convert returned result data (in [](*C.char)) to []byte
func ConvertResultData(ed [](*C.char), blockSize int) []byte {
	data := []byte{}

	for _, d := range ed {
		data = append(data, C.GoBytes(unsafe.Pointer(d), C.int(blockSize))...)
	}
	return data
}
