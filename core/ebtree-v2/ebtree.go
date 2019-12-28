package ebtree_v2

var (
	//MaxLeafNodeCapability=uint8(32)
	//MaxInternalNodeCapability=uint64(256)
	MaxLeafNodeCapability     = uint8(5)
	MaxInternalNodeCapability = uint64(5)
)

type EBTree struct {
	Sequence  uint64
	Root      EBTreen
	LastPath  Path
	FirstLeaf *LeafNode
}
type Path struct {
	Leaf      *LeafNode
	Internals []InternalNode
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

	//Start**********
	//process leaf node
	le := ebt.NewLeafNode()
	le.LeafDatas = d
	if ebt.LastPath.Leaf == nil {
		ebt.LastPath.Leaf = &le
		ebt.Root = &le
		ebt.FirstLeaf = &le
		return nil
	}
	ebt.LastPath.Leaf.NextPtr = &le
	if ebt.LastPath.Internals == nil {
		in, err := ebt.CreateInternalNode(ebt.LastPath.Leaf, &le)
		ebt.FirstLeaf.NextPtr = &le
		if err != nil {
			return nil
		}
		ebt.LastPath.Internals = append(ebt.LastPath.Internals, in)
		ebt.Root = &in
		ebt.LastPath.Leaf = &le
		return nil
	}

	//process leaf node
	//End**********

	//Start**********
	//process internal node
	err := ebt.AdjustNodeInPath(0, ebt.LastPath.Leaf, &le)
	if err != nil {
		return err
	}
	//process internal node
	//End**********
	ebt.LastPath.Leaf = &le
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
		for j := 0; j < len(le.LeafDatas); j++ {
			ds = append(ds, le.LeafDatas[j])
			if int64(len(ds)) >= k {
				return ds
			}
		}
		le = le.NextPtr
	}
	return ds
}

func (ebt *EBTree) RangeSearch(begin int64, end int64) []ResultD {
	var ds []ResultD
	return ds
}

func (ebt *EBTree) FindFirstLeaf(value int64) *LeafNode {
	var le *LeafNode

	return le
}

// Search functions in EBTree
//End*****************************
