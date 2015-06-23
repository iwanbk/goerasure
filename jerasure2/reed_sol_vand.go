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

// Encode encodes data using reed solomon
func (rsv ReedSolVand) Encode(data []byte) ([][]byte, [][]byte, int, error) {
	ed, ep, blockSize := utils.PrepareDataForEncode(rsv.k, rsv.m, rsv.w, data)

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
	ed, ep := utils.PrepareDataForDecode(rsv.k, rsv.m, encodedData, encodedParity)

	C.jerasure_matrix_decode(C.int(rsv.k), C.int(rsv.m), C.int(rsv.w),
		rsv.matrix, 1,
		(*C.int)(unsafe.Pointer(&missingIDs[0])),
		(**C.char)(unsafe.Pointer(&ed[0])),
		(**C.char)(unsafe.Pointer(&ep[0])),
		C.int(blockSize))

	return utils.ConvertResultData(ed, blockSize)
}
