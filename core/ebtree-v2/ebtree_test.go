package ebtree_v2

import (
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"testing"
)

func TestInsert(t *testing.T) {
	var ebt *EBTree
	var err error
	ebt, err = NewEBTree()
	if err != nil {
		fmt.Print(err)
		return
	}
	for i := 10000000000; i >= 10; i = i / 10 {
		ds := ReturnResultD(i)
		ebt.InsertDatasToTree(ds)
	}
	ebt.FirstLeaf.fstring("hello")
	fmt.Print("hello")
	//test topk search
	ds := ebt.TopkVSearch(100)
	fmt.Println(len(ds))

	dt, err := ebt.RangeSearch(IntToBytes(1000), IntToBytes(100001))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(len(dt))

	var de ResultD
	de, err = ebt.EquivalentSearch(IntToBytes(100005))
	if err != nil {
		fmt.Println(err)
		return
	}
	if de.Value != nil {
		fmt.Println("found")
	} else {
		fmt.Println("not found")
	}
}

func ReturnResultD(t int) []ResultD {
	var ds []ResultD
	var d ResultD
	for j := int(MaxLeafNodeCapability); uint8(j) >= 1; j-- {
		d.Value = IntToBytes(uint64(j + t))
		var td TD
		td.IdentifierData = IntToBytes(uint64(t))
		d.ResultData = append(d.ResultData, td)
		ds = append(ds, d)
	}
	return ds
}

type Test interface {
	fstring()
}

type (
	tt struct {
		t int
		u int
	}
	ee struct {
		e int
		f int
	}
)

func (te *tt) fstring() {
	fmt.Println("hello")
}
func (te *ee) fstring() {
	fmt.Println("hello")
}

type pp struct {
	p int
	t []byte
	x []byte
}
type mm struct {
	te Test
}

func TestH(t *testing.T) {
	var t1 tt
	t1.t = 3
	t1.u = 5

	var m1 mm
	m1.te = &t1
	t1.t = 8
	fmt.Println("hello")
	changem(&t1)
	fmt.Println("hello")
}

func changem(s *tt) {
	s.t = 100
	s.u = 1000
}

func TestEncode(t *testing.T) {
	var tp pp
	tp.p = 4
	tp.t = []byte("bye-bye")
	tp.x = []byte("hello")
	b1, err := rlp.EncodeToBytes(tp)
	fmt.Println(b1)
	t0 := tt{
		t: 100,
		u: 100,
	}
	r1, err := rlp.EncodeToBytes(t0)
	fmt.Println(err)
	fmt.Println(r1)

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

	test := InternalNode{Id: IntToBytes(3)}
	var vr1 ChildData
	var vr2 ChildData
	var vr3 ChildData

	var pb1 IdNode
	pb1 = IntToBytes(87)
	vr1.NodePtr = &pb1
	vr1.Value = IntToBytes(100)

	var pb2 IdNode
	pb2 = IntToBytes(88)
	vr2.NodePtr = &pb2
	vr2.Value = IntToBytes(30)

	var pb3 IdNode
	pb3 = IntToBytes(89)
	vr3.NodePtr = &pb3
	vr3.Value = IntToBytes(3000)

	test.Children = append(test.Children, vr1)
	test.Children = append(test.Children, vr2)
	test.Children = append(test.Children, vr3)
	result1, _ := rlp.EncodeToBytes(test)
	fmt.Println(result1)
	in, err := DecodeInternal(result1)
	fmt.Println(err)
	fmt.Println(in.Id)

	testle := LeafNode{Id: IntToBytes(5)}
	var ntid IdNode
	ntid = IntToBytes(8)
	testle.NextPtr = &ntid

	var rd1 ResultD
	var rd2 ResultD
	var rd3 ResultD

	var td11 TD
	var td12 TD
	var td13 TD
	td11.IdentifierData = []byte("hello")
	td12.IdentifierData = []byte("world")
	td13.IdentifierData = []byte("!")

	var td21 TD
	var td22 TD
	var td23 TD
	td21.IdentifierData = []byte("My")
	td22.IdentifierData = []byte("Lover")
	td23.IdentifierData = []byte("biu")

	var td31 TD
	var td32 TD
	var td33 TD
	td31.IdentifierData = []byte("I")
	td32.IdentifierData = []byte("Love")
	td33.IdentifierData = []byte("China!")

	rd1.Value = IntToBytes(1000)
	rd1.ResultData = append(rd1.ResultData, td11)
	rd1.ResultData = append(rd1.ResultData, td12)
	rd1.ResultData = append(rd1.ResultData, td13)

	rd2.Value = IntToBytes(2000)
	rd2.ResultData = append(rd2.ResultData, td21)
	rd2.ResultData = append(rd2.ResultData, td22)
	rd2.ResultData = append(rd2.ResultData, td23)

	rd3.Value = IntToBytes(3000)
	rd3.ResultData = append(rd3.ResultData, td31)
	rd3.ResultData = append(rd3.ResultData, td32)
	rd3.ResultData = append(rd3.ResultData, td33)

	testle.LeafDatas = append(testle.LeafDatas, rd1)
	testle.LeafDatas = append(testle.LeafDatas, rd2)
	testle.LeafDatas = append(testle.LeafDatas, rd3)

	result2, _ := rlp.EncodeToBytes(testle)
	fmt.Println(result2)
	le, err := DecodeLeafNode(result2)
	fmt.Println(err)
	fmt.Println(le.Id)
}
func DecodeInternal(encode []byte) (InternalNode, error) {
	var in InternalNode
	var err error
	elems, _, _ := rlp.SplitList(encode)
	//the number of fields in internal node
	c, _ := rlp.CountValues(elems)
	fmt.Println(c)

	kbuf, rest, _ := rlp.SplitString(elems)
	in.Id = kbuf
	fmt.Println(kbuf)
	fmt.Println(rest)
	elems = rest
	bbuf, rest, _ := rlp.SplitString(elems)
	fmt.Println(bbuf)
	fmt.Println(rest)
	elems, _, _ = rlp.SplitList(elems)
	//the number of children
	c, _ = rlp.CountValues(elems)
	fmt.Println(elems)
	fmt.Println(c)
	for i := 0; i < c; i++ {
		var rest2 []byte
		elems, rest2, _ = rlp.SplitList(elems)

		//the number of fields in childData
		//c, _ = rlp.CountValues(elems)
		//fmt.Println(elems)
		//fmt.Println(c)
		var child ChildData

		bd1uf, rest, _ := rlp.SplitString(elems)
		fmt.Println(bd1uf)
		child.Value = bd1uf
		fmt.Println(rest)
		elems = rest
		bd2uf, _, _ := rlp.SplitString(elems)
		fmt.Println(bd2uf)
		var npid IdNode
		npid = bd2uf
		child.NodePtr = &npid
		in.Children = append(in.Children, child)
		fmt.Println(rest2)
		elems = rest2
	}

	return in, err
}

//todo:test the decode leaf node
func DecodeLeafNode(encode []byte) (LeafNode, error) {
	var le LeafNode
	var err error
	elems, _, _ := rlp.SplitList(encode)
	//the number of fields in internal node
	c, _ := rlp.CountValues(elems)
	fmt.Println(c)

	//get the id
	kbuf, rest, _ := rlp.SplitString(elems)
	le.Id = kbuf
	fmt.Println(kbuf)
	fmt.Println(rest)
	elems = rest

	//get the nextptr
	kbuf, rest, _ = rlp.SplitString(elems)
	var ntid IdNode
	ntid = kbuf
	le.NextPtr = &ntid
	fmt.Println(kbuf)
	fmt.Println(rest)
	elems = rest

	//get the data
	elems, _, _ = rlp.SplitList(elems)
	//the number of data
	c, _ = rlp.CountValues(elems)
	fmt.Println(elems)
	fmt.Println(c)
	for i := 0; i < c; i++ {
		var rest2 []byte
		elems, rest2, _ = rlp.SplitList(elems)

		//the number of fields in childData
		//c, _ = rlp.CountValues(elems)
		//fmt.Println(elems)
		//fmt.Println(c)
		var rd ResultD

		//get the value of resultd
		bd1uf, rest, _ := rlp.SplitString(elems)
		fmt.Println(bd1uf)
		rd.Value = bd1uf
		fmt.Println(rest)
		elems = rest

		//get the tds of resultd
		elems, _, _ = rlp.SplitList(elems)
		//the number of td
		tdc, _ := rlp.CountValues(elems)
		fmt.Println(elems)
		fmt.Println(tdc)
		for i := 0; i < tdc; i++ {
			var rest3 []byte
			elems, rest3, _ = rlp.SplitList(elems)

			var td TD
			//get the tds of td
			bd2uf, _, _ := rlp.SplitString(elems)
			fmt.Println(bd2uf)
			td.IdentifierData = bd2uf
			rd.ResultData = append(rd.ResultData, td)

			fmt.Println(rest3)
			elems = rest3
		}
		le.LeafDatas = append(le.LeafDatas, rd)
		elems = rest2
	}

	return le, err
}
