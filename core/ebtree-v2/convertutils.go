package ebtree_v2

import (
	"encoding/binary"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
)

//todo fill TD struct
func Convert2IdentifierData(blockno int, txIndex int) []byte {
	return []byte(string(blockno) + "," + string(txIndex))
}

func StringToBig(a string) hexutil.Big {
	Inta, _ := new(big.Int).SetString(a, 10)
	return hexutil.Big(*Inta)
}

func IntToBytes(i uint64) []byte {
	var buf = PutUint64(i)
	return buf
}

func IntToBytesOrigin(i uint64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, i)
	return buf
}

func PutUint64(v uint64) []byte {
	if v == 0 {
		return make([]byte, 1)
	}

	var b = make([]byte, 8)
	_ = b[7] // early bounds check to guarantee safety of writes below
	b[0] = byte(v >> 56)
	b[1] = byte(v >> 48)
	b[2] = byte(v >> 40)
	b[3] = byte(v >> 32)
	b[4] = byte(v >> 24)
	b[5] = byte(v >> 16)
	b[6] = byte(v >> 8)
	b[7] = byte(v)

	var length int
	for i := 0; i < 8; i++ {
		if b[i] == 0 {
			length++
		} else {
			break
		}
	}

	return b[length:]
}

func BytesToInt(data []byte) (num uint64) {
	//return binary.BigEndian.Uint64(b)
	var ret uint64 = 0
	var len int = len(data)
	var i uint = 0
	for i = 0; i < uint(len); i++ {
		ret = ret<<(8) | uint64(data[i])
	}
	return ret
}

func BytesToIntOrigin(data []byte) (num uint64) {
	return binary.BigEndian.Uint64(data)
}
