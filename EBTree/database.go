// Copyright 2018 The go-ethereum Authors
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

package EBTree

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"sync"
	"time"

	"github.com/allegro/bigcache"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
)

var (
	memcacheCleanHitMeter   = metrics.NewRegisteredMeter("trie/memcache/clean/hit", nil)
	memcacheCleanMissMeter  = metrics.NewRegisteredMeter("trie/memcache/clean/miss", nil)
	memcacheCleanReadMeter  = metrics.NewRegisteredMeter("trie/memcache/clean/read", nil)
	memcacheCleanWriteMeter = metrics.NewRegisteredMeter("trie/memcache/clean/write", nil)

	memcacheFlushTimeTimer  = metrics.NewRegisteredResettingTimer("trie/memcache/flush/time", nil)
	memcacheFlushNodesMeter = metrics.NewRegisteredMeter("trie/memcache/flush/nodes", nil)
	memcacheFlushSizeMeter  = metrics.NewRegisteredMeter("trie/memcache/flush/size", nil)

	memcacheGCTimeTimer  = metrics.NewRegisteredResettingTimer("trie/memcache/gc/time", nil)
	memcacheGCNodesMeter = metrics.NewRegisteredMeter("trie/memcache/gc/nodes", nil)
	memcacheGCSizeMeter  = metrics.NewRegisteredMeter("trie/memcache/gc/size", nil)

	memcacheCommitTimeTimer  = metrics.NewRegisteredResettingTimer("trie/memcache/commit/time", nil)
	memcacheCommitNodesMeter = metrics.NewRegisteredMeter("trie/memcache/commit/nodes", nil)
	memcacheCommitSizeMeter  = metrics.NewRegisteredMeter("trie/memcache/commit/size", nil)
)

// secureKeyPrefix is the database key prefix used to store trie node preimages.
var secureKeyPrefix = []byte("secure-key-")

// secureKeyLength is the length of the above prefix + 32byte hash.
const secureKeyLength = 11 + 32

// DatabaseReader wraps the Get and Has method of a backing store for the trie.
type DatabaseReader interface {
	// Get retrieves the value associated with key from the database.
	Get(key []byte) (value []byte, err error)

	// Has retrieves whether a key is present in the database.
	Has(key []byte) (bool, error)
}

// Database is an intermediate write layer between the trie data structures and
// the disk database. The aim is to accumulate trie writes in-memory and only
// periodically flush a couple tries to disk, garbage collecting the remainder.
type Database struct {
	diskdb  ethdb.Database           // Persistent storage for matured trie nodes
	dirties map[string]*collapseNode // Data and references relationships of dirty nodes
	cleans  *bigcache.BigCache       // GC friendly memory cache of clean node RLPs
	oldest  []byte                   // Oldest tracked node, flush-list head
	newest  []byte                   // Newest tracked node, flush-list tail

	seckeybuf [secureKeyLength]byte // Ephemeral buffer for calculating preimage keys

	gctime  time.Duration      // Time spent on garbage collection since last commit
	gcnodes uint64             // Nodes garbage collected since last commit
	gcsize  common.StorageSize // Data storage garbage collected since last commit

	flushtime  time.Duration      // Time spent on data flushing since last commit
	flushnodes uint64             // Nodes flushed since last commit
	flushsize  common.StorageSize // Data storage flushed since last commit

	dirtiesSize common.StorageSize // Storage size of the dirty node cache (exc. flushlist)

	lock sync.RWMutex
}

// cachedNode is all the information we know about a single cached node in the
// memory database write layer.
type cachedNode struct {
	node idNode // Cached collapsed trie node, or raw rlp data
	size uint16 // Byte size of the useful cached data

	children [][]byte // External children referenced by this node

	flushPrev []byte // Previous node in the flush-list
	flushNext []byte // Next node in the flush-list
}

type collapseNode struct {
	node      EBTreen
	flushPrev []byte
	flushNext []byte
	size      uint16 // Byte size of the useful cached data
}

// NewDatabase creates a new trie database to store ephemeral trie content before
// its written out to disk or garbage collected. No read cache is created, so all
// data retrievals will hit the underlying disk database.
func NewDatabase(diskdb ethdb.Database) *Database {
	return NewDatabaseWithCache(diskdb, 256)
}

// NewDatabaseWithCache creates a new trie database to store ephemeral trie content
// before its written out to disk or garbage collected. It also acts as a read cache
// for nodes loaded from disk.
func NewDatabaseWithCache(diskdb ethdb.Database, cache int) *Database {
	var cleans *bigcache.BigCache
	if cache > 0 {
		cleans, _ = bigcache.NewBigCache(bigcache.Config{
			Shards:             1024,
			LifeWindow:         time.Hour,
			MaxEntriesInWindow: cache * 1024,
			MaxEntrySize:       512,
			HardMaxCacheSize:   cache,
		})
	}
	return &Database{
		diskdb:  diskdb,
		cleans:  cleans,
		dirties: map[string]*collapseNode{"": {}},
	}
}

// insert inserts a collapsed tree node into the memory database.
func (db *Database) insert(id []byte, blob []byte, node EBTreen) {
	// If the node's already cached, skip
	fmt.Println("into insert,the id is:")
	fmt.Println(id)
	if _, ok := db.dirties[string(id)]; ok {
		log.Info("the node has been cached, skip")
		return
	}
	// Create the cached entry for this node
	entry := &collapseNode{
		node:      node,
		size:      uint16(len(blob)),
		flushPrev: db.newest,
	}
	fmt.Printf("the size of dirty befor insert is %d", len(db.dirties))
	db.dirties[string(id)] = entry
	fmt.Printf("the size of dirty after insert is %d", len(db.dirties))
	// Update the flush-list endpoints
	if db.oldest == (nil) {
		fmt.Println(id)
		db.oldest, db.newest = id, id
	} else {
		if db.newest == nil {
			fmt.Println("wrong in insert! db.newest==nil")
			return
		}
		//find bugs
		if db.dirties[string(db.newest)] == nil {
			fmt.Println("wrong in insert! db.dirties[db.newest]==nil")
			return
		}
		db.dirties[string(db.newest)].flushNext, db.newest = id, id
	}

	//todo:make sure the dirtiesSize of db
	db.dirtiesSize += common.StorageSize(2*common.HashLength + entry.size)
}

// Commit iterates over all the children of a particular node, writes them out
// to disk, forcefully tearing down all references in both directions.
//
func (db *Database) Commit(node []byte, report bool) error {
	// Create a database batch to flush persistent data out. It is important that
	// outside code doesn't see an inconsistent state (referenced data removed from
	// memory cache during commit but not yet in persistent storage). This is ensured
	// by only uncaching existing data when the database write finalizes.
	db.lock.RLock()

	start := time.Now()
	batch := db.diskdb.NewBatch()
	// Move the trie itself into the batch, flushing if enough data is accumulated
	nodes, storage := len(db.dirties), db.dirtiesSize
	if err := db.commit(node, batch); err != nil {
		log.Error("Failed to commit trie from trie database", "err", err)
		fmt.Println(node)
		db.lock.RUnlock()
		return err
	}
	// Write batch ready, unlock for readers during persistence
	if err := batch.Write(); err != nil {
		log.Error("Failed to write trie to disk", "err", err)
		db.lock.RUnlock()
		return err
	}
	db.lock.RUnlock()

	// Write successful, clear out the flushed data
	db.lock.Lock()
	defer db.lock.Unlock()

	//db.uncache(node)

	memcacheCommitTimeTimer.Update(time.Since(start))
	memcacheCommitSizeMeter.Mark(int64(storage - db.dirtiesSize))
	memcacheCommitNodesMeter.Mark(int64(nodes - len(db.dirties)))

	logger := log.Info
	if !report {
		logger = log.Debug
	}
	logger("Persisted trie from memory database", "nodes", nodes-len(db.dirties)+int(db.flushnodes), "size", storage-db.dirtiesSize+db.flushsize, "time", time.Since(start)+db.flushtime,
		"gcnodes", db.gcnodes, "gcsize", db.gcsize, "gctime", db.gctime, "livenodes", len(db.dirties), "livesize", db.dirtiesSize)

	// Reset the garbage collection statistics
	db.gcnodes, db.gcsize, db.gctime = 0, 0, 0
	db.flushnodes, db.flushsize, db.flushtime = 0, 0, 0

	return nil
}

// commit is the private locked version of Commit.
func (db *Database) commit(id []byte, batch ethdb.Batch) error {
	// If the node does not exist, it's a previously committed node
	log.Info("into commit func")
	node, ok := db.dirties[string(id)]
	if !ok {
		//todo:why this node is not in dirty(maybe related to the limited size of dirty)
		fmt.Println("this node is not dirty, should not be commit")
		err := errors.New("this node is not dirty, should not be commit")
		return err
	}

	//enode the node
	var result []byte
	result = nil

	switch nt := (node.node).(type) {
	case *leafNode:
		log.Info("into node type:leafnode")
		var enode leafNode
		for _, d := range nt.Data {
			var cd data
			switch dt := (d).(type) {
			case dataEncode:
				err := errors.New("wrong data type")
				return err
			case data:
				cd.Value = dt.Value
				cd.Keylist = dt.Keylist
			case *dataEncode:
				err := errors.New("wrong data type")
				return err
			case *data:
				cd.Value = dt.Value
				cd.Keylist = dt.Keylist
			default:
				err := errors.New("wrong data")
				return err
			}
			enode.Data = append(enode.Data, cd)
		}
		enode.Id = nt.Id
		enode.Next = nt.Next
		if err := encodeLeaf(&result, &enode); err != nil {
			panic("encode error: " + err.Error())
		}
		//fmt.Print(result)
	case *internalNode:
		log.Info("into node type:internal node")
		var enode internalNode
		enode.Id = nt.Id
		for i := 0; i < len(nt.Children); i++ {
			//for _, c := range nt.Children {
			var ec child
			switch ct := (nt.Children[i]).(type) {
			case childEncode:
				log.Info("into childrens:child encode")
				err := errors.New("wrong type:childEncode")
				return err
			case child:
				log.Info("into childrens:child")
				switch cpt := (ct.Pointer).(type) {
				case *ByteNode:
					ec.Value = ct.Value
					ec.Pointer = cpt

				case *leafNode:
					var cptid ByteNode
					cptid = cpt.Id
					ct.Pointer = &cptid
					nt.Children[i] = ct
					err := db.commit(cpt.Id, batch)
					if err != nil {
						err := wrapError(err, "something wrong in child pointer commit as leafnode")
						return err
					}
				case *internalNode:
					var cptid ByteNode
					cptid = cpt.Id
					ct.Pointer = &cptid
					nt.Children[i] = ct
					err := db.commit(cpt.Id, batch)
					if err != nil {
						err := wrapError(err, "something wrong in child pointer commit as internalnode")
						return err
					}
				default:
					log.Info("into childres:child:pointer type:default")
					err := errors.New("wrong child pointer type:default")
					return err
				}
			default:
				log.Info("into childrens:wrong child type")
				err := errors.New("wrong child type")
				return err
			}
			enode.Children = append(enode.Children, ec)
		}
		if err := encodeInternal(&result, &enode); err != nil {
			panic("encode error: " + err.Error())
		}
	case *ByteNode:
		ntid, err := nt.cache()
		fmt.Println(err)
		fmt.Println("into node type:bytenode,id：")
		fmt.Println(ntid)
	default:
		log.Info("into node type:wrong type")
		err := errors.New("wrong node type")
		return err
	}

	if result == nil {
		err := errors.New("wrong encode")
		return err
	}
	if err := batch.Put(id[:], result); err != nil {
		return err
	}
	// If we've reached an optimal batch size, commit and start over
	if batch.ValueSize() >= ethdb.IdealBatchSize {
		if err := batch.Write(); err != nil {
			return err
		}
		batch.Reset()
	}
	return nil
}

// Size returns the current storage size of the memory cache in front of the
// persistent database layer.
func (db *Database) Size() common.StorageSize {
	db.lock.RLock()
	defer db.lock.RUnlock()

	// db.dirtiesSize only contains the useful data in the cache, but when reporting
	// the total memory consumption, the maintenance metadata is also needed to be
	// counted. For every useful node, we track 2 extra hashes as the flushlist.
	var flushlistSize = common.StorageSize((len(db.dirties) - 1) * 2 * common.HashLength)
	return db.dirtiesSize + flushlistSize
}
func (db *Database) GetTreeMetas(key []byte) ([]byte, error) {

	//todo:to improve the efficiency, we may use the cache someday
	return db.diskdb.Get(key)
}
func (db *Database) SetTreeMetas(key []byte, result []byte, batch ethdb.Batch) error {
	if err := batch.Put(key, result); err != nil {
		return err
	}

	if err := batch.Write(); err != nil {
		return err
	}
	batch.Reset()

	return nil
}

// node retrieves a cached trie node from memory, or returns nil if none can be
// found in the memory cache.
func (db *Database) node(id []byte, cachegen uint16) EBTreen {
	// todo:fix these problems:Retrieve the node from the clean cache if available
	/*if db.cleans != nil {
		if enc, err := db.cleans.Get(string(id[:])); err == nil && enc != nil {
			fmt.Printf("find the data from db.cleans for  %v",id)
			memcacheCleanHitMeter.Mark(1)
			memcacheCleanReadMeter.Mark(int64(len(enc)))
			return mustDecodeNode(id[:], enc)
		}
	}*/
	// Retrieve the node from the dirty cache if available
	db.lock.RLock()
	dirty := db.dirties[string(id)]
	db.lock.RUnlock()

	if dirty != nil {
		return dirty.node
		//return dirty.obj(id, cachegen)
	}
	// Content unavailable in memory, attempt to retrieve from disk
	fmt.Printf("missing the id:%v, in dirties\n", id[:])
	enc, err := db.diskdb.Get(id[:])
	if err != nil || enc == nil {
		fmt.Printf("not get the id from diskb, the error is:%s\n", err.Error())
		return nil
	}
	/*if db.cleans != nil {
		db.cleans.Set(string(id[:]), enc)
		memcacheCleanMissMeter.Mark(1)
		memcacheCleanWriteMeter.Mark(int64(len(enc)))
	}*/
	return mustDecodeNode(id[:], enc)
}

// obj returns the decoded and expanded trie node, either directly from the cache,
// or by regenerating it from the rlp encoded blob.
func (n *cachedNode) obj(id []byte, cachegen uint16) EBTreen {
	return mustDecodeNode(id[:], n.node.NodeData)
}

// uncache is the post-processing step of a commit operation where the already
// persisted trie is removed from the cache. The reason behind the two-phase
// commit is to ensure consistent data availability while moving from memory
// to disk.
func (db *Database) uncache(id []byte) {
	// If the node does not exist, we're done on this path
	node, ok := db.dirties[string(id)]
	if !ok {
		return
	}

	// Node still exists, remove it from the flush-list
	if bytes.Equal(id, db.oldest) {
		//当前节点为dirty中最老的节点
		db.oldest = node.flushNext
		db.dirties[string(node.flushNext)].flushPrev = nil
	} else if bytes.Equal(id, db.oldest) {
		//当前节点为dirty中最新的节点
		db.newest = node.flushPrev
		db.dirties[string(node.flushPrev)].flushNext = nil
	} else {
		db.dirties[string(node.flushPrev)].flushNext = node.flushNext
		db.dirties[string(node.flushNext)].flushPrev = node.flushPrev
	}

	// TODO:Uncache the node's subtries and remove the node itself too,don't sure if it is neccesary
	fmt.Printf("db.uncache:delete oldest of dirty:%s", string(db.oldest))
	delete(db.dirties, string(id))
	db.dirtiesSize -= common.StorageSize(2*common.HashLength + int(node.size))
}

// Cap iteratively flushes old but still referenced trie nodes until the total
// memory usage goes below the given threshold.
func (db *Database) Cap(limit common.StorageSize) error {
	// Create a database batch to flush persistent data out. It is important that
	// outside code doesn't see an inconsistent state (referenced data removed from
	// memory cache during commit but not yet in persistent storage). This is ensured
	// by only uncaching existing data when the database write finalizes.
	fmt.Println("into db.cap func")
	db.lock.RLock()

	nodes, storage, start := len(db.dirties), db.dirtiesSize, time.Now()
	batch := db.diskdb.NewBatch()

	// db.dirtiesSize only contains the useful data in the cache, but when reporting
	// the total memory consumption, the maintenance metadata is also needed to be
	// counted. For every useful node, we track 2 extra hashes as the flushlist.
	//todo:calculate the size of dirty
	size := db.dirtiesSize + common.StorageSize((len(db.dirties)-1)*2*common.HashLength)

	// Keep committing nodes from the flush-list until we're below allowance
	oldest := db.oldest
	//if the size of node is too large, the flush should happen ,too.
	for oldest != nil {
		// Fetch the oldest referenced node and push into the batch
		node := db.dirties[string(oldest)]

		//enode the node
		var result []byte
		switch nt := (node.node).(type) {
		case *leafNode:
			if err := encodeLeaf(&result, nt); err != nil {
				panic("encode error: " + err.Error())
			}
		case *internalNode:
			if err := encodeInternal(&result, nt); err != nil {
				panic("encode error: " + err.Error())
			}
		default:
			err := errors.New("wrong type")
			return err
		}

		if err := batch.Put(oldest[:], result); err != nil {
			db.lock.RUnlock()
			return err
		}
		// If we exceeded the ideal batch size, commit and reset
		if batch.ValueSize() >= ethdb.IdealBatchSize {
			if err := batch.Write(); err != nil {
				log.Error("Failed to write flush list to disk", "err", err)
				db.lock.RUnlock()
				return err
			}
			batch.Reset()
		}
		// Iterate to the next flush item, or abort if the size cap was achieved. Size
		// is the total size, including both the useful cached data (hash -> blob), as
		// well as the flushlist metadata (2*hash). When flushing items from the cache,
		// we need to reduce both.
		//todo: the size should be considered
		size -= common.StorageSize(3*common.HashLength + int(node.size))
		oldest = node.flushNext
	}
	// Flush out any remainder data from the last batch
	if err := batch.Write(); err != nil {
		log.Error("Failed to write flush list to disk", "err", err)
		db.lock.RUnlock()
		return err
	}
	db.lock.RUnlock()

	// Write successful, clear out the flushed data
	db.lock.Lock()
	defer db.lock.Unlock()

	for !bytes.Equal(db.oldest, oldest) {
		node := db.dirties[string(db.oldest)]
		fmt.Printf("db.cap:delete oldest of dirty:%s", string(db.oldest))
		delete(db.dirties, string(db.oldest))
		db.oldest = node.flushNext
		//todo:calculate the size of dirty
		db.dirtiesSize -= common.StorageSize(2*common.HashLength + int(node.size))
	}
	if db.oldest != nil {
		db.dirties[string(db.oldest)].flushPrev = nil
	}
	db.flushnodes += uint64(nodes - len(db.dirties))
	db.flushsize += storage - db.dirtiesSize
	db.flushtime += time.Since(start)

	memcacheFlushTimeTimer.Update(time.Since(start))
	memcacheFlushSizeMeter.Mark(int64(storage - db.dirtiesSize))
	memcacheFlushNodesMeter.Mark(int64(nodes - len(db.dirties)))

	log.Debug("Persisted nodes from memory database", "nodes", nodes-len(db.dirties), "size", storage-db.dirtiesSize, "time", time.Since(start),
		"flushnodes", db.flushnodes, "flushsize", db.flushsize, "flushtime", db.flushtime, "livenodes", len(db.dirties), "livesize", db.dirtiesSize)

	return nil
}
