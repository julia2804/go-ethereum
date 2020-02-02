package ebtree_v2

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"os"
	"strconv"
)

func minHeap(root int, end int, c []ResultD) {
	for {
		var child = 2*root + 1
		//判断是否存在child节点
		if child > end {
			break
		}
		//判断右child是否存在，如果存在则和另外一个同级节点进行比较
		if child+1 <= end && byteCompare((c)[child].Value, (c)[child+1].Value) > 0 {
			child += 1
		}
		if byteCompare(((c)[root].Value), (c)[child].Value) > 0 {
			(c)[root], (c)[child] = (c)[child], (c)[root]
			root = child
		} else {
			break
		}
	}
}

//降序排序
func HeapSortAndMergeSame(c []ResultD) []ResultD {
	var n = len(c) - 1
	for root := n / 2; root >= 0; root-- {
		minHeap(root, n, c)
	}
	//fmt.Println("堆构建完成")
	for end := n; end >= 0; end-- {
		if byteCompare((c)[0].Value, (c)[end].Value) < 0 {
			(c)[0], (c)[end] = (c)[end], (c)[0]
			minHeap(0, end-1, c)
		}
	}
	return mergeSamedata(c)
}

//heap sort response, 去重复
func mergeSamedata(array []ResultD) []ResultD {
	var hsrps []ResultD
	var size int
	pre := -1
	for i := 0; i < len(array); i++ {
		if pre == -1 || byteCompare((array)[i].Value, (array)[pre].Value) != 0 {
			hsrps = append(hsrps, (array)[i])
			pre = i
			size++
		} else {
			hsrps[size-1].ResultData = append(hsrps[size-1].ResultData, (array)[i].ResultData...)
		}
	}
	return hsrps
}

func simplemerge(a, b []ResultD) *[]ResultD {
	//判断数组的长度
	al := len(a)
	bl := len(b)
	cl := al + bl
	c := make([]ResultD, cl)
	ai := 0
	bi := 0
	ci := 0

	for ai < al && bi < bl {
		if byteCompare((a)[ai].Value, (b)[bi].Value) > 0 {
			c[ci] = (a)[ai]
			ci++
			ai++
		} else {
			c[ci] = (b)[bi]
			ci++
			bi++
		}
	}
	for ai < al {
		c[ci] = (a)[ai]
		ci++
		ai++
	}
	for bi < bl {
		c[ci] = (b)[bi]
		ci++
		bi++
	}
	return &c
}

//不再零散申请空间
func simplemergeV2(a []ResultD, sizea int, b []ResultD, sizeb int, c []ResultD, sizec int) int {
	rest := len(c) - sizec
	if rest < (sizea + sizeb) {
		fmt.Println(sizea, sizeb, sizec)
		fmt.Println(len(a), len(b), len(c))
		panic(errors.New("not enough from merge"))
	}

	ai := 0
	bi := 0
	ci := 0

	for ai < sizea && bi < sizeb {
		if byteCompare((a)[ai].Value, (b)[bi].Value) > 0 {
			//判断是否重复
			if byteCompare((c)[ci].Value, (a)[ai].Value) != 0 {
				(c)[ci] = (a)[ai]
				ci++
			} else {
				(c)[ci].ResultData = append((c)[ci].ResultData, (a)[ai].ResultData...)
			}
			ai++
		} else {
			if byteCompare((c)[ci].Value, (b)[bi].Value) != 0 {
				(c)[ci] = (b)[bi]
				ci++
			} else {
				(c)[ci].ResultData = append((c)[ci].ResultData, (b)[bi].ResultData...)
			}
			bi++
		}
	}
	for ai < sizea {
		if byteCompare((c)[ci].Value, (a)[ai].Value) != 0 {
			(c)[ci] = (a)[ai]
			ci++
		} else {
			(c)[ci].ResultData = append((c)[ci].ResultData, (a)[ai].ResultData...)
		}
		ai++
	}
	for bi < sizeb {
		if byteCompare((c)[ci].Value, (b)[bi].Value) != 0 {
			(c)[ci] = (b)[bi]
			ci++
		} else {
			(c)[ci].ResultData = append((c)[ci].ResultData, (b)[bi].ResultData...)
		}
		bi++
	}
	return ci + 1
}

func mergeSortAndMergeSame(matrix []TaskR) []ResultD {
	if len(matrix) <= 0 {
		panic(errors.New("not enough entity in taskR in mergesort"))
	}
	var length int
	for i := 0; i < len(matrix); i++ {
		length += len((matrix)[i].TaskResult)
	}

	b := make([]ResultD, length)
	c := make([]ResultD, length)

	var size int
	for i := 0; i < len(matrix); i++ {
		size = simplemergeV2((matrix)[i].TaskResult, len((matrix)[i].TaskResult), b, size, b, 0)
		tmp := b
		b = c
		c = tmp
	}
	return mergeSamedata(b)
}

func mergeFromTwoFiles(fileName1 string, fileName2 string, fileName3 string) {
	var array1length = 10
	var array2length = 10
	file1, _ := os.OpenFile(fileName1, os.O_RDWR|os.O_CREATE, 0644)
	file2, _ := os.OpenFile(fileName2, os.O_RDWR|os.O_CREATE, 0644)
	file3, _ := os.OpenFile(fileName3, os.O_RDWR|os.O_CREATE, 0644)
	defer file1.Close()
	defer file2.Close()
	defer file3.Close()

	reader1 := bufio.NewReader(file1)
	reader2 := bufio.NewReader(file2)

	array1 := make([]Entity, array1length)
	var index1 int = array1length
	array2 := make([]Entity, array2length)
	var index2 int = array2length

	var cache EBCache
	cache.size = 1024 * 16
	defer CloseCache(file3, &cache)
	//cache.data = make([]byte, cache.size)

	for {
		if index1 >= array1length {
			num := ReadFile(reader1, array1length, &array1)
			if num == 0 {
				AppendEntityArrayToFileByFile(array2, index2, file3)
				AppendFileToFileByFile(file2, reader2, file3)
				break
			} else {
				index1 = array1length - num
			}
		}
		if index2 >= array2length {
			num := ReadFile(reader2, array2length, &array2)
			if num == 0 {
				AppendEntityArrayToFileByFile(array1, index1, file3)
				AppendFileToFileByFile(file1, reader1, file3)
				break
			} else {
				index2 = array2length - num
			}
		}

		r := byteCompare(array1[index1].Value, array2[index2].Value)
		if r > 0 {
			WriteEntityToFileWithCache((array1)[index1], file3, &cache)
			index1++
		} else if r < 0 {
			WriteEntityToFileWithCache((array2)[index2], file3, &cache)
			index2++
		} else {
			tds1, _ := DecodeTds((array1)[index1].Data)
			tds2, _ := DecodeTds((array2)[index2].Data)
			bys, _ := EncodeTds(append(tds1, tds2...))
			(array1)[index1].Data = bys
			WriteEntityToFileWithCache((array1)[index1], file3, &cache)
			index1++
			index2++
		}
	}
}

func mergeFromFileAndMen(array1 []Entity, fileName2 string, fileName3 string) {
	file2, _ := os.Open(fileName2)
	file3, _ := os.Open(fileName3)
	reader := bufio.NewReader(file2)
	var array1L int = len(array1)
	var array2L int = 10000
	var index1 int = 0
	var index2 int = array2L

	var cache EBCache
	cache.size = 1024 * 16
	cache.data = make([]byte, cache.size)

	array2 := make([]Entity, array2L)
	for {
		if index2 >= array2L {
			num := ReadFile(reader, array2L, &array2)
			if num == 0 {
				AppendEntityArrayToFileByFile(array1, index1, file2)
				break
			} else {
				index2 = array2L - num
			}
		}
		if index1 >= array1L {
			AppendEntityArrayToFileByFile(array2, index2, file2)
			break
		}
		r := byteCompare((array1)[index1].Value, array2[index2].Value)
		if r > 0 {
			WriteEntityToFileWithCache((array1)[index1], file3, &cache)
			index1++
		} else if r < 0 {
			WriteEntityToFileWithCache((array2)[index2], file3, &cache)
			index2++
		} else {
			tds1, _ := DecodeTds((array1)[index1].Data)
			tds2, _ := DecodeTds((array2)[index2].Data)
			bys, _ := EncodeTds(append(tds1, tds2...))
			(array1)[index1].Data = bys
			WriteEntityToFileWithCache((array1)[index1], file3, &cache)
			index1++
			index2++
		}
	}

	if len(cache.data) != 0 {
		AppendToFileWithByteByFile(file3, cache.data)
	}
}

func MergeFromFiles(fileNames []string) string {
	var round int
	for len(fileNames) > 1 {
		round++
		var tmp []string
		var name string = "round_" + strconv.Itoa(round)
		if len(fileNames)%2 != 0 {
			tmp = append(tmp, fileNames[len(fileNames)-1])
		}
		for i := 0; i < len(fileNames)-1; i = i + 2 {
			fileName := homePath + name + "_" + strconv.Itoa(i)
			mergeFromTwoFiles(fileNames[i], fileNames[i+1], fileName)
			tmp = append(tmp, fileName)
		}
		fileNames = tmp
	}
	if len(fileNames) <= 0 {
		log.Error("files error")
	}
	return fileNames[0]

}
