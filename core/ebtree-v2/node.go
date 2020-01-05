package ebtree_v2

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
)

type EBTreen interface {
	fstring() []byte
}
type (
	LeafNode struct {
		Id        []byte
		NextPtr   EBTreen
		LeafDatas []ResultD
	}
	InternalNode struct {
		Id       []byte
		Children []ChildData
	}
	IdNode []byte
)
type ChildData struct {
	Value   []byte
	NodePtr EBTreen
}

//Start*****************************
// Extend the functions in EBTreen
func (n *InternalNode) fstring() []byte {
	return n.Id
}
func (n *LeafNode) fstring() []byte {
	return n.Id
}
func (n IdNode) fstring() []byte {
	return n
}

// Extend functions in EBTreen
//End*****************************

//Start*****************************
// Initial functions in EBTreen

func (ebt *EBTree) NewLeafNode() LeafNode {
	lid := ebt.NewSequence()
	le := LeafNode{
		Id:        lid,
		LeafDatas: nil,
		NextPtr:   nil,
	}
	return le
}

func (ebt *EBTree) NewInternalNode() InternalNode {
	iid := ebt.NewSequence()
	in := InternalNode{
		Id:       iid,
		Children: nil,
	}
	return in
}

func (ebt *EBTree) NewChildData(value []byte, node EBTreen) ChildData {
	chd := ChildData{
		Value:   value,
		NodePtr: node,
	}
	return chd
}

// Initial functions in EBTreen
//End*****************************

//Start*****************************
// insert functions in Nodes after tree inintialed

func (ebt *EBTree) InsertToFirstLeaf(d ResultD, n *LeafNode, i int, flag bool) error {
	var err error
	sle, err := ebt.InsertToLeaf(d, n, i, flag)
	if err != nil {
		return err
	}

	//the leaf node is not split, the result could be return
	if sle == nil {
		return err
	}

	//the leaf node is split, the internal node should be changed
	ile := len(ebt.LastPath.Internals)
	ebt.InsertToInternal(sle, ebt.LastPath.Internals[ile-1], ile-1)
	return nil
}

func (ebt *EBTree) InsertToInternal(n EBTreen, in *InternalNode, i int) (*InternalNode, error) {
	var v []byte
	switch nt := n.(type) {
	case *LeafNode:
		lle := len(nt.LeafDatas)
		v = nt.LeafDatas[lle-1].Value
		fmt.Println(v)
		//todo:now
		//i,err,_:=ebt.SearchInNode(v,in)
	case *InternalNode:
		lle := len(nt.Children)
		v = nt.Children[lle-1].Value
	default:
		err := errors.New("wrong node type in InsertToInternal")
		return nil, err
	}
	var sin *InternalNode
	var err error
	if (len(in.Children)) >= MaxInternalNodeCapability {
		//insert the node, and split it
	} else {
		//insert the node, and return
	}
	return sin, err
}

func (ebt *EBTree) InsertToLeaf(d ResultD, le *LeafNode, i int, flag bool) (*LeafNode, error) {

	var err error

	//if the d.value is already existed
	if flag {
		le.LeafDatas[i].ResultData = append(le.LeafDatas[i].ResultData, d.ResultData...)
		return nil, err
	}

	ile := len(le.LeafDatas)

	if (ile) <= MaxLeafNodeCapability {
		le.LeafDatas = append(le.LeafDatas, d)
		for j := ile - 1; j >= 0; j-- {
			if byteCompare(le.LeafDatas[j].Value, d.Value) < 0 {
				le.LeafDatas[j+1] = le.LeafDatas[j]
			}
		}
		return nil, err
	}

	//concate the new leaf nodes
	nle := ebt.NewLeafNode()
	mid := MaxLeafNodeCapability / 2
	nle.NextPtr = le.NextPtr
	le.NextPtr = &nle

	//allocate the leaf datas between two leafnodes
	if (i + 1) > mid-1 {
		//concate the leafdatas for new leaf node
		for j := mid; j <= (ile - 1); j++ {
			nle.LeafDatas = append(nle.LeafDatas, le.LeafDatas[j])
			if j == (i) {
				nle.LeafDatas = append(nle.LeafDatas, d)
			}
		}
		nle.LeafDatas = append(nle.LeafDatas, d)

		//concate the leafdata for original leaf node
		le.LeafDatas = le.LeafDatas[:mid]
		le.LeafDatas = append(le.LeafDatas, d)
		for j := mid - 1; j > (i); j-- {
			le.LeafDatas[j+1] = le.LeafDatas[j]
		}
	} else {
		//concate the leafdatas for new leaf node
		for j := (0); j <= MaxLeafNodeCapability-mid; j++ {
			nle.LeafDatas = append(nle.LeafDatas, le.LeafDatas[mid-1+j])
		}

		//concate the leafdata for original leaf node
		le.LeafDatas = le.LeafDatas[:mid-1]
		le.LeafDatas = append(le.LeafDatas, d)
		for j := mid - 1; j > (i); j-- {
			le.LeafDatas[j+1] = le.LeafDatas[j]
		}
	}

	return &nle, err
}

// insert functions in Nodes after tree inintialed
//End*****************************

//Start*****************************
// Update functions in InternalNode
func (ebt *EBTree) CreateInternalNode(first EBTreen, second EBTreen) (InternalNode, error) {
	in := ebt.NewInternalNode()
	var chd1, chd2 ChildData
	switch nt := (first).(type) {
	case *LeafNode:
		fvdl := len(nt.LeafDatas)
		fv := nt.LeafDatas[fvdl-1].Value
		chd1 = ebt.NewChildData(fv, nt)
		if second == nil {
			//chd1 = ebt.NewChildData(fv, nt)
			in.Children = append(in.Children, chd1)
			return in, nil
		} else {
			/*first node should be a IdNode
			var ntid IdNode
			ntid=nt.Id
			chd1=ebt.NewChildData(fv,&ntid)*/

			switch snt := (second).(type) {
			case *LeafNode:
				svdl := len(snt.LeafDatas)
				sv := snt.LeafDatas[svdl-1].Value
				chd2 = ebt.NewChildData(sv, snt)
			default:
				err := errors.New("wrong node type when first node is leaf node, while second node wrong")
				return InternalNode{}, err
			}
			in.Children = append(in.Children, chd1)
			in.Children = append(in.Children, chd2)
			return in, nil
		}
	case *InternalNode:
		fvdl := len(nt.Children)
		fv := nt.Children[fvdl-1].Value
		chd1 = ebt.NewChildData(fv, nt)
		if second == nil {
			//chd1 = ebt.NewChildData(fv, nt)
			in.Children = append(in.Children, chd1)
			return in, nil
		} else {

			/*first node should be a IdNode
			var ntid IdNode
			ntid=nt.Id
			chd1=ebt.NewChildData(fv,&ntid)*/

			switch snt := (second).(type) {
			case *InternalNode:
				svdl := len(snt.Children)
				sv := snt.Children[svdl-1].Value
				chd2 = ebt.NewChildData(sv, snt)
			default:
				err := errors.New("wrong node type when first node is leaf node, while second node wrong")
				return InternalNode{}, err
			}
			in.Children = append(in.Children, chd1)
			in.Children = append(in.Children, chd2)
			return in, nil
		}
	default:
		err := errors.New("wrong node type when first node is wrong")
		return InternalNode{}, err
	}

}

func (ebt *EBTree) AdjustNodeInPath(i int64, first EBTreen, second EBTreen) error {

	if int(i) == len(ebt.LastPath.Internals)-1 || (i < 0 && ebt.LastPath.Internals == nil) {
		//we reach to the root node of ebtree
		in, err := ebt.CreateInternalNode(first, second)
		if err != nil {
			return err
		}
		ebt.Root = &in
		ebt.LastPath.Internals = append(ebt.LastPath.Internals, &in)
		return nil
	} else {
		lin := len(ebt.LastPath.Internals[i+1].Children)
		if (lin) >= MaxInternalNodeCapability {
			//a new internal node needed to be created
			nin, err := ebt.CreateInternalNode(second, nil)
			if err != nil {
				return err
			}
			//the second node is put in new internal node
			err2 := ebt.AdjustNodeInPath(i+1, (ebt.LastPath.Internals[i+1]), &nin)
			if err2 != nil {
				return err2
			}

			ebt.LastPath.Internals[i+1] = &nin

			return nil
		} else {
			var v []byte
			switch snt := second.(type) {
			case *LeafNode:
				v = snt.LeafDatas[len(snt.LeafDatas)-1].Value
			case *InternalNode:
				v = snt.Children[len(snt.Children)-1].Value
				ebt.LastPath.Internals[i] = snt
			default:
				err := errors.New("wrong node type in UpdateInternalNodeInPath")
				return err
			}
			chd := ebt.NewChildData(v, second)
			ebt.LastPath.Internals[i+1].Children = append(ebt.LastPath.Internals[i+1].Children, chd)
			return nil
		}

	}
}

// Update functions in InternalNode
//End*****************************

//Start*****************************
// find functions in Node

func (ebt *EBTree) FindInNode(value []byte, n EBTreen, flag bool) (*LeafNode, error) {
	var le *LeafNode
	var err error
	switch nt := n.(type) {
	case *LeafNode:
		if flag {
			ebt.LastPath.Leaf = nt
		}
		return nt, nil
	case *InternalNode:
		if flag {
			ebt.LastPath.Internals = append(ebt.LastPath.Internals, nt)
		}
		i, err, _ := ebt.SearchInNode(value, nt)
		if err != nil {
			return nil, err
		}
		return ebt.FindInNode(value, nt.Children[i].NodePtr, flag)
	case *IdNode:
		nc, err := ebt.LoadNode(nt.fstring())
		if err != nil {
			return le, err
		}
		if nc == nil {
			return nil, nil
		}
		switch nct := nc.(type) {
		case *LeafNode:
			if flag {
				ebt.LastPath.Leaf = nct
			}
			return nct, nil
		case *InternalNode:
			return ebt.FindInNode(value, nct, flag)
		default:
			err = errors.New("wrong node type from leveldb")
		}
	default:
		err := errors.New("wrong node type in FindInNode")
		return nil, err
	}
	return le, err
}

func (ebt *EBTree) SearchInNode(value []byte, n EBTreen) (int64, error, bool) {
	//signal the value is existed in those ebtree
	var flag bool
	flag = false
	switch nt := n.(type) {
	case *LeafNode:
		lo, hi := 0, len(nt.LeafDatas)-1
		for lo <= hi {
			m := (lo + hi) >> 1
			if byteCompare(value, nt.LeafDatas[m].Value) < 0 {
				lo = m + 1
			} else if byteCompare(value, nt.LeafDatas[m].Value) > 0 {
				hi = m - 1
			} else {
				flag = true
				return int64(m), nil, flag
			}
		}
		//not found
		if hi < 0 {
			return int64(hi + 1), nil, flag
		}
		return int64(hi), nil, flag
	case *InternalNode:
		lo, hi := 0, len(nt.Children)-1
		for lo <= hi {
			m := (lo + hi) >> 1
			if byteCompare(value, nt.Children[m].Value) < 0 {
				lo = m + 1
			} else if byteCompare(value, nt.Children[m].Value) > 0 {
				hi = m - 1
			} else {
				flag = true
				return int64(m), nil, flag
			}
		}
		if hi < 0 {
			return int64(hi + 1), nil, flag
		}
		//not found
		return int64(hi), nil, flag
	default:
		err := errors.New("wrong node type in SearchInNode")
		return -1, err, flag
	}

}

// find functions in Node
//End*****************************

//Start*****************************
// commit prepare functions in Node

func (ebt *EBTree) CollapsedUnuseInternal(nt *InternalNode, j int) error {
	var err error
	for i := len(nt.Children) - 2; i >= 0; i-- {
		switch nct := (nt.Children[i].NodePtr).(type) {
		case *LeafNode:
			err = ebt.CollapseLeafNode(nct)

			if err != nil {
				return err
			}
			var idn IdNode
			idn = nct.Id
			nt.Children[i].NodePtr = &idn
		case *InternalNode:
			err = ebt.CollapseInternalNode(nct, false)
			if err != nil {
				return err
			}
			var idn IdNode
			idn = nct.Id
			nt.Children[i].NodePtr = &idn
		case *IdNode:
			break
		default:
			err = errors.New("wrong node type in CollapsedUnuseInternal")
			return err
		}
	}
	if byteCompare(nt.Id, ebt.LastPath.Internals[j].Id) == 0 {
		return err
	} else {
		return ebt.CollapseInternalNode(nt, false)
	}

}

func (ebt *EBTree) CollapseLeafNode(nt *LeafNode) error {
	var ntid []byte
	if nt.NextPtr != nil {
		switch nnt := (nt.NextPtr).(type) {
		case *LeafNode:
			ntid = nnt.Id
		case *IdNode:
			return nil
		default:
			err := errors.New("wrong node type in leaf node.nextptr")
			return err
		}
		var ntptr IdNode
		ntptr = ntid
		nt.NextPtr = &ntptr
	} else {
		var ntptr IdNode
		ntptr = NilNode
		nt.NextPtr = &ntptr
	}

	//ebt.Collapses = append(ebt.Collapses, nt)
	Pool.CacheChan <- nt
	return nil
}
func (ebt *EBTree) CollapseInternalNode(nt *InternalNode, final bool) error {
	nl := len(nt.Children)
	var ntid []byte
	var flag bool
	flag = false
	for i := 0; i <= nl-1; i++ {
		switch ntct := (nt.Children[i].NodePtr).(type) {
		case *IdNode:
			continue
		case *LeafNode:
			if !final && byteCompare(ntct.Id, ebt.LastPath.Leaf.Id) == 0 {
				return nil
			}
			flag = true
			err := ebt.CollapseLeafNode(ntct)
			if err != nil {
				return err
			}
			var ntptr IdNode
			ntid = ntct.Id
			ntptr = ntid
			nt.Children[i].NodePtr = &ntptr
		case *InternalNode:
			if !final && ebt.isInPath(ntct.Id) {
				return nil
			}
			flag = true
			err := ebt.CollapseInternalNode(ntct, final)
			if err != nil {
				return err
			}
			var ntptr IdNode
			ntid = ntct.Id
			ntptr = ntid
			nt.Children[i].NodePtr = &ntptr
		default:
			err := errors.New("the child of internalnode should be idnode in collaspNode")
			return err
		}
	}
	if flag {
		//ebt.Collapses = append(ebt.Collapses, nt)
		Pool.CacheChan <- nt
	}
	return nil
}

func (ebt *EBTree) isInPath(id []byte) bool {
	var flag bool
	flag = false
	for i := 0; i < len(ebt.LastPath.Internals); i++ {
		if byteCompare(id, ebt.LastPath.Internals[i].Id) == 0 {
			flag = true
			return flag
		}
	}
	return flag
}

// commit prepare functions in Node
//End*****************************

//Start*****************************
// encode/decode functions in Node
func EncodeNode(n EBTreen) ([]byte, error) {
	var encode []byte
	var err error
	encode, err = rlp.EncodeToBytes(n)
	return encode, err
}

func DecodeInternal(elems []byte) (InternalNode, error) {
	var in InternalNode
	var err error

	kbuf, rest, _ := rlp.SplitString(elems)
	in.Id = kbuf
	//fmt.Println(kbuf)
	//fmt.Println(rest)
	elems = rest
	_, rest, _ = rlp.SplitString(elems)
	//fmt.Println(bbuf)
	//fmt.Println(rest)
	elems, _, _ = rlp.SplitList(elems)
	//the number of children
	c, _ := rlp.CountValues(elems)
	//fmt.Println(elems)
	//fmt.Println(c)
	for i := 0; i < c; i++ {
		var rest2 []byte
		elems, rest2, _ = rlp.SplitList(elems)

		//the number of fields in childData
		//c, _ = rlp.CountValues(elems)
		//fmt.Println(elems)
		//fmt.Println(c)
		var child ChildData

		bd1uf, rest, _ := rlp.SplitString(elems)
		//fmt.Println(bd1uf)
		child.Value = bd1uf
		//fmt.Println(rest)
		elems = rest
		bd2uf, _, _ := rlp.SplitString(elems)
		//fmt.Println(bd2uf)
		var npid IdNode
		npid = bd2uf
		child.NodePtr = &npid
		in.Children = append(in.Children, child)
		//fmt.Println(rest2)
		elems = rest2
	}

	return in, err
}

//todo:test the decode leaf node
func DecodeLeafNode(elems []byte) (LeafNode, error) {
	var le LeafNode
	var err error

	//get the id
	kbuf, rest, _ := rlp.SplitString(elems)
	le.Id = kbuf
	//fmt.Println(kbuf)
	//fmt.Println(rest)
	elems = rest

	//get the nextptr
	kbuf, rest, _ = rlp.SplitString(elems)
	var ntid IdNode
	ntid = kbuf
	le.NextPtr = &ntid
	//fmt.Println(kbuf)
	//fmt.Println(rest)
	elems = rest

	//get the data
	elems, _, _ = rlp.SplitList(elems)
	//the number of data
	c, _ := rlp.CountValues(elems)
	//fmt.Println(elems)
	//fmt.Println(c)
	for i := 0; i < c; i++ {
		var rest2 []byte
		elems, rest2, _ = rlp.SplitList(elems)

		//the number of fields in childData
		//c, _ = rlp.CountValues(elems)
		//fmt.Println(elems)
		//fmt.Println(c)
		var rd ResultD

		//get the value of resultd
		bd1uf, rest, _ := rlp.SplitString(elems)
		//fmt.Println(bd1uf)
		rd.Value = bd1uf
		//fmt.Println(rest)
		elems = rest

		//get the tds of resultd
		elems, _, _ = rlp.SplitList(elems)
		//the number of td
		tdc, _ := rlp.CountValues(elems)
		//fmt.Println(elems)
		//fmt.Println(tdc)
		for i := 0; i < tdc; i++ {
			var rest3 []byte
			elems, rest3, _ = rlp.SplitList(elems)

			var td TD
			//get the tds of td
			bd2uf, _, _ := rlp.SplitString(elems)
			//fmt.Println(bd2uf)
			td.IdentifierData = bd2uf
			rd.ResultData = append(rd.ResultData, td)

			//fmt.Println(rest3)
			elems = rest3
		}
		le.LeafDatas = append(le.LeafDatas, rd)
		elems = rest2
	}

	return le, err
}
func DecodeNode(encode []byte) (EBTreen, error) {
	elems, _, _ := rlp.SplitList(encode)
	//the number of fields in internal node
	c, _ := rlp.CountValues(elems)
	//fmt.Println(c)
	if c == 2 {
		in, err := DecodeInternal(elems)
		return &in, err
	} else if c == 3 {
		le, err := DecodeLeafNode(elems)
		return &le, err
	} else {
		err := errors.New("some thing wrong in encode")
		return nil, err
	}
}

// encode/decode functions in Node
//End*****************************
