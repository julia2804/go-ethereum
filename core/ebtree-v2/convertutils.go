package ebtree_v2

import (
	"encoding/binary"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
	"strconv"
)

//todo fill TD struct
func Convert2IdentifierData(blockno int, txIndex int) []byte {
	return []byte(strconv.Itoa(blockno) + "," + strconv.Itoa(txIndex))
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
func EncodeTds(tds []TD) ([]byte, error) {
	var encode []byte
	var err error
	encode, err = rlp.EncodeToBytes(tds)
	return encode, err
}

func DecodeTds(elems []byte) ([]TD, error) {
	var tds []TD

	elems, _, _ = rlp.SplitList(elems)
	c, _ := rlp.CountValues(elems)
	for i := 0; i < c; i++ {
		var rest2 []byte
		var td TD
		elems, rest2, _ = rlp.SplitList(elems)
		bdluf, _, _ := rlp.SplitString(elems)
		td.IdentifierData = bdluf
		elems = rest2
		tds = append(tds, td)
	}
	return tds, nil
}

func StringToEntity(v string, d string) Entity {
	var entity Entity
	entity.Value = []byte(v)
	entity.Data = []byte(d)
	return entity
}

func IntToBytes2(v int) []byte {
	if v > 256*256 {
		fmt.Println("v value", v)
	}
	var b = make([]byte, 4)
	_ = b[3] // early bounds check to guarantee safety of writes below
	b[0] = byte(v >> 24)
	b[1] = byte(v >> 16)
	b[2] = byte(v >> 8)
	b[3] = byte(v)

	return b
}
func BytesToInt2(data []byte) int {
	//return binary.BigEndian.Uint64(b)
	var ret int = 0
	var len int = len(data)
	var i uint = 0
	for i = 0; i < uint(len); i++ {
		ret = ret<<(8) | int(data[i])
	}
	return ret
}
