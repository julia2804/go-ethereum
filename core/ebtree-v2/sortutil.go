package ebtree_v2

func minHeap(root int, end int, c []ResultD) {
	for {
		var child = 2*root + 1
		//判断是否存在child节点
		if child > end {
			break
		}
		//判断右child是否存在，如果存在则和另外一个同级节点进行比较
		if child+1 <= end && byteCompare(c[child].Value, c[child+1].Value) > 0 {
			child += 1
		}
		if byteCompare((c[root].Value), c[child].Value) > 0 {
			c[root], c[child] = c[child], c[root]
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
		if byteCompare(c[0].Value, c[end].Value) < 0 {
			c[0], c[end] = c[end], c[0]
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
		if pre == -1 || byteCompare(array[i].Value, array[pre].Value) != 0 {
			hsrps = append(hsrps, array[i])
			pre = i
			size++
		} else {
			hsrps[size-1].ResultData = append(hsrps[size-1].ResultData, array[i].ResultData...)
		}
	}
	return hsrps
}

func simplemerge(a, b []ResultD) []ResultD {
	//判断数组的长度
	al := len(a)
	bl := len(b)
	cl := al + bl
	c := make([]ResultD, cl)
	ai := 0
	bi := 0
	ci := 0

	for ai < al && bi < bl {
		if byteCompare(a[ai].Value, b[bi].Value) > 0 {
			c[ci] = a[ai]
			ci++
			ai++
		} else {
			c[ci] = b[bi]
			ci++
			bi++
		}
	}
	for ai < al {
		c[ci] = a[ai]
		ci++
		ai++
	}
	for bi < bl {
		c[ci] = b[bi]
		ci++
		bi++
	}
	return c
}

func mergeSortAndMergeSame(matrix []TaskR) []ResultD {
	var rps []ResultD
	for i := 0; i < len(matrix); i++ {
		rps = simplemerge(rps, matrix[i].TaskResult)
	}
	return mergeSamedata(rps)
}
