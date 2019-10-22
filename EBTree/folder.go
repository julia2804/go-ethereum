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
	"errors"
	"github.com/ethereum/go-ethereum/rlp"
	"sync"
)

type folder struct {
	tmp        sliceBuffer
	cachegen   uint16
	cachelimit uint16
	onleaf     LeafCallback
}

type sliceBuffer []byte

func (b *sliceBuffer) Write(data []byte) (n int, err error) {
	*b = append(*b, data...)
	return len(data), nil
}

func (b *sliceBuffer) Reset() {
	*b = (*b)[:0]
}

//TODO:folder pool process
// folders live in a global db.
var folderPool = sync.Pool{
	New: func() interface{} {
		return &folder{
			tmp: make(sliceBuffer, 0, 550), // cap is as large as a full fullNode.

		}
	},
}

func newFolder(cachegen, cachelimit uint16, onleaf LeafCallback) *folder {
	h := folderPool.Get().(*folder)
	h.cachegen, h.cachelimit, h.onleaf = cachegen, cachelimit, onleaf
	return h
}

func returnFolderToPool(h *folder) {
	folderPool.Put(h)
}

// fold folds a node , also returning a copy of the
// original node without next field to replace the original one.
func (f *folder) fold(n EBTreen, db *Database, force bool) (EBTreen, error) {
	//log.Info("into fold func")
	// If we're not storing the node, just folding, use available cached data
	if db == nil {
		err := errors.New("the db given in fold is nil")
		return nil, err
	}

	//var collapsedNode idNode
	switch nt := (n).(type) {
	case *leafNode:
		//log.Info("encode a leaf node")
		var collapsed leafNode
		if nt.Id == nil {
			err := errors.New("empty node")
			return nil, err
		}
		collapsed.Id = nt.Id
		da, err := CopyData(nt.Data)
		collapsed.Data = da
		if nt.Next != nil {
			switch cnt := (nt.Next).(type) {
			case *leafNode:
				//log.Info("fold:collapsedNode:leafnode")
				var nb ByteNode
				nb = cnt.Id
				collapsed.Next = &nb
			case *internalNode:
				//log.Info("fold:collapsedNode:internalnode")
				var nb ByteNode
				nb = cnt.Id
				collapsed.Next = &nb
			case *ByteNode:
				//log.Info("fold:collapsedNode:bytenode")
				var nb ByteNode
				nb, _ = cnt.cache()
				collapsed.Next = &nb
			default:
				err := errors.New("fold: wrong collapsed node type")
				return nil, err
			}
		}
		//fmt.Println("we are going to store this node")
		_, err = f.store(&collapsed, db, force)
		if err != nil {
			return nil, err
		}
		return &collapsed, nil
	case *internalNode:
		var collapsed internalNode
		if nt.Id == nil {
			err := errors.New("empty node")
			return nil, err
		}
		collapsed.Id = nt.Id

		for i := 0; i < len(nt.Children); i++ {
			if nt.Children[i] == nil {
				err := errors.New("n.children is nil")
				return nil, err
			}
			switch ct := (nt.Children[i]).(type) {
			case childEncode:
				err := errors.New("child is encoded in fold function")
				return nil, err
			case child:
				if ct.Pointer != nil {
					pEncode, err := f.fold(ct.Pointer, db, false)
					var pet ByteNode
					switch pt := (pEncode).(type) {
					case *leafNode:
						pet = pt.Id
						var cchild child
						cchild.Pointer = &pet
						cchild.Value = ct.Value
						collapsed.Children = append(collapsed.Children, cchild)
						if err != nil {
							return &collapsed, err
						}
					case *internalNode:
						pet = pt.Id
						var cchild child
						cchild.Pointer = &pet
						cchild.Value = ct.Value
						collapsed.Children = append(collapsed.Children, cchild)
						if err != nil {
							return &collapsed, err
						}
					case *ByteNode:
						return pEncode,nil
					default:
						err := errors.New("wrong type")
						return nil, err
					}
				}
			default:
				err := errors.New("wrong type in child")
				return nil, err
			}

		}


		/*_,err := f.foldChildren(collapsed, db)
		if err != nil {
			return nil, err
		}*/
		/* Generate the RLP encoding of the node
		var result []byte
		if err := encodeInternal(&result, nt); err != nil {
			panic("encode error: " + err.Error())
		}
		collapsedNode.NodeData=result*/
		//TODO:process error

		_, _ = f.store(&collapsed, db, force)
		return &collapsed, nil

		// Trie not processed yet or needs storage, walk the children
	case *ByteNode:
		return n, nil


	}

	err := errors.New("there is something wrong in  fold when swich the type of node")
	return nil, err
}

// foldChildren replaces the children of a node with their id .
func (f *folder) foldChildren(original EBTreen, db *Database) (EBTreen, error) {
	switch n := original.(type) {
	case *leafNode:
		var collapsed leafNode
		switch cot := (collapsed.Next).(type) {
		case *leafNode:
			var cnb ByteNode
			cnb = cot.Id
			collapsed.Next = &cnb
		case *internalNode:
			var cnb ByteNode
			cnb = cot.Id
			collapsed.Next = &cnb
		default:
			err := errors.New("wrong type")
			return nil, err
		}
		return &collapsed, nil

	case *internalNode:
		// fold the full node's children, caching the newly hashed subtrees
		var collapsed internalNode
		collapsed.Id = n.Id
		for i := 0; i < len(n.Children); i++ {
			if n.Children[i] == nil {
				err := errors.New("n.children is nil")
				return nil, err
			}
			switch ct := (n.Children[i]).(type) {
			case childEncode:
				return nil, nil
			case child:
				if ct.Pointer != nil {
					pEncode, err := f.fold(ct.Pointer, db, false)
					var pet ByteNode
					switch pt := (pEncode).(type) {
					case *leafNode:
						pet = pt.Id
						ct.Pointer = &pet
						n.Children[i] = ct
						if err != nil {
							return &collapsed, err
						}
					case *internalNode:
						pet = pt.Id
						ct.Pointer = &pet
						n.Children[i] = ct
						if err != nil {
							return &collapsed, err
						}
					default:
						err := errors.New("wrong type")
						return nil, err
					}
				}
			default:
				err := errors.New("wrong type in child")
				return nil, err
			}

		}

		return &collapsed, nil

	default:
		// Value and hash nodes don't have children so they're left as were
		return nil, nil
	}
}
func (f *folder) EncodeNode(n EBTreen) []byte {
	f.tmp.Reset()
	if err := rlp.Encode(&f.tmp, n); err != nil {
		panic("encode error: " + err.Error())
	}
	return f.tmp
}

// store stores the node n and if we have a storage layer specified, it writes
// the key/value pair to it and tracks any node->child references as well as any
// node->external trie references.
func (f *folder) store(n EBTreen, db *Database, force bool) ([]byte, error) {
	// Don't store hashes or empty nodes.
	//db.Cap(1024*1024*64)
	//fmt.Println("into folder.store")
	if db != nil {
		// We are pooling the trie nodes into an intermediate memory cache
		// Generate the RLP encoding of the node
		f.tmp.Reset()
		if err := rlp.Encode(&f.tmp, n); err != nil {
			panic("encode error: " + err.Error())
		}
		db.lock.Lock()
		switch nt := (n).(type) {
		case *leafNode:
			//log.Info("store:into leafnode")
			db.insert(nt.Id, f.tmp, n)
			db.lock.Unlock()
			return nt.Id, nil
		case *internalNode:
			//log.Info("store:into internalnode")
			db.insert(nt.Id, f.tmp, n)
			db.lock.Unlock()
			return nt.Id, nil
		default:
			err := errors.New("wrong type of node")
			db.lock.Unlock()
			return nil, err
		}

	}
	err := errors.New("the db is nil in store functin")
	return nil, err
}
