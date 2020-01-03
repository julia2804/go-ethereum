package ebtree_v2

import (
	"fmt"
	"github.com/ethereum/go-ethereum/core"
	"time"
)

func ConstructTree(outerbc *core.BlockChain, outblocksnum int) (int, error) {

	//cpuf, err := os.Create("cpu_profile")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//pprof.StartCPUProfile(cpuf)
	//defer pprof.StopCPUProfile()

	Initial(outerbc, outblocksnum)
	trps := GetTrans()
	t := time.Now()
	n, err := InsertToTree(trps)
	fmt.Printf("insert to ebtree, timeElapsed: %f s\n", time.Now().Sub(t).Seconds())
	return n, err
}

func InsertToTree(trps []TaskR) (int, error) {
	results := mergeSortAndMergeSame(trps)
	tree, err := NewEBTree()
	err = tree.InsertDatasToTree(results)

	topkrps := tree.TopkVSearch(100000000)
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
		fmt.Print(i, ",", len(array[i].Value), ":")
		fmt.Println(array[i])
	}
}

func compareResult(array1, array2 []ResultD) {
	if len(array1) != len(array2) {
		fmt.Println("array1 is not array2 : the length")
		fmt.Println(len(array1), len(array2))
	} else {
		for i := 0; i < len(array1); i++ {
			if !ResultDIsSame(array1[i], array2[i]) {
				fmt.Println("array1 is not array2, i : ", i)
				fmt.Println(array1[i], array2[i])
			}

			if i != 0 {
				r := ResultCompare(array1[i], array1[i-1])
				if r > 0 {
					fmt.Println("array is ascend, i :", i)
					fmt.Println(array1[i], array1[i-1])
				}
			}
		}
	}
}
