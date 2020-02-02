package ebtree_v2

import (
	"fmt"
	"strconv"
	"testing"
)

func TestHeapSort(t *testing.T) {
	//arr := []int{3,4,3,4,1,6}
	//HeapSort(arr)
	//fmt.Println()
}
func TestMergerSort(t *testing.T) {
	merge([]int{2, 5, 7, 9}, 4, []int{1, 3, 6, 9}, 4)

}

func merge(nums1 []int, m int, nums2 []int, n int) {
	//把nums1复制到temp中
	temp := make([]int, m)
	copy(temp, nums1)

	t, j := 0, 0 //t为temp的索引，j为nums2的索引
	for i := 0; i < len(nums1); i++ {
		//当t大于temp的长度，那就是说temp全部放进去了nums1中，那剩下的就是放nums2剩余的值了
		if t >= len(temp) {
			nums1[i] = nums2[j]
			j++
			continue
		}
		//当j大于nums2的长度的时候，那就是说明nums2全部都放进去了nums1中，那剩下的就是放temp剩余的值了
		if j >= n {
			nums1[i] = temp[t]
			t++
			continue
		}
		//比较nums2与temp对应值的大小，小的那个就放进nums1中
		if nums2[j] <= temp[t] {
			nums1[i] = nums2[j]
			j++
		} else {
			nums1[i] = temp[t]
			t++
		}
	}
	fmt.Println(nums1)
}

func TestSortArry(t *testing.T) {
	fmt.Println("Hello World!")

	var a = []int{9, 7, 5, 3, 1, 0}
	//var a  []int
	var b = []int{8, 6, 4, 2, 0}

	c := sortArr(a, b)
	for i, v := range c {
		fmt.Println(i, ":", v)
	}
}

func sortArr(a, b []int) []int {

	//判断数组的长度
	al := len(a)
	bl := len(b)
	cl := al + bl

	fmt.Println(cl)
	//var c [cl]int // non-constant array bound cl
	c := make([]int, cl)

	fmt.Println(len(c))
	fmt.Println(cap(c))
	ai := 0
	bi := 0
	ci := 0

	for ai < al && bi < bl {

		if a[ai] > b[bi] {
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

	/*	for i, v := range c {
			fmt.Println(i, ":", v)
		}
	*/
	return c
}

func TestMergeFromFile(t *testing.T) {
	fileName1 := "/home/mimota/savetest" + strconv.Itoa(1) + "_" + strconv.Itoa(500000)
	fileName2 := "/home/mimota/savetest" + strconv.Itoa(500001) + "_" + strconv.Itoa(1000000)
	mergeFromTwoFiles(fileName1, fileName2, "/home/mimota/file3")
	//mergeFromTwoFiles("/home/mimota/file1.txt", "/home/mimota/file2.txt", "/home/mimota/file3.txt")
}
