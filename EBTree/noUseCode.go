package EBTree

/*test topkValueSearch
var k []byte
k = IntToBytes(uint64(100))
_, result, _ := tree.TopkValueSearch(k, true)

for i:=0;i<len(result);i++{


	for j:=0;j<len(result[i].data);j++{
		fmt.Println(result[i].data[j])
	}
}*/

/*var s1 []byte
value2, e :=new(big.Int).SetString("15532600000000000000000000000000000",10)
if !e{
	fmt.Println("error")
} else {
	bv2 := value2.Bytes()
	dif := 8 - len(bv2)
	b0 := byte(0)
	for {
		if dif <= 0 {
			break
		} else {
			s1 = append(s1, b0)
			dif = dif - 1
		}
	}
	for i := 0; i < len(bv2); i++ {
		s1 = append(s1, bv2[i])
	}
	result2, err := SearchNode(s1, tree.Root, tree)
		if err != nil {
			fmt.Printf("somethine wrong in search node")
			return
		}
		fmt.Printf("the result for 552:\n")
		for i, r := range result2 {
			fmt.Printf("the %dth:\n", i)
			fmt.Printf("%v",r)
			fmt.Println()
		}
		fmt.Println()
	}
	var k []byte
	k = IntToBytes(uint64(500000))

	_, result, _ := tree.RangeValueSearch(s0, s1, k)
	if(len(result)==0){
		fmt.Println("no data")
	}
	for i:=0;i<len(result);i++{
		fmt.Println("%d value:",i)
		fmt.Println(result[i].value)
		fmt.Println("data:")
		for j:=0;j<len(result[i].data);j++{
			fmt.Println(result[i].data[j])
		}
	}
	//combineAndPrintSearchValue(result, s0, tree, k, false)
	fmt.Println("first find")
}*/
/*_, _ = tree.Commit(nil)
switch rt := (tree.Root).(type) {
case *leafNode:
	//todo:
	tree.Db.Commit(rt.Id, true)

	triedb := tree.Db
	//TODO:
	tree, _ = New(rt.Id, triedb)
case *internalNode:
	//todo:
	tree.Db.Commit(rt.Id, true)

	triedb := tree.Db
	//TODO:
	tree, _ = New(rt.Id, triedb)
default:
	return

}
fmt.Println(tree.sequence)*/

/*	var s0 []byte
	value, e :=new(big.Int).SetString("245431000000000000",10)
	if !e{
		fmt.Println("error")
	} else {
		bv := value.Bytes()
		dif := 8 - len(bv)
		b0 := byte(0)
		for {
			if dif <= 0 {
				break
			} else {
				s0 = append(s0, b0)
				dif = dif - 1
			}
		}
		for i := 0; i < len(bv); i++ {
			s0 = append(s0, bv[i])
		}
	}
	result1, err := SearchNode(s0, tree.Root, tree)
		if err != nil {
			fmt.Printf("somethine wrong in search node")
			return
		}
		fmt.Printf("the result for 129:\n")
		for i, r := range result1 {
			fmt.Printf("the %dth:\n", i)
			fmt.Printf("%v",r)
			fmt.Println()
		}
		fmt.Println()
	}
}*/

//bak-code1
/*如果更新的是最大值，应该同时更新children.value
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
}*/

//当前节点为高频值节点
/*if nt.special {
	//获取当前节点在parent节点中的位置，对于根节点来说，此时需要重新创建一个根节点
	_,parent,pos,err:=getLeafNodePosition(nt,parent,t)
	if err!=nil {
		wrapError(err,"insertData:when nt.special is true,get leaf node position wrong")
	}

	if value == nt.Data[0].value {
		nt.Data[0].keylist=append(nt.Data[0].keylist,da)
		return true, parent, nil
	}
	var datalist []data
	dai,err:=createData(value,da)
	datalist=append(datalist,dai)
	if err!=nil {
		wrapError(err,"insertData:when nt.special is true,create data wrong")
	}
	newn,err:=createLeafNode(t,datalist)
	tmp:=nt.next
	tmpid:=nt.nextid
	nt.next=&newn
	nt.nextid=newn.id
	newn.next=tmp
	newn.nextid=tmpid

	//此时，需要为data新创建一个叶子节点
	if value < nt.Data[0].value {
		//将新节点插入到当前节点的前一个位置
		_,parent,err=t.insertLeafNode(&newn,pos,parent,value,newn.id)
		return true, parent, nil
	} else {
		//将新节点插入到当前节点的后一个位置
		_,parent,err=t.insertLeafNode(&newn,pos+1,parent,value,newn.id)
		return true, parent, nil
	}

}*/

//为value创建一个新的叶子节点
/*var datalist []data
	dai,err := createData(value, da)
	if err!=nil {
		return false,parent,err
	}
	datalist = append(datalist, dai)
	newn, err := createLeafNode(t, datalist)
	if err!=nil {
		return false,parent,err
	}

	//将特殊节点链接到tree上


	//查找当前节点的位置
	var pos uint8
	_,parent,pos,err=getLeafNodePosition(nt,parent,t)
	if err!=nil {
		wrapError(err,"insertData:when nt.special is true,get leaf node position wrong")
	}

	//如果value<当前节点的最小值，则将新叶子节点放在pos-1的位置，并链接到EBTree上
	//否则，创建一个新到叶子节点，存放被value截断之后到所有值，则将新创建到两个叶子节点放在pos+1，pos+2的位置，并链接到EBTree上
	if value < nt.Data[0].value {
		//TODO:6/5出现bug
		if pos==0 {
			newn.nextid=nt.id
			newn.next=nt
		} else {
			pre:=parent.Children[pos-1].pointer
			switch pret:=(pre).(type) {
			case *leafNode:
				pret.nextid=newn.id
				pret.next=&newn
				newn.next=nt
				newn.nextid=nt.id
			}

		}
		_,_,err=t.insertLeafNode(&newn, uint8(pos), parent, value, newn.id)
		if err!=nil {
			return false,parent,err
		}
		return true, parent, nil
	} else if value > nt.Data[nt.count-1].value {
		//将叶子节点插入到当前节点之后
		nt.nextid=newn.id
		nt.next=&newn
		_,_,err=t.insertLeafNode(&newn, pos+1, parent, value, newn.id)
		if err!=nil {
			return false,parent,err
		}
		return true, parent, nil
	} else {
		//叶子节点被插入到节点中间，当前节点被分裂成两个节点
		var i uint8
		i=0
		for _,d:=range nt.Data {
			if d.value>value {
				break
			}
			i++
		}
		_,nt,newnt,err:=t.splitIntoTwoLeaf(nt,i)
		tmp:=nt.next
		tmpid:=nt.nextid
		nt.next=&newn
		nt.nextid=newn.id
		newn.next=newnt
		newn.nextid=newnt.id
		newnt.next=tmp
		newnt.nextid=tmpid
		if err!=nil {
			wrapError(err,"when value is special, spilt into two leaf wrong")
			return false,parent,err
		}
		_,_,err=t.insertLeafNode(&newn, pos+1, parent, value, newn.id)
		if err!=nil {
			wrapError(err,"when value is special, spilt insert leaf wrong")
			return false,parent,err
		}
		_,parent,err=t.insertLeafNode(newnt,pos+2,parent,value,newn.id)
		if err!=nil {
			wrapError(err,"when value is special, spilt insert leaf wrong 2")
			return false,parent,err
		}

	}
	return true, parent, nil

}*/

//如果当前节点为Root节点，当发生分裂时，需要添加一个root节点
/*if parent == nil {
	var childlist []child
	var chi child
	chi.id = nt.id
	se, err := t.newSequence()
	if err!=nil {
		err=wrapError(err,"split internal node: when parent is nil, new sequence wrong")
		return false,parent,err
	}
	newr := &internalNode{se, childlist, 2, true}
	child1, err:= createChild(nt.Children[nt.count-1].value, nt, nt.id)
	if err!=nil {
		err=wrapError(err,"split internal node: when parent is nil, create nt child wrong")
		return false,parent,err
	}
	child2, err := createChild(newn.Children[newn.count-1].value, newn, newn.id)
	if err!=nil {
		err=wrapError(err,"split internal node: when parent is nil, create newn child wrong")
		return false,parent,err
	}
	newr.Children = append(newr.Children, child1)
	newr.Children = append(newr.Children, child2)
	t.root = newr
	return true,newr,nil
}*/

/*elems, _, err := rlp.SplitList(buf)
if err != nil {
	return nil, fmt.Errorf("decode error: %v", err)
}
switch c, _ := rlp.CountValues(elems); c {
case 1:
	n, err := decodeLeaf(id, elems, cachegen)
	return n, wrapError(err, "short")
case 8:
	n, err := decodeInternal(id, elems, cachegen)
	return n, wrapError(err, "full")
default:
	return nil, fmt.Errorf("invalid number of list elements: %v", c)
}*/
