package ebtree_v2

import (
	"fmt"
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

	dt, err := ebt.RangeSearch(IntToBytes(100), IntToBytes(10001))
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
