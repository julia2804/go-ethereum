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
	for i := 10; i <= 100000000; i = i * 10 {
		ds := ReturnResultD(i)
		ebt.InsertDatasToTree(ds)
	}
	ebt.FirstLeaf.fstring("hello")
	fmt.Print("hello")
	//test topk search
	ds := ebt.TopkVSearch(100)
	fmt.Println(len(ds))
}

func ReturnResultD(t int) []ResultD {
	var ds []ResultD
	var d ResultD
	for j := 1; uint8(j) <= MaxLeafNodeCapability; j++ {
		d.Value = IntToBytes(uint64(j + t))
		var td TD
		td.IdentifierData = IntToBytes(uint64(t))
		d.ResultData = append(d.ResultData, td)
		ds = append(ds, d)
	}
	return ds
}
