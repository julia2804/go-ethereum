package ebtree_v2

import (
	"fmt"
	"github.com/allegro/bigcache"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"sync"
	"time"
)

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
	diskdb ethdb.Database     // Persistent storage for matured trie nodes
	cleans *bigcache.BigCache // GC friendly memory cache of clean node RLPs
	oldest []byte             // Oldest tracked node, flush-list head
	newest []byte             // Newest tracked node, flush-list tail

	gctime  time.Duration      // Time spent on garbage collection since last commit
	gcnodes uint64             // Nodes garbage collected since last commit
	gcsize  common.StorageSize // Data storage garbage collected since last commit

	flushtime  time.Duration      // Time spent on data flushing since last commit
	flushnodes uint64             // Nodes flushed since last commit
	flushsize  common.StorageSize // Data storage flushed since last commit

	dirtiesSize common.StorageSize // Storage size of the dirty node cache (exc. flushlist)

	lock sync.RWMutex
}

// NewDatabase creates a new ebtree database to store ephemeral ebtree content before
// its written out to disk or garbage collected. No read cache is created, so all
// data retrievals will hit the underlying disk database.
func NewDatabase(diskdb ethdb.Database) *Database {
	return NewDatabaseWithCache(diskdb, 256)
}

// NewDatabaseWithCache creates a new ebtree database to store ephemeral ebtree content
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
		diskdb: diskdb,
		cleans: cleans,
	}
}

// 提交保存以node为根节点的树
func (db *Database) commit(id []byte, result []byte, batch ethdb.Batch) error {

	//将result放入batch中
	if err := batch.Put(id, result); err != nil {
		return err
	}
	// If we've reached an optimal batch size, commit and start over
	//if batch.ValueSize() >= ethdb.IdealBatchSize {
	if batch.ValueSize() >= ethdb.IdealBatchSize {
		if err := batch.Write(); err != nil {
			return err
		}
		batch.Reset()
	}

	return nil
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

func (db *Database) node(id []byte) (EBTreen, error) {
	enc, err := db.diskdb.Get(id[:])
	if err != nil || enc == nil {
		fmt.Printf("not get the id from diskb, the error is:%s\n", err.Error())
		return nil, err
	}
	return DecodeNode(enc)
}
