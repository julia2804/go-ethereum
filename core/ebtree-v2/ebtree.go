package ebtree_v2

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	MaxLeafNodeCapability     = 32
	MaxInternalNodeCapability = 256
	//MaxLeafNodeCapability     = uint8(3)
	//MaxInternalNodeCapability = uint64(3)
	MaxCollapseCapbility = uint64(100)
	MetaKey              = []byte("metas")
	NilNode              = IntToBytes(0)
)
var Pool *WorkerPool2

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
type Meta struct {
	Sequence  []byte
	Root      []byte
	FirstLeaf []byte
}

//Start*****************************
// Initial functions in EBTree
func NewEBTreeFromDb(db *Database) (*EBTree, error) {
	var ebt *EBTree
	var err error
	ebt = &EBTree{
		Db:        db,
		Sequence:  0,
		Root:      nil,
		FirstLeaf: nil,
	}
	var me Meta
	var metae []byte
	metae, err = db.GetTreeMetas(MetaKey)
	if err != nil {
		fmt.Println(err)
	} else {
		me, err = DecodeMeta(metae)
		if err != nil {
			return nil, err
		}
		var rid IdNode
		rid = me.Root
		var lid IdNode
		lid = me.FirstLeaf
		ebt.Root = rid
		ebt.FirstLeaf = lid

		if len(me.Root) != 0 {

			rootNode, err := ebt.LoadNode(me.Root)
			if err != nil {
				return ebt, err
			}
			ebt.Root = rootNode
		}
	}

	return ebt, err
}

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

func (ebt *EBTree) InsertDataToEBTree(d ResultD) error {
	var err error
	//le,err:=ebt.FindFirstLeaf(d.Value,true)

	return err
}

func (ebt *EBTree) Inserts(d []ResultD) error {
	var err error
	groupsNum := (len(d)) / (MaxLeafNodeCapability)
	rest := len(d) % MaxLeafNodeCapability

	for i := (0); i < groupsNum; i++ {
		//var ds []ResultD
		err = ebt.InsertDatasToTree(d[MaxLeafNodeCapability*i : MaxLeafNodeCapability*(i+1)])

		if err != nil {
			return err
		}
	}
	if rest != 0 {
		err = ebt.InsertDatasToTree(d[MaxLeafNodeCapability*groupsNum:])
	}
	return err
}

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
func (ebt *EBTree) TopkVSearch(k int64) ([]ResultD, error) {
	var ds []ResultD
	le := ebt.FirstLeaf
	for i := 0; int64(len(ds)) < k; i++ {
		if le == nil {
			return ds, nil
		}
		var lle int
		var ft LeafNode
		switch lt := le.(type) {
		case *LeafNode:
			lle = len(lt.LeafDatas)
			ft = *lt
		case IdNode:
			nt, err := ebt.LoadNode(lt.fstring())
			if err != nil {
				return ds, err
			}
			if nt == nil {
				return ds, nil
			}
			switch ntt := nt.(type) {
			case *LeafNode:
				lle = len(ntt.LeafDatas)
				ft = *ntt
			default:
				err := errors.New("wrong node type in firstleaf after load from levledb")
				return ds, err
			}
		case *IdNode:
			nt, err := ebt.LoadNode(lt.fstring())
			if err != nil {
				return ds, err
			}
			if nt == nil {
				return ds, nil
			}
			switch ntt := nt.(type) {
			case *LeafNode:
				lle = len(ntt.LeafDatas)
				ft = *ntt
			default:
				err := errors.New("wrong node type in firstleaf after load from levledb")
				return ds, err
			}
		default:
			err := errors.New("wrong node type in firstleaf")
			return ds, err
		}
		for j := 0; j < lle; j++ {
			ds = append(ds, ft.LeafDatas[j])
			if int64(len(ds)) >= k {
				return ds, nil
			}
		}
		le = ft.NextPtr
	}
	return ds, nil
}

func (ebt *EBTree) RangeSearch(begin []byte, end []byte) ([]ResultD, error) {
	var ds []ResultD
	le, err := ebt.FindFirstLeaf(end, false)
	if err != nil {
		return ds, err
	}
	if le == nil {
		return ds, nil
	}
	for le != nil && len(le.LeafDatas) > 0 && byteCompare(&le.LeafDatas[0].Value, &begin) >= 0 {
		for j := 0; j < len(le.LeafDatas); j++ {
			if byteCompare(&begin, &le.LeafDatas[j].Value) <= 0 && byteCompare(&end, &le.LeafDatas[j].Value) >= 0 {
				ds = append(ds, le.LeafDatas[j])
			} else if byteCompare(&begin, &le.LeafDatas[j].Value) > 0 {
				return ds, nil
			}
		}
		if le.NextPtr == nil {
			return ds, nil
		}
		switch nt := (le.NextPtr).(type) {
		case *LeafNode:
			le = nt
		case *IdNode:
			n, err := ebt.LoadNode(nt.fstring())
			if err != nil {
				return ds, err
			}
			if n == nil {
				return ds, nil
			}
			switch lnt := n.(type) {
			case *LeafNode:
				le = lnt
			default:
				err := errors.New("wrong node type from leveldb in  RangeSearch")
				return ds, err
			}
		default:
			err := errors.New("wrong node type of leaf.nextprt in RangeSearch")
			return ds, err
		}
	}
	return ds, err
}

func (ebt *EBTree) EquivalentSearch(value []byte) (ResultD, error) {
	var d ResultD
	le, err := ebt.FindFirstLeaf(value, false)
	if err != nil {
		return d, err
	}
	if le == nil {
		return d, nil
	}
	for le != nil && len(le.LeafDatas) > 0 && byteCompare(&le.LeafDatas[0].Value, &value) >= 0 {
		for j := 0; j < len(le.LeafDatas); j++ {
			if byteCompare(&value, &le.LeafDatas[j].Value) == 0 {
				d = le.LeafDatas[j]
				return d, nil
			} else if byteCompare(&value, &le.LeafDatas[j].Value) > 0 {
				return d, nil
			}
		}
		if le.NextPtr == nil {
			return d, nil
		}
		switch nt := (le.NextPtr).(type) {
		case *LeafNode:
			le = nt
		case *IdNode:
			n, err := ebt.LoadNode(nt.fstring())
			if err != nil {
				return d, err
			}
			if n == nil {
				return d, nil
			}
			switch lnt := n.(type) {
			case *LeafNode:
				le = lnt
			default:
				err := errors.New("wrong node type from leveldb in  RangeSearch")
				return d, err
			}
		default:
			err := errors.New("wrong node type of leaf.nextprt in EquivalentSearch")
			return d, err
		}
	}
	return d, err
}

func (ebt *EBTree) FindFirstLeaf(value []byte, flag bool) (*LeafNode, error) {
	var le *LeafNode
	var err error
	le, err = ebt.FindInNode(value, ebt.Root, flag)
	if err != nil {
		return nil, err
	}
	if le == nil {
		fmt.Println("nil leaf node")
		return nil, nil
	}
	//i,err,flag:=ebt.SearchInNode(value,le)
	_, err, _ = ebt.SearchInNode(value, le)
	return le, err
}

// Search functions in EBTree
//End*****************************

//Start*****************************
// commit functions in ebtree

//after insert the data into ebtree, we need to

func (ebt *EBTree) FinalCollapse() error {
	il := len(ebt.LastPath.Internals)
	if il == 0 {
		return ebt.CollapseLeafNode(ebt.LastPath.Leaf)
	}
	return ebt.CollapseInternalNode((ebt.LastPath.Internals[il-1]), true)

}

func (ebt *EBTree) CollapseEBTree() error {
	var err error
	if len(ebt.LastPath.Internals) == 0 {
		return err
	}
	l := len(ebt.LastPath.Internals) - 1
	err = ebt.CollapsedUnuseInternal(ebt.LastPath.Internals[l], l)
	return err
}

func (ebt *EBTree) CommitMeatas() error {
	var me Meta

	var lid IdNode
	switch lt := ebt.FirstLeaf.(type) {
	case *LeafNode:
		lid = lt.Id
	case *IdNode:
		lid = *lt
	default:
		err := errors.New("wrong node type of ebt.FirstLeaf")
		return err
	}

	var rid IdNode
	switch rt := ebt.Root.(type) {
	case *LeafNode:
		rid = rt.Id
	case *InternalNode:
		rid = rt.Id
	case *IdNode:
		rid = *rt
	default:
		err := errors.New("wrong node type of ebt.FirstLeaf")
		return err
	}

	me.FirstLeaf = lid
	me.Root = rid
	me.Sequence = IntToBytes(ebt.Sequence)

	en, err := EncodeMeata(me)
	if err != nil {
		return err
	}

	batch := ebt.Db.diskdb.NewBatch()
	err = ebt.Db.SetTreeMetas(MetaKey, en, batch)
	if err != nil {
		return err
	}
	return nil
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
	if err := batch.Write(); err != nil {
		return err
	}
	batch.Reset()
	return err
}

func (ebt *EBTree) CommitNode(treen EBTreen, batch ethdb.Batch) error {
	var err error
	//batch := ebt.Db.diskdb.NewBatch()
	encode, err := EncodeNode(treen)
	if err != nil {
		return err
	}
	switch et := (treen).(type) {
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
	if err := batch.Write(); err != nil {
		return err
	}
	batch.Reset()
	return err
}

// commit functions in Node
//End*****************************

//Start*****************************
// load data from leveldb functions in ebtree

func (ebt *EBTree) LoadNode(id []byte) (EBTreen, error) {
	var n EBTreen
	var err error
	if byteCompare(&id, &NilNode) == 0 {
		fmt.Println("the last leaf node")
		return nil, nil
	}
	enodes, err := ebt.Db.diskdb.Get(id)
	if err != nil {
		return nil, err
	}
	n, err = DecodeNode(enodes)
	return n, err
}

// load data from leveldb functions in Node
//End*****************************

//Start*****************************
// encode/decode functions in Node

func EncodeMeata(me Meta) ([]byte, error) {
	var encode []byte
	var err error
	encode, err = rlp.EncodeToBytes(me)
	return encode, err
}

func DecodeMeta(elems []byte) (Meta, error) {
	var me Meta
	var err error
	elems, _, _ = rlp.SplitList(elems)
	//the number of fields in internal node
	//c, _ := rlp.CountValues(elems)
	//fmt.Println(c)

	kbuf, rest, err := rlp.SplitString(elems)
	if err != nil {
		return me, err
	}
	me.Sequence = kbuf
	//fmt.Println(kbuf)
	//fmt.Println(rest)
	elems = rest

	kbuf, rest, err = rlp.SplitString(elems)
	if err != nil {
		return me, err
	}
	me.Root = kbuf
	//fmt.Println(kbuf)
	//fmt.Println(rest)
	elems = rest

	kbuf, rest, err = rlp.SplitString(elems)
	if err != nil {
		return me, err
	}
	me.FirstLeaf = kbuf
	//fmt.Println(kbuf)
	//fmt.Println(rest)
	elems = rest

	return me, err
}
