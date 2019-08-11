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
	"strings"

	"github.com/ethereum/go-ethereum/rlp"
)

var indices = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f", "[17]"}

type EBTreen interface {
	fstring(string) string
	cache() ([]byte, bool)
}

type (
	internalNode struct {
		Children []ChildInterface // Actual ebtree internal node data to encode/decode (needs custom encoder)
		Id       []byte
	}
	leafNode struct {
		Data []DataInterface
		Next EBTreen
		Id   []byte
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

func (n *internalNode) cache() ([]byte, bool) { return n.Id, true }
func (n *leafNode) cache() ([]byte, bool)     { return n.Id, true }
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

type searchValue struct {
	value []byte
	data  [][]byte
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
	log.Info("into create leaf node")
	var empty []byte
	se, err := tree.newSequence()
	log.Info(string(se))
	if err != nil {
		err = wrapError(err, "CreateLeafNode")
		return leafNode{}, err
	}
	newn, err := constructLeafNode(se, uint8(len(datalist)), datalist, false, true, nil, empty, 0)
	return newn, err
}

func constructLeafNode(id []byte, count uint8, datalist []data, special bool, dirty bool, next EBTreen, nexid []byte, gen uint16) (leafNode, error) {
	log.Info("into construct leaf node")
	//newn := &leafNode{}
	var dataInter []DataInterface
	for i := uint8(0); i < count; i++ {
		dataInter = append(dataInter, &datalist[i])
	}
	newn := leafNode{Data: dataInter, Id: id, Next: next}
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
	newn := internalNode{Id: id, Children: childlist}
	return newn, nil
}

func addData(leaf leafNode, da data, position int) {
	var i int
	for i = len(leaf.Data); i > position; i-- {
		leaf.Data[i] = leaf.Data[i-1]
	}
	leaf.Data[i] = da
}
func add(b []byte, i uint64) []byte {
	f := BytesToInt(b)
	return IntToBytes(f + i)
}

func minus(b []byte, i uint64) []byte {
	f := BytesToInt(b)
	return IntToBytes(f - i)
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
	log.Info("into moveData")
	n.Data = append(n.Data, data{})

	for i := len(n.Data) - 1; i > pos; i-- {
		n.Data[i] = n.Data[i-1]
	}
	return true, n, nil
}

func moveChildren(n *internalNode, pos int) (bool, *internalNode, error) {

	n.Children = append(n.Children, child{})

	if pos > len(n.Children) {
		err := errors.New("the length of n.children is smaller than count,something wrong")
		return false, n, err
	}

	//如果需要在最后一个节点插入，那么，不需要移动元素
	if pos == len(n.Children) {
		return true, n, nil
	}
	for i := len(n.Children) - 1; i > pos; i-- {
		n.Children[i] = n.Children[i-1]
	}
	return true, n, nil
}

func createChild(val []byte, po EBTreen) (child, error) {
	ch := &child{Value: val, Pointer: po}
	if ch == nil {
		err := errors.New("create child failed")
		return *ch, err
	}
	return *ch, nil
}

//TODO:将定位当前叶子节点在parent节点中的位置，便于后期的插入、查找
func getLeafNodePosition(n *leafNode, parent *internalNode, t *EBTree) (bool, *internalNode, uint8, error) {
	if parent == nil {
		datainter := n.Data[len(n.Data)-1]
		switch dt := (datainter).(type) {
		case dataEncode:
			//TODO:
			_ = dt
			return false, nil, 0, nil
		case data:
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
			return false, nil, 0, nil
		}
	}
	var re uint8
	flag := false
	if parent == nil {
		return false, nil, 0, nil
	}
	for i, pc := range parent.Children {
		switch ct := (pc).(type) {
		case childEncode:
			return false, nil, 0, nil
		case child:
			if ct.Pointer == n {
				switch dt := (n.Data[len(n.Data)-1]).(type) {
				case dataEncode:
					//TODO:
					_ = dt
					return false, nil, 0, nil
				case data:
					re = uint8(i)
					ct.Value = dt.Value
					parent.Children[i] = pc
					flag = true
					return true, parent, re, nil
				default:
					return false, nil, 0, nil
				}
				//fmt.Println(pc.value)
			}
		default:
			return false, nil, 0, nil
		}

	}
	if !flag {
		err := errors.New("there is not such leaf node in this parent node")
		return false, parent, re, err
	}
	return true, parent, re, nil
}

//TODO:将定位当前节点在parent节点中的位置，便于后期的插入、查找
func getInternalNodePosition(n *internalNode, parent *internalNode, t *EBTree) (bool, *internalNode, uint8, error) {
	switch ct := (n.Children[len(n.Children)-1]).(type) {
	case childEncode:
		return false, nil, 0, nil
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
				return false, nil, 0, nil
			case child:
				if cpt.Pointer == n {
					re = uint8(i)
					cpt.Value = ct.Value
					parent.Children[i] = pc
					flag = true
					//fmt.Println(pc.value)
				}
			}

		}
		if !flag {
			err := errors.New("there is not such leaf node in this parent node")
			return false, parent, re, err
		}
		return true, parent, re, nil
	}
	return false, nil, 0, nil
}

func createData(value []byte, da []byte) (data, error) {
	log.Info("into createData")
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
			if bytes.Compare(ct.Value, value) >= 0 {
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

func findFirstNode(n EBTreen) (*leafNode, error) {
	if n == nil {
		err := errors.New("find first node wrong: t.root is nil")
		return nil, err
	}
	switch nt := (n).(type) {
	case *leafNode:
		return nt, nil
	case *internalNode:
		if len(nt.Children) <= 0 {
			err := errors.New("find first node wrong: when node is internal node, nt.count <=0")
			return nil, err
		}
		switch ct := (nt.Children[0]).(type) {
		case childEncode:
			return nil, nil
		case child:
			re, err := findFirstNode(ct.Pointer)
			if err != nil {
				wrapError(err, "find first node wrong:when node is internal node")
				return nil, err
			}
			return re, nil
		}

	}
	err := errors.New("find first node wrong: the type of node is wrong")
	return nil, err
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

		//遍历当前节点的所有data
		//value一定小于或等于当前节点到最大值
		for i := 0; i < len(nt.Data); i++ {
			switch dt := (nt.Data[i]).(type) {
			case dataEncode:
				//TODO:
				return false, nt, 0, nil
			case data:
				if bytes.Compare(dt.Value, value) < 0 {
					//EBTree叶子节点按升序排列，继续向后查找
					continue
				} else {
					//找到节点

					return true, nt, uint8(i), nil
				}
			default:
				return false, nt, 0, nil
			}

		}
	case *internalNode:
		var i int
		for i = 0; i < len(nt.Children); i++ {
			switch ct := (nt.Children[i]).(type) {
			case childEncode:
				return false, nil, 0, nil
			case child:
				if bytes.Compare(ct.Value, value) < 0 {
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
	default:
		err := errors.New("the node is not leaf or internal, something wrong")
		return false, nil, 0, err
	}
	err := errors.New("find node wrong, something wrong")
	return false, nil, 0, err
}

//top-k value search
func (t *EBTree) TopkValueSearch(k []byte, max bool) (bool, []searchValue, error) {
	var result []searchValue
	if max {
		n, err := findFirstNode(t.Root)
		if err != nil {
			wrapError(err, "top-k search value wrong:find first node wrong")
			return false, nil, err
		}
		for {
			if n == nil {
				break
			}
			if bytes.Compare(IntToBytes(uint64(len(result))), k) >= 0 {
				break
			}
			flag := false
			for i := 0; i < len(n.Data); i++ {
				if bytes.Compare(IntToBytes(uint64(len(result))), k) < 0 {
					switch dt := (n.Data[i]).(type) {
					case dataEncode:
						//TODO:
						return false, nil, nil
					case data:
						r := searchValue{dt.Value, dt.Keylist}
						result = append(result, r)
					default:
						return false, nil, nil
					}

				} else {
					flag = true
					break
				}
			}

			if !flag {
				switch nnt := (n.Next).(type) {
				case *leafNode:
					n = nnt
				case *ByteNode:
					//todo: load from cache or database
					return false, nil, nil
				default:
					err := errors.New("wrong type")
					return false, nil, err
				}

			} else {
				break
			}
		}
		if bytes.Compare(IntToBytes(uint64(len(result))), k) < 0 {
			wrapError(err, "top-k value search wrong:not enough data")
			return false, result, err
		} else if bytes.Compare(IntToBytes(uint64(len(result))), k) > 0 {
			wrapError(err, "top-k value search wrong:get too much data")
			return false, result, err
		} else {
			return true, result, nil
		}
	}
	return false, nil, nil

}

//top-k data search
func (t *EBTree) TopkDataSearch(k []byte, max bool) (bool, []searchData, error) {
	var result []searchData
	if max {
		n, err := findFirstNode(t.Root)
		if err != nil {
			wrapError(err, "top-k search data wrong:find first node wrong")
			return false, nil, err
		}
		for {
			if n == nil {
				break
			}
			if bytes.Compare(IntToBytes(uint64(len(result))), k) >= 0 {
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
						if bytes.Compare(IntToBytes(uint64(len(result))), k) < 0 {
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
				case *leafNode:
					n = nnt
				case *ByteNode:
					//todo: load from cache or database
					return false, nil, nil
				default:
					err := errors.New("wrong type")
					return false, nil, err
				}
			} else {
				break
			}
		}
		if bytes.Compare(IntToBytes(uint64(len(result))), k) < 0 {
			wrapError(err, "top-k search data wrong:not enough data")
			return false, result, err
		} else if bytes.Compare(IntToBytes(uint64(len(result))), k) > 0 {
			wrapError(err, "top-k search data wrong:get too much data")
			return false, result, err
		} else {
			return true, result, nil
		}

	} else {

	}
	return false, nil, nil
}

//range value search
func (t *EBTree) RangeValueSearch(min []byte, max []byte, k []byte) (bool, []searchValue, error) {
	var result []searchValue

	su, n, pos, err := t.findNode(t.Root, min)
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
			if bytes.Compare(IntToBytes(uint64(len(result))), k) < 0 && bytes.Compare(dt.Value, max) <= 0 {
				r := searchValue{dt.Value, dt.Keylist}
				result = append(result, r)
			} else {
				return true, result, nil
			}
		default:
			return false, nil, nil
		}
	}
	switch nnt := (n.Next).(type) {
	case *leafNode:
		n = nnt
	case *ByteNode:
		//todo: load from cache or database
		return false, nil, nil
	default:
		err := errors.New("wrong type")
		return false, nil, err
	}
	for {
		if n == nil {
			break
		}
		if bytes.Compare(IntToBytes(uint64(len(result))), k) >= 0 {
			break
		}
		flag := false
		for i := 0; i < len(n.Data); i++ {
			switch dt := (n.Data[i]).(type) {
			case dataEncode:
				//TODO:
				return false, nil, nil
			case data:
				if bytes.Compare(IntToBytes(uint64(len(result))), k) < 0 && bytes.Compare(dt.Value, max) <= 0 {
					r := searchValue{dt.Value, dt.Keylist}
					result = append(result, r)
				} else {
					flag = true
					break
				}
			}
		}

		if !flag {
			switch nnt := (n.Next).(type) {
			case *ByteNode:
				return false, nil, nil
			case *leafNode:
				n = nnt
			default:
				err := errors.New("wrong type")
				return false, nil, err
			}
		} else {
			break
		}
	}
	if bytes.Compare(IntToBytes(uint64(len(result))), k) > 0 {
		wrapError(err, "top-k value search wrong:get too much data")
		return false, result, err
	} else {
		return true, result, nil
	}

	return false, nil, nil

}

//top-k data search
func (t *EBTree) RangeDataSearch(k []byte, min []byte, max []byte) (bool, []searchData, error) {
	var result []searchData

	su, n, pos, err := t.findNode(t.Root, min)
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
				if bytes.Compare(IntToBytes(uint64(len(result))), k) < 0 && bytes.Compare(dt.Value, max) <= 0 {
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
		if bytes.Compare(IntToBytes(uint64(len(result))), k) >= 0 {
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
					if bytes.Compare(IntToBytes(uint64(len(result))), k) < 0 && bytes.Compare(dt.Value, max) <= 0 {
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
	if bytes.Compare(IntToBytes(uint64(len(result))), k) > 0 {
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
		if bytes.Compare(t.special[i].value, max) <= 0 && bytes.Compare(t.special[i].value, min) >= 0 {
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
func (tree *EBTree) CombineAndPrintSearchData(result []searchData, pos []byte, k []byte, top bool) {
	if pos == nil {
		pos = IntToBytes(uint64(0))
	}
	_, result, err := tree.CombineSearchDataResult(result, pos, k, top)
	if err != nil {
		fmt.Printf("something wrong in combine search data result\n")
		return
	}
	for i, r := range result {
		fmt.Printf("the %dth value is %d,the data is:\n", i, r.value)
		fmt.Printf(string(r.data))
		fmt.Println()
	}
}
func (tree *EBTree) CombineSearchDataResult(result []searchData, min []byte, k []byte, top bool) (bool, []searchData, error) {

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
			if bytes.Compare(result[i].value, tree.special[pos].value) > 0 {
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
			if i >= uint64(len(result)) || bytes.Compare(IntToBytes(uint64(i)), k) >= 0 {
				break
			}

		}
		if bytes.Compare(IntToBytes(uint64(len(finalR))), k) < 0 {
			err := errors.New("there is less data than k")
			return false, finalR, err
		}
	}
	return true, finalR, nil
}

func (tree *EBTree) CombineSearchValueResult(result []searchValue, min []byte, k []byte, top bool) (bool, []searchValue, error) {

	var finalR []searchValue
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
			sd := searchValue{}
			if pos >= max {
				sd = result[i]
				finalR = append(finalR, sd)
				i++
				continue
			}
			if bytes.Compare(result[i].value, tree.special[pos].value) > 0 {
				sd.value = tree.special[pos].value
				sd.data = tree.special[pos].data
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
			if i >= uint64(len(result)) || bytes.Compare(IntToBytes(i), k) >= 0 {
				break
			}

		}
		if bytes.Compare(IntToBytes(uint64(len(finalR))), k) < 0 && top {
			err := errors.New("there is less data than k")
			return false, finalR, err
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
	if len(buf) == 0 {
		return nil, io.ErrUnexpectedEOF
	}
	elems, _, err := rlp.SplitList(buf)
	if err != nil {
		return nil, fmt.Errorf("decode error: %v", err)
	}
	switch c, _ := rlp.CountValues(elems); c {

	case 3:
		n, err := decodeLeaf(id, elems)
		return n, wrapError(err, "short")
	case 2:
		n, err := decodeInternal(id, elems)
		return n, wrapError(err, "full")
	default:
		return nil, fmt.Errorf("invalid number of list elements: %v", c)
	}

}

func decodeData(buf []byte) (data, error) {
	log.Info("into decoded data")
	d := data{}
	elems, _, _ := rlp.SplitList(buf)
	c, _ := rlp.CountValues(elems)
	fmt.Println(elems)
	fmt.Println(c)

	elems, rest, _ := rlp.SplitList(elems)
	cc, _ := rlp.CountValues(elems)
	fmt.Println(elems)
	fmt.Println(cc)

	for i := 0; i < cc; i++ {
		kbuf, rest1, _ := rlp.SplitString(elems)
		fmt.Print(i)
		fmt.Println(kbuf)
		d.Keylist = append(d.Keylist, kbuf)
		fmt.Println(rest1)
		elems = rest1
	}
	elems = rest
	value, rest3, _ := rlp.SplitString(elems)
	d.Value = value
	fmt.Print("value:")
	fmt.Println(value)
	fmt.Println(rest3)

	return d, nil

}

func decodeChild(buf []byte) (child, []byte, error) {
	elems, _, _ := rlp.SplitList(buf)
	c, _ := rlp.CountValues(elems)
	fmt.Println(elems)
	fmt.Println(c)

	evalue, rest, err := rlp.SplitString(elems)
	if err != nil {
		return child{}, nil, err
	}

	epointer, _, err := rlp.SplitString(rest)
	if err != nil {
		return child{}, nil, err
	}

	value := evalue
	var ebpointer ByteNode
	ebpointer = epointer
	cd := constructChild(value, &ebpointer)
	return cd, buf, err

}

func decodeLeaf(id, buf []byte) (EBTreen, error) {
	log.Info("decode leafnode:")
	fmt.Println(id)
	le := leafNode{}
	le.Id = id
	elems, rest, _ := rlp.SplitList(buf)
	c, _ := rlp.CountValues(elems)
	fmt.Println(elems)
	fmt.Println(c)

	for i := 0; i < c; i++ {
		kbuf, rest1, _ := rlp.SplitString(elems)
		d, _ := decodeData(kbuf)
		fmt.Println(d.Value)
		fmt.Print(i)
		fmt.Println(kbuf)
		le.Data = append(le.Data, d)
		fmt.Println(rest1)
		elems = rest1
	}

	elems = rest
	nextid, rest5, _ := rlp.SplitString(elems)
	var nextByteNode ByteNode
	nextByteNode = nextid
	le.Next = &nextByteNode
	fmt.Print("nextid:")
	fmt.Println(nextid)
	fmt.Println(rest5)
	return &le, nil

	return nil, nil
}

func decodeInternal(id, buf []byte) (EBTreen, error) {
	in := internalNode{}
	in.Id = id
	elems, _, _ := rlp.SplitList(buf)
	c, _ := rlp.CountValues(elems)
	fmt.Println(elems)
	fmt.Println(c)

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
		fmt.Print(i)
		fmt.Println(kbuf)
		in.Children = append(in.Children, cd)
		fmt.Println(rest1)
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

//TODO:6/19
//将节点转化为[]数组，将节点元素转化为byte,byte[],或者有byte[]组成的列表
func encodeNode(oriNode EBTreen) ([]byte, error) {
	switch n := (oriNode).(type) {
	case *leafNode:
		//TODO:encode leafNode
		le, _ := rlp.EncodeToBytes(n)
		fmt.Printf("bytes is %v", le)
		return nil, nil
	case *internalNode:
		//TODO:encode internalNode
		return nil, nil
	}
	return nil, nil
}

func (*leafNode) encodeData(d data) ([]byte, error) {
	//TODO:encode the data for leaf node
	return nil, nil
}

func (*leafNode) encodeKeList(k [][]byte) ([]byte, error) {
	//TODO:encode the keylist for data
	return nil, nil
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
	in.Id = nil
	buff3 := bytes.Buffer{}
	rlp.Encode(&buff3, &in)
	b1 := buff3.Bytes()
	for _, i := range b1 {
		*bb = append(*bb, i)
	}
	return nil
}

func encodeLeaf(result *[]byte, le *leafNode) error {
	for i := 0; i < len(le.Data); i++ {
		switch dt := (le.Data[i]).(type) {
		case dataEncode:
			continue
		case data:
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

	le.Id = nil
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
