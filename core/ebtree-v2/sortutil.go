package ebtree_v2

import (
	"errors"
	"fmt"
)

func minHeap(root int, end int, c *[]ResultD) {
	for {
		var child = 2*root + 1
		//判断是否存在child节点
		if child > end {
			break
		}
		//判断右child是否存在，如果存在则和另外一个同级节点进行比较
		if child+1 <= end && byteCompare(&(*c)[child].Value, &(*c)[child+1].Value) > 0 {
			child += 1
		}
		if byteCompare((&(*c)[root].Value), &(*c)[child].Value) > 0 {
			(*c)[root], (*c)[child] = (*c)[child], (*c)[root]
			root = child
		} else {
			break
		}
	}
}

//降序排序
func HeapSortAndMergeSame(c *[]ResultD) *[]ResultD {
	var n = len(*c) - 1
	for root := n / 2; root >= 0; root-- {
		minHeap(root, n, c)
	}
	//fmt.Println("堆构建完成")
	for end := n; end >= 0; end-- {
		if byteCompare(&(*c)[0].Value, &(*c)[end].Value) < 0 {
			(*c)[0], (*c)[end] = (*c)[end], (*c)[0]
			minHeap(0, end-1, c)
		}
	}
	return mergeSamedata(c)
}

//heap sort response, 去重复
func mergeSamedata(array *[]ResultD) *[]ResultD {
	var hsrps []ResultD
	var size int
	pre := -1
	for i := 0; i < len(*array); i++ {
		if pre == -1 || byteCompare(&(*array)[i].Value, &(*array)[pre].Value) != 0 {
			hsrps = append(hsrps, (*array)[i])
			pre = i
			size++
		} else {
			hsrps[size-1].ResultData = append(hsrps[size-1].ResultData, (*array)[i].ResultData...)
		}
	}
	return &hsrps
}

func simplemerge(a, b *[]ResultD) *[]ResultD {
	//判断数组的长度
	al := len(*a)
	bl := len(*b)
	cl := al + bl
	c := make([]ResultD, cl)
	ai := 0
	bi := 0
	ci := 0

	for ai < al && bi < bl {
		if byteCompare(&(*a)[ai].Value, &(*b)[bi].Value) > 0 {
			c[ci] = (*a)[ai]
			ci++
			ai++
		} else {
			c[ci] = (*b)[bi]
			ci++
			bi++
		}
	}
	for ai < al {
		c[ci] = (*a)[ai]
		ci++
		ai++
	}
	for bi < bl {
		c[ci] = (*b)[bi]
		ci++
		bi++
	}
	return &c
}

//不再零散申请空间
func simplemergeV2(a *[]ResultD, sizea int, b *[]ResultD, sizeb int, c *[]ResultD, sizec int) int {
	rest := len(*c) - sizec
	if rest < (sizea + sizeb) {
		fmt.Println(sizea, sizeb, sizec)
		fmt.Println(len(*a), len(*b), len(*c))
		panic(errors.New("not enough from merge"))
	}

	ai := 0
	bi := 0
	ci := 0

	for ai < sizea && bi < sizeb {
		if byteCompare(&(*a)[ai].Value, &(*b)[bi].Value) > 0 {
			//判断是否重复
			if byteCompare(&(*c)[ci].Value, &(*a)[ai].Value) != 0 {
				(*c)[ci] = (*a)[ai]
				ci++
			} else {
				(*c)[ci].ResultData = append((*c)[ci].ResultData, (*a)[ai].ResultData...)
			}
			ai++
		} else {
			if byteCompare(&(*c)[ci].Value, &(*b)[bi].Value) != 0 {
				(*c)[ci] = (*b)[bi]
				ci++
			} else {
				(*c)[ci].ResultData = append((*c)[ci].ResultData, (*b)[bi].ResultData...)
			}
			bi++
		}
	}
	for ai < sizea {
		if byteCompare(&(*c)[ci].Value, &(*a)[ai].Value) != 0 {
			(*c)[ci] = (*a)[ai]
			ci++
		} else {
			(*c)[ci].ResultData = append((*c)[ci].ResultData, (*a)[ai].ResultData...)
		}
		ai++
	}
	for bi < sizeb {
		if byteCompare(&(*c)[ci].Value, &(*b)[bi].Value) != 0 {
			(*c)[ci] = (*b)[bi]
			ci++
		} else {
			(*c)[ci].ResultData = append((*c)[ci].ResultData, (*b)[bi].ResultData...)
		}
		bi++
	}
	return ci + 1
}

func mergeSortAndMergeSame(matrix *[]TaskR) *[]ResultD {
	if len(*matrix) <= 0 {
		panic(errors.New("not enough entity in taskR in mergesort"))
	}
	var length int
	for i := 0; i < len(*matrix); i++ {
		length += len((*matrix)[i].TaskResult)
	}

	b := make([]ResultD, length)
	c := make([]ResultD, length)
	b_point := &b
	c_point := &c

	var size int
	for i := 0; i < len(*matrix); i++ {
		size = simplemergeV2(&(*matrix)[i].TaskResult, len((*matrix)[i].TaskResult), b_point, size, c_point, 0)
		tmp := b_point
		b_point = c_point
		c_point = tmp
	}
	return mergeSamedata(b_point)
}
