package ebtree_v2

import (
	"bufio"
	"fmt"
	"github.com/ethereum/go-ethereum/core"
	"log"
	"os"
	"runtime/pprof"
	"time"
)

func ConstructTree(outerbc *core.BlockChain, begin int, end int) (int, error) {

	cpuf, err := os.Create("cpu_profile")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(cpuf)
	defer pprof.StopCPUProfile()

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
	t := time.Now()
	fileName := ReadDirAndMerge(dir)
	fmt.Println("final fileName", fileName)
	fmt.Printf("merge finish, timeElapsed: %f s\n", time.Now().Sub(t).Seconds())

	t1 := time.Now()
	//results := TestReadResultDs(fileName)
	var db *Database
	db = NewDatabase(*bc.GetDB())
	n, err := InsertToTreeWithDbByFile(fileName, db)
	//n, err := InsertToTree(trps)
	fmt.Printf("insert to ebtree, timeElapsed: %f s\n", time.Now().Sub(t1).Seconds())
	fmt.Println("dir ", dir)
	return n, err
}

func constructTreeHelper(outerbc *core.BlockChain, begin int, end int) (int, error) {
	Initial(outerbc, begin, end)
	defer CloseParams()
	trps := GetTransAndSort()

	t := time.Now()
	var db *Database
	db = NewDatabase(*bc.GetDB())
	n, err := InsertToTreeWithDb(trps, db)
	fmt.Printf("insert to ebtree, timeElapsed: %f s\n", time.Now().Sub(t).Seconds())

	return n, err
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
	datas := make([]ResultD, MaxLeafNodeCapability*1024)

	var number int
	var threod int
	threod = 1
	for {
		num := ReadResultDs(reader, len(datas), &datas)
		number += num
		if number/10000 > threod {
			fmt.Println("tmp number ", number)
			threod++
		}
		if num < len(datas) {
			tree.Inserts(datas[:num])
			//tree.InsertDatasToTree(datas[:num])
			break
		}
		tree.Inserts(datas)
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
	reader := bufio.NewReader(f)
	var final []ResultD
	datas := make([]ResultD, MaxLeafNodeCapability)

	var number int
	for {
		num := ReadResultDs(reader, MaxLeafNodeCapability, &datas)
		number += num
		if num < MaxLeafNodeCapability {
			//InsertToTreeWithDb(datas[:num], db)
			final = append(final, datas[:num]...)
			break
		}
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

func TransNumInResultDArray(array []ResultD) int {
	var num int
	for i := 0; i < len(array); i++ {
		num += len((array)[i].ResultData)
	}
	return num
}

func GetTransAndSort() []ResultD {
	fmt.Println("get, merge, sort start")
	t := time.Now()

	fmt.Println("getblocks from db start")
	t3 := time.Now()
	prepool := AssembleTaskAndStart(pretasknum, prethreadnum, ToChannel, nil)
	results := prepool.Results(pretasknum)
	fmt.Printf("getblocks finished, timeElapsed: %f s\n", time.Now().Sub(t3).Seconds())

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

	fmt.Println("heapsort start")
	t2 := time.Now()
	trps := HeapSortAndMergeSame(data)
	fmt.Printf("heapsort finished, timeElapsed: %f s\n", time.Now().Sub(t2).Seconds())

	//takenum = pretasknum / aftertasknum
	//afterpool := AssembleTaskAndStart(aftertasknum, afterthreadnum, FromChannel, prepool)

	//trps := afterpool.Results(aftertasknum)

	fmt.Printf("get, merge, sort finished, timeElapsed: %f s\n", time.Now().Sub(t).Seconds())

	return trps
}
