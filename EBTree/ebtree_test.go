// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package EBTree

import (
	"bytes"
	"encoding/binary"
	"encoding/csv"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"strconv"

	"math/rand"
	"strings"
	"testing"
	"time"
)

var r *rand.Rand
var DataNum int
var pre []byte
var wcount int
var wrong [1500]uint64

func init() {
	r = rand.New(rand.NewSource(time.Now().Unix()))
	spew.Config.Indent = "    "
	spew.Config.DisableMethods = false
	DataNum = 0
	pre = nil
	wcount = 0
}

// Used for testing :create a new tree without root
func newEmpty() (*EBTree, error) {
	se := IntToBytes(uint64(1))
	root, _ := constructLeafNode(se, 0, nil, false, true, nil, nil, 0)
	tree := &EBTree{NewDatabase(ethdb.NewMemDatabase()), &root, se, nil, 0, 0}
	tree.special = SetSpecialData(tree)
	return tree, nil
}

func SetSpecialData(tree *EBTree) []SpecialData {
	a := []uint64{0, 200000, 2100000, 499999}
	var result []SpecialData
	for _, i := range a {
		sData := SpecialData{}
		sData.value = IntToBytes(i)
		/*var da [][]byte
		da=append(da,[]byte("hello"))
		da=append(da,[]byte("world"))
		sData.data=da*/
		result = append(result, sData)
	}
	return result
}

func TestIntToBytes(t *testing.T) {
	i := uint64(1000)
	se := IntToBytes(i)
	fmt.Println(se)
	ib := BytesToInt(se)
	fmt.Println(ib)
}

// RandString 生成随机字符串
func RandString(len int) []byte {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return bytes
}

func updateString(tree *EBTree) {
	var a [509]int

	var v []byte
	for i := 1; i < 500; i++ {
		var j int
		v = []byte("qwerqwerqwerqwerqwerqwerqwerqwer")
		j = rand.Intn(30000)

		a[i-1] = j

		err := tree.InsertDataToTree(IntToBytes(uint64(j)), v)
		if err != nil {
			fmt.Sprintf("the error is not nil,%v", err)
		}

	}

	fmt.Println()

}
func updateString2(tree *EBTree) {
	var a [509]int

	var v []byte
	for i := 1; i < 500; i++ {
		var j int
		v = []byte("qwerqwerqwerqwerqwerqwerqwerqwer")
		j = rand.Intn(30000 - 10)
		if j < 0 {
			j = 0 - j
		}
		a[i-1] = j

		err := tree.InsertDataToTree(IntToBytes(uint64(j)), v)
		if err != nil {
			fmt.Sprintf("the error is not nil,%v", err)
		}

	}

	fmt.Println()

}

func TestInsert(t *testing.T) {
	tree, _ := newEmpty()
	v := []byte("qwermytecautqwerqwerqwerqwerqwer")
	tree.InsertDataToTree(IntToBytes(uint64(99)), v)
	tree.InsertDataToTree(IntToBytes(uint64(10)), v)
	tree.InsertDataToTree(IntToBytes(uint64(70)), v)
	tree.InsertDataToTree(IntToBytes(uint64(90)), v)
	tree.InsertDataToTree(IntToBytes(uint64(95)), v)
	tree.InsertDataToTree(IntToBytes(uint64(100)), v)
	tree.InsertDataToTree(IntToBytes(uint64(60)), v)
	tree.InsertDataToTree(IntToBytes(uint64(78)), v)
	printTree(tree.Root)

	//fmt.Println(BytesToInt(tree.sequence))

}

func combineAndPrintSearchValue(result []searchValue, pos []byte, tree *EBTree, k []byte, top bool) {
	_, result, err := tree.CombineSearchValueResult(result, pos, k, top)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for i, r := range result {
		fmt.Printf("the %dth value is %d,the data is:\n", i, r.value)
		for i, k := range r.data {
			fmt.Printf("the %dth data for value is:", i)
			fmt.Printf(hex.EncodeToString(k))
			fmt.Println()
		}
		fmt.Println()
	}
}

func TestTopkDataSearch(t *testing.T) {
	tree, _ := newEmpty()
	updateString(tree)
	var k []byte
	k = IntToBytes(uint64(3000))
	/*su, result, err := tree.TopkDataSearch(k, true)
	if !su {
		fmt.Printf("something may be wrong in top-k search\n")
		if err == nil {
			tree.CombineAndPrintSearchData(result, IntToBytes(uint64(0)), k, true)
			return
		} else {
			fmt.Println(err.Error())
			return
		}
	}
	tree.CombineAndPrintSearchData(result, IntToBytes(uint64(0)), k, true)*/
	_, _ = tree.Commit(nil)
	var triedb *Database
	tree.DBCommit()
	rid, _ := tree.Root.cache()
	triedb = tree.Db
	//triedb.Cap(1024)
	tri, _ := New(rid, triedb)
	updateString2(tri)
	_, _ = tri.Commit(nil)
	var tridb *Database
	tri.DBCommit()
	rid, _ = tri.Root.cache()
	tridb = tri.Db
	//triedb.Cap(1024)
	tri2, _ := New(rid, tridb)
	su, result, _ := tri2.TopkValueSearch(k, true)
	if !su {
		fmt.Printf("something wrong in top-k search")
	}

	combineAndPrintSearchValue(result, IntToBytes(uint64(0)), tri2, k, true)
}

func TestTopkValueSearch(t *testing.T) {
	tree, _ := newEmpty()
	updateString(tree)
	var k []byte
	k = IntToBytes(uint64(100))
	su, result, _ := tree.TopkValueSearch(k, true)
	if !su {
		fmt.Printf("something wrong in top-k search")
	}

	combineAndPrintSearchValue(result, IntToBytes(uint64(0)), tree, k, true)

}

func TestRangeValueSearch(t *testing.T) {
	tree, _ := newEmpty()
	updateString(tree)
	var k []byte
	k = IntToBytes(uint64(5000))
	min := IntToBytes(uint64(12))
	max := IntToBytes(uint64(900))
	su, result, err := tree.RangeValueSearch(min, max, k)
	if !su {
		fmt.Printf("something may be wrong in top-k search\n")
		if err == nil {
			combineAndPrintSearchValue(result, min, tree, k, false)
		} else {
			fmt.Println(err.Error())
		}
	}
	combineAndPrintSearchValue(result, min, tree, k, false)
}

func TestRangeDataSearch(t *testing.T) {
	tree, _ := newEmpty()
	updateString(tree)
	var k []byte
	k = IntToBytes(uint64(5000))
	min := IntToBytes(uint64(3))
	max := IntToBytes(uint64(900))
	_, result, _ := tree.RangeDataSearch(k, min, max)

	tree.CombineAndPrintSearchData(result, min, k, false)
}

func TestSearch(t *testing.T) {
	tree, _ := newEmpty()
	updateString(tree)
	result1, err := SearchNode(IntToBytes(uint64(25)), tree.Root, tree)
	if err != nil {
		fmt.Printf("somethine wrong in search node")
		return
	}
	fmt.Printf("the result for 31:\n")
	for i, r := range result1 {
		fmt.Printf("the %dth:\n", i)
		fmt.Printf(string(r))
		fmt.Println()
	}
	fmt.Println()
	result2, err := SearchNode(IntToBytes(uint64(18)), tree.Root, tree)
	if err != nil {
		fmt.Printf("somethine wrong in search node")
		return
	}
	fmt.Printf("the result for 552:\n")
	for i, r := range result2 {
		fmt.Printf("the %dth:\n", i)
		fmt.Printf(string(r))
		fmt.Println()
	}
	fmt.Println()

}

//TODO:BUG:when node is small, there is no problem
func TestMutiLeveInsert(t *testing.T) {
	//fmt.Sprintf("tree sequence")
	tree, _ := newEmpty()
	var arr [9000000][]byte
	for i := 0; i < 9000000; i++ {
		arr[i] = IntToBytes(uint64(rand.Intn(500000000)))
	}
	for i, a := range arr {
		su, _, err := tree.InsertData(tree.Root, uint8(i), nil, a, RandString(64))
		if !su {
			fmt.Sprintf("the error is not nil,%v", err)
		}
	}
	printTree(tree.Root)
	fmt.Printf("wrong number count:%d\n", wcount)
	fmt.Printf("end of output")
	fmt.Printf(string(tree.sequence))
}
func printTree(n EBTreen) {
	switch nt := (n).(type) {
	case *internalNode:
		if len(nt.Children) == 0 {
			fmt.Printf("empty tree")
			return
		}
		switch ct := (nt.Children[0]).(type) {
		case childEncode:
			return
		case child:
			printTree(ct.Pointer)
		default:
			return
		}

	case *leafNode:
		printNode(nt)
	}

}

func printNode(nt *leafNode) {
	for _, d := range nt.Data {
		switch dt := (d).(type) {
		case dataEncode:

			return
		case data:
			if dt.Keylist == nil {
				break
			}

			for _, k := range dt.Keylist {
				DataNum++

				fmt.Printf("%d,%d,%s\n", DataNum, dt.Value, k)
			}

			pre = dt.Value
		}

	}
	if nt.Next != nil {
		switch nnt := (nt.Next).(type) {
		case *leafNode:
			printNode(nnt)
		case *ByteNode:
			//todo: load from cache or database
			return
		default:
			err := errors.New("wrong type")
			fmt.Println(err)
			return
		}
	} else {
		fmt.Printf("end of print")
	}

}

func TestMissingNodeDisk(t *testing.T)    { testMissingNode(t, false) }
func TestMissingNodeMemonly(t *testing.T) { testMissingNode(t, true) }

//TODO：6/20，跟踪该测试并返回正确结果
func testMissingNode(t *testing.T, memonly bool) {
	src := []byte("hello")
	encodeStr := hex.EncodeToString(src)
	test, _ := hex.DecodeString(encodeStr)
	fmt.Println(bytes.Compare(test, src))
	diskdb := ethdb.NewMemDatabase()
	triedb := NewDatabase(diskdb)

	tree, _ := newEmpty()

	updateString(tree)

	_, _ = tree.Commit(nil)
	switch rt := (tree.Root).(type) {
	case *leafNode:

		tree.Db.Commit(rt.Id, true)

		triedb = tree.Db

		tree, _ = New(rt.Id, triedb)
	case *internalNode:

		tree.Db.Commit(rt.Id, true)

		triedb = tree.Db

		tree, _ = New(rt.Id, triedb)
	default:
		return

	}
	triedb.Cap(1*1024*1024 - ethdb.IdealBatchSize)
	//获取value对应的数据列表
	keylist, err := tree.TryGet(IntToBytes(uint64(25)))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	fmt.Printf("get the keylist1 as follows:\n")
	//输出返回结果
	for i, k := range keylist {
		fmt.Printf("get the %dth key:%s", i, string(k))
	}
	keylist2, err := tree.TryGet(IntToBytes(uint64(18)))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	fmt.Printf("get the keylist2 as follows:\n")
	//输出返回结果
	if keylist2 != nil {
		for i, k := range keylist2 {
			fmt.Printf("get the %dth key:%s", i, string(k))
		}
	}

}

func TestDecodeNode(t *testing.T) {
	var le leafNode
	//var in internalNode
	s1 := "hello2345"
	k1 := []byte(s1)
	s2 := "world87872345"
	k2 := []byte(s2)
	var keylist [][]byte
	keylist = append(keylist, k1)
	keylist = append(keylist, k2)
	dai := constructData(IntToBytes(uint64(100)), keylist)
	dai2 := constructData(IntToBytes(uint64(101)), keylist)
	var da []data
	da = append(da, dai)
	da = append(da, dai2)
	le, _ = constructLeafNode(IntToBytes(1), 1, da, false, false, nil, nil, 0)
	_ = le
	//6/24
	//b:=IntToBytes(2)
	//encode k1

	for i := 0; i < len(le.Data); i++ {
		switch dt := (le.Data[i]).(type) {
		case dataEncode:

			return
		case data:
			for j := 0; j < len(dt.Keylist); j++ {

			}
		}
	}

	r1, _ := rlp.EncodeToBytes([]uint{1, 2})
	_ = r1
	var i []uint
	rlp.DecodeBytes(r1, &i)
	lenew, _ := constructLeafNode(IntToBytes(1), 1, da, false, false, nil, nil, 0)
	_, r, _ := rlp.EncodeToReader(lenew)
	var rle leafNode
	rlp.Decode(r, &rle)

	buff := bytes.Buffer{}
	rlp.Encode(&buff, &le)
	//dle:=mustDecodeNode(IntToBytes(1), result, 0)
	//_=dle
	fmt.Println(buff.Bytes())
}

type TestRlpStruct struct {
	Id    []byte
	Count []byte
	Gen   []byte
	//C      []byte
	//BigInt *big.Int
	BM []ChildInterface
}

type TChildInterface interface {
	childString(string) string
}
type (
	tchild struct {
		Value   []byte
		Pointer EBTreen
		Id      []byte
	}
)
type tchildEncode []byte

func (n tchildEncode) childString(ind string) string {
	return fmt.Sprintf("<%x> ", string(n))
}

func (n tchild) childString(ind string) string {
	return fmt.Sprintf("<%x> ", string(n.Value))
}

type valN interface {
	fstring(string) string
}
type bcN struct {
	F string
	T DataInterface
}
type bcK struct {
	F byte
}

type bcV []byte

func (b bcV) fstring(ind string) string {
	return fmt.Sprintf("for %s, %x\n", ind, []byte(b))
}

func (b *bcN) fstring(ind string) string {
	return fmt.Sprintf("for %s,b:%x\n", ind, &b.F)
}

func TestInternalNode(t *testing.T) {
	/*internal:=internalNode{}
	leaf1:=leafNode{}
	leaf2:=leafNode{}
	leaf1.next=&leaf2
	leaf1.gen=0
	leaf1.dirty=true
	leaf1.special=false
	leaf1.id=IntToBytes(1)
	leaf1.count=2*/

}

//rlp用法
func TestRlp(t *testing.T) {
	//1.将一个整数数组序列化

	/*var rbc2 bcN
	var d child
	d.Id=IntToBytes(87)
	d.Value=IntToBytes(100)
	rbc1:=bcN{F:"mimtjyuytuhrtuyhurgyuttttnmtryutnfmtyujftyu meyuunj35tyjhrtbbw46555yrty otajulia",T:d}
	//rbc3:=bcK{B:"hello",A:8}

	//rbc2.a=9
	rbc2.F="test"
	test:=TestRlpStruct{A:3, B:"4"}

	buff:=bytes.Buffer{}
	rlp.Encode(&buff,&rbc1)
	bb:=buff.Bytes()
	var vr1 bcV
	vr1=bb[2:]
	fmt.Println(buff.Bytes())
	buff2:=bytes.Buffer{}
	rlp.Encode(&buff2,&rbc2)
	bb=buff2.Bytes()
	var vr2 bcV
	vr2=bb[2:]
	fmt.Println(buff.Bytes())
	/*r1,_:=rlp.EncodeToBytes(rbc1)
	r2,_:=rlp.EncodeToBytes(rbc2)
	var vr1 bcV
	var vr2 bcV
	r1 = r1[1:]
	r2 = r2[1:]
	vr1=r1
	vr2=r2*/

	test := internalNode{Id: IntToBytes(3)}
	var vr1 child
	var vr2 child
	var pb1 ByteNode
	pb1 = IntToBytes(87)
	vr1.Pointer = &pb1
	vr1.Value = IntToBytes(100)
	var pb2 ByteNode
	pb2 = IntToBytes(87)
	vr2.Pointer = &pb2
	vr2.Value = IntToBytes(30)
	//test.BM[0]=vr1
	//test.BM[1]=vr2
	var c1 childEncode
	var c2 childEncode
	result1, _ := rlp.EncodeToBytes(vr1)
	fmt.Println(result1)
	c1 = result1
	result2, _ := rlp.EncodeToBytes(vr2)
	fmt.Println(result2)
	c2 = result2
	test.Children = append(test.Children, c1)
	test.Children = append(test.Children, c2)
	//result,_:=rlp.EncodeToBytes(test)
	//fmt.Println(result)
	buff3 := bytes.Buffer{}
	rlp.Encode(&buff3, &test)
	bb := buff3.Bytes()
	fmt.Println(bb)
	//elems, _, err := rlp.SplitList(bb)
	//c, _ := rlp.CountValues(elems)
	//fmt.Println(c)

	elems, _, _ := rlp.SplitList(bb)
	c, _ := rlp.CountValues(elems)
	fmt.Println(elems)
	fmt.Println(c)
	kbuf, rest, _ := rlp.SplitString(elems)
	fmt.Println(kbuf)
	fmt.Println(rest)
	elems = rest
	bbuf, rest, _ := rlp.SplitString(elems)
	fmt.Println(bbuf)
	fmt.Println(rest)
	elems, _, _ = rlp.SplitList(rest)
	c, _ = rlp.CountValues(elems)
	fmt.Println(elems)
	fmt.Println(c)

	bd1uf, rest, _ := rlp.SplitString(elems)
	fmt.Println(bd1uf)
	fmt.Println(rest)
	elems = rest
	bd2uf, rest, _ := rlp.SplitString(elems)
	fmt.Println(bd2uf)
	fmt.Println(rest)
	/*var teststruct TestRlpStruct
	err := rlp.Decode(bytes.NewReader(bb), &teststruct)
	_=err
	//{A:0x3, B:"44", C:[]uint8{0x12, 0x32}, BigInt:32}
	fmt.Printf("teststruct=%#v\n", teststruct)*/

	//5.将任意一个struct序列化
	//将一个struct序列化到reader中
	/*s1:="hello2345"
	k1:=[]byte(s1)
	s2:="world87872345"
	k2:=[]byte(s2)
	var keylist [][]byte
	keylist=append(keylist,k1)
	keylist=append(keylist,k2)
	dai:=constructData(100,keylist,2)
	dai2:=constructData(101,keylist,2)
	var da []data
	da=append(da,dai)
	da=append(da,dai2)
	te:=leafNode{IntToBytes(1),1,da,false,false,nil,nil,0}
	buff=bytes.Buffer{}
	//var teststruct leafNode
	rlp.Encode(&buff,&te)
	fmt.Println(buff.Bytes())
	elems, _, err = rlp.SplitList(buff.Bytes())
	fmt.Println(elems)

	kbuf, rest, err = rlp.SplitString(elems)
	fmt.Println(kbuf)
	fmt.Println(rest)
	//{A:0x3, B:"44", C:[]uint8{0x12, 0x32}, BigInt:32}
	fmt.Printf("teststruct=%#v\n", teststruct)*/

}

func TestDataInterface(t *testing.T) {
	leaf1 := createTestLeaf(IntToBytes(1), 100)
	/*switch dt:=(leaf1.Data[0]).(type) {
	case data:
		fmt.Println("data,%s",dt.dataString("data"))
	case dataEncode:
		fmt.Println("dataEncode,%s",dt.dataString("dataEncode"))
	default:
		fmt.Println("wrong")
	}*/
	for i := 0; i < len(leaf1.Data); i++ {
		switch dt := (leaf1.Data[i]).(type) {
		case data:
			leaf1.Data[i] = encodeData(dt)
		case dataEncode:
			return
		default:
			return
		}

	}
	fmt.Println()

}

type DataStruct struct {
	A uint
	B string
}

func encodeData(d data) dataEncode {
	buff := bytes.Buffer{}
	rlp.Encode(&buff, d)
	bb := buff.Bytes()
	var vr1 dataEncode
	vr1 = bb
	return vr1
}

func TestEncodeData(t *testing.T) {
	k1 := []byte("hello")
	k2 := []byte("world")
	var k [][]byte
	k = append(k, k1)
	k = append(k, k2)
	d := data{Value: IntToBytes(1000000000000000), Keylist: k}
	bb, _ := rlp.EncodeToBytes(d)
	fmt.Println(bb)

	da, _ := decodeData(bb)
	_ = da
	fmt.Println()

}

func encodeTestLeaf(result *[]byte, le *leafNode) error {
	for i := 0; i < len(le.Data); i++ {
		bb, _ := rlp.EncodeToBytes(le.Data[i])
		var dataE dataEncode
		dataE = bb
		le.Data[i] = dataE
		le.Id = nil
	}
	r1, err := rlp.EncodeToBytes(le)
	if err != nil {
		err := wrapError(err, "encode leaf wrong")
		return err
	}
	for _, i := range r1 {
		*result = append(*result, i)
	}
	return nil
}
func TestEncodeTestLeafNode(t *testing.T) {
	s := "abcdfe233456"
	v, _ := rlp.EncodeToBytes(s)
	fmt.Println(v)
	le := createTestLeaf(IntToBytes(6), 1000)
	var result []byte
	err := encodeTestLeaf(&result, &le)
	if err != nil {
		return
	}
	fmt.Println(result)
	rle, _ := decodeNode(IntToBytes(6), result)
	_ = rle
	fmt.Println()
}

type tleafNode struct {
	Data   []DataInterface
	Nextid []byte
	Id     []byte
}
type tDataInterface interface {
	dataString(string) string
}

type (
	tdata struct {
		Keylist [][]byte
		Value   []byte
	}
)

type tdataEncode []byte

func (n tdataEncode) dataString(ind string) string {
	return fmt.Sprintf("<%x> ", string(n))
}

func (n tdata) dataString(ind string) string {
	return fmt.Sprintf("<%x> ", string(n.Value))
}

func createTestLeaf(id []byte, dataLow uint64) leafNode {
	leaf1 := leafNode{}
	leaf1.Id = id
	var nextb ByteNode
	nextb = IntToBytes(10)
	leaf1.Next = &nextb
	for i := uint8(0); i < maxLeafNodeCount; i++ {
		d1 := data{}
		d1.Value = add(IntToBytes(dataLow), 1)
		k1 := "0x994840d01d8c60b3a1a52b9119865dcbae683660482175038cd22c1cbbec679c"
		d1.Keylist = append(d1.Keylist, hexutil.MustDecode(k1))
		k2 := []byte("asdfasdfasdfasdfasdfasdfasdfasdf")
		d1.Keylist = append(d1.Keylist, k2)
		leaf1.Data = append(leaf1.Data, d1)
	}
	return leaf1
}

func TestEncodeInternalNode(t *testing.T) {
	in := createTestInternal(IntToBytes(1))
	var bb []byte
	err := encodeInternal(&bb, &in)
	if err != nil {
		wrapError(err, "wrong in encodeInternal")
		return
	}
	fmt.Println(bb)
	rle, _ := decodeNode(IntToBytes(1), bb)
	_ = rle
	fmt.Println()
}

func decodeTestLeaf(id, buf []byte) (leafNode, error) {
	if len(buf) == 0 {
		return leafNode{}, nil
	}
	elems, _, err := rlp.SplitList(buf)
	if err != nil {
		return leafNode{}, fmt.Errorf("decode error: %v", err)
	}

	buf = elems
	le := leafNode{}
	elems, rest, _ := rlp.SplitList(buf)
	c, _ := rlp.CountValues(elems)
	fmt.Println(elems)
	fmt.Println(c)

	for i := 0; i < c; i++ {
		kbuf, rest1, _ := rlp.SplitString(elems)
		d, _ := decodeData(kbuf)
		fmt.Print(i)
		fmt.Println(kbuf)
		le.Data = append(le.Data, d)
		fmt.Println(rest1)
		elems = rest1
	}
	elems = rest
	nextid, rest5, _ := rlp.SplitString(elems)
	var nextb ByteNode
	le.Next = &nextb
	fmt.Print("next:")
	fmt.Println(nextid)
	fmt.Println(rest5)
	return le, nil
}

func createTestInternal(id []byte) internalNode {
	internal1 := internalNode{}
	internal1.Id = id
	for i := uint8(0); i < maxInternalNodeCount; i++ {
		c1 := child{}
		c1.Value = IntToBytes(uint64(10000 + int(i)))
		le1 := createTestLeaf(IntToBytes(uint64(100+int(i))), 1000)
		var pb ByteNode
		pb = le1.Id
		c1.Pointer = &pb
		//c1.Pointer=&le1
		internal1.Children = append(internal1.Children, c1)
	}
	return internal1
}
func Float64ToByte(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)

	return bytes
}

func TestReadCsv(t *testing.T) {
	tree, _ := newEmpty()
	tree.special = SetSpecialData(tree)
	dat, err := ioutil.ReadFile("transaction2.csv")
	if err != nil {
		log.Fatal(err)
	}
	r := csv.NewReader(strings.NewReader(string(dat[:])))
	i := 0
	for {
		//process data
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		datastr := record[0][2:]
		//fmt.Printf("%T\n", datastr)
		datav, _ := rlp.EncodeToBytes(datastr)
		//fmt.Println(datav)

		var s string
		s = strings.TrimSpace(record[1])
		if s == "0" {
			//value=0的情况，首先定位，什么数字的插入会导致错误发生
			//发现一个错误，0被插入到tree中,hash是0x01cde6e47904c689e183ee5c3cc39167ab61b6e368d9a158ff87075ee4ea75c1,0
			//fmt.Println(bytes.Compare(data,value))
			tree.InsertData(tree.Root, uint8(i), nil, IntToBytes(0), datav)
			fmt.Println()
			if err != nil {
				fmt.Println(err)
			}

			i = i + 1
		} else if len(s) == 0 {
			fmt.Println("len is zero")
			i = i + 1
			continue
		} else if s[0] == '(' {
			//fmt.Println("wrong")
			i = i + 1
			continue
		} else {
			value, err := strconv.ParseFloat(s, 64)

			if err != nil {
				fmt.Println(err)
				continue
			}
			//value=0的情况，首先定位，什么数字的插入会导致错误发生
			//发现一个错误，0被插入到tree中,hash是0x01cde6e47904c689e183ee5c3cc39167ab61b6e368d9a158ff87075ee4ea75c1,0
			//fmt.Println(bytes.Compare(data,value))
			bv := Float64ToByte(value)
			dif := 8 - len(bv)
			b0 := byte(0)
			var s0 []byte
			for {
				if dif <= 0 {
					break
				} else {
					s0 = append(s0, b0)
					dif = dif - 1
				}
			}
			for i := 0; i < len(bv); i++ {
				s0 = append(s0, bv[i])
			}
			tree.InsertData(tree.Root, uint8(i), nil, s0, datav)
			fmt.Println()
			if err != nil {
				fmt.Println(err)
			}
			//fmt.Println(data)
			//fmt.Println(value)
			i = i + 1
		}

	}

	_, _ = tree.Commit(nil)
	var triedb *Database
	var rid []byte
	switch rt := (tree.Root).(type) {
	case *leafNode:

		tree.Db.Commit(rt.Id, true)

		triedb = tree.Db

		rid = rt.Id
	case *internalNode:

		tree.Db.Commit(rt.Id, true)

		triedb = tree.Db

		rid = rt.Id

	default:
		return

	}
	triedb.Cap(1024)
	tree, _ = New(rid, triedb)
	var s0 []byte
	value, e := new(big.Int).SetString("0", 10)
	if !e {
		fmt.Println("error")
	} else {
		bv := value.Bytes()
		dif := 8 - len(bv)
		b0 := byte(0)
		for {
			if dif <= 0 {
				break
			} else {
				s0 = append(s0, b0)
				dif = dif - 1
			}
		}
		for i := 0; i < len(bv); i++ {
			s0 = append(s0, bv[i])
		}
	}

	var s1 []byte
	value2, e := new(big.Int).SetString("20000000000000000000000000", 10)
	if !e {
		fmt.Println("error")
	} else {
		bv2 := value2.Bytes()
		dif := 8 - len(bv2)
		b0 := byte(0)
		for {
			if dif <= 0 {
				break
			} else {
				s1 = append(s1, b0)
				dif = dif - 1
			}
		}
		for i := 0; i < len(bv2); i++ {
			s1 = append(s1, bv2[i])
		}
		/*result2, err := SearchNode(s1, tree.Root, tree)
		if err != nil {
			fmt.Printf("somethine wrong in search node")
			return
		}
		fmt.Printf("the result for 552:\n")
		for i, r := range result2 {
			fmt.Printf("the %dth:\n", i)
			fmt.Printf("%v", r)
			fmt.Println()
		}
		fmt.Println()*/
	}
	var k []byte
	k = IntToBytes(uint64(500000))
	//_,result,err:=tree.TopkDataSearch(k,true)
	_, result, err := tree.RangeValueSearch(s0, s1, k)
	if len(result) == 0 {
		fmt.Println("no data")
	}
	for i := 0; i < len(result); i++ {
		fmt.Printf("%d value:", i)
		fmt.Println(result[i].value)
		fmt.Println("data:")
		for j := 0; j < len(result[i].data); j++ {
			fmt.Println(result[i].data[j])
		}
	}

	updateString(tree)
	_, result, err = tree.RangeValueSearch(s0, s1, k)
	if len(result) == 0 {
		fmt.Println("no data")
	}
	for i := 0; i < len(result); i++ {
		fmt.Printf("%d value:", i)
		fmt.Println(result[i].value)
		fmt.Println("data:")
		for j := 0; j < len(result[i].data); j++ {
			fmt.Println(result[i].data[j])
		}
	}

}
