package ebtree_v2

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/log"
	"runtime"
	"strconv"
	"time"
)

var bc *core.BlockChain
var interval int
var blocksnum int

var (
	pretasknum, _   = strconv.Atoi(ethereum.GetValueFromDefaultPath("thread", "pretasknum"))
	aftertasknum, _ = strconv.Atoi(ethereum.GetValueFromDefaultPath("thread", "aftertasknum"))

	prethreadnum, _   = strconv.Atoi(ethereum.GetValueFromDefaultPath("thread", "prethreadnum"))
	afterthreadnum, _ = strconv.Atoi(ethereum.GetValueFromDefaultPath("thread", "afterthreadnum"))
)

var takenum int

type Task struct {
	Id      int
	Err     error
	Prepool *WorkerPool
	f       func(id int, prepool *WorkerPool) TaskR
}

func (task *Task) Do() TaskR {
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
		re := task.Do()
		pool.resultsChan <- re
	}
}

func (pool *WorkerPool) results(e int) []TaskR {
	results := make([]TaskR, e)
	for i := 0; i < e; i++ {
		results[i] = <-pool.resultsChan
	}
	return results
}

func ToChannel(id int, prepool *WorkerPool) TaskR {
	//single task repsonse
	var strps TaskR
	var tmprps []ResultD

	for i := 1; i <= interval; i++ {
		block := bc.GetBlockByNumber(uint64((id * interval) + i))
		if block != nil {
			trans := block.Transactions()
			for j := 0; j < trans.Len(); j++ {
				var tmprd ResultD
				tmprd.Value = trans[j].Value().Bytes()
				var tmptd TD
				tmptd.IdentifierData = Convert2IdentifierData(i, j)
				tmprd.ResultData = append(tmprd.ResultData, tmptd)
				tmprps = append(tmprps, tmprd)
			}
		}
	}
	strps.TaskResult = HeapSortAndMergeSame(tmprps)
	return strps
}

func FromChannel(id int, prepool *WorkerPool) TaskR {
	var rps TaskR
	trps := prepool.Results(takenum)
	//把这些排序
	rps.TaskResult = mergeSortAndMergeSame(trps)
	return rps
}

func Initial(outerbc *core.BlockChain, outblocksnum int) {
	maxProces := runtime.NumCPU()
	if maxProces > 1 {
		maxProces -= 1
	}
	runtime.GOMAXPROCS(maxProces)

	if pretasknum == 0 {
		pretasknum = 1
	}

	if prethreadnum == 0 {
		prethreadnum = 1
	}

	if aftertasknum == 0 {
		aftertasknum = 1
	}

	if afterthreadnum == 0 {
		afterthreadnum = 1
	}

	bc = outerbc
	blocksnum = outblocksnum
	interval = blocksnum / pretasknum
	log.Info("initial over, the final blocknum is :", interval*pretasknum)
}

func AssembleTaskAndStart(tasknum int, threadnum int, f func(id int, prepool *WorkerPool) TaskR, prepool *WorkerPool) *WorkerPool {
	tasks := make([]Task, tasknum)
	for i := 0; i < tasknum; i++ {
		tasks[i] = *new(Task)
		tasks[i].Id = i
		tasks[i].Prepool = prepool
		tasks[i].f = f
	}

	pool := NewWorkerPool(tasks, threadnum)
	pool.Start()
	return pool
}

func GetTrans() []TaskR {
	t := time.Now()

	prepool := AssembleTaskAndStart(pretasknum, prethreadnum, ToChannel, nil)

	takenum = pretasknum / aftertasknum
	afterpool := AssembleTaskAndStart(aftertasknum, afterthreadnum, FromChannel, prepool)

	trps := afterpool.Results(aftertasknum)

	fmt.Printf("all tasks finished, timeElapsed: %f s\n", time.Now().Sub(t).Seconds())

	b, err := json.Marshal(trps)
	if err != nil {
		fmt.Println("error:", err)
	}
	AppendToFileWithByte("/home/mimota/sss.txt", b)
	return trps
}

func InsertToTree(trps []TaskR) (int, error) {
	results := mergeSortAndMergeSame(trps)
	tree, err := NewEBTree()
	err = tree.InsertDatasToTree(results)
	return len(results), err
}

func ConstructTree(outerbc *core.BlockChain, outblocksnum int) (int, error) {
	Initial(outerbc, outblocksnum)
	trps := GetTrans()
	t := time.Now()
	n, err := InsertToTree(trps)
	fmt.Printf("insert to ebtree, timeElapsed: %f s\n", time.Now().Sub(t).Seconds())
	return n, err
}
