package ebtree_v2

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestWrite(t *testing.T) {
	AppendToFileWithString("/home/mimota/test2.txt", "bbb")

}

func TestRead(t *testing.T) {
	//fileName := "/home/mimota/savetest" + strconv.Itoa(1) +"_" + strconv.Itoa(200000) + ".txt"
	//file, _  := os.Open(fileName)
	////scanner := bufio.NewScanner(file)
	//r := bufio.NewReader(file)
	//line, _ := r.ReadString(byte('\n'))
	//fmt.Println(line)

	fileName2 := "/home/mimota/save.txt"
	file2, _ := os.Open(fileName2)
	//scanner := bufio.NewScanner(file)
	r2 := bufio.NewReader(file2)
	line2, _ := r2.ReadString(byte('\n'))
	fmt.Println(line2)

}

func TestRead2(t *testing.T) {
	f, err := os.Open("/home/mimota/file1.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行

		if err != nil || io.EOF == err {
			break
		}
		fmt.Println(line)
	}
}

func TestRead3(t *testing.T) {
	f, err := os.Open("/home/mimota/file3.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	var num int
	var knum int
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行

		if err != nil || io.EOF == err {
			break
		}
		if len(line) != 0 {
			//fmt.Println(line)
			if len(line) == 1 {
				//fmt.Println("line", line)
				knum++

			}
			num++
		}
		//fmt.Println(line)
	}
	fmt.Println("knum", knum)
	fmt.Println("num", num)
}
func TestCount(t *testing.T) {
	//CountNum("/home/mimota/savetest1_500000")
	//CountNum("/home/mimota/savetest500001_1000000")
	//CountNum("/home/mimota/savetest1_2000000")
	//CountNum("/home/mimota/file3")
	CountNum("/home/mimota/round_2_0")

}

func TestReadDirAndMerge(t *testing.T) {
	fmt.Println(ReadDirAndMerge("/home/mimota/testread/"))
}
