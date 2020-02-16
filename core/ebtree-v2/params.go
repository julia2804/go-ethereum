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
	gettasknum, _ = strconv.Atoi(ethereum.GetValueFromDefaultPath("thread", "gettasknum"))

	getthreadnum, _ = strconv.Atoi(ethereum.GetValueFromDefaultPath("thread", "getthreadnum"))

	insertthreadnum, _ = strconv.Atoi(ethereum.GetValueFromDefaultPath("thread", "insertthreadnum"))

	insertbuffer, _ = strconv.Atoi(ethereum.GetValueFromDefaultPath("thread", "insertbuffer"))

	treesize, _ = strconv.Atoi(ethereum.GetValueFromDefaultPath("thread", "treesize"))

	constructSavePath = ethereum.GetValueFromDefaultPath("thread", "constructSavePath")

	MaxLeafNodeCapability, _     = strconv.Atoi(ethereum.GetValueFromDefaultPath("thread", "leaf_num"))
	MaxInternalNodeCapability, _ = strconv.Atoi(ethereum.GetValueFromDefaultPath("thread", "internal_num"))

	recordPath = ethereum.GetValueFromDefaultPath("thread", "recordPath")

	insert_begin_end_Path = ethereum.GetValueFromDefaultPath("thread", "SingleAppendPath")

	//experSavePath = ethereum.GetValueFromDefaultPath("thread", "experSavePath")
	duplicate, _     = strconv.Atoi(ethereum.GetValueFromDefaultPath("thread", "duplicate"))
	topkSavePath     = ethereum.GetValueFromDefaultPath("thread", "topkSavePath")
	specificSavePath = ethereum.GetValueFromDefaultPath("thread", "specificSavePath")
	rangeSavePath    = ethereum.GetValueFromDefaultPath("thread", "rangeSavePath")
)

func Initial(outerbc *core.BlockChain, outbegin int, outend int) {
	maxProces := runtime.NumCPU()
	if maxProces > 1 {
		maxProces -= 1
	}
	runtime.GOMAXPROCS(maxProces)

	if gettasknum == 0 {
		gettasknum = 1
	}

	if getthreadnum == 0 {
		if gettasknum == 1 {
			getthreadnum = 1
		} else {
			getthreadnum = maxProces
		}
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

	if constructSavePath == "" {
		constructSavePath = "/root/experConstruct/"
	}

	if MaxLeafNodeCapability == 0 {
		MaxLeafNodeCapability = 32
	}
	if MaxInternalNodeCapability == 0 {
		MaxInternalNodeCapability = 512
	}

	if recordPath == "" || len(recordPath) == 0 {
		recordPath = "/root/record.txt"
		//recordPath = "/root/record" + time.Now().Format("2006-01-02 15:04:05.9999") + ".txt"
	}

	if insert_begin_end_Path == "" || len(insert_begin_end_Path) == 0 {
		insert_begin_end_Path = "/root/insert_begin_end.txt"
	}

	if topkSavePath == "" || len(topkSavePath) == 0 {
		topkSavePath = "/root/expersave.txt"
	}

	if rangeSavePath == "" || len(rangeSavePath) == 0 {
		rangeSavePath = "/root/expersave.txt"
	}

	if specificSavePath == "" || len(specificSavePath) == 0 {
		specificSavePath = "/root/expersave.txt"
	}

	if duplicate == 0 {
		duplicate = 5
	}

	begin = outbegin
	end = outend
	interval = (end - begin + 1) / gettasknum

	log.Info("initial over, the final blocknum is :", "begin", begin, "fn", (interval*gettasknum + begin - 1), "gettasknum", gettasknum,
		"threadnum", getthreadnum, "maxProces", maxProces, "constructPath", constructSavePath,
		"leaf_num", MaxLeafNodeCapability, "internal_num", MaxInternalNodeCapability, "recordpath",
		recordPath, "insert_begin_end_path", insert_begin_end_Path, "duplicate", duplicate,
		"topkSavePath", topkSavePath, "specificSavePath", specificSavePath, "rangeSavePath", rangeSavePath)
}

func CloseParams() {
	bc.Stop()
	bc = nil
	begin = 0
	end = 0
}
