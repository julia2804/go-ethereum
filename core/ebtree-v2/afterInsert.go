package ebtree_v2

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
)

func (t *EBTree) AfterInsertDataToTree(value []byte, da []byte) error {
	if len(value) == 0 {
		value = IntToBytes(0)
	}
	//fmt.Print("start to insert value :")
	//fmt.Println(value)
	err := t.AfterInsertDataToNode(&t.Root, value, da)
	if err != nil {
		return err
	}
	return t.AftersplitNode(&t.Root, nil, 0)
}

/**************start****************
insert to internal
***********************************/
//insert data into internalnode child, not split（reconstruct)
func (t *EBTree) AfterInsertToInternalChild(cd *ChildData, value []byte, da []byte) error {
	//if the child pointer is bytenode,we should construct a leaf/internal node from it first
	switch pt := (cd.NodePtr).(type) {
	case *LeafNode, *InternalNode:
		//call the insert function to
		err := t.AfterInsertDataToNode(&cd.NodePtr, value, da)
		if err != nil {
			fmt.Println("insert data: when the data was added into appropriate child, something wrong")
			return err
		}
		if byteCompare(cd.Value, value) > 0 {
			cd.Value = value
		}
		return nil
	case *IdNode:
		//decode the pointer
		ptid := pt.fstring()
		decoden, err := t.LoadNode(ptid)
		if err != nil {
			fmt.Println("decoden error")
			return err
		}
		//replace the bytenode with leaf/internal node
		cd.NodePtr = decoden
		//call the insert function to
		err = t.AfterInsertDataToNode(&cd.NodePtr, value, da)
		if err != nil {
			fmt.Println("insert data: when the data was added into appropriate child, something wrong")
			return err
		}
		if byteCompare(cd.Value, value) > 0 {
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
func (t *EBTree) AfterInsertDataToInternalNode(nt *InternalNode, value []byte, da []byte) error {
	var j int
	for j = 0; j < len(nt.Children); j++ {
		if byteCompare(nt.Children[j].Value, value) <= 0 {
			//将数据插入到对应的child中
			err := t.AfterInsertToInternalChild(&nt.Children[j], value, da)
			if err != nil {
				return err
			}
			//判断child对应的节点是否需要split
			err = t.AftersplitNode(&nt.Children[j].NodePtr, nt, j)
			//返回结果
			if err != nil {
				fmt.Println("insert data: when the data was added into appropriate child, something wrong")
				return err
			}
			return nil
		} else {
			continue
		}
	}

	//将该值插入到节点末尾
	//update the value of children
	//call the insert function to
	ct := (nt.Children[j-1])

	err := t.AfterInsertToInternalChild(&ct, value, da)
	nt.Children[j-1] = ct
	if err != nil {
		return err
	}
	//判断child对应的节点是否需要split
	err = t.AftersplitNode(&ct.NodePtr, nt, j-1)
	//返回结果
	if err != nil {
		fmt.Println("insert data: when the data was added into appropriate child, something wrong")
		return err
	}
	return nil

}

/**************end****************
insert to internal
***********************************/
/**************start****************
insert to leaf
***********************************/
//insert data into leafNode（reconstruct)
func moveData(n *LeafNode, pos int) (bool, *LeafNode, error) {
	//log.Info("into moveData")
	n.LeafDatas = append(n.LeafDatas, ResultD{})

	for i := len(n.LeafDatas) - 1; i > pos; i-- {
		n.LeafDatas[i] = n.LeafDatas[i-1]
	}
	return true, n, nil
}

//insert into dataNode（reconstruct)
func (t *EBTree) AfterInsertToLeafData(nt *LeafNode, i int, d ResultD, value []byte, da []byte) error {
	if byteCompare(d.Value, value) == 0 {
		//EBTree中已经存储了该value，因此，只要把data加入到对应到datalist中即可
		var td TD
		td.IdentifierData = da
		d.ResultData = append(d.ResultData, td)
		nt.LeafDatas[i] = d
		return nil
	} else {
		//说明EBTree中不存在value值，此时，需要构建data，并将其加入到节点中
		sucess, nt, err := moveData(nt, i)
		if !sucess {
			fmt.Println("insert data: move data wrong")
			return err
		}
		nt.LeafDatas[i] = NewLeafData(value, da)
		return nil
	}
}

func (t *EBTree) AfterInsertDataToLeafNode(nt *LeafNode, value []byte, da []byte) error {
	//向叶子节点插入数据
	//若当前节点为空时，直接插入节点。
	if len(nt.LeafDatas) == 0 {
		//create a data item for da
		dai := NewLeafData(value, da)
		nt.LeafDatas = append(nt.LeafDatas, dai)
		return nil
	}

	//遍历当前节点的所有data，将da插入合适的位置
	//value一定小于或等于当前节点到最大值
	for i := 0; i < len(nt.LeafDatas); i++ {
		//log.Info("find the appropriate position in nt datas")

		if byteCompare(nt.LeafDatas[i].Value, value) > 0 {
			continue
		} else {
			err := t.AfterInsertToLeafData(nt, i, nt.LeafDatas[i], value, da)
			if err != nil {
				return err
			}
			return nil
		}

	}
	//将该值插入到节点末尾
	//log.Info("the data should be put in the last ")
	dai := NewLeafData(value, da)
	nt.LeafDatas = append(nt.LeafDatas, dai)
	return nil

}

//将value插入到该节点或节点的子节点中
/**************end****************
insert to leaf
***********************************/
func (t *EBTree) AfterInsertDataToNode(n *EBTreen, value []byte, da []byte) error {

	switch nt := (*n).(type) {
	case *LeafNode:
		//insert into leafNode
		return t.AfterInsertDataToLeafNode(nt, value, da)
		//不进行split
	case *InternalNode:
		//insert into internal node
		return t.AfterInsertDataToInternalNode(nt, value, da)
	case *IdNode:
		//insert into byte node,need to use real node to replace it
		ntid := nt.fstring()
		decoden, err := t.LoadNode(ntid)
		if err != nil {
			return err
		}
		n = &decoden
		return t.AfterInsertDataToNode(n, value, da)
	case nil:
		dai := NewLeafData(value, da)

		var da []ResultD
		da = append(da, dai)
		newn := t.NewLeafNode()
		newn.LeafDatas = da
		t.Root = &newn
		return nil
	default:
		log.Info("n with wrong node type")
		err := errors.New("the node is not leaf or internal, something wrong")
		return err
	}
	err := errors.New("the function reach to the bottom in InsertDataToNode, something wrong")
	return err
}

/**************start****************
split node
***********************************/
//split leafnode into two leaf nodes(recontruct)
func (t *EBTree) AftersplitIntoTwoLeafNode(n *LeafNode, pos int) (*LeafNode, error) {
	//fmt.Println("split leafnode into two leaf nodes")
	newn := t.NewLeafNode()
	for j := len(n.LeafDatas) - 1; j >= pos; j-- {
		newn.LeafDatas = append(newn.LeafDatas, ResultD{})
	}
	for i := len(n.LeafDatas) - 1; i >= pos; i-- {

		newn.LeafDatas[i-pos] = n.LeafDatas[i]
		n.LeafDatas = append(n.LeafDatas[:i])
	}
	return &newn, nil
}
func AfteraddChild(internal InternalNode, chil ChildData, position int) (bool, InternalNode, error) {

	internal.Children = append(internal.Children, ChildData{})
	var i int
	for i = len(internal.Children) - 1; i > position; i-- {
		internal.Children[i] = internal.Children[i-1]
	}
	internal.Children[i] = chil
	return true, internal, nil
}

//split node(recontruct)
func (t *EBTree) AftersplitNode(n *EBTreen, parent *InternalNode, i int) error {
	switch nt := (*n).(type) {
	case *LeafNode:
		if (len(nt.LeafDatas)) <= MaxLeafNodeCapability {
			return nil
		}
		pos := (len(nt.LeafDatas) + 1) / 2
		//split the leaf node into two
		newn, err := t.AftersplitIntoTwoLeafNode(nt, pos)
		if err != nil {
			return err
		}
		if nt.NextPtr != nil {
			temp := nt.NextPtr
			nt.NextPtr = newn
			newn.NextPtr = temp
		} else {
			nt.NextPtr = newn
		}
		//对于根节点为叶子节点的情况,需要单独讨论
		if parent == nil {
			//需要为上级根节点确定value的值
			dt := (nt.LeafDatas[len(nt.LeafDatas)-1])
			chil, err := t.NewChildData(*n)
			chil.Value = dt.Value
			var chil2 ChildData
			if err != nil {
				fmt.Println("get leaf node position wrong: when parent is nil, create child wrong")
				return err
			}
			dtt := (newn.LeafDatas[len(newn.LeafDatas)-1])
			chil2, err = t.NewChildData(newn)
			chil2.Value = dtt.Value

			var children []ChildData
			children = append(children, chil)
			children = append(children, chil2)
			pnode := t.NewInternalNode()
			parent = &pnode
			parent.Children = children
			t.Root = parent
			return nil
		}

		//当前节点的元素被split，对应的parent中的children的值也要修改
		dt := (nt.LeafDatas[len(nt.LeafDatas)-1])
		ct := (parent.Children[i])
		ct.Value = dt.Value
		parent.Children[i] = ct
		if err != nil {
			return err
		}
		dtt := (newn.LeafDatas[len(newn.LeafDatas)-1])
		child2, err := t.NewChildData(newn)
		child2.Value = dtt.Value
		if err != nil {
			fmt.Println("split leaf node :create child to connect the new node to root")
			return err
		}
		su, presult, err := AfteraddChild(*parent, child2, int(i+1))
		if !su {
			fmt.Println("split leaf node: add the new child to root")
			return err
		}
		parent.Children = presult.Children
		return nil
	case *InternalNode:
		if (len(nt.Children)) <= MaxInternalNodeCapability {
			return nil
		}
		//carry the child node to new node
		var childList []ChildData
		pos := (len(nt.Children) + 1) / 2
		newn := t.NewInternalNode()
		newn.Children = childList
		for j := len(nt.Children) - 1; j >= pos; j-- {
			newn.Children = append(newn.Children, ChildData{})
		}
		for i := len(nt.Children) - 1; i >= pos; i-- {
			newn.Children[i-pos] = nt.Children[i]
			nt.Children = append(nt.Children[:i])
		}
		//直接将新节点链接到当前节点到后面，并链接到父节点上
		//为新创建到节点，创建一个child对象
		ct := (newn.Children[len(newn.Children)-1])
		chil, err := t.NewChildData(&newn)
		chil.Value = ct.Value
		if err != nil {
			fmt.Println("split internal node: create newn child wrong")
			return err
		}
		//对于根节点为叶子节点的情况,需要单独讨论
		if parent == nil {
			//需要为上级根节点确定value的值
			dt := (nt.Children[len(nt.Children)-1])

			chil2, err := t.NewChildData(*n)
			chil2.Value = dt.Value
			if err != nil {
				fmt.Println("get leaf node position wrong: when parent is nil, create child wrong")
				return err
			}
			var children []ChildData
			children = append(children, chil2)
			children = append(children, chil)
			pnode := t.NewInternalNode()
			parent = &pnode
			parent.Children = children
			if err != nil {
				fmt.Println("get leaf node position wrong: when parent is nil, create root")
				return err
			}
			t.Root = parent
			return nil
		}
		cpt := (parent.Children[i])
		cnt := (nt.Children[len(nt.Children)-1])
		cpt.Value = cnt.Value
		parent.Children[i] = cpt
		if err != nil {
			return err
		}
		su, presult, err := AfteraddChild(*parent, chil, int(i+1))
		if !su {
			fmt.Println("split internal node: add the new child to root")
			return err
		}
		parent.Children = presult.Children

	default:
		err := errors.New("node is defalut in splitNode")
		return err
	}
	err := errors.New("something wrong")
	return err
}

/**************end****************
split node
***********************************/
