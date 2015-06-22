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

func PrepareDataForEncode(k, m, w int, data []byte) ([](*C.char), [](*C.char), int) {
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

func PrepareDataForDecode(k, m int, encodedData, encodedParity [][]byte) ([](*C.char), [](*C.char)) {
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

func ConvertResultData(ed [](*C.char), blockSize int) []byte {
	data := []byte{}

	for _, d := range ed {
		data = append(data, C.GoBytes(unsafe.Pointer(d), C.int(blockSize))...)
	}
	return data
}