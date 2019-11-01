// Copyright 2014 The go-ethereum Authors
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

// Package trie implements Merkle Patricia Tries.
package EBTree

import (
	"bytes"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/pkg/errors"
)

var (
	// emptyRoot is the known root hash of an empty trie.
	emptyRoot = common.HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")

	// emptyState is the known hash of an empty state trie entry.
	emptyState = crypto.Keccak256Hash(nil)

	//maxInternalNodeCount
	maxInternalNodeCount = uint8(32)

	//maxLeafNodeCount
	maxLeafNodeCount = uint8(16)

	cacheMissCounter   = metrics.NewRegisteredCounter("trie/cachemiss", nil)
	cacheUnloadCounter = metrics.NewRegisteredCounter("trie/cacheunload", nil)
)

// LeafCallback is a callback type invoked when a trie operation reaches a leaf
// node. It's used by state sync and commit to allow handling external references
// between account and storage tries.
type LeafCallback func(leaf []byte, parent common.Hash) error

type SpecialData struct {
	value []byte
	dirty bool
	data  [][]byte
}

type EBTree struct {
	Db                   *Database
	Root                 EBTreen
	sequence             []byte
	special              []SpecialData
	cachegen, cachelimit uint16
}

// SetCacheLimit sets the number of 'cache generations' to keep.
// A cache generation is created by a call to Commit.
func (t *EBTree) SetCacheLimit(l uint16) {
	t.cachelimit = l
}

func (t *EBTree) isSpecial(value []byte) (bool, uint8) {
	for j, i := range t.special {
		if bytes.Equal(i.value, value) {
			return true, uint8(j)
		}
	}
	return false, 0
}

func (tree *EBTree) DBCommit() ([]byte, error) {
	//store the metas for tree
	batch := tree.Db.diskdb.NewBatch()
	err := tree.Db.SetTreeMetas([]byte("sequence"), tree.sequence, batch)
	if err != nil {
		wrapError(err, "something wrong in store tree.sequence")
	}
	//首先拿到root
	//调用递归commit操作

	switch rt := (tree.Root).(type) {
	case *leafNode:
		err := tree.Db.Commit(rt, true)
		if err != nil {
			wrapError(err, "error in db.commit in func DBCommit")
			return nil, err
		}

		return rt.Id, nil
	case *internalNode:
		err := tree.Db.Commit(rt, true)
		if err != nil {
			wrapError(err, "error in db.commit in func DBCommit")
			return nil, err
		}

		return rt.Id, nil
	default:
		err := errors.New("wong root node type")
		return nil, err
	}

}

// New creates a trie with an existing root node from db.
// If root is the zero hash or the sha3 hash of an empty string, the
// trie is initially empty and does not require a database. Otherwise,
// New will panic if db is nil and returns a MissingNodeError if root does
// not exist in the database. Accessing the trie loads nodes from db on demand.
func New(rid []byte, db *Database) (*EBTree, error) {
	if db == nil {
		panic("trie.New called without a database")
	}
	se, err := db.GetTreeMetas([]byte("sequence"))
	if err != nil {
		log.Info(err.Error())
		se = IntToBytes(0)
	}
	ebt := &EBTree{
		Db:       db,
		sequence: se,
	}

	ebt.Db = db

	if len(rid) != 0 {

		rootNode, err := ebt.resolveHash(rid[:])
		if err != nil {
			return ebt, err
		}

		switch rt := (rootNode).(type) {
		case *idNode:
			rt.Id = rid
		case *leafNode:
			rt.Id = rid
		case *internalNode:
			rt.Id = rid
		default:
			err := errors.New("wrong type")
			return nil, err
		}
		ebt.Root = rootNode
	}
	return ebt, nil
}

//split leafnode into two leaf nodes
func (t *EBTree) splitIntoTwoLeaf(n *leafNode, pos int) (bool, *leafNode, *leafNode, error) {
	var datalist []data
	//fmt.Println("split leafnode into two leaf nodes")
	newn, err := CreateLeafNode(t, datalist)
	if err != nil {
		err = wrapError(err, "split into two leaf node: create leaf node error")
		return false, nil, nil, err
	}
	for j := len(n.Data) - 1; j >= pos; j-- {
		newn.Data = append(newn.Data, data{})
	}
	for i := len(n.Data) - 1; i >= pos; i-- {

		newn.Data[i-pos] = n.Data[i]
		n.Data = append(n.Data[:i])
	}
	return true, n, &newn, nil
}

//split node
func (t *EBTree) split(n EBTreen, parent *internalNode) (bool, *internalNode, error) {
	//fmt.Println("into split ebtree")
	switch nt := n.(type) {
	case *leafNode:
		pos := (len(nt.Data) + 1) / 2

		//split the leaf node into two
		_, _, newn, err := t.splitIntoTwoLeaf(nt, pos)
		if err != nil {
			return false, nil, err
		}
		if nt.Next != nil {
			temp := nt.Next
			nt.Next = newn
			newn.Next = temp
		} else {
			nt.Next = newn
		}
		//fix bug：10/12 err：there is no such leaf node in this parent node
		_, parent, re, err := getLeafNodePosition(nt, parent, t)
		if err != nil {
			return false, nil, err
		}
		//当前节点的元素被split，对应的parent中的children的值也要修改
		switch dt := (nt.Data[len(nt.Data)-1]).(type) {
		case dataEncode:
			_ = dt
			err := errors.New("wrong data type")
			return false, nil, err
		case data:
			switch ct := (parent.Children[re]).(type) {
			case childEncode:
				err := errors.New("wrong data  child type:childEncode")
				return false, nil, err
			case child:
				ct.Value = dt.Value
				parent.Children[re] = ct
				if err != nil {
					return false, parent, err
				}
				switch dtt := (newn.Data[len(newn.Data)-1]).(type) {
				case dataEncode:
					_ = dt
					err := errors.New("wrong data type:dataEncoded")
					return false, nil, err
				case data:
					child2, err := createChild(dtt.Value, newn)
					if err != nil {
						err = wrapError(err, "split leaf node :create child to connect the new node to root")
						return false, parent, err
					}
					su, presult, err := addChild(*parent, child2, int(re+1))
					if !su {
						err = wrapError(err, "split leaf node: add the new child to root")
						return false, parent, err
					}
					parent.Children = presult.Children
					return true, parent, nil
				case *dataEncode:
					_ = dt
					err := errors.New("wrong data type:dataEncoded")
					return false, nil, err
				case *data:
					child2, err := createChild(dtt.Value, newn)
					if err != nil {
						err = wrapError(err, "split leaf node :create child to connect the new node to root")
						return false, parent, err
					}
					su, presult, err := addChild(*parent, child2, int(re+1))
					if !su {
						err = wrapError(err, "split leaf node: add the new child to root")
						return false, parent, err
					}
					parent.Children = presult.Children
					return true, parent, nil
				default:
					err := errors.New("wrong data type:default")
					return false, parent, err
				}
			default:
				err := errors.New("wrong data  child type:default")
				return false, nil, err
			}
		case *data:
			switch ct := (parent.Children[re]).(type) {
			case childEncode:
				err := errors.New("wrong data type:dataEncoded in getLeafNodePosition")
				return false, nil, err
			case child:
				ct.Value = dt.Value
				parent.Children[re] = ct
				if err != nil {
					return false, parent, err
				}
				switch dtt := (newn.Data[len(newn.Data)-1]).(type) {
				case dataEncode:
					err := errors.New("wrong data type")
					return false, nil, err
				case data:
					child2, err := createChild(dtt.Value, newn)
					if err != nil {
						err = wrapError(err, "split leaf node :create child to connect the new node to root")
						return false, parent, err
					}
					su, presult, err := addChild(*parent, child2, int(re+1))
					if !su {
						err = wrapError(err, "split leaf node: add the new child to root")
						return false, parent, err
					}
					parent.Children = presult.Children
					return true, parent, nil
				default:
					err := errors.New("wrong data type:dataEncoded in getLeafNodePosition")
					return false, parent, err
				}
			}
		default:
			err := errors.New("node wrong  type")
			return false, nil, err

		}

	case *internalNode:
		//carry the child node to new node
		var childList []ChildInterface
		pos := (len(nt.Children) + 1) / 2
		newn, err := createInternalNode(t, childList)
		if err != nil {
			err = wrapError(err, "split internal node: create internal node error")
			return false, parent, err
		}
		for j := len(nt.Children) - 1; j >= pos; j-- {
			newn.Children = append(newn.Children, child{})
		}
		for i := len(nt.Children) - 1; i >= pos; i-- {
			newn.Children[i-pos] = nt.Children[i]
			nt.Children = append(nt.Children[:i])
		}
		//直接将新节点链接到当前节点到后面，并链接到父节点上
		//为新创建到节点，创建一个child对象
		switch ct := (newn.Children[len(newn.Children)-1]).(type) {
		case childEncode:
			err := errors.New("wrong data type")
			return false, nil, err
		case child:
			chil, err := createChild(ct.Value, newn)
			if err != nil {
				err = wrapError(err, "split internal node: create newn child wrong")
				return false, parent, err
			}
			//查找当前节点在父节点到位置，新节点放在当前节点后面
			_, parent, re, err := getInternalNodePosition(nt, parent, t)
			switch cpt := (parent.Children[re]).(type) {
			case childEncode:
				err := errors.New("wrong data type")
				return false, nil, err
			case child:
				switch cnt := (nt.Children[len(nt.Children)-1]).(type) {
				case childEncode:
					err := errors.New("wrong data type")
					return false, nil, err
				case child:
					cpt.Value = cnt.Value
					parent.Children[re] = cpt
					if err != nil {
						return false, parent, err
					}
					su, presult, err := addChild(*parent, chil, int(re+1))
					if !su {
						err = wrapError(err, "split internal node: add the new child to root")
						return false, parent, err
					}
					parent.Children = presult.Children
					return true, parent, nil
				default:
					err := errors.New("wrong data type")
					return false, nil, err
				}
			default:
				err := errors.New("wrong data type")
				return false, nil, err
			}

		}

	}
	return false, parent, nil
}

//向EBtree中插入叶子节点
func (t *EBTree) insertLeafNode(n *leafNode, pos int, parent *internalNode, value []byte, id []byte) (bool, *internalNode, error) {
	_, parent, err := moveChildren(parent, pos)
	if err != nil {
		err = wrapError(err, "insert leaf node: move child wrong")
		return false, nil, err
	}
	chil, err := createChild(value, n)
	if err != nil {
		err = wrapError(err, "insert leaf node: create child wrong")
		return false, nil, err
	}
	parent.Children[pos] = chil
	return true, parent, nil
}

//向EBtree中插入internal节点
func (t *EBTree) insertInternalNode(n *internalNode, pos int, parent *internalNode, value []byte, id []byte) (bool, EBTreen, error) {
	_, parent, err := moveChildren(parent, pos)
	if err != nil {
		err = wrapError(err, "insert internal node: move child wrong")
		return false, parent, err
	}
	chil, err := createChild(value, n)
	if err != nil {
		err = wrapError(err, "insert leaf node: create child wrong")
		return false, parent, err
	}
	parent.Children[pos] = chil
	return true, parent, nil
}

//insert into dataNode
func (t *EBTree) InsertToDataNode(i int, nt *leafNode, d *data, value []byte, da []byte, flag bool, parent *internalNode) (bool, bool, *internalNode, error) {
	if Compare(d.Value, value) == 0 {
		//EBTree中已经存储了该value，因此，只要把data加入到对应到datalist中即可
		d.Keylist = append(d.Keylist, da)
		flag = true
		nt.Data[i] = d
		return flag, true, parent, nil
	} else if Compare(d.Value, value) > 0 {

		flag = true
		//说明EBTree中不存在value值，此时，需要构建data，并将其加入到节点中
		su, nt, err := moveData(nt, i)
		if !su {
			err = wrapError(err, "insert data: move data wrong")
			return flag, false, parent, err
		}
		nt.Data[i], err = createData(value, da)
		if err != nil {
			err = wrapError(err, "insert data: create data wrong")
			return flag, false, parent, err
		}
		//split the leaf node
		if len(nt.Data) > int(maxLeafNodeCount) {
			su, parent, err := t.split(nt, parent)
			if !su {
				err = wrapError(err, "insert data: split the leaf node wrong")
				return flag, false, parent, err
			}
			return flag, true, parent, nil
		}

		return flag, true, parent, nil
	}
	return false, false, nil, nil
}

//insert data into internalnode child
func (t *EBTree) InsertToChild(i int, nt *internalNode, cd *child, value []byte, da []byte, parent *internalNode, decoden EBTreen) (bool, bool, *internalNode, EBTreen, error) {
	if Compare(cd.Value, value) >= 0 {

		//call the insert function to
		su, _, err := t.InsertData(decoden, uint8(i), nt, value, da)
		if !su {
			err = wrapError(err, "insert data: when the data was added into appropriate child, something wrong")
			return true, false, parent, decoden, err
		}
		//split internal node
		if len(nt.Children) > int(maxInternalNodeCount) {
			_, parent, err := t.split(nt, parent)
			if err != nil {
				err = wrapError(err, "insert data: when the data was added into appropriate child, split internal node wrong")
				return true, false, parent, decoden, err
			}
		}

		return true, true, parent, decoden, nil
	} else {
		return false, false, parent, decoden, nil
	}
}

func (t *EBTree) InsertToChildEncode(ct childEncode, nt *internalNode, i int, value []byte, da []byte, parent *internalNode) (bool, bool, *internalNode, *child, error) {
	cd, _, err := decodeChild(ct)
	decoden, err := t.resolveHash(ct)
	if err != nil {
		return true, false, nil, nil, err
	}
	cd.Pointer = decoden
	nt.Children[i] = cd
	//insert the data into specific child
	flag, su, parent, decoden, err := t.InsertToChild(i, nt, &cd, value, da, parent, decoden)

	return flag, su, parent, &cd, err
}

//insert into sepcial
func (t *EBTree) InsertToSpecial(p uint8, parent *internalNode, da []byte) (bool, *internalNode, error) {
	//将对应data存入special中
	t.special[p].data = append(t.special[p].data, da)
	t.special[p].dirty = true
	return true, parent, nil
}

//insert data into internalNode
func (t *EBTree) InsertDataToInternal(nt *internalNode, pos uint8, parent *internalNode, value []byte, da []byte) (bool, *internalNode, error) {
	flag := false
	var j int
	for j = 0; j < len(nt.Children); j++ {
		switch ct := (nt.Children[j]).(type) {
		case childEncode:

			//insert the data into specific child
			flag, su, parent, cd, err := t.InsertToChildEncode(ct, nt, j, value, da, parent)
			nt.Children[j] = cd
			if err != nil {
				return false, nil, err
			}
			if !flag {
				continue
			} else {
				return su, parent, err
			}
		case child:
			if Compare(ct.Value, value) >= 0 {
				_, su, parent, cn, err := t.InsertToChild(j, nt, &ct, value, da, parent, ct.Pointer)
				ct.Pointer = cn
				//nt.Children[j] = ct
				return su, parent, err
			} else {
				continue
			}
		default:
			log.Info("wrong child type")
			err := errors.New("wrong child type:default in  InsertDataToInternal")
			return false, nil, err
		}

	}
	if !flag {
		//将该值插入到节点末尾
		//call the insert function to
		switch ct := (nt.Children[j-1]).(type) {
		case childEncode:
			//insert the data into specific child
			_, su, parent, cd, err := t.InsertToChildEncode(ct, nt, j-1, value, da, parent)
			nt.Children[j-1] = cd
			if err != nil {
				return false, nil, err
			}

			return su, parent, err

		case child:
			su, _, err := t.InsertData(ct.Pointer, uint8(j-1), nt, value, da)
			if !su {
				err = wrapError(err, "insert data: when the data was added into last child, something wrong")
				return false, parent, err
			}
			ct.Value = value
			nt.Children[j-1] = ct
			if len(nt.Children) > int(maxInternalNodeCount) {
				su, parent, err := t.split(nt, parent)
				if !su {
					err = wrapError(err, "insert data: when the data was added into last child, split internal node wrong")
					return false, parent, err
				}
				return true, parent, nil
			}

			return true, parent, nil
		default:
			err := errors.New("wrong child type:default in InsertDataToInternal")
			return false, nil, err
		}

	}
	return false, nil, nil
}

//split leafnode into two leaf nodes(recontruct)
func (t *EBTree) splitIntoTwoLeafNode(n *leafNode, pos int) (*leafNode, error) {
	var datalist []data
	//fmt.Println("split leafnode into two leaf nodes")
	newn, err := CreateLeafNode(t, datalist)
	if err != nil {
		err = wrapError(err, "split into two leaf node: create leaf node error")
		return nil, err
	}
	for j := len(n.Data) - 1; j >= pos; j-- {
		newn.Data = append(newn.Data, data{})
	}
	for i := len(n.Data) - 1; i >= pos; i-- {

		newn.Data[i-pos] = n.Data[i]
		n.Data = append(n.Data[:i])
	}
	return &newn, nil
}

//split node(recontruct)
func (t *EBTree) splitNode(n *EBTreen, parent *internalNode, i int) error {
	switch nt := (*n).(type) {
	case *leafNode:
		if uint8(len(nt.Data)) <= maxLeafNodeCount {
			return nil
		}
		pos := (len(nt.Data) + 1) / 2
		//split the leaf node into two
		newn, err := t.splitIntoTwoLeafNode(nt, pos)
		if err != nil {
			return err
		}
		if nt.Next != nil {
			temp := nt.Next
			nt.Next = newn
			newn.Next = temp
		} else {
			nt.Next = newn
		}
		//对于根节点为叶子节点的情况,需要单独讨论
		if parent == nil {
			//需要为上级根节点确定value的值
			switch dt := (nt.Data[len(nt.Data)-1]).(type) {
			case dataEncode, *dataEncode:
				err := errors.New("data is encoded")
				return err
			case data:
				chil, err := createChild(dt.Value, *n)
				var chil2 child
				if err != nil {
					err = wrapError(err, "get leaf node position wrong: when parent is nil, create child wrong")
					return err
				}
				switch dtt := (newn.Data[len(newn.Data)-1]).(type) {
				case dataEncode, *dataEncode:
					err := errors.New("data is encoded")
					return err
				case data:
					chil2, err = createChild(dtt.Value, newn)
					if err != nil {
						err = wrapError(err, "get leaf node position wrong: when parent is nil, create child wrong")
						return err
					}
				case *data:
					chil2, err = createChild(dtt.Value, newn)
					if err != nil {
						err = wrapError(err, "get leaf node position wrong: when parent is nil, create child wrong")
						return err
					}
				default:
					err := errors.New("wrong data type:default")
					return err
				}
				var children []ChildInterface
				children = append(children, chil)
				children = append(children, chil2)
				parent, err = createInternalNode(t, children)
				if err != nil {
					err = wrapError(err, "get leaf node position wrong: when parent is nil, create root")
					return err
				}
				t.Root = parent
				return nil
			case *data:
				chil, err := createChild(dt.Value, *n)
				var chil2 child
				if err != nil {
					err = wrapError(err, "get leaf node position wrong: when parent is nil, create child wrong")
					return err
				}
				switch dtt := (newn.Data[len(newn.Data)-1]).(type) {
				case dataEncode, *dataEncode:
					err := errors.New("data is encoded")
					return err
				case data:
					chil2, err = createChild(dtt.Value, newn)
					if err != nil {
						err = wrapError(err, "get leaf node position wrong: when parent is nil, create child wrong")
						return err
					}
				case *data:
					chil2, err = createChild(dtt.Value, newn)
					if err != nil {
						err = wrapError(err, "get leaf node position wrong: when parent is nil, create child wrong")
						return err
					}
				default:
					err := errors.New("wrong data type:default")
					return err
				}
				var children []ChildInterface
				children = append(children, chil)
				children = append(children, chil2)
				parent, err = createInternalNode(t, children)
				if err != nil {
					err = wrapError(err, "get leaf node position wrong: when parent is nil, create root")
					return err
				}
				t.Root = parent
				return nil
			default:
				err := errors.New("wrong data type:default")
				return err
			}

		}

		//当前节点的元素被split，对应的parent中的children的值也要修改
		switch dt := (nt.Data[len(nt.Data)-1]).(type) {
		case dataEncode:
			fmt.Println("wrong data type")
			err := errors.New("wrong data type")
			return err
		case data:
			switch ct := (parent.Children[i]).(type) {
			case childEncode:
				err := errors.New("wrong data  child type:childEncode")
				return err
			case child:
				ct.Value = dt.Value
				parent.Children[i] = ct
				if err != nil {
					return err
				}
				switch dtt := (newn.Data[len(newn.Data)-1]).(type) {
				case dataEncode:
					err := errors.New("wrong data type:dataEncoded")
					return err
				case data:
					child2, err := createChild(dtt.Value, newn)
					if err != nil {
						err = wrapError(err, "split leaf node :create child to connect the new node to root")
						return err
					}
					su, presult, err := addChild(*parent, child2, int(i+1))
					if !su {
						err = wrapError(err, "split leaf node: add the new child to root")
						return err
					}
					parent.Children = presult.Children
					return nil
				case *dataEncode:
					err := errors.New("wrong data type:dataEncoded")
					return err
				case *data:
					child2, err := createChild(dtt.Value, newn)
					if err != nil {
						err = wrapError(err, "split leaf node :create child to connect the new node to root")
						return err
					}
					su, presult, err := addChild(*parent, child2, int(i+1))
					if !su {
						err = wrapError(err, "split leaf node: add the new child to root")
						return err
					}
					parent.Children = presult.Children
					return nil
				default:
					err := errors.New("wrong data type:default")
					return err
				}
			default:
				err := errors.New("wrong data  child type:default")
				return err
			}
		case *data:
			switch ct := (parent.Children[i]).(type) {
			case childEncode:
				err := errors.New("wrong data type:dataEncoded in getLeafNodePosition")
				return err
			case child:
				ct.Value = dt.Value
				parent.Children[i] = ct
				if err != nil {
					return err
				}
				switch dtt := (newn.Data[len(newn.Data)-1]).(type) {
				case dataEncode:
					err := errors.New("wrong data type")
					return err
				case data:
					child2, err := createChild(dtt.Value, newn)
					if err != nil {
						err = wrapError(err, "split leaf node :create child to connect the new node to root")
						return err
					}
					su, presult, err := addChild(*parent, child2, int(i+1))
					if !su {
						err = wrapError(err, "split leaf node: add the new child to root")
						return err
					}
					parent.Children = presult.Children
					return nil
				default:
					err := errors.New("wrong data type:dataEncoded in getLeafNodePosition")
					return err
				}
			}
		default:
			err := errors.New("node wrong  type")
			return err

		}

	case *internalNode:
		if uint8(len(nt.Children)) <= maxInternalNodeCount {
			return nil
		}
		//carry the child node to new node
		var childList []ChildInterface
		pos := (len(nt.Children) + 1) / 2
		newn, err := createInternalNode(t, childList)
		if err != nil {
			err = wrapError(err, "split internal node: create internal node error")
			return err
		}
		for j := len(nt.Children) - 1; j >= pos; j-- {
			newn.Children = append(newn.Children, child{})
		}
		for i := len(nt.Children) - 1; i >= pos; i-- {
			newn.Children[i-pos] = nt.Children[i]
			nt.Children = append(nt.Children[:i])
		}
		//直接将新节点链接到当前节点到后面，并链接到父节点上
		//为新创建到节点，创建一个child对象
		switch ct := (newn.Children[len(newn.Children)-1]).(type) {
		case childEncode:
			err := errors.New("wrong data type")
			return err
		case child:
			chil, err := createChild(ct.Value, newn)
			if err != nil {
				err = wrapError(err, "split internal node: create newn child wrong")
				return err
			}
			//对于根节点为叶子节点的情况,需要单独讨论
			if parent == nil {
				//需要为上级根节点确定value的值
				switch dt := (nt.Children[len(nt.Children)-1]).(type) {
				case childEncode, *childEncode:
					err := errors.New("child is encoded")
					return err
				case child:
					chil2, err := createChild(dt.Value, *n)
					if err != nil {
						err = wrapError(err, "get leaf node position wrong: when parent is nil, create child wrong")
						return err
					}
					var children []ChildInterface

					children = append(children, chil2)
					children = append(children, chil)
					parent, err = createInternalNode(t, children)
					if err != nil {
						err = wrapError(err, "get leaf node position wrong: when parent is nil, create root")
						return err
					}
					t.Root = parent
					return nil
				case *child:
					chil2, err := createChild(dt.Value, *n)
					if err != nil {
						err = wrapError(err, "get leaf node position wrong: when parent is nil, create child wrong")
						return err
					}
					var children []ChildInterface
					children = append(children, chil)
					children = append(children, chil2)
					parent, err = createInternalNode(t, children)
					if err != nil {
						err = wrapError(err, "get leaf node position wrong: when parent is nil, create root")
						return err
					}
					t.Root = parent
					return nil
				default:
					err := errors.New("wrong data type:default")
					return err
				}

			}
			switch cpt := (parent.Children[i]).(type) {
			case childEncode:
				err := errors.New("wrong data type")
				return err
			case child:
				switch cnt := (nt.Children[len(nt.Children)-1]).(type) {
				case childEncode:
					err := errors.New("wrong data type")
					return err
				case child:
					cpt.Value = cnt.Value
					parent.Children[i] = cpt
					if err != nil {
						return err
					}
					su, presult, err := addChild(*parent, chil, int(i+1))
					if !su {
						err = wrapError(err, "split internal node: add the new child to root")
						return err
					}
					parent.Children = presult.Children
					return nil
				default:
					err := errors.New("wrong data type")
					return err
				}
			default:
				err := errors.New("wrong data type")
				return err
			}
		case *child:
			chil, err := createChild(ct.Value, newn)
			if err != nil {
				err = wrapError(err, "split internal node: create newn child wrong")
				return err
			}
			//对于根节点为叶子节点的情况,需要单独讨论
			if parent == nil {
				//需要为上级根节点确定value的值
				switch dt := (nt.Children[len(nt.Children)-1]).(type) {
				case childEncode, *childEncode:
					err := errors.New("child is encoded")
					return err
				case child:
					chil2, err := createChild(dt.Value, *n)
					if err != nil {
						err = wrapError(err, "get leaf node position wrong: when parent is nil, create child wrong")
						return err
					}
					var children []ChildInterface
					children = append(children, chil)
					children = append(children, chil2)
					parent, err = createInternalNode(t, children)
					if err != nil {
						err = wrapError(err, "get leaf node position wrong: when parent is nil, create root")
						return err
					}
					t.Root = parent
					return nil
				case *child:
					chil2, err := createChild(dt.Value, *n)
					if err != nil {
						err = wrapError(err, "get leaf node position wrong: when parent is nil, create child wrong")
						return err
					}
					var children []ChildInterface
					children = append(children, chil)
					children = append(children, chil2)
					parent, err = createInternalNode(t, children)
					if err != nil {
						err = wrapError(err, "get leaf node position wrong: when parent is nil, create root")
						return err
					}
					t.Root = parent
					return nil
				default:
					err := errors.New("wrong data type:default")
					return err
				}

			}
			switch cpt := (parent.Children[i]).(type) {
			case childEncode:
				err := errors.New("wrong data type")
				return err
			case child:
				switch cnt := (nt.Children[len(nt.Children)-1]).(type) {
				case childEncode:
					err := errors.New("wrong data type")
					return err
				case child:
					cpt.Value = cnt.Value
					parent.Children[i] = cpt
					if err != nil {
						return err
					}
					su, presult, err := addChild(*parent, chil, int(i+1))
					if !su {
						err = wrapError(err, "split internal node: add the new child to root")
						return err
					}
					parent.Children = presult.Children
					return nil
				default:
					err := errors.New("wrong child type")
					return err
				}
			default:
				err := errors.New("wrong child type")
				return err
			}
		default:
			err := errors.New("wrong child type")
			return err
		}
	default:
		err := errors.New("node is defalut in splitNode")
		return err
	}
	err := errors.New("something wrong")
	return err
}

//insert data into internalnode child, not split（reconstruct)
func (t *EBTree) InsertToInternalChild(cd *child, value []byte, da []byte) error {
	//if the child pointer is bytenode,we should construct a leaf/internal node from it first
	switch pt := (cd.Pointer).(type) {
	case *leafNode, *internalNode:
		//call the insert function to
		err := t.InsertDataToNode(&cd.Pointer, value, da)
		if err != nil {
			wrapError(err, "insert data: when the data was added into appropriate child, something wrong")
			return err
		}
		if Compare(cd.Value, value) > 0 {
			cd.Value = value
		}
		return nil
	case *ByteNode:
		//decode the pointer
		ptid, _ := pt.cache()
		decoden, err := t.resolveHash(ptid)
		if err != nil {
			wrapError(err, "decoden error")
			return err
		}
		//replace the bytenode with leaf/internal node
		cd.Pointer = decoden
		//call the insert function to
		err = t.InsertDataToNode(&cd.Pointer, value, da)
		if err != nil {
			wrapError(err, "insert data: when the data was added into appropriate child, something wrong")
			return err
		}
		if Compare(cd.Value, value) > 0 {
			cd.Value = value
		}
		return nil
	default:
		log.Info("wrong pointer type：default")
		err := errors.New("wrong pointer type:default in  InsertToInternalChild")
		return err
	}
}

//insert data into internalNode（reconstruct)
func (t *EBTree) InsertDataToInternalNode(nt *internalNode, value []byte, da []byte) error {
	var j int
	for j = 0; j < len(nt.Children); j++ {
		switch ct := (nt.Children[j]).(type) {
		case childEncode:
			log.Info("wrong child type：childEncode")
			err := errors.New("child is encoded in InsertDataToInternalNode")
			return err
		case child:
			if Compare(ct.Value, value) <= 0 {
				//将数据插入到对应的child中
				err := t.InsertToInternalChild(&ct, value, da)
				nt.Children[j] = ct
				if err != nil {
					return err
				}
				//判断child对应的节点是否需要split
				err = t.splitNode(&ct.Pointer, nt, j)
				//返回结果
				if err != nil {
					wrapError(err, "insert data: when the data was added into appropriate child, something wrong")
					return err
				}
				return nil
			} else {
				continue
			}

		default:
			log.Info("wrong child type：default")
			err := errors.New("wrong child type:default in  InsertDataToInternal")
			return err
		}

	}

	//将该值插入到节点末尾
	//update the value of children
	//call the insert function to
	switch ct := (nt.Children[j-1]).(type) {
	case childEncode:
		log.Info("wrong child type：childEncode in the last")
		err := errors.New("child is encoded in InsertDataToInternalNode")
		return err
	case child:
		err := t.InsertToInternalChild(&ct, value, da)
		nt.Children[j-1] = ct
		if err != nil {
			return err
		}
		//判断child对应的节点是否需要split
		err = t.splitNode(&ct.Pointer, nt, j-1)
		//返回结果
		if err != nil {
			wrapError(err, "insert data: when the data was added into appropriate child, something wrong")
			return err
		}
		return nil
	default:
		log.Info("wrong child type：default in the last")
		err := errors.New("wrong child type:default in  InsertDataToInternal")
		return err
	}

}

//insert into dataNode（reconstruct)
func (t *EBTree) InsertToLeafData(nt *leafNode, i int, d *data, value []byte, da []byte) error {
	if Compare(d.Value, value) == 0 {
		//EBTree中已经存储了该value，因此，只要把data加入到对应到datalist中即可
		d.Keylist = append(d.Keylist, da)
		return nil
	} else {

		//说明EBTree中不存在value值，此时，需要构建data，并将其加入到节点中
		sucess, nt, err := moveData(nt, i)
		if !sucess {
			err = wrapError(err, "insert data: move data wrong")
			return err
		}
		nt.Data[i], err = createData(value, da)
		if err != nil {
			err = wrapError(err, "insert data: create data wrong")
			return err
		}
		return nil
	}
}

//insert data into leafNode（reconstruct)
func (t *EBTree) InsertDataToLeafNode(nt *leafNode, value []byte, da []byte) error {
	//向叶子节点插入数据
	//若当前节点为空时，直接插入节点。
	if len(nt.Data) == 0 {
		//log.Info("the data is nil")
		//create a data item for da
		dai, err := createData(value, da)
		if err != nil {
			err = wrapError(err, "insert data: create data wrong")
			return err
		}
		nt.Data = append(nt.Data, dai)

		return nil
	}

	//遍历当前节点的所有data，将da插入合适的位置
	//value一定小于或等于当前节点到最大值
	for i := 0; i < len(nt.Data); i++ {
		//log.Info("find the appropriate position in nt datas")
		switch dt := (nt.Data[i]).(type) {
		case dataEncode:
			//decode the data
			log.Info("data is encoded！")
			err := errors.New("data type is *dataEncode")
			return err
		case data:
			if Compare(dt.Value, value) > 0 {
				continue
			} else {
				err := t.InsertToLeafData(nt, i, &dt, value, da)
				if err != nil {
					return err
				}
				return nil
			}
		case *data:
			if Compare(dt.Value, value) > 0 {
				continue
			} else {
				err := t.InsertToLeafData(nt, i, dt, value, da)
				if err != nil {
					return err
				}
				return nil
			}
		case *dataEncode:
			log.Info("pointer data encode type")
			err := errors.New("data type is *dataEncode")
			return err
		default:
			log.Info("data type is not appropriate")
			err := errors.New("data type is not appropriate")
			return err
		}
	}

	//将该值插入到节点末尾
	//log.Info("the data should be put in the last ")
	dai, err := createData(value, da)
	if err != nil {
		err = wrapError(err, "insert data: when the data was added into the end of node, create data wrong")
		return err
	}
	nt.Data = append(nt.Data, dai)
	return nil

}

//insert data into leafNode
func (t *EBTree) InsertDataToLeaf(nt *leafNode, pos uint8, parent *internalNode, value []byte, da []byte) (bool, *internalNode, error) {
	//向叶子节点插入数据
	//若当前节点为空时，直接插入节点。
	if len(nt.Data) == 0 {
		//log.Info("the data is nil")
		//create a data item for da
		dai, err := createData(value, da)
		if err != nil {
			err = wrapError(err, "insert data: create data wrong")
			return false, parent, err
		}
		nt.Data = append(nt.Data, dai)

		return true, parent, nil
	}

	//用于标记该value是否被插入成功
	flag := false
	//遍历当前节点的所有data，将da插入合适的位置
	//value一定小于或等于当前节点到最大值

	for i := 0; i < len(nt.Data); i++ {
		//log.Info("find the appropriate position in nt datas")

		switch dt := (nt.Data[i]).(type) {
		case dataEncode:
			//decode the data
			log.Info("data is encoded,dt is:")
			fmt.Println(dt)
			d, err := decodeData(dt)
			if err != nil {
				return false, parent, err
			}
			flag, su, parent, err := t.InsertToDataNode(i, nt, &d, value, da, flag, parent)
			if flag {
				return su, parent, err
			} else {
				continue
			}
		case data:
			flag, su, parent, err := t.InsertToDataNode(i, nt, &dt, value, da, flag, parent)
			if flag {
				return su, parent, err
			} else {
				continue
			}

		case *data:
			flag, su, parent, err := t.InsertToDataNode(i, nt, dt, value, da, flag, parent)
			if flag {
				return su, parent, err
			} else {
				continue
			}
		case *dataEncode:
			log.Info("pointer data encode type")

		default:
			log.Info("data type is not appropriate")
			err := errors.New("data type is not appropriate")
			return false, nil, err
		}

	}

	//将该值插入到节点末尾
	if !flag {
		//log.Info("the data should be put in the last ")
		dai, err := createData(value, da)
		if err != nil {
			err = wrapError(err, "insert data: when the data was added into the end of node, create data wrong")
			return false, parent, err
		}

		nt.Data = append(nt.Data, dai)

		//bak_code 1
		return true, nil, nil
	}
	return true, nil, nil
}

func (t *EBTree) InsertDataToTree(value []byte, da []byte) error {
	//fmt.Print("start to insert value :")
	//fmt.Println(value)
	err := t.InsertDataToNode(&t.Root, value, da)
	if err != nil {
		return err
	}
	return t.splitNode(&t.Root, nil, 0)
}

//将value插入到该节点或节点的子节点中
func (t *EBTree) InsertDataToNode(n *EBTreen, value []byte, da []byte) error {

	switch nt := (*n).(type) {
	case *leafNode:
		//insert into leafNode
		return t.InsertDataToLeafNode(nt, value, da)
		//不进行split
	case *internalNode:
		//insert into internal node
		return t.InsertDataToInternalNode(nt, value, da)
	case *ByteNode:
		//insert into byte node,need to use real node to replace it
		ntid, _ := (*n).cache()
		decoden, err := t.resolveHash(ntid)
		if err != nil {
			return err
		}
		n = &decoden
		return t.InsertDataToNode(n, value, da)
	case nil:
		dai, err := createData(value, da)
		if err != nil {
			err = wrapError(err, "insert data: create data wrong")
			return err
		}
		var da []data
		da = append(da, dai)
		newn, err := CreateLeafNode(t, da)
		t.Root = &newn
		if err != nil {
			log.Info("err in create leaf node")
			return err
		}
		return nil
	default:
		log.Info("n with wrong node type")
		err := errors.New("the node is not leaf or internal, something wrong")
		return err
	}
	err := errors.New("the function reach to the bottom in InsertDataToNode, something wrong")
	return err
}

//向EBTree中插入数据
//special value被存储在特定结构中
//其他值正常存储在tree中
func (t *EBTree) InsertData(n EBTreen, pos uint8, parent *internalNode, value []byte, da []byte) (bool, *internalNode, error) {
	//判断value是否special
	sp, p := t.isSpecial(value)
	if sp {
		return t.InsertToSpecial(p, parent, da)
	}
	switch nt := n.(type) {
	case *leafNode:
		return t.InsertDataToLeaf(nt, pos, parent, value, da)
	case *internalNode:
		return t.InsertDataToInternal(nt, pos, parent, value, da)
	case *ByteNode:
		nbb, _ := n.cache()
		if string(nbb) == "" {
			dai, err := createData(value, da)
			if err != nil {
				err = wrapError(err, "insert data: create data wrong")
				return false, parent, err
			}
			var da []data
			da = append(da, dai)
			newn, err := CreateLeafNode(t, da)
			t.Root = &newn
			if err != nil {
				log.Info("err in create leaf node")
				return false, nil, err
			}
			return true, nil, nil
		} else {
			var nb []byte
			nb, _ = n.cache()
			decoden, err := t.resolveHash(nb)
			if err != nil {
				return false, nil, err
			}
			n = decoden
			return t.InsertData(decoden, pos, parent, value, da)
		}
	default:
		log.Info("n with wrong node type")
		err := errors.New("the node is not leaf or internal, something wrong")
		return false, nil, err
	}
	err := errors.New("the function reach to the bottom, something wrong")
	return false, nil, err
}
func (t *EBTree) resolveLeaf(n []byte) (leafNode, error) {

	if BytesToInt(n) == uint64(20) {
		fmt.Println("something wrong")
	}
	cacheMissCounter.Inc(1)

	if node := t.Db.node(n, t.cachegen); node != nil {
		switch nt := (node).(type) {
		case *leafNode:
			var ds []data
			for i := 0; i < len(nt.Data); i++ {
				switch dt := (nt.Data[i]).(type) {
				case data:
					ds = append(ds, dt)
				}
			}
			var le leafNode
			if(nt.Next == nil){
				le, _ = constructLeafNode(nt.Id, uint8(len(nt.Data)), ds, false, true, nt.Next, nil, 0)
			}else{
				nextid, _ := nt.Next.cache()
				le, _ = constructLeafNode(nt.Id, uint8(len(nt.Data)), ds, false, true, nt.Next, nextid, 0)
			}
			return le, nil
		}

	}
	log.Info("not get the node from db")
	fmt.Println(n)
	err := errors.New("not get the leaf node %v from db")
	return leafNode{}, err
}
func (t *EBTree) resolveHash(n []byte) (EBTreen, error) {

	cacheMissCounter.Inc(1)

	if node := t.Db.node(n, t.cachegen); node != nil {
		return node, nil
	}
	log.Info("not get the node from db")
	fmt.Println(n)
	err := errors.New("not get the leaf node %v from db")
	return nil, err
}


func (t *EBTree) newSequence() ([]byte, error) {
	//log.Info("into new sequece")
	//log.Info(string(t.sequence))
	re := BytesToInt(t.sequence)
	re = re + 1

	if re < 0 {
		err := errors.New("BytesToInt return a negtive data")
		return nil, err
	}
	id := IntToBytes(uint64(re))
	t.sequence = id
	//fmt.Println(t.sequence)
	return id, nil
}
func (t *EBTree) OutputRoot() []byte {
	switch rt := (t.Root).(type) {
	case *leafNode:
		//log.Info("outputroot:root node type: leafnode")
		//log.Info(string(rt.Id))
		return rt.Id
	case *internalNode:
		//log.Info("outputroot:root node type: internalnode")
		return rt.Id
	case *ByteNode:
		//log.Info("outputroot:root node type: bytenode")
		id, _ := rt.cache()
		return id
	default:
		log.Info("Output  root: wrong root node type")
	}
	return nil
}

// Commit writes all nodes to the trie's memory database, tracking the internal
// and external (for account tries) references.
func (t *EBTree) Commit(onleaf LeafCallback) ([]byte, error) {
	if t == nil {
		panic("nil tree")
	}
	if t.Db == nil {
		panic("commit called on tree with nil database")
	}
	collapsedNode, err := t.foldRoot(t.Db, onleaf)
	if err != nil {
		return nil, err
	}
	//t.root = collapsedNode
	t.cachegen++
	if collapsedNode == nil {
		err := errors.New("collased node is nil")
		return nil, err
	}
	rt, _ := collapsedNode.cache()
	//log.Info(string(rt))
	return rt, nil
}

func (t *EBTree) foldRoot(db *Database, onleaf LeafCallback) (EBTreen, error) {
	//log.Info("into fold root ")
	if t.Root == nil {
		err := errors.New("tree is nil")
		return nil, err
	}
	f := newFolder(t.cachegen, t.cachelimit, onleaf)
	defer returnFolderToPool(f)

	//todo:fold the sequence and special

	return f.fold(t.Root, db, true)
}

// TryGet returns the value for key stored in the trie.
// The value bytes must not be modified by the caller.
// If a node was not found in the database, a MissingNodeError is returned.
func (t *EBTree) TryGet(value []byte) ([][]byte, error) {

	data, _, didResolve, err := t.tryGet(t.Root, value, 0)
	if err == nil && didResolve {
		//t.root = newroot
	}
	return data, err
}

//bool is used to mark whether the node is decoded right
func (t *EBTree) tryGet(origNode EBTreen, value []byte, pos int) ([][]byte, EBTreen, bool, error) {
	switch n := (origNode).(type) {
	case *ByteNode:
		//decode this node
		nc, _ := n.cache()
		decoden, err := t.resolveHash(nc)
		if err != nil {
			return nil, decoden, true, err
		} else {
			//得到该节点之后，需要将其中包含的data或者children解析出来
			switch ct := (decoden).(type) {
			//对于internal 节点
			case *internalNode:
				//先恢复ID
				ct.Id, _ = n.cache()
				//ct.Id=[]byte(n.fstring("a"))
				//再将该节点的子节点恢复出来。这里主要是要将pointer恢复成对应的子节点
				for i := 0; i < len(ct.Children); i++ {
					switch ctt := (ct.Children[i]).(type) {
					//首先判断该child是否已经被解码
					case childEncode:
						err := errors.New("child is encoded")
						return nil, nil, false, err
					case child:
						if Compare(ctt.Value, value) < 0 {
							continue
						} else {
							switch cpt := (ctt.Pointer).(type) {
							case *ByteNode:
								//将pointer对应的子节点解析出来
								var cb ByteNode
								cb, _ = cpt.cache()
								decodechild, err := t.resolveHash(cb)
								if err != nil {
									err = wrapError(err, "some thing wrong in resolve hash")
									return nil, nil, false, err
								}
								chi := child{}
								chi.Pointer = decodechild
								chi.Value = ctt.Value

								//继续去子节点中搜索对应value
								data, encodeNode, su, err := t.tryGet(chi.Pointer, value, 0)
								chi.Pointer = encodeNode
								//解析之后将其放回到对应域中
								ct.Children[i] = chi
								decoden = ct

								return data, decoden, su, err
							case *leafNode:
								continue
							case *internalNode:
								continue
							default:
								err := errors.New("wrong type")
								return nil, nil, false, err
							}
						}
					}
				}
				decoden = ct
				err := errors.New("not found the data for the value")
				return nil, decoden, false, err
			case *leafNode:
				ct.Id, _ = n.cache()
				data, newnode, didResolve, err := t.tryGet(ct, value, pos)
				if err != nil {
					return nil, decoden, true, err
				}
				return data, newnode, didResolve, nil
			}
		}
	case *leafNode:
		//判断value是否在当前data的范围中
		switch dt := (n.Data[0]).(type) {
		case dataEncode:
			err := errors.New("data is encoded for data 0")
			return nil, nil, false, err
		case data:
			switch dmt := (n.Data[len(n.Data)-1]).(type) {
			case dataEncode:
				err := errors.New("data is encoded for data len")
				return nil, nil, false, err
			case data:
				if Compare(value, dt.Value) < 0 || Compare(value, dmt.Value) > 0 {
					// key not found in trie
					err := errors.New("key not found")
					return nil, n, false, err
				}
				//在data中查找
				for i := 0; i < len(n.Data); i++ {
					switch dt := (n.Data[i]).(type) {
					case dataEncode:
						err := errors.New("data is encoded for data 0")
						return nil, nil, false, err
					case data:
						if Compare(dt.Value, value) == 0 {
							return dt.Keylist, n, true, nil
						}
					default:
						err := errors.New("wrong type")
						return nil, nil, false, err
					}
				}
				err := errors.New("no such value in leaf node!")
				return nil, n, false, err
			default:
				return nil, nil, false, nil

			}
		default:
			err := errors.New("wrong  type")
			return nil, nil, false, err
		}
	case *internalNode:
		for i := 0; i < len(n.Children); i++ {
			switch ct := (n.Children[i]).(type) {
			case childEncode:
				err := errors.New("data is encoded for data len")
				return nil, nil, false, err
			case child:
				if Compare(ct.Value, value) >= 0 {
					result, decodeChild, su, err := t.tryGet(ct.Pointer, value, 0)
					ct.Pointer = decodeChild
					n.Children[i] = ct
					return result, n, su, err
				}
			default:
				err := errors.New("wrong type")
				return nil, nil, false, err
			}
		}

		err := errors.New("no such value in internal node!")
		return nil, n, false, err
	default:
		panic(fmt.Sprintf("%T: invalid node: %v", origNode, origNode))
	}
	return nil, nil, false, nil

}


