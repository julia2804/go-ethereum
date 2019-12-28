package ebtree_v2

import (
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