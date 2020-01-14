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
	ebt.FirstLeaf.fstring()
	fmt.Print("hello")
	//test topk search
	ds, err := ebt.TopkVSearch(100)
	if err != nil {
		return
	}
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

func TestInsertToNodes(t *testing.T) {
	le := ConstructTestLeaf()
	var rd ResultD
	rd.Value = IntToBytes(10099)
	var td11 TD
	var td12 TD
	var td13 TD
	td11.IdentifierData = []byte("hoho")
	td12.IdentifierData = []byte("piupiu")
	td13.IdentifierData = []byte("rhfur")
	rd.ResultData = append(rd.ResultData, td11)
	rd.ResultData = append(rd.ResultData, td12)
	rd.ResultData = append(rd.ResultData, td13)

	tebt, err := NewEBTree()
	if err != nil {
		fmt.Println(err)
		return
	}
	i, err, flag := tebt.SearchInNode(rd.Value, &le)
	tebt.InsertToLeaf(rd, &le, int(i), flag)

}

func TestPer(t *testing.T) {
	var i int
	var e int
	i = 3
	e = 7
	fmt.Println(float32(i) / float32(e) * 100)
}

func ConstructTestLeaf() LeafNode {
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

	return testle
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
	result1, _ := EncodeNode(&test)
	fmt.Println(result1)
	in, err := DecodeNode(result1)
	fmt.Println(err)
	in.fstring()

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

	result2, _ := EncodeNode(&testle)
	fmt.Println(result2)
	le, err := DecodeNode(result2)
	fmt.Println(err)
	le.fstring()

	/*var me Meta
	me.FirstLeaf = []byte("hello")
	me.Root = []byte("world")
	me.Sequence = IntToBytes(12)

	result, err := EncodeMeata(me)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)

	me2, err := DecodeMeta(result)
	fmt.Println(me2.Sequence)
	fmt.Println(me2.Root)
	fmt.Println(me2.FirstLeaf)*/
}
