package jerasure2

// #include <stdlib.h>
// #include <jerasure.h>
// #include <jerasure/reed_sol.h>
// #cgo CFLAGS: -I /usr/include/jerasure/
// #cgo LDFLAGS: -lJerasure
import "C"

import (
	"errors"
	"unsafe"
)

var ErrInvalidBlocksize = errors.New("Invalid blockSize")

type Context struct {
	k            int
	m            int
	w            int
	dataSlices   [][]byte
	data         [](*C.char)
	codingSlices [][]byte
	coding       [](*C.char)
	inputSize    int
	blockSize    int
	codingSize   int
	matrix       (*C.int)
}

func fill(slice []byte, value byte) {
	for i := range slice {
		slice[i] = value
	}
}

func NewContext(k int, m int) *Context {
	self := Context{}
	self.k = k
	self.m = m
	self.w = 8
	self.dataSlices = make([][]byte, k)
	self.data = make([](*C.char), k)
	self.codingSlices = make([][]byte, m)
	self.coding = make([](*C.char), m)
	self.matrix = C.reed_sol_vandermonde_coding_matrix(C.int(self.k), C.int(self.m), C.int(self.w))
	return &self
}

func (self *Context) DeleteContext() {
	C.free(unsafe.Pointer(self.matrix))
}

func (self *Context) SetDataSize(dataSize int) {
	self.inputSize = dataSize
	self.blockSize = (dataSize + (self.k - 1)) / self.k
	self.codingSize = self.blockSize * self.m
}

func (self *Context) RoundDataSize(inputSize int) int {
	return ((inputSize + (self.k - 1)) / self.k) * self.k
}

func (self *Context) GetCodingSize() int {
	return self.codingSize
}

func (self *Context) GetBlockSize() int {
	return self.blockSize
}

func (self *Context) SetDataSlice(slice []byte, index int) {
	self.dataSlices[index] = slice
	self.data[index] = (*C.char)(unsafe.Pointer(&slice[0]))
}

func (self *Context) SetCodingSlice(slice []byte, index int) {
	self.codingSlices[index] = slice
	self.coding[index] = (*C.char)(unsafe.Pointer(&slice[0]))
}

func (self *Context) PutData(data []byte, offset int) int {
	chunkLen := (len(data) + (self.k - 1)) / self.k
	for i := range self.dataSlices {
		slice := self.dataSlices[i]
		copied := copy(slice[offset:offset+chunkLen], data)
		fill(slice[copied:chunkLen], 0)
		data = data[copied:]
	}
	return chunkLen
}

func (self *Context) GetData(offset int, length int) []byte {
	data := make([]byte, length)
	retval := data
	chunkLen := (length + (self.k - 1)) / self.k
	for i := range self.dataSlices {
		slice := self.dataSlices[i]
		copied := copy(data, slice[offset:offset+chunkLen])
		data = data[copied:]
	}
	return retval
}

func (self *Context) Encode(blockSize int) error {
	if blockSize > self.blockSize {
		return ErrInvalidBlocksize
	}
	if blockSize == 0 {
		blockSize = self.blockSize
	}
	C.jerasure_matrix_encode(C.int(self.k), C.int(self.m), C.int(self.w), self.matrix, (**C.char)(unsafe.Pointer(&self.data[0])), (**C.char)(unsafe.Pointer(&self.coding[0])), C.int(blockSize))
	return nil
}

func (self *Context) Decode(erasures []int, blockSize int) error {
	if blockSize > self.blockSize {
		return ErrInvalidBlocksize
	}
	if blockSize == 0 {
		blockSize = self.blockSize
	}
	C.jerasure_matrix_decode(C.int(self.k), C.int(self.m), C.int(self.w), self.matrix, 1, (*C.int)(unsafe.Pointer(&erasures[0])), (**C.char)(unsafe.Pointer(&self.data[0])), (**C.char)(unsafe.Pointer(&self.coding[0])), C.int(blockSize))
	return nil
}
