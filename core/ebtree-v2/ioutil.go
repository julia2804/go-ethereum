package ebtree_v2

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

//var file1 *os.File
//var file2 *os.File
//
//var reader1 *bufio.Reader
//var reader2 *bufio.Reader

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

func AppendToFileWithStringByFile(file *os.File, content string) {

	// 查找文件末尾的偏移量
	n, _ := file.Seek(0, 2)

	// 从末尾的偏移量开始写入内容
	file.WriteAt([]byte(content), n)

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

func AppendToFileWithByteByFile(file *os.File, content []byte) {
	// 查找文件末尾的偏移量
	n, _ := file.Seek(0, 2)

	// 从末尾的偏移量开始写入内容
	file.WriteAt(content, n)
}

func WriteResultD(file *os.File, re ResultD, cache *EBCache) {
	var entity Entity
	entity.Value = re.Value
	entity.Data, _ = EncodeTds(re.ResultData)
	//fmt.Println("value", entity.Value)
	//fmt.Println("data", entity.Data)
	WriteEntityToFileWithCache(entity, file, cache)
}

func WriteResultDArray(fileName string, res []ResultD) {
	f, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	var cache EBCache
	cache.size = 1024 * 1024 * 16
	defer CloseCache(f, &cache)
	//cache.data = make([]byte, cache.size)
	for i := 0; i < len(res); i++ {
		WriteResultD(f, (res)[i], &cache)
	}
}

//func initial(fileName1 string, fileName2 string){
//	file1, _  = os.Open(fileName1)
//	file2, _  = os.Open(fileName2)
//	reader1 = bufio.NewReader(file1)
//	reader2 = bufio.NewReader(file2)
//}

func ReadEntitys(reader *bufio.Reader, num int, array *[]Entity) int {
	var length = len(*array)
	var i int
	for i = 0; i < num; i++ {
		entity := ReadOneEntity(reader)
		if entity.Data == nil {
			break
		}
		(*array)[i] = entity
	}
	if i != num {
		copy((*array)[length-i:], (*array)[:i])
	}
	return i
}

func ReadOneResultD(reader *bufio.Reader) ResultD {
	var result ResultD

	value, data, err := ReadHelper(reader)
	tds, _ := DecodeTds(data)
	if data != nil && err == nil {
		result.Value = value
		result.ResultData = tds
	}
	return result
}

func ReadOneEntity(reader *bufio.Reader) Entity {
	var entity Entity

	value, data, err := ReadHelper(reader)
	if data != nil && err == nil {
		entity.Value = value
		entity.Data = data
	}
	return entity
}

func ReadHelper(reader *bufio.Reader) ([]byte, []byte, error) {
	var num int
	var err error

	sizeArray := make([]byte, 4)
	num, err = io.ReadFull(reader, sizeArray)
	if num != 4 || err != nil {
		return nil, nil, err
	}
	size := BytesToInt2(sizeArray)

	value := make([]byte, size)
	num, err = io.ReadFull(reader, value)
	if num != size || err != nil {
		return nil, nil, err
	}

	sizeArray = make([]byte, 4)
	num, err = io.ReadFull(reader, sizeArray)
	if num != 4 || err != nil {
		return nil, nil, err
	}
	size = BytesToInt2(sizeArray)
	if size == 0 {
		return nil, nil, err
	}

	data := make([]byte, size)
	num, err = io.ReadFull(reader, data)
	if num != size || err != nil {
		return nil, nil, err
	}
	return value, data, nil
}

func ReadResultDs(reader *bufio.Reader, num int) []ResultD {
	var i int
	var results []ResultD
	for i = 0; i < num; i++ {
		result := ReadOneResultD(reader)
		if result.ResultData == nil {
			break
		}
		results = append(results, result)
	}
	return results
}

func TestReadResultDs(fileName string) []ResultD {
	file, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	defer file.Close()

	reader := bufio.NewReader(file)
	var results []ResultD
	for {
		result := ReadOneResultD(reader)
		if result.ResultData == nil {
			break
		}
		results = append(results, result)
	}
	return results
}

func AppendEntityArrayToFile(array *[]Entity, index int, fileName string) {
	file, _ := os.Open(fileName)
	for i := index; i < len(*array); i++ {
		AppendToFileWithByteByFile(file, (*array)[i].Value)
		AppendToFileWithStringByFile(file, ",")
	}
}

func AppendEntityArrayToFileByFile(array []Entity, index int, file *os.File) {
	for i := index; i < len(array); i++ {
		AppendToFileWithByteByFile(file, IntToBytes2(len((array)[i].Value)))
		AppendToFileWithByteByFile(file, (array)[i].Value)
		AppendToFileWithByteByFile(file, IntToBytes2(len((array)[i].Data)))
		AppendToFileWithByteByFile(file, (array)[i].Data)
	}
}

func AppendFileToFileByFile(src *os.File, reader *bufio.Reader, dst *os.File) {
	buf := make([]byte, 256)
	for {
		n, err := reader.Read(buf)
		if err == io.EOF {
			fmt.Println("read the file finished")
			break
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		AppendToFileWithByteByFile(dst, buf[:n])
	}
}

func WriteEntityToFileWithCache(entity Entity, file *os.File, cache *EBCache) {
	if len(entity.Data) != 0 {
		cache.data = append(cache.data, IntToBytes2(len(entity.Value))...)
		cache.data = append(cache.data, entity.Value...)

		cache.data = append(cache.data, IntToBytes2(len(entity.Data))...)
		cache.data = append(cache.data, entity.Data...)
	}

	if len(cache.data) >= cache.size {
		AppendToFileWithByteByFile(file, cache.data)
		cache.data = nil
	}
}

func CountNum(fileName string) int {
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err.Error())
	}
	reader := bufio.NewReader(f)
	var i int
	for {
		result := ReadOneResultD(reader)
		if result.ResultData == nil {
			break
		}
		i++
	}

	fmt.Println("count nums", i)
	return i
}

func ReadDirAndMerge(dirPath string) string {
	flist, _ := ioutil.ReadDir(dirPath)
	fileNames := make([]string, len(flist))
	var i int
	for _, f := range flist {
		fileNames[i] = dirPath + f.Name()
		i++
	}
	return MergeFromFiles(fileNames)
}
