package ebtree_v2

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/core"
	"os"
	"runtime"
	"time"
)

var bc *core.BlockChain
var interval int
var blocksnum int

var(
	pretasknum   = 100
	aftertasknum = 10

	prethreadnum = 10
	afterthreadnum = 10

	//每次从channel中几个数组来排序
	takenum  = 10
)

type Task struct {
	Id  int
	Err error
	Prepool *WorkerPool
	f   func(id int, prepool *WorkerPool) (TaskR, error)
}

func (task *Task) Do() (TaskR, error) {
	return task.f(task.Id, task.Prepool)
}

type WorkerPool struct {
	PoolSize    int
	tasksSize   int
	tasksChan   chan Task
	resultsChan chan TaskR
	Results     func(e int) []TaskR
}

func NewWorkerPool(tasks []Task, poolsize int) *WorkerPool {
	tasksChan := make(chan Task, len(tasks))
	resultsChan := make(chan TaskR, len(tasks))
	for _, task := range tasks {
		tasksChan <- task
	}
	close(tasksChan)
	pool := &WorkerPool{PoolSize: poolsize, tasksSize: len(tasks), tasksChan: tasksChan, resultsChan: resultsChan}
	pool.Results = pool.results
	return pool
}

func (pool *WorkerPool) Start() {
	for i := 0; i < pool.PoolSize; i++ {
		go pool.worker()
	}
}

func (pool *WorkerPool) worker() {
	for task := range pool.tasksChan {
		re, _ := task.Do()
		pool.resultsChan <- re
	}
}

func (pool *WorkerPool) results(e int) []TaskR {
	results := make([]TaskR, e)
	for i := 0; i < e; i++ {
		results[i] = <- pool.resultsChan
	}
	return results
}


func ToChannel(id int, prepool *WorkerPool) (TaskR, error) {
	//single task repsonse
	var strps TaskR
	var tmprps []ResultD

	for i := 0; i < interval; i++ {
		block := bc.GetBlockByNumber(uint64((id*interval)+i))
		trans := block.Transactions()

		for j:=0; j < trans.Len(); j++{
			var tmprd ResultD
			tmprd.Value = trans[j].Value().Bytes()
			var tmptd TD
			tmptd.IdentifierData = Convert2IdentifierData(i, j)
			tmprd.ResultData = append(tmprd.ResultData, tmptd)
			tmprps = append(tmprps, tmprd)
		}
	}
	strps.TaskResult = HeapSort(tmprps)
	return strps, nil
}

func FromChannel(id int, prepool *WorkerPool) (TaskR, error){
	var rps TaskR
	trps := prepool.Results(takenum)
	//把这些排序
	rps.TaskResult = mergeSort(trps)
	return rps, nil
}

func Initial(outerbc *core.BlockChain, outblocksnum int) {
	bc = outerbc
	blocksnum = outblocksnum
	interval = blocksnum / pretasknum
}

func GetAll() []TaskR {
	maxProces := runtime.NumCPU()
	if maxProces > 1 {
		maxProces -= 1
	}
	runtime.GOMAXPROCS(maxProces)

	t := time.Now()

	pretasks := make([]Task, pretasknum)
	for i := 0; i < pretasknum; i++ {
		pretasks[i] = *new(Task)
		pretasks[i].Id = i
		pretasks[i].f = ToChannel
	}

	prepool := NewWorkerPool(pretasks, prethreadnum)
	prepool.Start()


	takenum = pretasknum / aftertasknum
	aftertasks := make([]Task, aftertasknum)
	for i := 0; i < aftertasknum; i++{
		aftertasks[i] = *new(Task)
		//aftertasks[i].Id = i
		aftertasks[i].Prepool = prepool
		aftertasks[i].f = FromChannel
	}

	afterpool := NewWorkerPool(aftertasks, afterthreadnum)
	afterpool.Start()

	results := afterpool.Results(aftertasknum)
	fmt.Printf("all pretasks finished, timeElapsed: %f s\n", time.Now().Sub(t).Seconds())


	var f *os.File
	f, err := os.Create("/home/mimota/sss.txt")
	b, err := json.Marshal (results)
	if err != nil {
		fmt. Println ( "error:" , err )
	}
	f.Write(b)
	return results
}
