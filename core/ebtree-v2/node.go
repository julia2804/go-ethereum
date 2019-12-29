package ebtree_v2

import "errors"

type EBTreen interface {
	fstring(string)
}
type (
	LeafNode struct {
		Id        []byte
		LeafDatas []ResultD
		NextPtr   *LeafNode
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
func (n *InternalNode) fstring(ind string) {
}
func (n *LeafNode) fstring(ind string) {
}
func (n *IdNode) fstring(ind string) {
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
			in.Children = append(in.Children, chd1)
			return in, nil
		} else {
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
			in.Children = append(in.Children, chd1)
			return in, nil
		} else {
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
		if uint64(lin) >= MaxInternalNodeCapability {
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
			if i != -1 {
				ebt.LastPath.Internals[i] = &nin
			}

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

func (ebt *EBTree) FindInNode(value []byte, n EBTreen) (*LeafNode, error) {
	var le *LeafNode
	var err error
	switch nt := n.(type) {
	case *LeafNode:
		return nt, nil
	case *InternalNode:
		i, err := ebt.SearchInNode(value, nt)
		if err != nil {
			return nil, err
		}
		return ebt.FindInNode(value, nt.Children[i].NodePtr)
	default:
		err := errors.New("wrong node type in FindInNode")
		return nil, err
	}

	return le, err
}

func (ebt *EBTree) SearchInNode(value []byte, n EBTreen) (int64, error) {
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
				return int64(m), nil
			}
		}
		//not found
		if hi < 0 {
			return int64(hi + 1), nil
		}
		return int64(hi), nil
	case *InternalNode:
		lo, hi := 0, len(nt.Children)-1
		for lo <= hi {
			m := (lo + hi) >> 1
			if byteCompare(value, nt.Children[m].Value) < 0 {
				lo = m + 1
			} else if byteCompare(value, nt.Children[m].Value) > 0 {
				hi = m - 1
			} else {
				return int64(m), nil
			}
		}
		if hi < 0 {
			return int64(hi + 1), nil
		}
		//not found
		return int64(hi), nil
	default:
		err := errors.New("wrong node type in SearchInNode")
		return -1, err
	}

}

// find functions in Node
//End*****************************
