package ebtree_v2

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestConvertInt2Byte(t *testing.T) {
	size := 10000000

	array := make([]uint64, size)
	for i := 0; i < size; i++ {
		array[i] = rand.Uint64()
	}

	tt := time.Now()
	for i := 0; i < size; i++ {
		before := array[i]
		buf := IntToBytes(before)
		if len(buf) > 0 && buf[0] == 0 {
			fmt.Println(buf)
		}
		after := BytesToInt(buf)
		if before != after {
			fmt.Println("error", before, after)
			fmt.Println("-----------------")
		}
	}
	fmt.Printf("finished, timeElapsed: %f s\n", time.Now().Sub(tt).Seconds())
}

func TestConvertInt2ByteV2(t *testing.T) {
	size := 100

	array := make([]uint64, size)
	for i := 0; i < size; i++ {
		array[i] = rand.Uint64()
	}

	tt := time.Now()
	for i := 0; i < size; i++ {
		before := array[i]
		fmt.Println("before", before)
		buf := IntToBytes(before)
		fmt.Println("buffer", buf)
		fmt.Println("origin", IntToBytesOrigin(before))
		after := BytesToInt(buf)
		fmt.Println(after)
	}
	fmt.Printf("finished, timeElapsed: %f s\n", time.Now().Sub(tt).Seconds())
}

func TestNilByteArray(t *testing.T) {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, 0)
	fmt.Println(buf)

	var buf2 = PutUint64(0)
	fmt.Println(buf2)

	tt := time.Now()
	var size = 1000000
	for i := 0; i < size; i++ {
		var array []int
		array = append(array, 1)
	}
	fmt.Printf("finished, timeElapsed: %f s\n", time.Now().Sub(tt).Seconds())

}
func TestIntToBytes2(t *testing.T) {
	a1 := IntToBytes2(78762)
	a2 := IntToBytes(78762)
	fmt.Println(a1)
	fmt.Println(a2)
	fmt.Println(BytesToInt2(a1))
}
