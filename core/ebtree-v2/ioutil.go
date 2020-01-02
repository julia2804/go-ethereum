package ebtree_v2

import (
	"fmt"
	"io/ioutil"
	"os"
)

func ReadDir(dirPath string) int64 {
	var dirSize int64
	flist, _ := ioutil.ReadDir(dirPath)
	for _, f := range flist {
		if f.IsDir() {
			dirSize = ReadDir(dirPath+"/"+f.Name()) + dirSize
		} else {
			dirSize = f.Size() + dirSize
		}
	}
	return dirSize
}

func WriteFile(name string, content []byte) {
	if ioutil.WriteFile(name, content, 0644) == nil {
		fmt.Println("WriteFile")
	} else {
		fmt.Println("NOT WriteFile")
	}
}

func getFileSize(filename string) int64 {
	fileInfo, _ := os.Stat(filename)
	//文件大小
	filesize := fileInfo.Size()
	return filesize
}

func AppendToFileWithString(fileName string, content string) {
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		// 打开文件失败处理
		fmt.Println("append to file error")
	} else {
		// 查找文件末尾的偏移量
		n, _ := f.Seek(0, 2)

		// 从末尾的偏移量开始写入内容
		_, err = f.WriteAt([]byte(content), n)
	}
	defer f.Close()
}

func AppendToFileWithByte(fileName string, content []byte) {
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		// 打开文件失败处理
		fmt.Println("append to file error")
	} else {
		// 查找文件末尾的偏移量
		n, _ := f.Seek(0, 2)

		// 从末尾的偏移量开始写入内容
		_, err = f.WriteAt(content, n)
	}
	defer f.Close()
}
