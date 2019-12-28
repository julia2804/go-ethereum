package ebtree_v2

type TD struct {
	//IdentifierData=blockNo,txIndex
	IdentifierData []byte
}
type ResultD struct {
	//value stores the transaction value
	value int64
	//ResultData stores all IdentifierData of transactions whose value is equal to value
	ResultData []TD
}
type TaskR struct {
	//TaskResult stores all results returned by one task
	TaskResult []ResultD
}
