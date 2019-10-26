package EBTree

import (
	"encoding/binary"
	"reflect"
)

func IntToBytes(i uint64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, i)
	return buf
}

func BytesToInt(b []byte) (i uint64) {
	return binary.BigEndian.Uint64(b)
}

func typeof(v interface{}) string {
	return reflect.TypeOf(v).String()
}