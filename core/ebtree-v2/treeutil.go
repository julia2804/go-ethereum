package ebtree_v2

import (
	"bufio"
	"fmt"
	"github.com/ethereum/go-ethereum/core"
	"os"
	"strconv"
	"time"
)

func ConstructTree(outerbc *core.BlockChain, begin int, end int) (int, error) {

	/*
		cpuf, err := os.Create("cpu_profile")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(cpuf)
		defer pprof.StopCPUProfile()

	*/

	if treesize == 0 {
		treesize = 10000000
	}
	fmt.Println("treesize", treesize)
	nums := end - begin + 1
	if nums > treesize {
		var err error
		n1, err := constructTreeHelper(outerbc, begin, begin+treesize-1)
		n2, err := constructTreeHelper(nil, begin+treesize, end)
		return n1 + n2, err
	} else {
		return constructTreeHelper(outerbc, begin, end)
	}
}

func MergeFilesAndInsert(outerbc *core.BlockChain, dir string) (int, error) {
	Initial(outerbc, 0, 0)
	defer CloseParams()

	f, _ := os.OpenFile(recordPath, os.O_RDWR|os.O_CREATE, 0644)
	defer f.Close()

	t := time.Now()
	fileName := ReadDirAndMerge(dir)
	fmt.Println("final fileName :", fileName)
	fmt.Printf("merge finish, timeElapsed: %f s\n", time.Now().Sub(t).Seconds())
	AppendToFileWithStringByFile(f, strconv.FormatInt(time.Now().Sub(t).Milliseconds(), 10)+",")

	t1 := time.Now()
	//results := TestReadResultDs(fileName)
	var db *Database
	db = NewDatabase(*bc.GetDB())
	n, err := InsertToTreeWithDbByFile(fileName, db)
	//n, err := InsertToTree(trps)
	fmt.Printf("insert to ebtree, timeElapsed: %f s\n", time.Now().Sub(t1).Seconds())
	fmt.Println("dir ", dir)
	AppendToFileWithStringByFile(f, strconv.FormatInt(time.Now().Sub(t1).Milliseconds(), 10))
	AppendToFileWithStringByFile(f, "\n")

	AppendToFileWithStringByFile(f, "\n\n\n")
	return n, err
}

func constructTreeHelper(outerbc *core.BlockChain, begin int, end int) (int, error) {
	Initial(outerbc, begin, end)
	f, _ := os.OpenFile(recordPath, os.O_RDWR|os.O_CREATE, 0644)
	defer f.Close()

	AppendToFileWithStringByFile(f, strconv.Itoa(begin)+","+strconv.Itoa(end)+",")

	defer CloseParams()
	trps := GetTransAndSort(f)
	var fileName string
	fileName = constructSavePath + "save" + strconv.Itoa(begin) + "_" + strconv.Itoa(end)
	t1 := time.Now()
	WriteResultDArray(fileName, trps)
	//n := CountNum(fileName)
	//if(n != )
	fmt.Printf("write finished, timeElapsed: %f s\n", time.Now().Sub(t1).Seconds())
	AppendToFileWithStringByFile(f, strconv.FormatInt(time.Now().Sub(t1).Milliseconds(), 10))
	AppendToFileWithStringByFile(f, "\n")

	/*
		t := time.Now()
		//results := TestReadResultDs(fileName)
		var db *Database
		db = NewDatabase(*bc.GetDB())
		n, err := InsertToTreeWithDbByFile(fileName, db)
		//n, err := InsertToTree(trps)
		fmt.Printf("insert to ebtree, timeElapsed: %f s\n", time.Now().Sub(t).Seconds())

	*/

	return len(trps), nil
}

func AppendInsert(bc *core.BlockChain, begin int, end int) (int64, error) {
	Initial(bc, begin, end)
	defer CloseParams()
	f, _ := os.OpenFile(insert_begin_end_Path, os.O_RDWR|os.O_CREATE, 0644)
	defer f.Close()

	e := end - begin + 1
	var p float32
	p = 5
	var transNum int64
	db := NewDatabase(*bc.GetDB())
	tree, _ := NewEBTreeFromDb(db)
	t1 := time.Now()
	for i := begin; i <= end; i++ {
		block := bc.GetBlockByNumber(uint64(i))
		if block != nil {
			trans := block.Transactions()
			for j := 0; j < trans.Len(); j++ {
				tree.AfterInsertDataToTree(trans[j].Value().Bytes(), Convert2IdentifierData(i, j))
			}
			transNum += int64(trans.Len())
			per := float32(i) / float32(e) * 100
			if per >= p {
				fmt.Println("finish task ", per, "%")
				p = p + 5
			}
		}
	}
	err := tree.FinalCollapse()
	if err != nil {
		fmt.Println(err.Error())
		return transNum, err
	}
	err = tree.CommitMeatas()
	if err != nil {
		fmt.Println(err.Error())
		return transNum, err
	}
	fmt.Printf("after insert, timeElapsed: %d ms\n", time.Now().Sub(t1).Milliseconds())
	AppendToFileWithStringByFile(f, strconv.Itoa(begin)+",")
	AppendToFileWithStringByFile(f, strconv.Itoa(end)+",")
	AppendToFileWithStringByFile(f, strconv.FormatInt(time.Now().Sub(t1).Milliseconds(), 10)+",")
	AppendToFileWithStringByFile(f, strconv.FormatInt(transNum, 10)+"\n")
	return transNum, nil
}

func InsertToTreeWithDb(trps []ResultD, db *Database) (int, error) {
	//results := mergeSortAndMergeSame(trps)
	tree, err := NewEBTreeFromDb(db)
	Pool = CreatPoolAndRun(tree, insertthreadnum, insertbuffer)
	defer Pool.Close()
	err = tree.Inserts(trps)
	err = tree.FinalCollapse()
	if err != nil {
		return len(trps), err
	}
	err = tree.CommitMeatas()
	if err != nil {
		return len(trps), err
	}

	return len(trps), err
}

func InsertToTreeWithDbByFile(fileName string, db *Database) (int, error) {
	tree, err := NewEBTreeFromDb(db)
	Pool = CreatPoolAndRun(tree, insertthreadnum, insertbuffer)
	defer Pool.Close()

	f, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	defer f.Close()
	reader := bufio.NewReader(f)

	pi := MaxLeafNodeCapability * 16

	var number int
	for {
		datas := ReadResultDs(reader, pi)
		number += len(datas)
		if len(datas) > 0 {
			tree.Inserts(datas)
		}
		if len(datas) < (pi) {
			//tree.Inserts(datas[:num])
			//tree.InsertDatasToTree(datas[:num])
			break
		}
		//tree.InsertDatasToTree(datas)
	}
	err = tree.FinalCollapse()
	if err != nil {
		return number, err
	}
	err = tree.CommitMeatas()
	if err != nil {
		return number, err
	}

	return number, err
}

func TestInsertToTreeWithDbByFile(fileName string, db *Database) (int, error) {
	f, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	defer f.Close()
	//reader := bufio.NewReader(f)
	var final []ResultD
	datas := make([]ResultD, MaxLeafNodeCapability)

	//var number int
	for {
		//num := ReadResultDs(reader, MaxLeafNodeCapability, &datas)
		//number += num
		//if num < MaxLeafNodeCapability {
		//	//InsertToTreeWithDb(datas[:num], db)
		//	final = append(final, datas[:num]...)
		//	break
		//}
		//InsertToTreeWithDb(datas, db)
		final = append(final, datas...)
	}
	fmt.Println("final", len(final))
	InsertToTreeWithDb(final, db)
	return 0, nil
}

/*
func InsertToTree(trps []TaskR) (int, error) {
	results := mergeSortAndMergeSame(trps)
	tree, err := NewEBTree()
	err = tree.InsertDatasToTree(results)

	topkrps, err := tree.TopkVSearch(100000000)
	if err != nil {
		fmt.Println(err)
		return len(results), err
	}
	compareResult(results, topkrps)
	fmt.Println("topk num : ", len(topkrps))

	return len(results), err
}

*/

func TestResult(n int, array []ResultD) {
	var tmp int

	if len(array) >= n {
		tmp = n
	} else {
		tmp = len(array)
	}

	for i := 0; i < tmp; i++ {
		fmt.Print(i, ",", len((array)[i].Value), ":")
		fmt.Println((array)[i])
	}
}

func compareResult(array1, array2 []ResultD) {
	if len(array1) != len(array2) {
		fmt.Println("array1 is not array2 : the length")
		fmt.Println(len(array1), len(array2))
	} else {
		for i := 0; i < len(array1); i++ {
			if !ResultDIsSame((array1)[i], (array2)[i]) {
				fmt.Println("array1 is not array2, i : ", i)
				fmt.Println((array1)[i], (array2)[i])
			}

			if i != 0 {
				r := ResultCompare((array1)[i], (array2)[i-1])
				if r > 0 {
					fmt.Println("array is ascend, i :", i)
					fmt.Println((array1)[i], (array2)[i-1])
				}
			}
		}
	}
}

func TransNumInResultDArray(array []ResultD) int64 {
	var num int64
	for i := 0; i < len(array); i++ {
		num += int64(len((array)[i].ResultData))
	}
	return num
}

func GetTransAndSort(file *os.File) []ResultD {
	fmt.Println("get, merge, sort start")
	t := time.Now()

	fmt.Println("getblocks from db start")
	t3 := time.Now()
	prepool := AssembleTaskAndStart(gettasknum, getthreadnum, ToChannel, nil)
	results := prepool.Results(gettasknum)
	fmt.Printf("getblocks finished, timeElapsed: %f s\n", time.Now().Sub(t3).Seconds())
	AppendToFileWithStringByFile(file, strconv.FormatInt(time.Now().Sub(t3).Milliseconds(), 10)+",")

	var length int
	for i := 0; i < len(results); i++ {
		length += len((results)[i].TaskResult)
	}

	fmt.Println("merge tasks start")
	t1 := time.Now()
	data := make([]ResultD, length)
	var size int
	for i := 0; i < len(results); i++ {
		copy(data[size:], (results)[i].TaskResult)
		size += len((results)[i].TaskResult)
	}
	fmt.Printf("merge tasks finished, timeElapsed: %f s\n", time.Now().Sub(t1).Seconds())
	AppendToFileWithStringByFile(file, strconv.FormatInt(time.Now().Sub(t1).Milliseconds(), 10)+",")

	fmt.Println("heapsort start")
	t2 := time.Now()
	trps := HeapSortAndMergeSame(data)
	fmt.Printf("heapsort finished, timeElapsed: %f s\n", time.Now().Sub(t2).Seconds())
	AppendToFileWithStringByFile(file, strconv.FormatInt(time.Now().Sub(t2).Milliseconds(), 10)+",")

	//takenum = gettasknum / aftertasknum
	//afterpool := AssembleTaskAndStart(aftertasknum, afterthreadnum, FromChannel, prepool)

	//trps := afterpool.Results(aftertasknum)

	fmt.Printf("get, merge, sort finished, timeElapsed: %f s\n", time.Now().Sub(t).Seconds())
	AppendToFileWithStringByFile(file, strconv.FormatInt(time.Now().Sub(t).Milliseconds(), 10)+",")

	return trps
}
