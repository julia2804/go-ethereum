package ebtree_v2

import (
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
		treesize = 4000000
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

func constructTreeHelper(outerbc *core.BlockChain, begin int, end int) (int, error) {
	Initial(outerbc, begin, end)
	defer CloseParams()
	GetTransAndSort()
	//var fileName string
	//fileName = "/home/mimota/savetest" + strconv.Itoa(begin) + "_" + strconv.Itoa(end) + ".txt"
	t1 := time.Now()
	//WriteResultDArray(fileName, trps)
	fmt.Printf("write finished, timeElapsed: %f s\n", time.Now().Sub(t1).Seconds())
	t := time.Now()
	//var db *Database
	//db = NewDatabase(*bc.GetDB())
	//n, err := InsertToTreeWithDb(trps, db)
	//n, err := InsertToTree(trps)
	fmt.Printf("insert to ebtree, timeElapsed: %f s\n", time.Now().Sub(t).Seconds())
	return 0, nil
}

func InsertToTreeWithDb(trps []ResultD, db *Database) (int, error) {
	//results := mergeSortAndMergeSame(trps)
	tree, err := NewEBTreeFromDb(db)
	Pool = CreatPoolAndRun(tree, insertthreadnum, insertbuffer)
	err = tree.Inserts(trps)
	err = tree.FinalCollapse()
	if err != nil {
		return len(trps), err
	}
	Pool.Close()
	err = tree.CommitMeatas()
	if err != nil {
		return len(trps), err
	}

	return len(trps), err
}

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
