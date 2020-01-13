package ebtree_v2

import (
	"bytes"
	"fmt"
	"testing"
)

func TestByteArrayCompare(t *testing.T) {
	var a = []byte("hello")
	var b = []byte("hello")
	fmt.Println(a, b)
	fmt.Println(byteCompare(&a, &b))
	fmt.Println(byteCompare(&b, &a))
	fmt.Println(bytes.Compare(a, b))
}
