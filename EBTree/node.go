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

package EBTree

import (
	"bytes"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"github.com/pkg/errors"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/rlp"
)

var indices = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f", "[17]"}

type EBTreen interface {
	fstring(string) string
	cache() ([]byte, bool)
}

type (
	//中间节点
	internalNode struct {
		Children []ChildInterface // Actual ebtree internal node data to encode/decode (needs custom encoder)
		Id       []byte
		Dirty    bool
	}
	//叶子节点
	leafNode struct {
		Data  []DataInterface
		Next  EBTreen
		Id    []byte
		Dirty bool
	}
	idNode struct {
		Id       []byte
		NodeData []byte
	}
	ByteNode []byte
)

func (n *internalNode) fstring(ind string) string {
	resp := fmt.Sprintf("[\n%s  ", ind)
	for i := 0; i < len(n.Children); i++ {
		//TODO:fstring the node of ebtree
		resp += fmt.Sprintf("%s", indices[i])
	}

	return resp + fmt.Sprintf("\n%s] ", ind)
}
func (n *leafNode) fstring(ind string) string {

	resp := fmt.Sprintf("[\n%s  ", ind)
	for i := 0; i < len(n.Data); i++ {
		//TODO:fstring the  leaf node of ebtree
		resp += fmt.Sprintf("%s", indices[i])
	}

	return resp + fmt.Sprintf("\n%s] ", ind)
}
func (n *idNode) fstring(ind string) string {
	return fmt.Sprintf("<%x> ", string(n.NodeData))
}
func (n *ByteNode) fstring(ind string) string {
	return fmt.Sprintf("<%x> ", n)
}

func (n *leafNode) copy() *leafNode         { copy := *n; return &copy }
func (n *internalNode) copy() *internalNode { copy := *n; return &copy }
func (n *child) copy() *child               { copy := *n; return &copy }
func (n *data) copy() *data                 { copy := *n; return &copy }

//获取节点ID
func (n *internalNode) cache() ([]byte, bool) { return n.Id, n.Dirty }
func (n *leafNode) cache() ([]byte, bool)     { return n.Id, n.Dirty }
func (n *idNode) cache() ([]byte, bool)       { return n.Id, true }
func (n ByteNode) cache() ([]byte, bool)      { return n, true }

type ChildInterface interface {
	childString(string) string
}
type (
	child struct {
		Value   []byte
		Pointer EBTreen
	}
)
type childEncode []byte

func (n childEncode) childString(ind string) string {
	return fmt.Sprintf("<%x> ", string(n))
}

func (n child) childString(ind string) string {
	return fmt.Sprintf("<%x> ", string(n.Value))
}

type searchData struct {
	value []byte
	data  []byte
}

type SearchValue struct {
	Value []byte
	Data  [][]byte
}

type DataInterface interface {
	dataString(string) string
}

type (
	data struct {
		Keylist [][]byte
		Value   []byte
	}
)

type dataEncode []byte

func (n dataEncode) dataString(ind string) string {
	return fmt.Sprintf("<%x> ", string(n))
}

func (n data) dataString(ind string) string {
	return fmt.Sprintf("<%x> ", string(n.Value))
}

// wraps a decoding error with information about the path to the
// invalid child node (for debugging encoding issues).
type decodeError struct {
	what  error
	stack []string
}

func wrapError(err error, ctx string) error {
	if err == nil {
		return nil
	}
	if decErr, ok := err.(*decodeError); ok {
		decErr.stack = append(decErr.stack, ctx)
		return decErr
	}
	return &decodeError{err, []string{ctx}}
}

func CreateLeafNode(tree *EBTree, datalist []data) (leafNode, error) {
	//log.Info("into create leaf node")
	var empty []byte
	se, err := tree.newSequence()
	//log.Info(string(se))
	if err != nil {
		err = wrapError(err, "CreateLeafNode")
		return leafNode{}, err
	}
	newn, err := constructLeafNode(se, uint8(len(datalist)), datalist, false, true, nil, empty, 0)
	return newn, err
}

func constructLeafNode(id []byte, count uint8, datalist []data, special bool, dirty bool, next EBTreen, nexid []byte, gen uint16) (leafNode, error) {
	//log.Info("into construct leaf node")
	//newn := &leafNode{}
	var dataInter []DataInterface
	for i := uint8(0); i < count; i++ {
		dataInter = append(dataInter, &datalist[i])
	}
	newn := leafNode{Data: dataInter, Id: id, Next: next, Dirty: false}
	return newn, nil
}

func createInternalNode(tree *EBTree, children []ChildInterface) (*internalNode, error) {
	se, err := tree.newSequence()
	newn, _ := constructInternalNode(se, uint8(len(children)), children, 0)
	if err != nil {
		err = wrapError(err, "createInternalNode")
		return &newn, err
	}
	return &newn, err
}

func constructInternalNode(id []byte, count uint8, childlist []ChildInterface, gen uint16) (internalNode, error) {
	//newn := &leafNode{}
	newn := internalNode{Id: id, Children: childlist, Dirty: false}
	return newn, nil
}

func addChild(internal internalNode, chil child, position int) (bool, internalNode, error) {

	internal.Children = append(internal.Children, child{})

	var i int

	for i = len(internal.Children) - 1; i > position; i-- {
		internal.Children[i] = internal.Children[i-1]
	}
	internal.Children[i] = chil

	return true, internal, nil
}

func moveData(n *leafNode, pos int) (bool, *leafNode, error) {
	//log.Info("into moveData")
	n.Data = append(n.Data, data{})

	for i := len(n.Data) - 1; i > pos; i-- {
		n.Data[i] = n.Data[i-1]
	}
	return true, n, nil
}

func collapsedLeafNode(nt *leafNode) (*leafNode, error) {
	//log.Info("encode a leaf node")
	var collapsed leafNode
	if nt.Id == nil {
		err := errors.New("empty node")
		return nil, err
	}
	collapsed.Id = nt.Id
	da, err := CopyData(nt.Data)
	if err != nil {
		return nil, err
	}
	collapsed.Data = da

	if nt.Next != nil {
		switch cnt := (nt.Next).(type) {
		case *leafNode:
			//log.Info("fold:collapsedNode:leafnode")
			var nb ByteNode
			nb = cnt.Id
			if len(nb) == 0 {
				fmt.Println("wrong in func : collapsedLeafNode.297")
			}
			collapsed.Next = &nb
		case *internalNode:
			//log.Info("fold:collapsedNode:internalnode")
			var nb ByteNode
			nb = cnt.Id
			if len(nb) == 0 {
				fmt.Println("wrong in func : collapsedLeafNode.305")
			}
			collapsed.Next = &nb
		case *ByteNode:
			//log.Info("fold:collapsedNode:bytenode")
			var nb ByteNode
			nb, _ = cnt.cache()
			if len(nb) == 0 {
				fmt.Println("wrong in func : collapsedLeafNode.313")
			}
			collapsed.Next = &nb
		default:
			err := errors.New("fold: wrong collapsed node type")
			return nil, err
		}
	}
	return &collapsed, nil
}

func createChild(val []byte, po EBTreen) (child, error) {
	ch := &child{Value: val, Pointer: po}
	if ch == nil {
		err := errors.New("create child failed")
		return *ch, err
	}
	return *ch, nil
}

//将定位当前叶子节点在parent节点中的位置，便于后期的插入、查找
func getLeafNodePosition(n *leafNode, parent *internalNode, t *EBTree) (bool, *internalNode, uint8, error) {
	//说明此时分割的是作为树根节点的叶子节点
	if parent == nil {
		datainter := n.Data[len(n.Data)-1]
		switch dt := (datainter).(type) {
		case dataEncode:
			err := errors.New("wrong data type:dataEncoded in getLeafNodePosition")
			return false, nil, 0, err
		case data:
			//需要为上级根节点确定value的值
			chil, err := createChild(dt.Value, n)
			if err != nil {
				err = wrapError(err, "get leaf node position wrong: when parent is nil, create child wrong")
				return false, parent, uint8(0), err
			}
			var children []ChildInterface
			children = append(children, chil)
			parent, err = createInternalNode(t, children)
			if err != nil {
				err = wrapError(err, "get leaf node position wrong: when parent is nil, create root")
				return false, parent, uint8(0), err
			}
			t.Root = parent
			return true, parent, uint8(0), nil
		default:
			fmt.Println(dt)
			err := errors.New("wrong data type:dataEncoded in getLeafNodePosition")
			return false, nil, 0, err
		}
	}
	var re uint8
	flag := false
	//此时需要确定当前叶子节点是parent的第几个子节点
	for i, pc := range parent.Children {
		switch ct := (pc).(type) {
		case childEncode:
			err := errors.New("wrong children type:childEncode in getLeafNodePosition")
			return false, nil, 0, err
		case child:
			//找到这个节点之后，需要更新value值，并插入新节点
			switch ctpt := (ct.Pointer).(type) {
			case *leafNode:
				//fmt.Printf("this node id is:%v, n id is %v\n", ctpt.Id, n.Id)
				if ctpt == n {
					switch dt := (n.Data[len(n.Data)-1]).(type) {
					case dataEncode:
						err := errors.New("data is encoded in getLeafNodePosition")
						return false, nil, 0, err
					case data:
						//找到之后返回对应的节点序号
						re = uint8(i)
						ct.Value = dt.Value
						parent.Children[i] = pc
						flag = true
						return true, parent, re, nil
					case *dataEncode:
						err := errors.New("data is encoded in getLeafNodePosition")
						return false, nil, 0, err
					case *data:
						//找到之后返回对应的节点序号
						re = uint8(i)
						ct.Value = dt.Value
						parent.Children[i] = pc
						flag = true
						return true, parent, re, nil
					default:
						fmt.Println("data type error in getLeafNodePosition")
						err := errors.New("data type error in getLeafNodePosition")
						return false, nil, 0, err
					}
				}
			case *internalNode:
				err := errors.New("wrong child type:internalnode error in getLeafNodePosition")
				return false, nil, 0, err
			case *ByteNode:
				le, err := t.ResolveByteNode(ctpt)
				ct.Pointer = le
				parent.Children[i] = ct
				//replace the bytenode into leafode
				if err != nil {
					wrapError(err, "wrong ini resolve byte node when get leafnode position")
					return false, nil, 0, err
				}
				//change the parent's child value
				//fmt.Printf("this node id is:%v, n id is %v\n", le.Id, n.Id)
				//fix bug.when n is split, the length of n'data is different from le. to get the position, we just to compare the node id
				if Compare(le.Id, n.Id) == 0 {
					switch dt := (n.Data[len(n.Data)-1]).(type) {
					case dataEncode:
						err := errors.New("data is encoded in getLeafNodePosition")
						return false, nil, 0, err
					case data:
						//找到之后返回对应的节点序号
						re = uint8(i)
						ct.Value = dt.Value
						parent.Children[i] = ct
						flag = true
						return true, parent, re, nil
					case *dataEncode:
						err := errors.New("data is encoded in getLeafNodePosition")
						return false, nil, 0, err
					case *data:
						//找到之后返回对应的节点序号
						re = uint8(i)
						ct.Value = dt.Value
						parent.Children[i] = ct
						flag = true
						return true, parent, re, nil
					default:
						fmt.Println("data type error in getLeafNodePosition")
						err := errors.New("data type error in getLeafNodePosition")
						return false, nil, 0, err
					}
				}
			default:
				err := errors.New("wrong child type:default error in getLeafNodePosition")
				return false, nil, 0, err
			}

		default:
			fmt.Println("child type error in getLeafNodePosition")
			err := errors.New("child type error in getLeafNodePosition")
			return false, nil, 0, err
		}

	}
	if !flag {
		err := errors.New("there is no such leaf node in this parent node")
		return false, parent, re, err
	}
	return true, parent, re, nil
}

//
func getInternalNodePosition(n *internalNode, parent *internalNode, t *EBTree) (bool, *internalNode, uint8, error) {
	switch ct := (n.Children[len(n.Children)-1]).(type) {
	case childEncode:
		err := errors.New("wrong child type:childEncode in getInternalNodePosition")
		return false, parent, 0, err
	case child:
		if parent == nil {
			chil, err := createChild(ct.Value, n)
			if err != nil {
				err = wrapError(err, "get internal node position wrong: when parent is nil, create child wrong")
				return false, parent, uint8(0), err
			}
			var children []ChildInterface
			children = append(children, chil)
			parent, err = createInternalNode(t, children)
			if err != nil {
				err = wrapError(err, "get internal node position wrong: when parent is nil, create root")
				return false, parent, uint8(0), err
			}
			t.Root = parent
			return true, parent, uint8(0), nil

		}
		var re uint8
		flag := false
		for i, pc := range parent.Children {
			switch cpt := (pc).(type) {
			case childEncode:
				err := errors.New("wrong child type:childEncode in  getInternalNodePosition")
				return false, nil, 0, err
			case child:
				switch cptct := (cpt.Pointer).(type) {
				case *leafNode:
					err := errors.New("wrong child type:leafnode error in getInternalNodePosition")
					return false, nil, 0, err
				case *internalNode:
					if cpt.Pointer == n {
						re = uint8(i)
						cpt.Value = ct.Value
						parent.Children[i] = pc
						flag = true
						//fmt.Println(pc.value)
					}
				case *ByteNode:
					nid, _ := cptct.cache()
					in, err := t.resolveHash(nid)
					ct.Pointer = in
					parent.Children[i] = ct
					//replace the bytenode into leafode
					if err != nil {
						wrapError(err, "wrong ini resolve byte node when get internalNode position")
						return false, nil, 0, err
					}
					//change the parent's child value
					//fmt.Printf("n id is %v\n", n.Id)
					if in == n {
						re = uint8(i)
						cpt.Value = ct.Value
						parent.Children[i] = pc
						flag = true
						//fmt.Println(pc.value)
					}
				default:
					err := errors.New("wrong child pointer type:default in  getInternalNodePosition")
					return false, nil, 0, err
				}

			default:
				err := errors.New("wrong child type:default in  getInternalNodePosition")
				return false, nil, 0, err
			}

		}
		if !flag {
			err := errors.New("there is no such internal node in this parent node")
			return false, parent, re, err
		}
		return true, parent, re, nil
	default:
		err := errors.New("wrong child type:default in getInternalNodePosition")
		return false, parent, 0, err
	}

}

func createData(value []byte, da []byte) (data, error) {
	//log.Info("into createData")
	//create a data item for da
	var kel [][]byte
	//create a key list for data item
	kel = append(kel, da)
	dai := constructData(value, kel)
	/*if dai == data{} {
		err := errors.New("create data failed")
		if err != nil {
			return dai, err
		}
	}*/
	return dai, nil
}

func (err *decodeError) Error() string {
	return fmt.Sprintf("%v (decode path: %s)", err.what, strings.Join(err.stack, "<-"))
}

//在叶子节点中搜索value对应到节点
func SearchLeafNode(value []byte, n *leafNode) ([][]byte, error) {

	for i := 0; i < len(n.Data); i++ {
		switch dt := (n.Data[i]).(type) {
		case dataEncode:
			//TODO
			return nil, nil
		case data:
			if bytes.Equal(dt.Value, value) {
				return dt.Keylist, nil
			}

		case *data:
			if bytes.Equal(dt.Value, value) {
				return dt.Keylist, nil
			}
		default:
			return nil, nil
		}

	}
	err := errors.New("none data matches!")
	return nil, err

}

//在中间节点中搜索value对应到节点
func SearchInternalNode(value []byte, n *internalNode, t *EBTree) ([][]byte, error) {

	for i := 0; i < len(n.Children); i++ {
		switch ct := (n.Children[i]).(type) {
		case childEncode:
			return nil, nil
		case child:
			if Compare(ct.Value, value) >= 0 {
				switch dt := (ct.Pointer).(type) {
				case *leafNode:
					return SearchLeafNode(value, dt)
				case *internalNode:
					return SearchInternalNode(value, dt, t)
				case *ByteNode:
					dtc, _ := dt.cache()
					decoden, err := t.resolveHash(dtc)
					if err != nil {
						return nil, err
					}
					switch det := (decoden).(type) {
					case *leafNode:
						return SearchLeafNode(value, det)
					case *internalNode:
						return SearchInternalNode(value, det, t)
					case *ByteNode:
						err := errors.New("wrong det type")
						return nil, err
					default:
						err := errors.New("wrong det node type")
						return nil, err
					}
				}
			}
		}

	}

	err := errors.New("none data matches!")
	return nil, err

}

//在中间节点中搜索value对应到节点
func SearchNode(value []byte, n EBTreen, t *EBTree) ([][]byte, error) {
	su, p := t.isSpecial(value)
	if su {
		return t.special[p].data, nil
	}
	switch nt := (n).(type) {
	case *leafNode:
		result, err := SearchLeafNode(value, nt)
		if err != nil {
			err = wrapError(err, "search node: search leaf node wrong")
			return nil, err
		}
		return result, err
	case *internalNode:
		result2, err := SearchInternalNode(value, nt, t)
		if err != nil {
			err = wrapError(err, "search node: search leaf node wrong")
			return nil, err
		}
		return result2, err
	}
	return nil, nil

}

func findFirstNode(n EBTreen, tree *EBTree) (*leafNode, EBTreen, error) {
	if n == nil {
		err := errors.New("find first node wrong: t.root is nil")
		return nil, nil, err
	}
	switch nt := (n).(type) {
	case *leafNode:
		return nt, nt, nil
	case *internalNode:
		if len(nt.Children) <= 0 {
			err := errors.New("find first node wrong: when node is internal node, nt.count <=0")
			return nil, nil, err
		}
		switch ct := (nt.Children[0]).(type) {
		case childEncode:
			cd, _, err := decodeChild(ct)
			if err != nil {
				wrapError(err, "find first node :decode wrong:when node is internal node")
				return nil, nil, err
			}
			switch cpt := (cd.Pointer).(type) {
			case *ByteNode:
				cptid, _ := cpt.cache()
				decoden, err := tree.resolveHash(cptid)
				if err != nil {
					return nil, nil, err
				}
				cd.Pointer = decoden
				re, decoden, err := findFirstNode(decoden, tree)
				if err != nil {
					wrapError(err, "find first node wrong:when node is internal node")
					return nil, nil, err
				}
				nt.Children[0] = cd
				return re, nt, nil
			case *leafNode, *internalNode:
				cd.Pointer = cpt
				re, _, err := findFirstNode(cpt, tree)
				if err != nil {
					wrapError(err, "find first node wrong:when node is internal node")
					return nil, nil, err
				}
				nt.Children[0] = cd
				return re, nt, nil
			}

		case child:
			re, decoden, err := findFirstNode(ct.Pointer, tree)
			if err != nil {
				wrapError(err, "find first node wrong:when node is internal node")
				return nil, nil, err
			}
			ct.Pointer = decoden
			nt.Children[0] = ct
			return re, nt, nil
		}
	case *ByteNode:
		ntid, _ := nt.cache()
		decoden, _ := tree.resolveHash(ntid)

		return findFirstNode(decoden, tree)

	default:
		err := errors.New("wrong type ")
		return nil, nil, err
	}
	err := errors.New("wrong in func")
	return nil, nil, err
}

func (t *EBTree) findNode(n EBTreen, value []byte) (bool, *leafNode, uint8, error) {
	if n == nil {
		err := errors.New("find node wrong: t.root is nil")
		return false, nil, 0, err
	}

	switch nt := n.(type) {
	case *leafNode:
		//若当前节点为空时,返回空
		if len(nt.Data) == 0 {
			err := errors.New("find node wrong: n.count==0")
			return false, nil, uint8(0), err
		}

		for i := 0; i < len(nt.Data); i++ {
			switch dt := (nt.Data[i]).(type) {
			case data:
				if Compare(dt.Value, value) > 0 {
					//EBTree叶子节点按升序排列，继续向后查找
					continue
				} else {
					//找到节点
					return true, nt, uint8(i), nil
				}
			case *data:
				if Compare(dt.Value, value) > 0 {
					//EBTree叶子节点按升序排列，继续向后查找
					continue
				} else {
					//找到节点
					return true, nt, uint8(i), nil
				}
			default:
				err := errors.New("error in func findnode()")
				return false, nt, 0, err
			}
		}
	case *internalNode:
		var i int
		for i = 0; i < len(nt.Children); i++ {
			switch ct := (nt.Children[i]).(type) {
			case childEncode:
				return false, nil, 0, nil
			case child:
				if Compare(ct.Value, value) > 0 {
					continue
				} else {
					//call the find node function to
					if ct.Pointer != nil {
						su, re, po, err := t.findNode(ct.Pointer, value)
						if !su {
							err = wrapError(err, "insert data: when the data was added into appropriate child, something wrong")
							return false, re, 0, err
						}
						return true, re, po, nil
					} else {
						//there is no child in this position,which should be error
						return false, nil, 0, nil

					}
				}
			}
		}
	case *ByteNode:
		ntid, _ := nt.cache()
		decoden, _ := t.resolveHash(ntid)

		return t.findNode(decoden, value)
	default:
		err := errors.New("the node is not leaf or internal, something wrong")
		return false, nil, 0, err
	}
	//err := errors.New("find node wrong, something wrong")
	return false, nil, 0, nil
}

func (t *EBTree) ResolveByteNode(dt *ByteNode) (*leafNode, error) {
	dtc, _ := dt.cache()
	decoden, err := t.resolveHash(dtc)
	if err != nil {
		return nil, err
	}
	switch det := (decoden).(type) {
	case *leafNode:
		return det, nil
	case *internalNode:
		err := errors.New("the next of leaf node should be leafnode, something wrong")
		return nil, err
	case *ByteNode:
		err := errors.New("wrong det type:there must be problems in resolve func")
		return nil, err
	default:
		err := errors.New("wrong det node type")
		return nil, err
	}
}

func AddToSearchValue(d DataInterface, bn []byte, value []byte) (SearchValue, error, bool) {
	var result SearchValue
	switch dt := d.(type) {
	case dataEncode:
		err := errors.New("topkvsearch:data is encoded")
		return result, err, false
	case data:
		if Compare(dt.Value, value) == 0 {
			var ds [][]byte
			for _, kl := range dt.Keylist {

				var ss string
				rlp.DecodeBytes(kl, &ss)
				sss := strings.Split(ss, ",")
				num, _ := strconv.Atoi(sss[0])
				if Compare(IntToBytes(uint64(num)), bn) <= 0 {
					ds = append(ds, kl)
				}

			}
			if len(ds) > 0 {
				result.Value = dt.Value
				result.Data = ds
				return result, nil, false
			} else {
				err := errors.New("no data found")
				return result, err, false
			}

		} else if Compare(dt.Value, value) < 0 {
			err := errors.New("no data found")
			return result, err, false
		} else {
			return result, nil, true
		}
	case *dataEncode:
		err := errors.New("topkvsearch:data is encoded")
		return result, err, false
	case *data:
		if Compare(dt.Value, value) == 0 {

			var ds [][]byte
			for _, kl := range dt.Keylist {

				var ss string
				rlp.DecodeBytes(kl, &ss)
				sss := strings.Split(ss, ",")
				num, _ := strconv.Atoi(sss[0])
				if Compare(IntToBytes(uint64(num)), bn) <= 0 {
					ds = append(ds, kl)
				}

			}
			if len(ds) > 0 {
				result.Value = dt.Value
				result.Data = ds
				return result, nil, false
			} else {
				err := errors.New("no data found")
				return result, err, false
			}
		} else if Compare(dt.Value, value) < 0 {
			err := errors.New("no data found")
			return result, err, false
		} else {
			return result, nil, true
		}
	default:
		err := errors.New("topkdatasearch:data in wrong type")
		return SearchValue{}, err, false
	}
	return result, nil, true

}

//top-k data search
func (t *EBTree) TopkVSearch(k []byte, bn []byte, max bool) (bool, []SearchValue, error) {
	var result []SearchValue
	if max {
		n, _, err := findFirstNode(t.Root, t)
		if err != nil {
			err = wrapError(err, "top-k search data wrong:find first node wrong")
			return false, nil, err
		}
		_ = n
		_ = result
		b := false

		notCompareBlockNum := false
		if Compare(bn, IntToBytes(0)) <= 0 {
			notCompareBlockNum = true
		}
		for {
			if b || n == nil || Compare(IntToBytes(uint64(len(result))), k) >= 0 {
				break
			}
			flag := false
			for i := 0; i < len(n.Data); i++ {
				switch dt := (n.Data[i]).(type) {
				case dataEncode:
					err := errors.New("topkvsearch:data is encoded")
					return false, nil, err
				case data:
					var tmpkl [][]byte
					for _, kl := range dt.Keylist {
						var ss string
						rlp.DecodeBytes(kl, &ss)
						sss := strings.Split(ss, ",")
						num, _ := strconv.Atoi(sss[0])
						if Compare(IntToBytes(uint64(len(result))), k) < 0 {
							if notCompareBlockNum || Compare(IntToBytes(uint64(num)), bn) <= 0 {
								tmpkl = append(tmpkl, kl)
							}
						} else {
							flag = true
							break
						}
					}
					if tmpkl != nil {
						r := SearchValue{dt.Value, tmpkl}
						result = append(result, r)
					}
				case *dataEncode:
					err := errors.New("topkvsearch:data is encoded")
					return false, nil, err
				case *data:
					var tmpkl [][]byte
					for _, kl := range dt.Keylist {
						var ss string
						rlp.DecodeBytes(kl, &ss)
						sss := strings.Split(ss, ",")
						num, _ := strconv.Atoi(sss[0])
						if Compare(IntToBytes(uint64(len(result))), k) < 0 {
							if notCompareBlockNum || Compare(IntToBytes(uint64(num)), bn) <= 0 {
								tmpkl = append(tmpkl, kl)
							}
						} else {
							flag = true
							break
						}
					}
					if tmpkl != nil {
						r := SearchValue{dt.Value, tmpkl}
						result = append(result, r)
					}
				default:
					err := errors.New("topkdatasearch:data in wrong type")
					return false, nil, err
				}

			}
			if n.Next == nil {
				break
			}
			if !flag {
				switch nnt := (n.Next).(type) {
				case *leafNode:
					n = nnt
				case *ByteNode:
					if nnt == nil {
						b = true
						fmt.Println("nnt is nil")
						break
					} else {
						nntid, _ := nnt.cache()
						if len(nntid) == 0 {
							b = true
							fmt.Println("nnt's length is 0")
							break
						} else {
							den, err := t.ResolveByteNode(nnt)
							if err != nil {
								return false, nil, err
							}
							n = den
						}
					}
				default:
					err := errors.New("topkdatasearch:wrong n.nexty type")
					return false, nil, err
				}
			} else {
				break
			}
		}
		//fmt.Println("out of for")
		if Compare(IntToBytes(uint64(len(result))), k) < 0 {
			err = wrapError(err, "top-k search data wrong:not enough data")
			return false, result, err
		} else if Compare(IntToBytes(uint64(len(result))), k) > 0 {
			err = wrapError(err, "top-k search data wrong:get too much data")
			return false, result, err
		} else {
			return true, result, nil
		}

	}
	fmt.Println("wrong in topkvSearch")
	err := errors.New("wrong in topkvsearch")
	return false, nil, err
}

var addToSearchValueTime int64

func GetAddToSearch() int64{
	return addToSearchValueTime
}
//find value in node
func (t *EBTree) SpecificValueSearchInNode(n EBTreen, value []byte, bn []byte) (SearchValue, error) {
	var result SearchValue
	switch nt := (n).(type) {
	case *leafNode:
		for i := 0; i < len(nt.Data); i++ {
			t1 := time.Now()
			result, err, flag := AddToSearchValue(nt.Data[i], bn, value)
			t2 := time.Now()
			t3 := t2.Sub(t1).Microseconds()
			addToSearchValueTime = addToSearchValueTime + t3

			if err != nil {
				return result, err
			}
			if !flag {
				return result, nil
			}
		}
		return result, nil

	case *internalNode:
		for i := 0; i < len(nt.Children); i++ {
			switch ct := (nt.Children[i]).(type) {
			case child:
				switch ctpt := (ct.Pointer).(type) {
				case *ByteNode:
					ctptid, _ := ctpt.cache()
					decoden, err := t.resolveHash(ctptid)
					if err != nil {
						return result, err
					}
					ct.Pointer = decoden
					nt.Children[i] = ct
				}
				if Compare(value, ct.Value) >= 0 {
					cpt := ct.Pointer
					return t.SpecificValueSearchInNode(cpt, value, bn)
				}
			case *child:
				switch ctpt := (ct.Pointer).(type) {
				case *ByteNode:
					ctptid, _ := ctpt.cache()
					decoden, err := t.resolveHash(ctptid)
					if err != nil {
						return result, err
					}
					ct.Pointer = decoden
					nt.Children[i] = ct
				}
				if Compare(value, ct.Value) >= 0 {
					cpt := ct.Pointer
					return t.SpecificValueSearchInNode(cpt, value, bn)
				}
			default:
				err := errors.New("wrong children type in SpecificValueSearchInNode")
				return result, err
			}

		}
	default:
		err := errors.New("wrong node type: default in SpecificValueSearchInNode")
		return result, err
	}

	return result, nil
}

//search in leafNode data
func (t *EBTree) RangeValueSearchLeaf(dt *data, max []byte, bn []byte, result []SearchValue) (bool, bool, []SearchValue, error) {
	if Compare(IntToBytes(uint64(len(result))), bn) < 0 && Compare(dt.Value, max) <= 0 {
		r := SearchValue{dt.Value, dt.Keylist}
		result = append(result, r)
		return false, true, result, nil
	} else {
		return true, true, result, nil
	}
}

//specific value search
func (t *EBTree) SpecificValueSearch(value []byte, bn []byte) (SearchValue, error) {
	var result SearchValue
	switch rt := (t.Root).(type) {
	case *ByteNode:
		rtid, _ := rt.cache()
		decoden, err := t.resolveHash(rtid)
		if err != nil {
			return result, err
		}
		t.Root = decoden

	}
	return t.SpecificValueSearchInNode(t.Root, value, bn)

}

//range value search
func (t *EBTree) RangeValueSearch(min []byte, max []byte, bn []byte) (bool, []SearchValue, error) {

	var result []SearchValue
	notCompareBlockNum := false
	if Compare(bn, IntToBytes(0)) <= 0 {
		notCompareBlockNum = true
	}

	su, n, pos, err := t.findNode(t.Root, max)
	if !su {
		wrapError(err, "range search value wrong:find node wrong")
		return false, nil, err
	}
	//find the result in range
	for i := int(pos); i < len(n.Data); i++ {
		switch dt := (n.Data[i]).(type) {

		case data:
			if Compare((&dt).Value, min) < 0 {
				return su, result, err
			}
			var tmpkl [][]byte
			for _, kl := range (&dt).Keylist {
				if !notCompareBlockNum {
					var ss string
					rlp.DecodeBytes(kl, &ss)
					sss := strings.Split(ss, ",")
					num, _ := strconv.Atoi(sss[0])
					//fmt.Println("num:", num)
					if Compare(IntToBytes(uint64(num)), bn) > 0 {
						continue
					}
				}
				tmpkl = append(tmpkl, kl)
			}
			if tmpkl != nil {
				r := SearchValue{(&dt).Value, tmpkl}
				result = append(result, r)
			}
		case *data:
			if Compare(dt.Value, min) < 0 {
				return su, result, err
			}
			var tmpkl [][]byte
			for _, kl := range dt.Keylist {
				if !notCompareBlockNum {
					var ss string
					rlp.DecodeBytes(kl, &ss)
					sss := strings.Split(ss, ",")
					num, _ := strconv.Atoi(sss[0])
					//fmt.Println("num:", num)
					if Compare(IntToBytes(uint64(num)), bn) > 0 {
						continue
					}
				}
				tmpkl = append(tmpkl, kl)
			}
			if tmpkl != nil {
				r := SearchValue{dt.Value, tmpkl}
				result = append(result, r)
			}
		default:
			err := errors.New("wrong data type")
			return false, nil, err
		}
	}
	switch nnt := (n.Next).(type) {
	case *leafNode:
		n = nnt
	case *ByteNode:

		nntid, _ := nnt.cache()
		if nntid == nil {
			return true, result, nil
		}

		decoden, err := t.resolveLeaf(nntid)
		if err != nil {
			return false, nil, err
		}
		n.Next = &decoden
		n = &decoden
	case nil:
		return true, result, err
	default:
		err := errors.New("wrong type")
		return false, nil, err
	}

	for {
		if n == nil {
			break
		}
		//if Compare(IntToBytes(uint64(len(result))), bn) >= 0 {
		//	break
		//}
		//flag := false
		for i := 0; i < len(n.Data); i++ {
			switch dt := (n.Data[i]).(type) {
			case dataEncode:
				fmt.Println("data is encoded")
				err := errors.New("data is encoded in RangeValueSearch")
				return false, nil, err
			case data:
				if Compare((&dt).Value, min) < 0 {
					return su, result, err
				}
				var tmpkl [][]byte
				for _, kl := range (&dt).Keylist {
					if !notCompareBlockNum {
						var ss string
						rlp.DecodeBytes(kl, &ss)
						sss := strings.Split(ss, ",")
						num, _ := strconv.Atoi(sss[0])
						//fmt.Println("num:", num)
						if Compare(IntToBytes(uint64(num)), bn) > 0 {
							continue
						}
					}
					tmpkl = append(tmpkl, kl)
				}
				if tmpkl != nil {
					r := SearchValue{(&dt).Value, tmpkl}
					result = append(result, r)
				}
			case *data:
				if Compare(dt.Value, min) < 0 {
					return su, result, err
				}
				var tmpkl [][]byte
				for _, kl := range dt.Keylist {
					if !notCompareBlockNum {
						var ss string
						rlp.DecodeBytes(kl, &ss)
						sss := strings.Split(ss, ",")
						num, _ := strconv.Atoi(sss[0])
						//fmt.Println("num:", num)
						if Compare(IntToBytes(uint64(num)), bn) > 0 {
							continue
						}
					}
					tmpkl = append(tmpkl, kl)
				}
				if tmpkl != nil {
					r := SearchValue{dt.Value, tmpkl}
					result = append(result, r)
				}
			default:
				err := errors.New("wrong data type")
				return false, result, err
			}
		}

		//if !flag {
		if n.Next == nil {
			return true, result, nil
		}
		switch nnt := (n.Next).(type) {
		case *ByteNode:
			nntid, _ := nnt.cache()
			if nntid == nil {
				return true, result, nil
			}
			decoden, err := t.resolveLeaf(nntid)
			if err != nil {
				return false, nil, err
			}
			n.Next = &decoden
			n = &decoden

		case *leafNode:
			n = nnt

		default:
			err := errors.New("wrong next node type")
			return false, nil, err
		}
		//} else {
		//	break
		//}
	}
	//if Compare(IntToBytes(uint64(len(result))), bn) > 0 {
	//	wrapError(err, "top-bn value search wrong:get too much data")
	//	return false, result, err
	//} else {
	return true, result, nil
	//}

}

//range data search
func (t *EBTree) RangeDataSearch(k []byte, min []byte, max []byte) (bool, []searchData, error) {
	var result []searchData

	su, n, pos, err := t.findNode(t.Root, max)
	if !su {
		wrapError(err, "range search value wrong:find node wrong")
		return false, nil, err
	}
	for i := int(pos); i < len(n.Data); i++ {
		switch dt := (n.Data[i]).(type) {
		case dataEncode:
			//TODO:
			return false, nil, nil
		case data:
			for _, kl := range dt.Keylist {
				if Compare(IntToBytes(uint64(len(result))), k) < 0 && Compare(dt.Value, max) <= 0 {
					r := searchData{dt.Value, kl}
					result = append(result, r)
				} else {
					return true, result, nil
				}
			}
		default:
			return false, nil, nil
		}
	}
	switch nnt := (n.Next).(type) {
	case *leafNode:
		n = nnt
	case *ByteNode:
		err := errors.New("the node is encoded")
		//todo: load from cache or database
		return false, nil, err
	default:
		err := errors.New("wrong type")
		return false, nil, err
	}
	for {
		if n == nil {
			break
		}
		if Compare(IntToBytes(uint64(len(result))), k) >= 0 {
			break
		}
		flag := false
		for i := 0; i < len(n.Data); i++ {
			switch dt := (n.Data[i]).(type) {
			case dataEncode:
				//TODO:
				return false, nil, nil
			case data:
				for _, kl := range dt.Keylist {
					if Compare(IntToBytes(uint64(len(result))), k) < 0 && Compare(dt.Value, max) <= 0 {
						r := searchData{dt.Value, kl}
						result = append(result, r)
					} else {
						flag = true
						break
					}
				}
			default:
				return false, nil, nil
			}

		}

		if !flag {
			switch nnt := (n.Next).(type) {
			case *ByteNode:
				//todo: resovle this field
				return false, nil, nil
			case *leafNode:
				n = nnt
			}
		} else {
			break
		}
	}
	if Compare(IntToBytes(uint64(len(result))), k) > 0 {
		wrapError(err, "top-k search data wrong:get too much data")
		return false, result, err
	} else {
		return true, result, nil
	}

	return false, nil, nil
}

//该判断方法适用于范围查询和top-k查询
//搜索时需要将搜索的结果和special值进行对比，保证输出正确结果
//输入为对tree进行搜索得到的结果
//返回值为最终的结果和错误
func (t *EBTree) CompareSpeacial(min []byte, max []byte) (bool, uint64, uint64, error) {
	flag := false
	var pos uint64
	var count uint64
	pos = 0
	count = 0
	//判断special值是否应该包含在结果集中
	for i := uint64(0); i < uint64(len(t.special)); i++ {
		if Compare(t.special[i].value, max) <= 0 && Compare(t.special[i].value, min) >= 0 {
			count++
			if !false {
				pos = i
				flag = true
			}
		}
	}
	return flag, pos, count, nil
	//如果结果集中应包含special值，special值应该放在哪里
}

func (tree *EBTree) CombineAndPrintSearchData(result []searchData, pos []byte, k []byte, top bool) error {
	log.Info("into comine and print searchdata")
	if pos == nil {
		pos = IntToBytes(uint64(0))
	}
	/*_, result, err := tree.CombineSearchDataResult(result, pos, k, top)
	if err != nil {
		fmt.Printf("something wrong in combine search data result\n")
		return err
	}*/
	for i, r := range result {
		fmt.Printf("the %dth value is %d,the data is:\n", i, r.value)
		fmt.Println(r.data)
	}
	return nil
}

func (tree *EBTree) CombineAndPrintSearchValue(result []SearchValue, pos []byte, k []byte, top bool) error {
	log.Info("into comine and print SearchValue")
	if pos == nil {
		pos = IntToBytes(uint64(0))
	}
	for i, r := range result {
		fmt.Printf("the %dth value is %d,the data is:\n", i, r.Value)
		fmt.Println(r.Data)
	}
	return nil
}

func (tree *EBTree) CombineSearchDataResult(result []searchData, min []byte, k []byte, top bool) (bool, []searchData, error) {
	log.Info("into CombineSearchDataResult in EBtree")
	var finalR []searchData
	var su bool
	var pos uint64
	var number uint64
	var err error
	if top {
		su, pos, number, err = tree.CompareSpeacial(IntToBytes(0), result[len(result)-1].value)
	} else {
		su, pos, number, err = tree.CompareSpeacial(min, result[len(result)-1].value)
	}

	if err != nil {
		fmt.Printf("something wrong in compare special\n")
	}
	if number == 0 {
		return true, result, nil
	}
	max := pos + number
	if su {
		i := uint64(0)
		for {
			if len(finalR) >= len(result) {
				break
			}
			sd := searchData{}
			if pos > max {
				sd = result[i]
				finalR = append(finalR, sd)
				i++
				continue
			}
			if Compare(result[i].value, tree.special[pos].value) > 0 {
				for i := uint64(0); i <= pos; i++ {
					if len(finalR) >= len(result) {
						break
					}
					for i := 0; i < len(tree.special[pos].data); i++ {
						sd.data = tree.special[pos].data[i]
						sd.value = tree.special[pos].value
						finalR = append(finalR, sd)
						if len(finalR) >= len(result) {
							break
						}
					}

				}

				pos = pos + 1
				if pos >= max {
					continue
				}
			} else {
				sd = result[i]
				finalR = append(finalR, sd)
				i++
			}
			if i >= uint64(len(result)) || Compare(IntToBytes(uint64(i)), k) >= 0 {
				break
			}

		}
		if Compare(IntToBytes(uint64(len(finalR))), k) < 0 {
			err := errors.New("there is less data than k")
			return false, finalR, err
		}
	}
	return true, finalR, nil
}

func (tree *EBTree) CombineSearchValueResult(result []SearchValue, min []byte, k []byte, top bool) (bool, []SearchValue, error) {

	var finalR []SearchValue
	var su bool
	var pos uint64
	var number uint64
	var err error
	if top {
		if result == nil {
			err := errors.New("result is nil, no data in range")
			return false, nil, err
		} else {
			su, pos, number, err = tree.CompareSpeacial(IntToBytes(0), result[len(result)-1].Value)
		}

	} else {
		if result == nil {
			err := errors.New("no data in such range")
			return false, nil, err
		} else {
			if result == nil {
				err := errors.New("result is nil, no data in range")
				return false, nil, err
			} else {
				su, pos, number, err = tree.CompareSpeacial(min, result[len(result)-1].Value)
			}
		}

	}

	if err != nil {
		fmt.Printf("something wrong in compare special\n")
	}
	if number == 0 {
		return true, result, nil
	}
	max := pos + number
	if su {
		i := uint64(0)
		for {
			if len(finalR) >= len(result) {
				break
			}
			sd := SearchValue{}
			if pos >= max {
				sd = result[i]
				finalR = append(finalR, sd)
				i++
				continue
			}
			if Compare(result[i].Value, tree.special[pos].value) > 0 {
				sd.Value = tree.special[pos].value
				sd.Data = tree.special[pos].data
				finalR = append(finalR, sd)
				pos = pos + 1
				if pos >= max {
					continue
				}
			} else {
				sd = result[i]
				finalR = append(finalR, sd)
				i++
			}
			if i >= uint64(len(result)) || Compare(IntToBytes(i), k) >= 0 {
				break
			}

		}
		if Compare(IntToBytes(uint64(len(finalR))), k) < 0 && top {
			fmt.Println("there is less data than k")

		}
	}
	return true, finalR, nil
}

func mustDecodeNode(id, buf []byte) EBTreen {
	n, err := decodeNode(id, buf)
	if err != nil {
		panic(fmt.Sprintf("node %x: %v", id, err))
	}
	return n
}

// decodeNode parses the RLP encoding of a tree node.
func decodeNode(id, buf []byte) (EBTreen, error) {
	//if BytesToInt(id) == uint64(33) {
	//	elems, _, err := rlp.SplitList(buf)
	//	n, err := decodeInternal(id, elems)
	//	return n, wrapError(err, "full")
	//}
	if len(buf) == 0 {
		return nil, io.ErrUnexpectedEOF
	}
	elems, _, err := rlp.SplitList(buf)
	if err != nil {
		return nil, fmt.Errorf("decode error: %v", err)
	}
	switch c, _ := rlp.CountValues(elems); c {

	case 4:
		n, err := decodeLeaf(id, elems)
		return n, wrapError(err, "short")
	case 3:
		n, err := decodeInternal(id, elems)
		return n, wrapError(err, "full")
	default:
		return nil, fmt.Errorf("invalid number of list elements: %v", c)
	}

}

func decodeData(buf []byte) (data, error) {
	//log.Info("into decoded data")
	d := data{}
	elems, _, _ := rlp.SplitList(buf)
	rlp.CountValues(elems)
	//fmt.Println(elems)
	//fmt.Println(c)

	elems, rest, _ := rlp.SplitList(elems)
	cc, _ := rlp.CountValues(elems)
	//fmt.Println(elems)
	//fmt.Println(cc)

	for i := 0; i < cc; i++ {
		kbuf, rest1, _ := rlp.SplitString(elems)
		//fmt.Print(i)
		//fmt.Println(kbuf)
		d.Keylist = append(d.Keylist, kbuf)
		//fmt.Println(rest1)
		elems = rest1
	}
	elems = rest
	value, _, _ := rlp.SplitString(elems)
	d.Value = value
	//fmt.Print("value:")
	//fmt.Println(value)
	//fmt.Println(rest3)

	return d, nil

}

func decodeChild(buf []byte) (child, []byte, error) {
	elems, _, _ := rlp.SplitList(buf)
	//c, _ := rlp.CountValues(elems)
	//fmt.Println("into decodeChild：")
	//fmt.Println(elems)
	//fmt.Println(c)

	evalue, rest, err := rlp.SplitString(elems)
	if err != nil {
		fmt.Println("into decodeChild：rlp.splitString evalue wrong")
		fmt.Println(err)
		wrapError(err, "into decodeChild：rlp.splitString evalue wrong")
		return child{}, nil, err
	}

	epointer, _, err := rlp.SplitString(rest)
	if err != nil {

		fmt.Println("into decodeChild：rlp.splitString epointer wrong")
		fmt.Println(err)
		wrapError(err, "into decodeChild：rlp.splitString epointer wrong")
		return child{}, nil, err
	}
	//fmt.Println("into decodeChild：rlp decode success！")
	value := evalue
	//fmt.Println(value)
	var ebpointer ByteNode
	ebpointer = epointer
	//fmt.Println(epointer)
	cd := constructChild(value, &ebpointer)
	return cd, buf, err

}

func decodeLeaf(id, buf []byte) (EBTreen, error) {
	//log.Info("to decode leafnode:")
	//fmt.Println(buf)
	le := leafNode{}
	le.Id = id
	elems, rest, _ := rlp.SplitList(buf)
	c, _ := rlp.CountValues(elems)

	//fmt.Printf("%d values in elems\n", c)

	for i := 0; i < c; i++ {
		kbuf, rest1, _ := rlp.SplitString(elems)
		d, _ := decodeData(kbuf)
		le.Data = append(le.Data, d)
		elems = rest1
	}

	elems = rest
	nextid, _, _ := rlp.SplitString(elems)
	if len(nextid) > 0 {
		var nextByteNode ByteNode
		nextByteNode = nextid
		le.Next = &nextByteNode
	}

	//fmt.Println(rest5)
	return &le, nil

}

func decodeInternal(id, buf []byte) (EBTreen, error) {
	in := internalNode{}
	in.Id = id
	elems, _, _ := rlp.SplitList(buf)
	c, _ := rlp.CountValues(elems)
	//fmt.Println("decodeInternal")
	//fmt.Println(id)
	for i := 0; i < c; i++ {
		kbuf, rest1, err := rlp.SplitString(elems)
		if err != nil {
			err = wrapError(err, "error in split string when decode children")
			return nil, err
		}
		cd, _, err := decodeChild(kbuf)
		if err != nil {
			wrapError(err, "decode internal node error when decode child")
			return nil, nil
		}
		in.Children = append(in.Children, cd)
		elems = rest1
	}

	return &in, nil
}

func constructData(value []byte, keylist [][]byte) data {
	d := data{Value: value, Keylist: keylist}
	return d
}

func constructChild(value []byte, pointer EBTreen) child {
	c := child{Value: value, Pointer: pointer}
	return c
}

func encodeInternal(bb *[]byte, in *internalNode) error {
	for i := 0; i < len(in.Children); i++ {
		switch ct := (in.Children[i]).(type) {
		case childEncode:
			err := errors.New("encode leaf error, it is encoded alredy")
			return err
		case child:
			bb, _ := rlp.EncodeToBytes(ct)
			var childE childEncode
			childE = bb
			in.Children[i] = childE
		default:
			err := errors.New("encode leaf error, wrong type")
			return err
		}
	}
	//in.Id = nil
	buff3 := bytes.Buffer{}
	rlp.Encode(&buff3, &in)
	b1 := buff3.Bytes()
	for _, i := range b1 {
		*bb = append(*bb, i)
	}
	return nil
}

func encodeLeaf(result *[]byte, le *leafNode) error {
	//fmt.Printf("into encode leaf,len of data is %d\n", len(le.Data))
	for i := 0; i < len(le.Data); i++ {
		switch dt := (le.Data[i]).(type) {
		case dataEncode:
		case data:
			bb, _ := rlp.EncodeToBytes(dt)
			var dataE dataEncode
			dataE = bb
			le.Data[i] = dataE
		case *data:
			bb, _ := rlp.EncodeToBytes(dt)
			var dataE dataEncode
			dataE = bb
			le.Data[i] = dataE
		default:
			err := errors.New("wrong type")
			return err
		}

	}
	if le.Next != nil {
		switch let := (le.Next).(type) {
		case *leafNode:
			var nt ByteNode
			nt = let.Id
			le.Next = &nt
		case *internalNode:
			var nt ByteNode
			nt = let.Id
			le.Next = &nt
		case *ByteNode:
			le.Next = let
		default:
			err := errors.New("wrong type")
			return err
		}
	}

	//le.Id = nil
	r1, err := rlp.EncodeToBytes(le)
	if err != nil {
		err := wrapError(err, "encode leaf wrong")
		return err
	}
	for _, i := range r1 {
		*result = append(*result, i)
	}

	return nil
}
