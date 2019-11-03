package EBTree

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"os"
	"reflect"
)

func IntToBytes(i uint64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, i)
	return buf
}

func BytesToInt(b []byte) (i uint64) {
	return binary.BigEndian.Uint64(b)
}

func typeof(v interface{}) string {
	return reflect.TypeOf(v).String()
}

//compare a and b
func Compare(a, b []byte) int {
	lev, lemax := IndentBytes(a, b)
	return bytes.Compare(lev, lemax)
}

//used to indent the bytes
func IndentBytes(f []byte, s []byte) ([]byte, []byte) {
	var lev []byte
	var lemax []byte
	if len(f) < len(s) {
		t := len(s)

		b0 := byte(0)
		for i := 0; i < t-len(f); i++ {
			lev = append(lev, b0)
		}
		for j := 0; j < len(f); j++ {
			lev = append(lev, f[j])
		}
		for j := 0; j < len(s); j++ {
			lemax = append(lemax, s[j])
		}
		return lev, lemax
	} else if len(f) > len(s) {
		t := len(f)

		b0 := byte(0)
		for i := 0; i < t-len(s); i++ {
			lemax = append(lemax, b0)
		}
		for j := 0; j < len(s); j++ {
			lemax = append(lev, s[j])
		}
		for j := 0; j < len(f); j++ {
			lev = append(lev, f[j])
		}
		return lev, lemax
	} else {
		for j := 0; j < len(f); j++ {
			lev = append(lev, f[j])
		}
		for j := 0; j < len(s); j++ {
			lemax = append(lemax, s[j])
		}
		return lev, lemax
	}
}

func add(b []byte, i uint64) []byte {
	f := BytesToInt(b)
	return IntToBytes(f + i)
}

func minus(b []byte, i uint64) []byte {
	f := BytesToInt(b)
	return IntToBytes(f - i)
}

func getFileSize(filename string) int64{
	fileInfo, _ := os.Stat(filename)
	//文件大小
	filesize:= fileInfo.Size()
	return filesize
}

func readDir(dirPath string) int64 {
	var dirSize int64
	flist, _ := ioutil.ReadDir(dirPath)
	for _, f := range flist {
		if f.IsDir() {
			dirSize = readDir(dirPath+"/"+f.Name()) + dirSize
		} else {
			dirSize= f.Size() + dirSize
		}
	}
	return dirSize
}
