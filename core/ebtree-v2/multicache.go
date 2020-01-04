package ebtree_v2

import (
	"github.com/ethereum/go-ethereum/log"
)

type WorkerPool2 struct {
	CacheChan chan EBTreen
	PoolSize  int
	ebt       *EBTree
}

func NewWorkerPool2(ebt *EBTree, threadSize int, buffersize int) *WorkerPool2 {
	cache := make(chan EBTreen, buffersize)
	pool := &WorkerPool2{CacheChan: cache, PoolSize: threadSize, ebt: ebt}
	return pool
}

func (pool *WorkerPool2) Start() {
	for i := 0; i < pool.PoolSize; i++ {
		go pool.consumer(pool.CacheChan)
	}
}

func (pool *WorkerPool2) worker() {
	for {
		tmp := <-pool.CacheChan
		pool.ebt.CommitNode(tmp)
	}
}

func (pool *WorkerPool2) consumer(ch chan EBTreen) {
	//可以循环 for i := range ch 来不断从 channel 接收值，直到它被关闭。
	for node := range ch {
		pool.ebt.CommitNode(node)
	}
}

func (pool *WorkerPool2) Close() {
	close(pool.CacheChan)
	log.Info("cache channel closed")
}

func CreatPoolAndRun(ebt *EBTree, threadSize int, buffersize int) *WorkerPool2 {
	pool := NewWorkerPool2(ebt, 10, 10)
	pool.Start()
	return pool
}
