package ebtree_v2

import (
	"fmt"
	"github.com/ethereum/go-ethereum/core"
	"strconv"
	"time"
)

func ExperStart(bc *core.BlockChain) {
	Initial(nil, 0, 0)
	var db *Database
	db = NewDatabase(*bc.GetDB())
	tree, err := NewEBTreeFromDb(db)
	if err != nil {
		fmt.Println("some errors happen", err.Error())
		fmt.Println("\n\n\n")
	}
	var content string

	content = ""
	value := "1000000000000000"
	for i := 0; i < 4; i++ {
		for j := 0; j < duplicate; j++ {
			t1 := time.Now()
			BigV := StringToBig(value)
			result, err := tree.EquivalentSearch(BigV.ToInt().Bytes())
			if err != nil {
				fmt.Println("some errors happen", err.Error())
				fmt.Println("\n\n\n")
			}
			t2 := time.Now()
			if j == 0 {
				content += strconv.Itoa(len(result.ResultData))
			}
			content += ","
			content += strconv.FormatInt(t2.Sub(t1).Microseconds(), 10)
		}
		content += "\n"
		value += "0"
	}
	content += "\n\n\n"
	AppendToFileWithString(specificSavePath, content)

	k := int64(10)
	content = ""
	for i := 0; i < 4; i++ {
		for j := 0; j < duplicate; j++ {
			t1 := time.Now()
			results, err := tree.TopkVSearch(k)
			if err != nil {
				fmt.Println("some errors happen", err.Error())
				fmt.Println("\n\n\n")
			}
			t2 := time.Now()
			transNum := TransNumInResultDArray(results)
			if j == 0 {
				content += strconv.Itoa(len(results))
				content += ","
				content += strconv.FormatInt(transNum, 10)
			}
			content += ","
			content += strconv.FormatInt(t2.Sub(t1).Microseconds(), 10)
		}
		k = k * 10
		content += "\n"
	}
	content += "\n\n\n"
	AppendToFileWithString(topkSavePath, content)

	content = ""
	start := "10000000000000000" //16个0
	Bigstart := StringToBig(start)
	span := "100000000000000" //14个0
	for i := 0; i < 4; i++ {
		for j := 0; j < duplicate; j++ {
			Bigend := BigAdd(start, span)
			t1 := time.Now()
			//fmt.Println(Bigstart.ToInt().String(),Bigend.ToInt().String())
			results, err := tree.RangeSearch(Bigstart.ToInt().Bytes(), Bigend.ToInt().Bytes())
			if err != nil {
				fmt.Println("some errors happen", err.Error())
				fmt.Println("\n\n\n")
			}
			t2 := time.Now()
			transNum := TransNumInResultDArray(results)
			if j == 0 {
				content += strconv.Itoa(len(results))
				content += ","
				content += strconv.FormatInt(transNum, 10)
			}
			content += ","
			content += strconv.FormatInt(t2.Sub(t1).Microseconds(), 10)
		}
		content += "\n"
		span += "0"
	}
	content += "\n\n\n"
	AppendToFileWithString(rangeSavePath, content)

}
