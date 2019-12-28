package ebtree_v2

import (
	"fmt"
	"github.com/ethereum/go-ethereum/core"
	"runtime"
	"time"
)

var bc *core.BlockChain
var interval int
var blocksnum int

type Task struct {
	Id  int
	Err error
	f   func(id int) ([]Data, error)
}

type Data struct {
	content []byte
}

func (task *Task) Do() ([]Data, error) {
	return task.f(task.Id)
}

type WorkerPool struct {
	PoolSize    int
	tasksSize   int
	tasksChan   chan Task
	resultsChan chan []Data
	Results     func() [][]Data
}

func NewWorkerPool(tasks []Task, size int) *WorkerPool {
	tasksChan := make(chan Task, len(tasks))
	resultsChan := make(chan []Data, len(tasks))
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

func (pool *WorkerPool) results() [][]Data {
	results := make([][]Data, pool.tasksSize)
	for i := 0; i < pool.tasksSize; i++ {
		results[i] = <-pool.resultsChan
	}
	return results
}

func get(i int, bc *core.BlockChain) []Data {
	b := bc.GetBlockByNumber(uint64(i))
	bts := b.Transactions()

	dts := make([]Data, bts.Len())
	for j := 0; j < len(bts); j++ {
		dts[j].content = bts[j].Value().Bytes()
	}
	return dts
}

func dosomething(id int) ([]Data, error) {
	var rs []Data
	for i := 0; i < interval; i++ {
		tmp := get((id*interval)+i, bc)
		rs = append(rs, tmp...)
	}
	return rs, nil
}

func Initial(outerbc *core.BlockChain, outinterval int, outblocksnum int) {
	bc = outerbc
	interval = outinterval
	blocksnum = outblocksnum
}

func GetAll() [][]Data {
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
	//for _, datalist := range results {
	//	fmt.Printf("Data of task is %v\n", datalist)
	//}
}
