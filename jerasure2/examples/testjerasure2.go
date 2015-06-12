package main

import (
	"fmt"
	"github.com/iwanbk/goerasure/jerasure2"
)

func fill(slice []byte, value byte) {
	for i := range slice {
		slice[i] = value
	}
}

const BUFFERSIZE int = 2 * 1024 * 1024
const K int = 16
const M int = 4

type ticket struct {
	offset int
	length int
}

func main() {
	fmt.Println("hello")
	ctx := jerasure2.NewContext(K, M)
	defer ctx.DeleteContext()
	ctx.SetDataSize(BUFFERSIZE)
	dataSize := ctx.RoundDataSize(BUFFERSIZE)
	codingSize := ctx.GetCodingSize()
	blockSize := ctx.GetBlockSize()

	// Make data-slices, coding-slices, bind slices to context
	data := make([]byte, dataSize)
	dataSlices := make([][]byte, K)
	for i := range dataSlices {
		begin := i * blockSize
		end := ((i + 1) * blockSize) - 1
		slice := data[begin:end]
		dataSlices[i] = slice
		ctx.SetDataSlice(slice, i)
	}
	coding := make([]byte, codingSize)
	codingSlices := make([][]byte, M)
	for i := range codingSlices {
		begin := i * blockSize
		end := ((i + 1) * blockSize) - 1
		slice := coding[begin:end]
		codingSlices[i] = slice
		ctx.SetCodingSlice(slice, i)
	}

	buffer1 := []byte("The quick brown fox jumps over the lazy dog 12345")
	buffer2 := []byte("Don't call us child, we'll call you!")
	buffer3 := []byte("This is one small step for a man, one giant leap for mankind.")

	var t ticket
	offset := ctx.PutData(buffer1, 0)
	t.offset = offset       // keep offset of second message
	t.length = len(buffer2) // keep real length of second message
	offset += ctx.PutData(buffer2, offset)
	offset += ctx.PutData(buffer3, offset)
	// ...
	for i := range dataSlices {
		slice := dataSlices[i]
		fmt.Printf("D%v <%v>\n", i, string(slice[:offset]))
	}

	fmt.Println("encode...")
	totalLen := offset
	ctx.Encode(totalLen)

	fmt.Println("corrupt...")
	fill(dataSlices[0], 'z')
	fill(codingSlices[0], 'z')
	for i := range dataSlices {
		slice := dataSlices[i]
		fmt.Printf("D%v <%v>\n", i, string(slice[:offset]))
	}

	fmt.Println("decode...")
	erasures := make([]int, 3)
	erasures[0] = 0
	erasures[1] = K
	erasures[2] = -1
	ctx.Decode(erasures, totalLen)
	for i := range dataSlices {
		slice := dataSlices[i]
		fmt.Printf("D%v <%v>\n", i, string(slice[:offset]))
	}

	message := ctx.GetData(t.offset, t.length)
	fmt.Printf("<%v>\n", string(message))

	fmt.Println("done")
}
