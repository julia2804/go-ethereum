package EBTree

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
