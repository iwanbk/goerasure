package jerasure2

// #include <jerasure.h>
// #include <jerasure/reed_sol.h>
// #cgo CFLAGS: -I /usr/include/jerasure/
// #cgo LDFLAGS: -lJerasure
import "C"

import (
	"unsafe"

	"github.com/iwanbk/goerasure/utils"
)

// ReedSolVand defines struct for jerasure2 reed solomon with vandermonde matrix
type ReedSolVand struct {
	matrix (*C.int)
	k      int
	m      int
	w      int
}

// NewReedSolVand create ReedSolVand object
func NewReedSolVand(k, m int) ReedSolVand {
	rsv := ReedSolVand{
		k: k,
		m: m,
		w: 8,
	}
	rsv.matrix = C.reed_sol_vandermonde_coding_matrix(C.int(k), C.int(m), C.int(rsv.w))
	return rsv
}

// Encode encodes data using reed solomon and vandermonde matrix
func (rsv ReedSolVand) Encode(data []byte) ([][]byte, [][]byte, int, error) {
	edBytes, epBytes, blockSize := utils.PrepareDataForEncode(rsv.k, rsv.m, rsv.w, data)

	ed := make([]*C.char, rsv.k)
	for i, d := range edBytes {
		ed[i] = (*C.char)(unsafe.Pointer(&d[0]))
	}

	ep := make([]*C.char, rsv.m)
	for i, d := range epBytes {
		ep[i] = (*C.char)(unsafe.Pointer(&d[0]))
	}

	C.jerasure_matrix_encode(C.int(rsv.k), C.int(rsv.m), C.int(rsv.w),
		rsv.matrix,
		(**C.char)(unsafe.Pointer(&ed[0])),
		(**C.char)(unsafe.Pointer(&ep[0])),
		C.int(blockSize))

	return edBytes, epBytes, blockSize, nil
}

// Decode decodes data
func (rsv ReedSolVand) Decode(encodedData, encodedParity [][]byte, blockSize int, missingIDs []int) []byte {
	ed, ep := utils.PrepareDataForDecode(rsv.k, rsv.m, encodedData, encodedParity)

	C.jerasure_matrix_decode(C.int(rsv.k), C.int(rsv.m), C.int(rsv.w),
		rsv.matrix, 1,
		(*C.int)(unsafe.Pointer(&missingIDs[0])),
		(**C.char)(unsafe.Pointer(&ed[0])),
		(**C.char)(unsafe.Pointer(&ep[0])),
		C.int(blockSize))

	return utils.ConvertResultData(ed, blockSize)
}
