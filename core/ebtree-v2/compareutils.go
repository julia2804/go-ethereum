package ebtree_v2

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
	"math/rand"
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

//compare a and b
func byteCompare(a, b []byte) int {
	lev, lemax := IndentBytes(a, b)
	return bytes.Compare(lev, lemax)
}

//used to indent the bytes
func IndentBytes(f []byte, s []byte) ([]byte, []byte) {
	var lev []byte
	var lemax []byte
	if len(f) < len(s) {
		t := len(s)

		b0 := byte(0)
		for i := 0; i < t-len(f); i++ {
			lev = append(lev, b0)
		}
		for j := 0; j < len(f); j++ {
			lev = append(lev, f[j])
		}
		for j := 0; j < len(s); j++ {
			lemax = append(lemax, s[j])
		}
		return lev, lemax
	} else if len(f) > len(s) {
		t := len(f)

		b0 := byte(0)
		for i := 0; i < t-len(s); i++ {
			lemax = append(lemax, b0)
		}
		for j := 0; j < len(s); j++ {
			lemax = append(lemax, s[j])
		}
		for j := 0; j < len(f); j++ {
			lev = append(lev, f[j])
		}
		return lev, lemax
	} else {
		for j := 0; j < len(f); j++ {
			lev = append(lev, f[j])
		}
		for j := 0; j < len(s); j++ {
			lemax = append(lemax, s[j])
		}
		return lev, lemax
	}
}

func add(b []byte, i uint64) []byte {
	f := BytesToInt(b)
	return IntToBytes(f + i)
}

func minus(b []byte, i uint64) []byte {
	f := BytesToInt(b)
	return IntToBytes(f - i)
}

//这里的n只是大致数量，输出的n才是真正的数量
func GenTestData(n int) {
	m1 := make(map[float32]string)
	var length int
	for i := 0; i < n; i++ {
		eth := rand.Float32() * 5
		if _, ok := m1[eth]; ok {
			//fmt.Println("is duclip",v)
		} else {
			fmt.Print("eth.sendTransaction({from:eth.coinbase,to:"+
				"\"0x4751c4cd1ef729afc3232b2064565f1d692a9346\",value:web3.toWei(", eth)
			fmt.Println(",'ether')})")
			m1[eth] = "1"
			length++
		}
	}
	fmt.Println(len(m1))
}


func BigAbs(a string, b string) hexutil.Big {
	Inta, _ := new(big.Int).SetString(a, 10)
	Intb, _ := new(big.Int).SetString(b, 10)
	return hexutil.Big(*Inta.Abs(Intb))
}

