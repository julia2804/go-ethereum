package ebtree_v2

import (
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/log"
	"runtime"
	"strconv"
)

var bc *core.BlockChain
var interval int
var begin int
var end int

var takenum int

var (
	pretasknum, _   = strconv.Atoi(ethereum.GetValueFromDefaultPath("thread", "pretasknum"))
	aftertasknum, _ = strconv.Atoi(ethereum.GetValueFromDefaultPath("thread", "aftertasknum"))

	prethreadnum, _    = strconv.Atoi(ethereum.GetValueFromDefaultPath("thread", "prethreadnum"))
	afterthreadnum, _  = strconv.Atoi(ethereum.GetValueFromDefaultPath("thread", "afterthreadnum"))
	insertthreadnum, _ = strconv.Atoi(ethereum.GetValueFromDefaultPath("thread", "insertthreadnum"))

	insertbuffer, _ = strconv.Atoi(ethereum.GetValueFromDefaultPath("thread", "insertbuffer"))

	treesize, _ = strconv.Atoi(ethereum.GetValueFromDefaultPath("thread", "treesize"))
)

func Initial(outerbc *core.BlockChain, outbegin int, outend int) {
	maxProces := runtime.NumCPU()
	if maxProces > 1 {
		maxProces -= 1
	}
	runtime.GOMAXPROCS(maxProces)

	if pretasknum == 0 {
		pretasknum = 1
	}

	if prethreadnum == 0 {
		if pretasknum == 1 {
			prethreadnum = 1
		} else {
			prethreadnum = maxProces
		}
	}

	if aftertasknum == 0 {
		aftertasknum = 1
	}

	if afterthreadnum == 0 {
		afterthreadnum = 1
	}

	if insertthreadnum == 0 {
		insertthreadnum = maxProces
	}

	if insertbuffer == 0 {
		insertbuffer = 10
	}

	if outerbc != nil {
		bc = outerbc
	}

	begin = outbegin
	end = outend
	interval = (end - begin + 1) / pretasknum

	log.Info("initial over, the final blocknum is :", "begin", begin, "fn", (interval*pretasknum + begin - 1), "tasknum", pretasknum, "threadnum", prethreadnum, "maxProces", maxProces)
}
