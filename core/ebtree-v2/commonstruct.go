package ebtree_v2

import "os"

type TD struct {
	//IdentifierData=blockNo,txIndex
	IdentifierData []byte
}
type ResultD struct {
	//value stores the transaction value
	Value []byte
	//ResultData stores all IdentifierData of transactions whose value is equal to value
	ResultData []TD
}
type TaskR struct {
	//TaskResult stores all results returned by one task
	TaskResult []ResultD
}
type Entity struct {
	Value []byte
	Data  []byte
}
type EBCache struct {
	size int
	data []byte
}

func CloseCache(file *os.File, cache *EBCache) {
	if cache != nil {
		AppendToFileWithByteByFile(file, cache.data)
	}
}
