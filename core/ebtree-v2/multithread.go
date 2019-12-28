package ebtree_v2

import (
	"fmt"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"runtime"
	"time"
)

var bc *core.BlockChain
var interval int
var blocksnum int

type Task struct {
	Id  int
	Err error
	f   func(id int) (TaskR, error)
}

func (task *Task) Do() (TaskR, error) {
	return task.f(task.Id)
}

type WorkerPool struct {
	PoolSize    int
	tasksSize   int
	tasksChan   chan Task
	resultsChan chan TaskR
	Results     func() []TaskR
}

func NewWorkerPool(tasks []Task, size int) *WorkerPool {
	tasksChan := make(chan Task, len(tasks))
	resultsChan := make(chan TaskR, len(tasks))
	for _, task := range tasks {
		tasksChan <- task
	}
	close(tasksChan)
	pool := &WorkerPool{PoolSize: size, tasksSize: len(tasks), tasksChan: tasksChan, resultsChan: resultsChan}
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

//todo 不需要for遍历所有的，可以设置阈值，每次拿出多少
func (pool *WorkerPool) results() []TaskR {
	results := make([]TaskR, pool.tasksSize)
	for i := 0; i < pool.tasksSize; i++ {
		results[i] = <- pool.resultsChan
	}
	return results
}


func dosomething(id int) (TaskR, error) {
	//single task repsonse
	var strps TaskR
	for i := 0; i < interval; i++ {
		block := bc.GetBlockByNumber(uint64((id*interval)+i))
		trans := block.Transactions()

		strps.TaskResult = simSortTrans(trans)
	}
	return strps, nil
}


//todo 对trans整合，合并重复value 并且排序。当然可以先排序，排序的过程中合并重复元素
func simSortTrans(trans types.Transactions) []ResultD {
	return nil
}

func Initial(outerbc *core.BlockChain, outinterval int, outblocksnum int) {
	bc = outerbc
	interval = outinterval
	blocksnum = outblocksnum
}

func GetAll() []TaskR {
	maxProces := runtime.NumCPU()
	if maxProces > 1 {
		maxProces -= 1
	}
	runtime.GOMAXPROCS(maxProces)

	tasknum := blocksnum / interval

	t := time.Now()

	tasks := make([]Task, tasknum)
	for i := 0; i < tasknum; i++ {
		tasks[i] = *new(Task)
		tasks[i].Id = i
		tasks[i].f = dosomething
	}

	pool := NewWorkerPool(tasks, maxProces*2)
	pool.Start()

	results := pool.Results()
	fmt.Printf("all tasks finished, timeElapsed: %f s\n", time.Now().Sub(t).Seconds())
	return results
}
