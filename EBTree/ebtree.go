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
	"encoding/binary"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/pkg/errors"
	"reflect"
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

// Trie is a Merkle Patricia Trie.
// The zero value is an empty trie with no database.
// Use New to create a trie that sits on top of a database.
//
// Trie is not safe for concurrent use.
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
func typeof(v interface{}) string {
	return reflect.TypeOf(v).String()
}
func DBCommit(tree *EBTree) ([]byte, error) {
	switch rt := (tree.Root).(type) {
	case *leafNode:
		//todo:
		log.Info("dbcommit:before commit:next is rt.id")
		log.Info(string(rt.Id))
		err := tree.Db.Commit(rt.Id, true)
		if err != nil {
			wrapError(err, "error in db.commit in func DBCommit")
			return nil, err
		}
		log.Info("dbcommit:next is rt.id")
		log.Info(string(rt.Id))
		return rt.Id, nil
	case *internalNode:
		//todo:
		err := tree.Db.Commit(rt.Id, true)
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
func New(root []byte, db *Database) (*EBTree, error) {
	log.Info("into new ebtree")
	if db == nil {
		panic("trie.New called without a database")
	}
	ebt := &EBTree{
		Db:       db,
		sequence: IntToBytes(0),
	}
	log.Info(string(len(root)))
	ebt.Db = db
	if len(root) != 0 {
		log.Info(string(root))
		rootNode, err := ebt.resolveHash(root[:])
		if err != nil {
			return ebt, err
		}

		switch rt := (rootNode).(type) {
		case *idNode:
			rt.Id = root
		case *leafNode:
			rt.Id = root
		case *internalNode:
			rt.Id = root
		default:
			err := errors.New("wrong type")
			return nil, err
		}
		ebt.Root = rootNode
	}
	return ebt, nil
}
func (t *EBTree) splitIntoTwoLeaf(n *leafNode, pos int) (bool, *leafNode, *leafNode, error) {
	var datalist []data
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
func (t *EBTree) split(n EBTreen, parent *internalNode) (bool, *internalNode, error) {
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

		_, parent, re, err := getLeafNodePosition(nt, parent, t)
		//当前节点的元素被split，对应的parent中的children的值也要修改
		switch dt := (nt.Data[len(nt.Data)-1]).(type) {
		case dataEncode:
			//TODO:
			_ = dt
			err := errors.New("wrong data type")
			return false, nil, err
		case data:
			switch ct := (parent.Children[re]).(type) {
			case childEncode:
				return false, nil, nil
			case child:
				ct.Value = dt.Value
				parent.Children[re] = ct
				if err != nil {
					return false, parent, err
				}
				switch dtt := (newn.Data[len(newn.Data)-1]).(type) {
				case dataEncode:
					//TODO:
					_ = dt
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
						err = wrapError(err, "split leaf node: add the new child to root")
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

func (t *EBTree) DoNothing() error {
	return nil
}

//向EBTree中插入数据
//special value被存储在特定结构中
//其他值正常存储在tree中
func (t *EBTree) InsertData(n EBTreen, pos uint8, parent *internalNode, value []byte, da []byte) (bool, *internalNode, error) {
	log.Info("into insert data,vaule is:")
	fmt.Println(value)
	//判断value是否special

	sp, p := t.isSpecial(value)

	if sp {
		//将对应data存入special中
		t.special[p].data = append(t.special[p].data, da)
		t.special[p].dirty = true
		return true, parent, nil
	}
	switch nt := n.(type) {
	case *leafNode:
		//向叶子节点插入数据
		//若当前节点为空时，直接插入节点。
		log.Info("insertdata:leafnode id:")
		fmt.Println(nt.Id)
		log.Info("insertdata:leafnode data length:")
		fmt.Println(len(nt.Data))
		outid := t.OutputRoot()
		log.Info("insertdata:leafnode:now tree root is")
		fmt.Println(outid)
		if len(nt.Data) == 0 {
			log.Info("the data is nil")
			//create a data item for da
			dai, err := createData(value, da)
			if err != nil {
				err = wrapError(err, "insert data: create data wrong")
				return false, parent, err
			}
			nt.Data = append(nt.Data, dai)
			n = nt
			log.Info("create data for leafnode,data length is:")
			l := len(nt.Data)
			fmt.Println(l)
			return true, parent, nil
		}

		//用于标记该value是否被插入成功
		flag := false
		//遍历当前节点的所有data，将da插入合适的位置
		//value一定小于或等于当前节点到最大值

		for i := 0; i < len(nt.Data); i++ {
			log.Info("find the appropriate position in nt datas")

			switch dt := (nt.Data[i]).(type) {
			case dataEncode:
				//decode the data
				log.Info("data is encoded,dt is:")
				fmt.Println(dt)
				d, err := decodeData(dt)
				if err != nil {
					return false, parent, err
				}
				fmt.Println("it's the %dth value", i)
				if bytes.Compare(d.Value, value) < 0 {
					//EBTree叶子节点按升序排列，应该data应该插入到nt.data[i]之后
					continue
				} else if bytes.Compare(d.Value, value) == 0 {
					//EBTree中已经存储了该value，因此，只要把data加入到对应到datalist中即可
					d.Keylist = append(d.Keylist, da)
					flag = true
					nt.Data[i] = d
					return true, parent, nil
				} else {
					nt.Data[i] = d
					//说明EBTree中不存在value值，此时，需要构建data，并将其加入到节点中
					su, nt, err := moveData(nt, i)
					if !su {
						err = wrapError(err, "insert data: move data wrong")
						return false, parent, err
					}
					nt.Data[i], err = createData(value, da)
					if err != nil {
						err = wrapError(err, "insert data: create data wrong")
						return false, parent, err
					}
					//split the leaf node
					if len(nt.Data) > int(maxLeafNodeCount) {
						su, parent, err := t.split(nt, parent)
						if !su {
							err = wrapError(err, "insert data: split the leaf node wrong")
							return false, parent, err
						}
						return true, parent, nil
					}
					flag = true
					return true, parent, nil
				}
				return true, parent, nil
			case data:
				log.Info("data is not encoded,dt value:")
				fmt.Println(dt.Value)
				log.Info(" the value will be inserted is::")
				fmt.Println(value)
				if bytes.Compare(dt.Value, value) < 0 {
					//EBTree叶子节点按升序排列，应该data应该插入到nt.data[i]之后
					continue
				} else if bytes.Compare(dt.Value, value) == 0 {
					//EBTree中已经存储了该value，因此，只要把data加入到对应到datalist中即可
					dt.Keylist = append(dt.Keylist, da)
					nt.Data[i] = dt
					n = nt
					flag = true
					return true, parent, nil
				} else {
					//说明EBTree中不存在value值，此时，需要构建data，并将其加入到节点中
					su, nt, err := moveData(nt, i)
					if !su {
						err = wrapError(err, "insert data: move data wrong")
						return false, parent, err
					}
					nt.Data[i], err = createData(value, da)
					if err != nil {
						err = wrapError(err, "insert data: create data wrong")
						return false, parent, err
					}
					//split the leaf node
					if len(nt.Data) > int(maxLeafNodeCount) {
						su, parent, err := t.split(nt, parent)
						if !su {
							err = wrapError(err, "insert data: split the leaf node wrong")
							return false, parent, err
						}
						return true, parent, nil
					}
					flag = true
					return true, parent, nil
				}
			case *data:
				log.Info("pointer data type")
				if bytes.Compare(dt.Value, value) < 0 {
					log.Info("find next value")
					//EBTree叶子节点按升序排列，应该data应该插入到nt.data[i]之后
					continue
				} else if bytes.Compare(dt.Value, value) == 0 {
					log.Info("find the same value")
					//EBTree中已经存储了该value，因此，只要把data加入到对应到datalist中即可
					dt.Keylist = append(dt.Keylist, da)
					nt.Data[i] = dt
					n = nt
					flag = true
					return true, parent, nil
				} else {
					log.Info("find a larger value")
					//说明EBTree中不存在value值，此时，需要构建data，并将其加入到节点中
					su, nt, err := moveData(nt, i)
					if !su {
						err = wrapError(err, "insert data: move data wrong")
						return false, parent, err
					}
					nt.Data[i], err = createData(value, da)
					if err != nil {
						err = wrapError(err, "insert data: create data wrong")
						return false, parent, err
					}
					//split the leaf node
					if len(nt.Data) > int(maxLeafNodeCount) {
						su, parent, err := t.split(nt, parent)
						if !su {
							err = wrapError(err, "insert data: split the leaf node wrong")
							return false, parent, err
						}
						return true, parent, nil
					}
					flag = true
					return true, parent, nil
				}
			case *dataEncode:
				log.Info("pointer data encode type")
			default:
				log.Info("data type is not appropriate")
				err := errors.New("data type is not appropriate")
				return false, nil, err
			}

		}
		if !flag { //将该值插入到节点末尾
			log.Info("the data should be put in the last ")
			dai, err := createData(value, da)
			if err != nil {
				err = wrapError(err, "insert data: when the data was added into the end of node, create data wrong")
				return false, parent, err
			}

			nt.Data = append(nt.Data, dai)

			//如果更新的是最大值，应该同时更新children.value
			//如果parent为空，那么不需要进行更新
			if parent != nil {
				_, parent, res, err := getLeafNodePosition(nt, parent, t)
				if err != nil {
					wrapError(err, "insert data: when the node is leaf node, get leaf node postion wrong")
					return false, parent, err
				}
				switch ct := (parent.Children[res]).(type) {
				case childEncode:
					return false, parent, err
				case child:
					ct.Value = value
				default:
					return false, nil, nil
				}
			}

			//split the leaf node
			if len(nt.Data) > int(maxLeafNodeCount) {
				su, parent, err := t.split(nt, parent)
				if !su {
					err = wrapError(err, "insert data: when the data was added into the end of node, split leaf node wrong")
					return false, parent, err
				}
				return true, parent, nil
			}
			return true, nil, nil
		}
	case *internalNode:
		log.Info("insertdata:next is internalnoode id")
		fmt.Println(nt.Id)
		flag := false
		var i int
		for i = 0; i < len(nt.Children); i++ {
			switch ct := (nt.Children[i]).(type) {
			case childEncode:
				//decode child
				cd, _, err := decodeChild(ct)
				decoden, err := t.resolveHash(ct)
				if err != nil {
					return false, nil, err
				}
				cd.Pointer = decoden
				nt.Children[i] = cd
				if bytes.Compare(cd.Value, value) < 0 {
					continue
				} else {
					//call the insert function to
					su, _, err := t.InsertData(decoden, uint8(i), nt, value, da)
					if !su {
						err = wrapError(err, "insert data: when the data was added into appropriate child, something wrong")
						return false, parent, err
					}

					if len(nt.Children) > int(maxInternalNodeCount) {
						_, parent, err := t.split(nt, parent)
						if err != nil {
							err = wrapError(err, "insert data: when the data was added into appropriate child, split internal node wrong")
							return false, parent, err
						}
					}
					flag = true
					return true, parent, nil
				}
			case child:
				if bytes.Compare(ct.Value, value) < 0 {
					continue
				} else {

					//call the insert function to
					su, _, err := t.InsertData(ct.Pointer, uint8(i), nt, value, da)
					if !su {
						err = wrapError(err, "insert data: when the data was added into appropriate child, something wrong")
						return false, parent, err
					}

					if len(nt.Children) > int(maxInternalNodeCount) {
						_, parent, err := t.split(nt, parent)
						if err != nil {
							err = wrapError(err, "insert data: when the data was added into appropriate child, split internal node wrong")
							return false, parent, err
						}
					}
					flag = true
					return true, parent, nil
				}
			default:
				return false, nil, nil
			}

		}
		if !flag {
			//TODO:将该值插入到节点末尾
			//call the insert function to
			switch ct := (nt.Children[i-1]).(type) {
			case childEncode:
				//decode child
				cd, _, err := decodeChild(ct)
				decoden, err := t.resolveHash(ct)
				if err != nil {
					return false, nil, err
				}
				cd.Pointer = decoden
				su, _, err := t.InsertData(decoden, uint8(i), nt, value, da)
				if !su {
					err = wrapError(err, "insert data: when the data was added into last child, something wrong")
					return false, parent, err
				}
				cd.Value = value
				nt.Children[i] = cd
				if len(nt.Children) > int(maxInternalNodeCount) {
					su, parent, err := t.split(nt, parent)
					if !su {
						err = wrapError(err, "insert data: when the data was added into last child, split internal node wrong")
						return false, parent, err
					}
					return true, parent, nil
				}

				return true, parent, nil
			case child:
				su, _, err := t.InsertData(ct.Pointer, uint8(i), nt, value, da)
				if !su {
					err = wrapError(err, "insert data: when the data was added into last child, something wrong")
					return false, parent, err
				}
				ct.Value = value
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
				return false, nil, nil
			}

		}
	case *ByteNode:
		nbb, _ := n.cache()
		log.Info("insertdata:next is bytenode id")
		fmt.Println(nbb)
		if string(nbb) == "" {
			dai, err := createData(value, da)
			if err != nil {
				err = wrapError(err, "insert data: create data wrong")
				return false, parent, err
			}
			var da []data
			da = append(da, dai)
			newn, err := CreateLeafNode(t, da)
			log.Info("next is newn leaf node id:")
			fmt.Println(newn.Id)
			t.Root = &newn
			outid := t.OutputRoot()
			log.Info("next is the root of tree:")
			fmt.Println(outid)
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

func (t *EBTree) resolveHash(n []byte) (EBTreen, error) {
	log.Info("into resolve hash func")
	fmt.Println(n)
	cacheMissCounter.Inc(1)

	if node := t.Db.node(n, t.cachegen); node != nil {
		return node, nil
	}
	log.Info("not get the node from db")
	return nil, &MissingNodeError{NodeId: n, Path: nil}
}

func BytesToInt(b []byte) (i uint64) {
	return binary.BigEndian.Uint64(b)
}

func IntToBytes(i uint64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, i)
	return buf
}

func (t *EBTree) newSequence() ([]byte, error) {
	log.Info("into new sequece")
	log.Info(string(t.sequence))
	re := BytesToInt(t.sequence)
	re = re + 1
	if re < 0 {
		err := errors.New("BytesToInt return a negtive data")
		return nil, err
	}
	id := IntToBytes(uint64(re))
	t.sequence = id
	log.Info(string(t.sequence))
	return id, nil
}
func (t *EBTree) OutputRoot() []byte {
	switch rt := (t.Root).(type) {
	case *leafNode:
		log.Info("outputroot:root node type: leafnode")
		log.Info(string(rt.Id))
		return rt.Id
	case *internalNode:
		log.Info("outputroot:root node type: internalnode")
		return rt.Id
	case *ByteNode:
		log.Info("outputroot:root node type: bytenode")
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

	if t.Db == nil {
		panic("commit called on trie with nil database")
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
	log.Info(string(rt))
	return rt, nil
}

func (t *EBTree) foldRoot(db *Database, onleaf LeafCallback) (EBTreen, error) {
	log.Info("into fold root ")
	if t.Root == nil {
		err := errors.New("tree is nil")
		return nil, err
	}
	f := newFolder(t.cachegen, t.cachelimit, onleaf)
	defer returnFolderToPool(f)
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
						if bytes.Compare(ctt.Value, value) < 0 {
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
				if bytes.Compare(value, dt.Value) < 0 || bytes.Compare(value, dmt.Value) > 0 {
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
						if bytes.Compare(dt.Value, value) == 0 {
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
				if bytes.Compare(ct.Value, value) >= 0 {
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

func testInsert(nt *leafNode, i uint64) (bool, error) {
	//test if the data is inserted right
	switch dt := (nt.Data[i]).(type) {
	case dataEncode:
		err := errors.New("insertData in  leaf node:data[i] is encoded.")
		return false, err
	case data:
		if i > 0 {
			switch ddt := (nt.Data[i-1]).(type) {
			case dataEncode:
				err := errors.New("insertData in  leaf node:data[i-1] is encoded.")
				return false, err
			case data:
				if bytes.Compare(ddt.Value, dt.Value) >= 0 {
					err := errors.New("insertData in leaf node: smaller than last data")
					return false, err
				}
				return true, nil
			default:
				err := errors.New("insertData in  leaf node:data[i-1] is in wrong format.")
				return false, err

			}
		}

	default:
		err := errors.New("insertData in  leaf node:data[i] is in wrong format.")
		return false, err
	}
	return true, nil
}
