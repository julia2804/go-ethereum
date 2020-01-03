package ebtree_v2

import (
	"errors"
	"fmt"
)

var (
	//MaxLeafNodeCapability=uint8(32)
	//MaxInternalNodeCapability=uint64(256)
	MaxLeafNodeCapability     = uint8(3)
	MaxInternalNodeCapability = uint64(3)
	MaxCollapseCapbility      = uint64(100)
)

type EBTree struct {
	Sequence  uint64
	Db        *Database
	Root      EBTreen
	LastPath  Path
	FirstLeaf EBTreen
	Collapses []EBTreen
}
type Path struct {
	Leaf      *LeafNode
	Internals []*InternalNode
}

//Start*****************************
// Initial functions in EBTree
func NewEBTree() (*EBTree, error) {
	pa := Path{
		Leaf:      nil,
		Internals: nil,
	}
	ebt := &EBTree{
		Sequence:  0,
		Root:      nil,
		FirstLeaf: nil,
		LastPath:  pa,
	}
	return ebt, nil
}

func (ebt *EBTree) NewSequence() []byte {
	re := ebt.Sequence
	re = re + 1
	id := IntToBytes(uint64(re))
	ebt.Sequence = re
	return id
}

// Initial functions in EBTree
//End*****************************

//Start*****************************
// Insert functions in EBTree
func (ebt *EBTree) InsertDatasToTree(d []ResultD) error {
	//todo:last path is wrong
	//Start**********
	//process leaf node
	le := ebt.NewLeafNode()
	le.LeafDatas = d
	if ebt.LastPath.Leaf == nil {
		//there are no node in ebtree

		ebt.LastPath.Leaf = &le
		ebt.Root = &le
		ebt.FirstLeaf = &le
		return nil
	}
	ebt.LastPath.Leaf.NextPtr = &le
	if ebt.LastPath.Internals == nil {
		//there are no internal node in ebtree

		//leaf nodes produce a internal node
		in, err := ebt.CreateInternalNode(ebt.LastPath.Leaf, &le)
		switch ft := (ebt.FirstLeaf).(type) {
		case *LeafNode:
			ft.NextPtr = &le
			ebt.FirstLeaf = ft
		default:
			err := errors.New("wrong node type in ebt.firstleaf")
			return err
		}
		if err != nil {
			return nil
		}
		/*collapse the leaf node
		ebt.CollapseLeafNode(ebt.LastPath.Leaf)
		//todo:send the collapse node to chanel to be processed*/

		ebt.LastPath.Internals = append(ebt.LastPath.Internals, &in)
		ebt.Root = &in
		ebt.LastPath.Leaf = &le
		err = ebt.CollapseEBTree()
		if err != nil {
			return err
		}
		if len(ebt.Collapses) > int(MaxCollapseCapbility) {
			return ebt.CommitNodes()
		}
		return nil
	}

	//process leaf node
	//End**********

	//Start**********
	//process internal node
	err := ebt.AdjustNodeInPath(-1, ebt.LastPath.Leaf, &le)
	if err != nil {
		return err
	}
	//process internal node
	//End**********

	/*collapse the leaf node
	ebt.CollapseLeafNode(ebt.LastPath.Leaf)
	//todo:send the collapse node to chanel to be processed*/
	ebt.LastPath.Leaf = &le
	err = ebt.CollapseEBTree()
	if err != nil {
		return err
	}
	if len(ebt.Collapses) > int(MaxCollapseCapbility) {
		return ebt.CommitNodes()
	}

	return nil

}

// Insert functions in EBTree
//End*****************************

//Start*****************************
// Search functions in EBTree
func (ebt *EBTree) TopkVSearch(k int64) []ResultD {
	var ds []ResultD
	le := ebt.FirstLeaf
	for i := 0; int64(len(ds)) < k; i++ {
		if le == nil {
			return ds
		}
		var lle int
		var ft LeafNode
		switch lt := le.(type) {
		case *LeafNode:
			lle = len(lt.LeafDatas)
			ft = *lt
		default:
			fmt.Println("wrong node type in firstleaf")
			return ds
		}
		for j := 0; j < lle; j++ {
			ds = append(ds, ft.LeafDatas[j])
			if int64(len(ds)) >= k {
				return ds
			}
		}
		le = ft.NextPtr
	}
	return ds
}

func (ebt *EBTree) RangeSearch(begin []byte, end []byte) ([]ResultD, error) {
	var ds []ResultD
	le, err := ebt.FindFirstLeaf(end)
	if err != nil {
		return ds, err
	}
	for le != nil && len(le.LeafDatas) > 0 && byteCompare(le.LeafDatas[0].Value, begin) >= 0 {
		for j := 0; j < len(le.LeafDatas); j++ {
			if byteCompare(begin, le.LeafDatas[j].Value) <= 0 && byteCompare(end, le.LeafDatas[j].Value) >= 0 {
				ds = append(ds, le.LeafDatas[j])
			} else if byteCompare(begin, le.LeafDatas[j].Value) > 0 {
				return ds, nil
			}
		}
		if le.NextPtr == nil {
			return ds, nil
		}
		switch nt := (le.NextPtr).(type) {
		case *LeafNode:
			le = nt
		default:
			err := errors.New("wrong node type of leaf.nextprt in RangeSearch")
			return ds, err
		}
	}
	return ds, err
}

func (ebt *EBTree) EquivalentSearch(value []byte) (ResultD, error) {
	var d ResultD
	le, err := ebt.FindFirstLeaf(value)
	if err != nil {
		return d, err
	}
	for le != nil && len(le.LeafDatas) > 0 && byteCompare(le.LeafDatas[0].Value, value) >= 0 {
		for j := 0; j < len(le.LeafDatas); j++ {
			if byteCompare(value, le.LeafDatas[j].Value) == 0 {
				d = le.LeafDatas[j]
				return d, nil
			} else if byteCompare(value, le.LeafDatas[j].Value) > 0 {
				return d, nil
			}
		}
		if le.NextPtr == nil {
			return d, nil
		}
		switch nt := (le.NextPtr).(type) {
		case *LeafNode:
			le = nt
		default:
			err := errors.New("wrong node type of leaf.nextprt in EquivalentSearch")
			return d, err
		}
	}
	return d, err
}

func (ebt *EBTree) FindFirstLeaf(value []byte) (*LeafNode, error) {
	var le *LeafNode
	var err error
	le, err = ebt.FindInNode(value, ebt.Root)
	return le, err
}

// Search functions in EBTree
//End*****************************

//Start*****************************
// commit functions in ebtree

//after insert the data into ebtree, we need to
func (ebt *EBTree) CollapseEBTree() error {
	var err error
	if len(ebt.LastPath.Internals) == 0 {
		return err
	}
	l := len(ebt.LastPath.Internals) - 1
	err = ebt.CollapsedUnuseInternal(ebt.LastPath.Internals[l], l)
	return err
}

func (ebt *EBTree) CommitNodes() error {
	var err error
	batch := ebt.Db.diskdb.NewBatch()
	for i := 0; i < len(ebt.Collapses); i++ {
		encode, err := EncodeNode((ebt.Collapses[i]))
		if err != nil {
			return err
		}
		switch et := (ebt.Collapses[i]).(type) {
		case *LeafNode:
			err = ebt.Db.commit(et.Id, encode, batch)
			if err != nil {
				return err
			}
		case *InternalNode:
			err = ebt.Db.commit(et.Id, encode, batch)
			if err != nil {
				return err
			}
		default:
			err = errors.New("wrong node type of ebtree")
			return err
		}

	}
	return err
}

// commit functions in Node
//End*****************************

//Start*****************************
// load data from leveldb functions in ebtree

func (ebt *EBTree) loadNode(id []byte) (EBTreen, error) {
	var n EBTreen
	var err error
	enodes, err := ebt.Db.diskdb.Get(id)
	if err != nil {
		return nil, err
	}
	n, err = DecodeNode(enodes)
	return n, err
}

// load data from leveldb functions in Node
//End*****************************
