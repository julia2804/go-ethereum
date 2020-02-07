package ebtree_v2

import (
	"fmt"
	"github.com/ethereum/go-ethereum/log"
)

type WorkerPool2 struct {
	CacheChan   chan EBTreen
	PoolSize    int
	ebt         *EBTree
	resultsChan chan int
}

func NewWorkerPool2(ebt *EBTree, threadSize int, buffersize int) *WorkerPool2 {
	cache := make(chan EBTreen, buffersize)
	results := make(chan int, threadSize)
	pool := &WorkerPool2{CacheChan: cache, PoolSize: threadSize, ebt: ebt, resultsChan: results}
	return pool
}

func (pool *WorkerPool2) Start() {
	for i := 0; i < pool.PoolSize; i++ {
		go pool.consumer(pool.CacheChan)
	}
}

//func (pool *WorkerPool2) worker() {
//	for {
//		tmp := <-pool.CacheChan
//		pool.ebt.CommitNode(tmp)
//	}
//}

func (pool *WorkerPool2) consumer(ch chan EBTreen) {
	//可以循环 for i := range ch 来不断从 channel 接收值，直到它被关闭。
	batch := pool.ebt.Db.diskdb.NewBatch()
	for node := range ch {
		pool.ebt.CommitNode(node, batch)
	}
	if err := batch.Write(); err != nil {
		log.Error(err.Error())
	}
	batch.Reset()
	pool.resultsChan <- 1
}

func (pool *WorkerPool2) Close() {
	close(pool.CacheChan)
	pool.results(pool.PoolSize)
	close(pool.resultsChan)
	log.Info("cache channel closed")
}

func (pool *WorkerPool2) results(e int) {
	var f float32
	f = 5
	results := make([]int, e)
	for i := 0; i < e; i++ {
		results[i] = <-pool.resultsChan
		per := float32(i) / float32(e) * 100
		if per >= f {
			fmt.Println("finish task ", per, "%")
			f = f + 5
		}
	}
}

func CreatPoolAndRun(ebt *EBTree, threadSize int, buffersize int) *WorkerPool2 {
	pool := NewWorkerPool2(ebt, threadSize, buffersize)
	pool.Start()
	return pool
}
