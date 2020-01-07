package ebtree_v2

import (
	"fmt"
	"github.com/ethereum/go-ethereum/core"
	"log"
	"os"
	"runtime/pprof"
	"time"
)

func ConstructTree(outerbc *core.BlockChain, outblocksnum int) (int, error) {

	cpuf, err := os.Create("cpu_profile")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(cpuf)
	defer pprof.StopCPUProfile()

	Initial(outerbc, outblocksnum)
	trps := GetTrans()
	t := time.Now()
	var db *Database
	db = NewDatabase(*outerbc.GetDB())
	n, err := InsertToTreeWithDb(trps, db)
	//n, err := InsertToTree(trps)
	fmt.Printf("insert to ebtree, timeElapsed: %f s\n", time.Now().Sub(t).Seconds())

	return n, err
}

func InsertToTreeWithDb(trps *[]TaskR, db *Database) (int, error) {
	results := mergeSortAndMergeSame(trps)
	tree, err := NewEBTreeFromDb(db)
	Pool = CreatPoolAndRun(tree, 10, 10)
	err = tree.Inserts(*results)
	err = tree.FinalCollapse()
	if err != nil {
		return len(*results), err
	}
	Pool.Close()
	err = tree.CommitMeatas()
	if err != nil {
		return len(*results), err
	}

	return len(*results), err
}

func InsertToTree(trps *[]TaskR) (int, error) {
	results := mergeSortAndMergeSame(trps)
	tree, err := NewEBTree()
	err = tree.InsertDatasToTree(*results)

	topkrps, err := tree.TopkVSearch(100000000)
	if err != nil {
		fmt.Println(err)
		return len(*results), err
	}
	compareResult(results, &topkrps)
	fmt.Println("topk num : ", len(topkrps))

	return len(*results), err
}

func TestResult(n int, array *[]ResultD) {
	var tmp int

	if len(*array) >= n {
		tmp = n
	} else {
		tmp = len(*array)
	}

	for i := 0; i < tmp; i++ {
		fmt.Print(i, ",", len((*array)[i].Value), ":")
		fmt.Println((*array)[i])
	}
}

func compareResult(array1, array2 *[]ResultD) {
	if len(*array1) != len(*array2) {
		fmt.Println("array1 is not array2 : the length")
		fmt.Println(len(*array1), len(*array2))
	} else {
		for i := 0; i < len(*array1); i++ {
			if !ResultDIsSame(&(*array1)[i], &(*array2)[i]) {
				fmt.Println("array1 is not array2, i : ", i)
				fmt.Println((*array1)[i], (*array2)[i])
			}

			if i != 0 {
				r := ResultCompare(&(*array1)[i], &(*array2)[i-1])
				if r > 0 {
					fmt.Println("array is ascend, i :", i)
					fmt.Println((*array1)[i], (*array2)[i-1])
				}
			}
		}
	}
}

func TransNumInResultDArray(array *[]ResultD) int {
	var num int
	for i := 0; i < len(*array); i++ {
		num += len((*array)[i].ResultData)
	}
	return num
}
