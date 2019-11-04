// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package eth

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/ethereum/go-ethereum/EBTree"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/bloombits"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth/downloader"
	"github.com/ethereum/go-ethereum/eth/gasprice"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
	"time"
)

// EthAPIBackend implements ethapi.Backend for full nodes
type EthAPIBackend struct {
	eth *Ethereum
	gpo *gasprice.Oracle
}

// ChainConfig returns the active chain configuration.
func (b *EthAPIBackend) ChainConfig() *params.ChainConfig {
	return b.eth.chainConfig
}

func (b *EthAPIBackend) CurrentBlock() *types.Block {
	return b.eth.blockchain.CurrentBlock()
}

func (b *EthAPIBackend) SetHead(number uint64) {
	b.eth.protocolManager.downloader.Cancel()
	b.eth.blockchain.SetHead(number)
}

func (b *EthAPIBackend) HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Header, error) {
	// Pending block is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block := b.eth.miner.PendingBlock()
		return block.Header(), nil
	}
	// Otherwise resolve and return the block
	if blockNr == rpc.LatestBlockNumber {
		return b.eth.blockchain.CurrentBlock().Header(), nil
	}
	return b.eth.blockchain.GetHeaderByNumber(uint64(blockNr)), nil
}

func (b *EthAPIBackend) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	return b.eth.blockchain.GetHeaderByHash(hash), nil
}

func (b *EthAPIBackend) BlockByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Block, error) {
	// Pending block is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block := b.eth.miner.PendingBlock()
		return block, nil
	}
	// Otherwise resolve and return the block
	if blockNr == rpc.LatestBlockNumber {
		return b.eth.blockchain.CurrentBlock(), nil
	}
	return b.eth.blockchain.GetBlockByNumber(uint64(blockNr)), nil
}

func (b *EthAPIBackend) StateAndHeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*state.StateDB, *types.Header, error) {
	// Pending state is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block, state := b.eth.miner.Pending()
		return state, block.Header(), nil
	}
	// Otherwise resolve the block number and return its state
	header, err := b.HeaderByNumber(ctx, blockNr)
	if header == nil || err != nil {
		return nil, nil, err
	}
	stateDb, err := b.eth.BlockChain().StateAt(header.Root)
	return stateDb, header, err
}

var specificValueSearchTime int64
var specificValueSearchNum int64

func (b *EthAPIBackend) SpecificValueSearch(ctx context.Context, v *hexutil.Big, bn uint64) (EBTree.SearchValue, error) {
	t1 := time.Now()
	fmt.Print("Specific Value search :")
	fmt.Println(v.ToInt().Bytes())
	root, err := b.GetEbtreeRoot(ctx)
	tree, err := EBTree.New(root, b.eth.blockchain.EbtreeCache())

	//var buf1 = make([]byte, 8)
	//binary.BigEndian.PutUint64(buf1, v)

	var buf2 = make([]byte, 8)
	binary.BigEndian.PutUint64(buf2, bn)

	data, err := tree.SpecificValueSearch(v.ToInt().Bytes(), buf2)
	fmt.Println("specific search data num:", len(data.Data))
	t2 := time.Now()
	t3 := t2.Sub(t1).Microseconds()
	specificValueSearchTime = specificValueSearchTime + t3
	specificValueSearchNum++
	return data, err
}

func (b *EthAPIBackend) SpecificValueSearchTime(ctx context.Context) {
	fmt.Println("SpecificValueSearchTime:", specificValueSearchTime, "us")
	fmt.Println("times：", specificValueSearchNum)
}

func (b *EthAPIBackend) ClearSpecificValueSearchTime(ctx context.Context) {
	specificValueSearchTime = 0
	specificValueSearchNum = 0
	fmt.Println("cleared SpecificValueSearchTime")
}

var topkVSearchTotalTime int64
var topkVSearchNum int64

func (b *EthAPIBackend) TopkVSearch(ctx context.Context, k uint64, bn uint64) ([][]byte, error) {
	t1 := time.Now()
	fmt.Print("top k search :")
	fmt.Println(k)
	root, err := b.GetEbtreeRoot(ctx)

	var buf1 = make([]byte, 8)
	binary.BigEndian.PutUint64(buf1, k)

	var buf2 = make([]byte, 8)
	binary.BigEndian.PutUint64(buf2, bn)

	data, err := b.eth.blockchain.TopkVSearch(buf1, buf2, root)
	t2 := time.Now()
	t3 := t2.Sub(t1).Microseconds()
	topkVSearchTotalTime = topkVSearchTotalTime + t3
	topkVSearchNum++
	return data, err
}

func (b *EthAPIBackend) TopkVSearchTime(ctx context.Context) {
	fmt.Println("topkVSearchTotalTime:", topkVSearchTotalTime, "us")
	fmt.Println("times：", topkVSearchNum)
}

func (b *EthAPIBackend) ClearTopkVSearchTime(ctx context.Context) {
	topkVSearchTotalTime = 0
	topkVSearchNum = 0
	fmt.Println("cleared topkVSearchTotalTime")
}

var rangeVSearchTotalTime int64
var rangeVSearchNum int64

func (b *EthAPIBackend) RangeVSearch(ctx context.Context, begin *hexutil.Big, end *hexutil.Big, bn uint64) ([][]byte, error) {
	t1 := time.Now()
	fmt.Print("starting range search : ")
	fmt.Print(begin.ToInt().Bytes())
	fmt.Print("--->")
	fmt.Println(end.ToInt().Bytes())

	root, err := b.GetEbtreeRoot(ctx)
	data, err := b.eth.blockchain.RangeVSearch(begin, end, bn, root)
	t2 := time.Now()
	t3 := t2.Sub(t1).Microseconds()
	rangeVSearchTotalTime = rangeVSearchTotalTime + t3
	rangeVSearchNum++
	return data, err
}

func (b *EthAPIBackend) RangeVSearchTime(ctx context.Context) {
	fmt.Println("rangeVSearchTotalTime:", rangeVSearchTotalTime, "us")
	fmt.Println("times: ", rangeVSearchNum)
}

func (b *EthAPIBackend) ClearRangeVSearchTime(ctx context.Context) {
	rangeVSearchTotalTime = 0
	rangeVSearchNum = 0
	fmt.Println("cleared rangeVSearchTotalTime")
}

func (b *EthAPIBackend) InsertTime(ctx context.Context) {
	b.eth.blockchain.InsertTime()
}

func (b *EthAPIBackend) CreateEbtree(ctx context.Context) (*EBTree.EBTree, error) {
	ebtree, err := b.eth.blockchain.CreateEbtree()
	return ebtree, err
}
func (b *EthAPIBackend) GetEbtreeRoot(ctx context.Context) ([]byte, error) {
	key := []byte("TEbtreeRoot")
	root, err := b.eth.chainDb.Get(key)
	fmt.Print("rid is :")
	fmt.Println(root)
	//root, err := b.eth.blockchain.GetEbtreeRoot()
	if err != nil {
		return nil, err
	}
	return root, nil
}

func (b *EthAPIBackend) GetBlock(ctx context.Context, hash common.Hash) (*types.Block, error) {
	return b.eth.blockchain.GetBlockByHash(hash), nil
}

func (b *EthAPIBackend) GetReceipts(ctx context.Context, hash common.Hash) (types.Receipts, error) {
	return b.eth.blockchain.GetReceiptsByHash(hash), nil
}

func (b *EthAPIBackend) GetLogs(ctx context.Context, hash common.Hash) ([][]*types.Log, error) {
	receipts := b.eth.blockchain.GetReceiptsByHash(hash)
	if receipts == nil {
		return nil, nil
	}
	logs := make([][]*types.Log, len(receipts))
	for i, receipt := range receipts {
		logs[i] = receipt.Logs
	}
	return logs, nil
}

func (b *EthAPIBackend) GetTd(blockHash common.Hash) *big.Int {
	return b.eth.blockchain.GetTdByHash(blockHash)
}

func (b *EthAPIBackend) GetEVM(ctx context.Context, msg core.Message, state *state.StateDB, header *types.Header) (*vm.EVM, func() error, error) {
	state.SetBalance(msg.From(), math.MaxBig256)
	vmError := func() error { return nil }

	context := core.NewEVMContext(msg, header, b.eth.BlockChain(), nil)
	return vm.NewEVM(context, state, b.eth.chainConfig, *b.eth.blockchain.GetVMConfig()), vmError, nil
}

func (b *EthAPIBackend) SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent) event.Subscription {
	return b.eth.BlockChain().SubscribeRemovedLogsEvent(ch)
}

func (b *EthAPIBackend) SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription {
	return b.eth.BlockChain().SubscribeChainEvent(ch)
}

func (b *EthAPIBackend) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	return b.eth.BlockChain().SubscribeChainHeadEvent(ch)
}

func (b *EthAPIBackend) SubscribeChainSideEvent(ch chan<- core.ChainSideEvent) event.Subscription {
	return b.eth.BlockChain().SubscribeChainSideEvent(ch)
}

func (b *EthAPIBackend) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return b.eth.BlockChain().SubscribeLogsEvent(ch)
}

func (b *EthAPIBackend) SendTx(ctx context.Context, signedTx *types.Transaction) error {
	return b.eth.txPool.AddLocal(signedTx)
}

func (b *EthAPIBackend) GetPoolTransactions() (types.Transactions, error) {
	pending, err := b.eth.txPool.Pending()
	if err != nil {
		return nil, err
	}
	var txs types.Transactions
	for _, batch := range pending {
		txs = append(txs, batch...)
	}
	return txs, nil
}

func (b *EthAPIBackend) GetPoolTransaction(hash common.Hash) *types.Transaction {
	return b.eth.txPool.Get(hash)
}

func (b *EthAPIBackend) GetPoolNonce(ctx context.Context, addr common.Address) (uint64, error) {
	return b.eth.txPool.State().GetNonce(addr), nil
}

func (b *EthAPIBackend) Stats() (pending int, queued int) {
	return b.eth.txPool.Stats()
}

func (b *EthAPIBackend) TxPoolContent() (map[common.Address]types.Transactions, map[common.Address]types.Transactions) {
	return b.eth.TxPool().Content()
}

func (b *EthAPIBackend) SubscribeNewTxsEvent(ch chan<- core.NewTxsEvent) event.Subscription {
	return b.eth.TxPool().SubscribeNewTxsEvent(ch)
}

func (b *EthAPIBackend) Downloader() *downloader.Downloader {
	return b.eth.Downloader()
}

func (b *EthAPIBackend) ProtocolVersion() int {
	return b.eth.EthVersion()
}

func (b *EthAPIBackend) SuggestPrice(ctx context.Context) (*big.Int, error) {
	return b.gpo.SuggestPrice(ctx)
}

func (b *EthAPIBackend) ChainDb() ethdb.Database {
	return b.eth.ChainDb()
}

func (b *EthAPIBackend) EventMux() *event.TypeMux {
	return b.eth.EventMux()
}

func (b *EthAPIBackend) AccountManager() *accounts.Manager {
	return b.eth.AccountManager()
}

func (b *EthAPIBackend) RPCGasCap() *big.Int {
	return b.eth.config.RPCGasCap
}

func (b *EthAPIBackend) BloomStatus() (uint64, uint64) {
	sections, _, _ := b.eth.bloomIndexer.Sections()
	return params.BloomBitsBlocks, sections
}

func (b *EthAPIBackend) ServiceFilter(ctx context.Context, session *bloombits.MatcherSession) {
	for i := 0; i < bloomFilterThreads; i++ {
		go session.Multiplex(bloomRetrievalBatch, bloomRetrievalWait, b.eth.bloomRequests)
	}
}
