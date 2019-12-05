package EBTree

import (
	"encoding/binary"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
	"strconv"
	"time"
)

var specificValueSearchTime int64
var specificValueSearchNum int64

func SpecificValueSearch(root []byte, db *Database, v *hexutil.Big, bn uint64) (SearchValue, int64, error) {
	t1 := time.Now()
	//fmt.Print("Specific Value search :")
	//fmt.Println(v.ToInt().Bytes())
	tree, err := New(root, db)

	var buf2 = make([]byte, 8)
	binary.BigEndian.PutUint64(buf2, bn)

	data, err := tree.SpecificValueSearch(v.ToInt().Bytes(), buf2)
	//fmt.Println("specific search data num:", len(data.Data))
	t2 := time.Now()
	t3 := t2.Sub(t1).Microseconds()
	specificValueSearchTime = specificValueSearchTime + t3
	specificValueSearchNum++

	var tmp SearchValue
	return tmp, int64(len(data.Data)), err
}

func SpecificValueSearchTime() {
	fmt.Println("SpecificValueSearchTime:", specificValueSearchTime, "us")
	fmt.Println("times：", specificValueSearchNum)
}

func ClearSpecificValueSearchTime() {
	specificValueSearchTime = 0
	specificValueSearchNum = 0
	//fmt.Println("cleared SpecificValueSearchTime")
}

var topkVSearchTotalTime int64
var topkVSearchNum int64

func TopkVSearch(root []byte, db *Database, k uint64, bn uint64) ([]SearchValue, int64, int64, error) {
	t1 := time.Now()
	//fmt.Print("top k search :")
	//fmt.Println(k)

	var buf1 = make([]byte, 8)
	binary.BigEndian.PutUint64(buf1, k)

	var buf2 = make([]byte, 8)
	binary.BigEndian.PutUint64(buf2, bn)

	tree, err := New(root, db)
	su, result, err := tree.TopkVSearch(buf1, buf2, true)
	if err != nil {
		fmt.Printf("something wrong in topk  search with error")
		return nil, 0, 0, err
	}
	if !su {
		//fmt.Println("something wrong in topk  search without error")
	}
	//fmt.Println("we totally find", len(result), "data")
	sum := int64(0)
	for i := 0; i < len(result); i++ {
		for j := 0; j < len(result[i].Data); j++ {
			sum++
		}
	}
	//fmt.Println("the total transactions number:", sum)
	t2 := time.Now()
	t3 := t2.Sub(t1).Microseconds()
	topkVSearchTotalTime = topkVSearchTotalTime + t3
	topkVSearchNum++
	return nil, int64(len(result)), sum, err
}

func TopkVSearchTime() {
	fmt.Println("topkVSearchTotalTime:", topkVSearchTotalTime, "us")
	fmt.Println("times：", topkVSearchNum)
}

func ClearTopkVSearchTime() {
	topkVSearchTotalTime = 0
	topkVSearchNum = 0
	//fmt.Println("cleared topkVSearchTotalTime")
}

var rangeVSearchTotalTime int64
var rangeVSearchNum int64

func RangeVSearch(root []byte, db *Database, begin *hexutil.Big, end *hexutil.Big, bn uint64) ([]SearchValue, int64, int64, error) {
	t1 := time.Now()
	//fmt.Print("starting range search : ")
	//fmt.Print(begin.ToInt().Bytes())
	//fmt.Print("--->")
	//fmt.Println(end.ToInt().Bytes())

	tree, err := New(root, db)
	if err != nil {
		return nil, 0, 0, err
	}

	var buf3 = make([]byte, 8)
	binary.BigEndian.PutUint64(buf3, bn)

	su, result, err := tree.RangeValueSearch(begin.ToInt().Bytes(), end.ToInt().Bytes(), buf3)
	if err != nil {
		fmt.Printf("something wrong in range search with error")
		return nil, 0, 0, err
	}
	if !su {
		fmt.Println("something wrong in range search without error")
	}
	//fmt.Println("range search num :", len(result))
	sum := int64(0)
	for i := 0; i < len(result); i++ {
		for j := 0; j < len(result[i].Data); j++ {
			sum++
		}
	}
	//fmt.Println("the total transactions number:", sum)
	t2 := time.Now()
	t3 := t2.Sub(t1).Microseconds()
	rangeVSearchTotalTime = rangeVSearchTotalTime + t3
	rangeVSearchNum++
	return nil, int64(len(result)), sum, err
}

func RangeVSearchTime() {
	fmt.Println("rangeVSearchTotalTime:", rangeVSearchTotalTime, "us")
	fmt.Println("times: ", rangeVSearchNum)
}

func ClearRangeVSearchTime() {
	rangeVSearchTotalTime = 0
	rangeVSearchNum = 0
	//fmt.Println("cleared rangeVSearchTotalTime")
}

var topkpath string
var rangepath string
var specificpath string

var addToSearchValuepath string

func ClearAddToSearchValueTime() {
	addToSearchValueTime = 0
}

func ExperStart(bn uint64, root []byte, db *Database) {
	if len(topkpath) == 0 {
		topkpath = ethereum.GetValueFromDefaultPath("experiment", "topkpath")
		AppendToFile(topkpath, time.Now().Format("2006-01-02 15:04:05")+"\n")
	}

	if len(rangepath) == 0 {
		rangepath = ethereum.GetValueFromDefaultPath("experiment", "rangepath")
		AppendToFile(rangepath, time.Now().Format("2006-01-02 15:04:05")+"\n")
	}

	if len(specificpath) == 0 {
		specificpath = ethereum.GetValueFromDefaultPath("experiment", "specificpath")
		AppendToFile(specificpath, time.Now().Format("2006-01-02 15:04:05")+"\n")
	}

	if len(addToSearchValuepath) == 0 {
		addToSearchValuepath = ethereum.GetValueFromDefaultPath("experiment", "addToSearchValuepath")
		AppendToFile(addToSearchValuepath, time.Now().Format("2006-01-02 15:04:05")+"\n")
	}

	if bn != 0 {
		bnStr := "block nums : "
		bnStr += strconv.FormatUint(bn, 10)
		bnStr += "\n"
		AppendToFile(topkpath, bnStr)
		AppendToFile(topkpath, "k, time(us), resultnum, sum\n")
		AppendToFile(rangepath, bnStr)
		AppendToFile(rangepath, "span, time(us), resultnum, sum\n")
		AppendToFile(specificpath, bnStr)
		AppendToFile(specificpath, "value, time(us), sum\n")

		AppendToFile(addToSearchValuepath, bnStr)
		AppendToFile(addToSearchValuepath, "value, time(us), sum\n")
	}

	k := uint64(1)
	for i := 0; i < 8; i++ {
		ClearTopkVSearchTime()
		var content string
		_, resultNum, sum, _ := TopkVSearch(root, db, k, uint64(0))
		content += strconv.FormatUint(k, 10)
		content += ","
		content += strconv.FormatInt(topkVSearchTotalTime, 10)
		content += ","
		content += strconv.FormatInt(resultNum, 10)
		content += ","
		content += strconv.FormatInt(sum, 10)
		content += "\n"
		AppendToFile(topkpath, content)
		k = k * 10
	}
	AppendToFile(topkpath, "\n")

	start := "10000000000000000"
	Intstart, _ := new(big.Int).SetString(start, 10)
	var Bigstart hexutil.Big
	Bigstart = hexutil.Big(*Intstart)
	span := "10000000000000000"
	for i := 0; i < 8; i++ {
		ClearRangeVSearchTime()
		var content string

		Bigend := BigAbs(start, span)
		_, resultNum, sum, _ := RangeVSearch(root, db, &Bigstart, &Bigend, uint64(0))

		content += span
		content += ","
		content += strconv.FormatInt(rangeVSearchTotalTime, 10)
		content += ","
		content += strconv.FormatInt(resultNum, 10)
		content += ","
		content += strconv.FormatInt(sum, 10)
		content += "\n"
		AppendToFile(rangepath, content)

		span += "0"
	}
	AppendToFile(rangepath, "\n")

	value := "10000000000000000"
	for i := 0; i < 6; i++ {
		ClearSpecificValueSearchTime()
		ClearAddToSearchValueTime()
		var content string

		BigV := StringToBig(value)
		_, sum, _ := SpecificValueSearch(root, db, &BigV, uint64(10000000))

		content += value
		content += ","
		content += strconv.FormatInt(specificValueSearchTime, 10)
		content += ","
		content += strconv.FormatInt(sum, 10)
		content += "\n"
		AppendToFile(specificpath, content)

		content = ""
		content += value
		content += ","
		content += strconv.FormatInt(GetAddToSearch(), 10)
		content += "\n"
		AppendToFile(addToSearchValuepath, content)

		value += "0"
	}
	AppendToFile(specificpath, "\n")

}
