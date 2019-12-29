package ebtree_v2

import (
	"encoding/binary"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
)

//todo fill TD struct
func Convert2IdentifierData(blockno int, txIndex int) []byte{
	return []byte(string(blockno) + "," + string(txIndex))
}

func StringToBig(a string) hexutil.Big {
	Inta, _ := new(big.Int).SetString(a, 10)
	return hexutil.Big(*Inta)
}

func IntToBytes(i uint64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, i)
	return buf
}

func BytesToInt(b []byte) (i uint64) {
	return binary.BigEndian.Uint64(b)
}