package ebtree_v2

import (
	"bufio"
	"fmt"
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
