package EBTree

/*




eth.topkVSearch()
miner.start()
personal.newAccount("ju")
personal.unlockAccount(eth.coinbase)
miner.stop()
./geth --port 30060 --rpcport 60060 -net.p-networkid 2805 --datadir /home/mimota/data/ console
./geth init /home/mimota/ethenv/genesis.json --datadir /home/mimota/data
build/bin/geth --port 30060 --rpcport 60060 --networkid 2804 --datadir /home/mimota/data/ console
build/bin/geth init /home/mimota/ethenv/genesis.json --datadir /home/mimota/data

eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:10})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:0})


eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(1,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(0.1,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(1.2,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(0.3,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(0.01,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(0.22,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(0.7,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(1.7,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(1.3,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(0.9,'ether')})

eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(1.6,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(0.01,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(0.82,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(0.99,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(1.8,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(1.001,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(1.3,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(2,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(0.81,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(1.92,'ether')})

eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(1.37,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(0.0401,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(0.272,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(1.77,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(1.07,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(2.3,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(0.9879,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(2.6,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(2.081,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(0.892,'ether')})

eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(1.899,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(2.8,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(1.7701,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(2.683,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(2.97,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(1.0081,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(1.09892,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(1.478737,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(1.79870401,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(0.687272,'ether')})

eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(1.78677,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(3.007,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(2.7868763,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(1.099879,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(2.87806,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(2.657081,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(3.892,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(1.0899,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(1.687678,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(2.997701,'ether')})

eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(3.6831,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(1.8991,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(2.81,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(1.77011,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(2.6183,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(2.917,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(1.00811,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(1.098192,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(1.4718737,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(1.719801,'ether')})

eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(0.617272,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(1.718677,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(3.1007,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(2.781663,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(1.099879,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(2.871806,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(2.651081,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(3.8192,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(1.08199,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(1.681678,'ether')})
eth.sendTransaction({from:eth.coinbase,to:"0x4751c4cd1ef729afc3232b2064565f1d692a9346",value:web3.toWei(2.971701,'ether')})

eth.topkVSearch()
miner.start()
personal.newAccount("ju")
personal.unlockAccount(eth.coinbase)
miner.stop()
./geth --port 30060 --rpcport 60060 --networkid 2805 --datadir /home/mimota/data/ console
./geth init /home/mimota/ethenv/genesis.json --datadir /home/mimota/data
build/bin/geth --port 30060 --rpcport 60060 --networkid 2805 --datadir /home/mimota/data/ console
build/bin/geth init /home/mimota/ethenv/genesis.json --datadir /home/mimota/data
*/
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
}

//向EBtree中插入叶子节点
func (t *EBTree) insertLeafNode(n *leafNode, pos int, parent *internalNode, value []byte, id []byte) (bool, *internalNode, error) {
	_, parent, err := moveChildren(parent, pos)
	if err != nil {
		err = wrapError(err, "insert leaf node: move child wrong")
		return false, nil, err
	}
	chil, err := createChild(value, n)
	if err != nil {
		err = wrapError(err, "insert leaf node: create child wrong")
		return false, nil, err
	}
	parent.Children[pos] = chil
	return true, parent, nil
}

//向EBtree中插入internal节点
func (t *EBTree) insertInternalNode(n *internalNode, pos int, parent *internalNode, value []byte, id []byte) (bool, EBTreen, error) {
	_, parent, err := moveChildren(parent, pos)
	if err != nil {
		err = wrapError(err, "insert internal node: move child wrong")
		return false, parent, err
	}
	chil, err := createChild(value, n)
	if err != nil {
		err = wrapError(err, "insert leaf node: create child wrong")
		return false, parent, err
	}
	parent.Children[pos] = chil
	return true, parent, nil
}


//向EBTree中插入数据
//special value被存储在特定结构中
//其他值正常存储在tree中
func (t *EBTree) InsertData(n EBTreen, pos uint8, parent *internalNode, value []byte, da []byte) (bool, *internalNode, error) {
	//判断value是否special
	sp, p := t.isSpecial(value)
	if sp {
		return t.InsertToSpecial(p, parent, da)
	}
	switch nt := n.(type) {
	case *leafNode:
		return t.InsertDataToLeaf(nt, pos, parent, value, da)
	case *internalNode:
		return t.InsertDataToInternal(nt, pos, parent, value, da)
	case *ByteNode:
		nbb, _ := n.cache()
		if string(nbb) == "" {
			dai, err := createData(value, da)
			if err != nil {
				err = wrapError(err, "insert data: create data wrong")
				return false, parent, err
			}
			var da []data
			da = append(da, dai)
			newn, err := CreateLeafNode(t, da)
			t.Root = &newn
			if err != nil {
				log.Info("err in create leaf node")
				return false, nil, err
			}
			return true, nil, nil
		} else {
			var nb []byte
			nb, _ = n.cache()
			decoden, err := t.resolveHash(nb)
			if err != nil {
				return false, nil, err
			}
			n = decoden
			return t.InsertData(decoden, pos, parent, value, da)
		}
	default:
		log.Info("n with wrong node type")
		err := errors.New("the node is not leaf or internal, something wrong")
		return false, nil, err
	}
	err := errors.New("the function reach to the bottom, something wrong")
	return false, nil, err
}

//insert data into internalnode child
func (t *EBTree) InsertToChild(i int, nt *internalNode, cd *child, value []byte, da []byte, parent *internalNode, decoden EBTreen) (bool, bool, *internalNode, EBTreen, error) {
	if Compare(cd.Value, value) >= 0 {

		//call the insert function to
		su, _, err := t.InsertData(decoden, uint8(i), nt, value, da)
		if !su {
			err = wrapError(err, "insert data: when the data was added into appropriate child, something wrong")
			return true, false, parent, decoden, err
		}
		//split internal node
		if len(nt.Children) > int(maxInternalNodeCount) {
			_, parent, err := t.split(nt, parent)
			if err != nil {
				err = wrapError(err, "insert data: when the data was added into appropriate child, split internal node wrong")
				return true, false, parent, decoden, err
			}
		}

		return true, true, parent, decoden, nil
	} else {
		return false, false, parent, decoden, nil
	}
}

func (t *EBTree) InsertToChildEncode(ct childEncode, nt *internalNode, i int, value []byte, da []byte, parent *internalNode) (bool, bool, *internalNode, *child, error) {
	cd, _, err := decodeChild(ct)
	decoden, err := t.resolveHash(ct)
	if err != nil {
		return true, false, nil, nil, err
	}
	cd.Pointer = decoden
	nt.Children[i] = cd
	//insert the data into specific child
	flag, su, parent, decoden, err := t.InsertToChild(i, nt, &cd, value, da, parent, decoden)

	return flag, su, parent, &cd, err
}

//insert data into internalNode
func (t *EBTree) InsertDataToInternal(nt *internalNode, pos uint8, parent *internalNode, value []byte, da []byte) (bool, *internalNode, error) {
	flag := false
	var j int
	for j = 0; j < len(nt.Children); j++ {
		switch ct := (nt.Children[j]).(type) {
		case childEncode:

			//insert the data into specific child
			flag, su, parent, cd, err := t.InsertToChildEncode(ct, nt, j, value, da, parent)
			nt.Children[j] = cd
			if err != nil {
				return false, nil, err
			}
			if !flag {
				continue
			} else {
				return su, parent, err
			}
		case child:
			if Compare(ct.Value, value) >= 0 {
				_, su, parent, cn, err := t.InsertToChild(j, nt, &ct, value, da, parent, ct.Pointer)
				ct.Pointer = cn
				//nt.Children[j] = ct
				return su, parent, err
			} else {
				continue
			}
		default:
			log.Info("wrong child type")
			err := errors.New("wrong child type:default in  InsertDataToInternal")
			return false, nil, err
		}

	}
	if !flag {
		//将该值插入到节点末尾
		//call the insert function to
		switch ct := (nt.Children[j-1]).(type) {
		case childEncode:
			//insert the data into specific child
			_, su, parent, cd, err := t.InsertToChildEncode(ct, nt, j-1, value, da, parent)
			nt.Children[j-1] = cd
			if err != nil {
				return false, nil, err
			}

			return su, parent, err

		case child:
			su, _, err := t.InsertData(ct.Pointer, uint8(j-1), nt, value, da)
			if !su {
				err = wrapError(err, "insert data: when the data was added into last child, something wrong")
				return false, parent, err
			}
			ct.Value = value
			nt.Children[j-1] = ct
			if len(nt.Children) > int(maxInternalNodeCount) {
				su, parent, err := t.split(nt, parent)
				if !su {
					err = wrapError(err, "insert data: when the data was added into last child, split internal node wrong")
					return false, parent, err
				}
				return true, parent, nil
			}

			return true, parent, nil
		default:
			err := errors.New("wrong child type:default in InsertDataToInternal")
			return false, nil, err
		}

	}
	return false, nil, nil
}


// Commit writes all nodes to the trie's memory database, tracking the internal
// and external (for account tries) references.
func (t *EBTree) Commit(onleaf LeafCallback) ([]byte, error) {
	if t == nil {
		panic("nil tree")
	}
	if t.Db == nil {
		panic("commit called on tree with nil database")
	}
	collapsedNode, err := t.foldRoot(t.Db, onleaf)
	if err != nil {
		return nil, err
	}
	//t.root = collapsedNode
	t.cachegen++
	if collapsedNode == nil {
		err := errors.New("collased node is nil")
		return nil, err
	}
	rt, _ := collapsedNode.cache()
	//log.Info(string(rt))
	return rt, nil
}

func (t *EBTree) foldRoot(db *Database, onleaf LeafCallback) (EBTreen, error) {
	//log.Info("into fold root ")
	if t.Root == nil {
		err := errors.New("tree is nil")
		return nil, err
	}
	f := newFolder(t.cachegen, t.cachelimit, onleaf)
	defer returnFolderToPool(f)

	//todo:fold the sequence and special

	return f.fold(t.Root, db, true)
}


// TryGet returns the value for key stored in the trie.
// The value bytes must not be modified by the caller.
// If a node was not found in the database, a MissingNodeError is returned.
func (t *EBTree) TryGet(value []byte) ([][]byte, error) {

	data, _, didResolve, err := t.tryGet(t.Root, value, 0)
	if err == nil && didResolve {
	}
	return data, err
}


//bool is used to mark whether the node is decoded right
func (t *EBTree) tryGet(origNode EBTreen, value []byte, pos int) ([][]byte, EBTreen, bool, error) {
	switch n := (origNode).(type) {
	case *ByteNode:
		//decode this node
		nc, _ := n.cache()
		decoden, err := t.resolveHash(nc)
		if err != nil {
			return nil, decoden, true, err
		} else {
			//得到该节点之后，需要将其中包含的data或者children解析出来
			switch ct := (decoden).(type) {
			//对于internal 节点
			case *internalNode:
				//先恢复ID
				ct.Id, _ = n.cache()
				//ct.Id=[]byte(n.fstring("a"))
				//再将该节点的子节点恢复出来。这里主要是要将pointer恢复成对应的子节点
				for i := 0; i < len(ct.Children); i++ {
					switch ctt := (ct.Children[i]).(type) {
					//首先判断该child是否已经被解码
					case childEncode:
						err := errors.New("child is encoded")
						return nil, nil, false, err
					case child:
						if Compare(ctt.Value, value) < 0 {
							continue
						} else {
							switch cpt := (ctt.Pointer).(type) {
							case *ByteNode:
								//将pointer对应的子节点解析出来
								var cb ByteNode
								cb, _ = cpt.cache()
								decodechild, err := t.resolveHash(cb)
								if err != nil {
									err = wrapError(err, "some thing wrong in resolve hash")
									return nil, nil, false, err
								}
								chi := child{}
								chi.Pointer = decodechild
								chi.Value = ctt.Value

								//继续去子节点中搜索对应value
								data, encodeNode, su, err := t.tryGet(chi.Pointer, value, 0)
								chi.Pointer = encodeNode
								//解析之后将其放回到对应域中
								ct.Children[i] = chi
								decoden = ct

								return data, decoden, su, err
							case *leafNode:
								continue
							case *internalNode:
								continue
							default:
								err := errors.New("wrong type")
								return nil, nil, false, err
							}
						}
					}
				}
				decoden = ct
				err := errors.New("not found the data for the value")
				return nil, decoden, false, err
			case *leafNode:
				ct.Id, _ = n.cache()
				data, newnode, didResolve, err := t.tryGet(ct, value, pos)
				if err != nil {
					return nil, decoden, true, err
				}
				return data, newnode, didResolve, nil
			}
		}
	case *leafNode:
		//判断value是否在当前data的范围中
		switch dt := (n.Data[0]).(type) {
		case dataEncode:
			err := errors.New("data is encoded for data 0")
			return nil, nil, false, err
		case data:
			switch dmt := (n.Data[len(n.Data)-1]).(type) {
			case dataEncode:
				err := errors.New("data is encoded for data len")
				return nil, nil, false, err
			case data:
				if Compare(value, dt.Value) < 0 || Compare(value, dmt.Value) > 0 {
					// key not found in trie
					err := errors.New("key not found")
					return nil, n, false, err
				}
				//在data中查找
				for i := 0; i < len(n.Data); i++ {
					switch dt := (n.Data[i]).(type) {
					case dataEncode:
						err := errors.New("data is encoded for data 0")
						return nil, nil, false, err
					case data:
						if Compare(dt.Value, value) == 0 {
							return dt.Keylist, n, true, nil
						}
					default:
						err := errors.New("wrong type")
						return nil, nil, false, err
					}
				}
				err := errors.New("no such value in leaf node!")
				return nil, n, false, err
			default:
				return nil, nil, false, nil

			}
		default:
			err := errors.New("wrong  type")
			return nil, nil, false, err
		}
	case *internalNode:
		for i := 0; i < len(n.Children); i++ {
			switch ct := (n.Children[i]).(type) {
			case childEncode:
				err := errors.New("data is encoded for data len")
				return nil, nil, false, err
			case child:
				if Compare(ct.Value, value) >= 0 {
					result, decodeChild, su, err := t.tryGet(ct.Pointer, value, 0)
					ct.Pointer = decodeChild
					n.Children[i] = ct
					return result, n, su, err
				}
			default:
				err := errors.New("wrong type")
				return nil, nil, false, err
			}
		}

		err := errors.New("no such value in internal node!")
		return nil, n, false, err
	default:
		panic(fmt.Sprintf("%T: invalid node: %v", origNode, origNode))
	}
	return nil, nil, false, nil

}


func addData(leaf leafNode, da data, position int) {
	var i int
	for i = len(leaf.Data); i > position; i-- {
		leaf.Data[i] = leaf.Data[i-1]
	}
	leaf.Data[i] = da
}

//ollapsed Node
func collapsedNode(n EBTreen) (EBTreen, error) {
	switch nt := (n).(type) {
	case *leafNode:
		return collapsedLeafNode(nt)
	case *internalNode:
		return collapsedInternalNode(nt)
	default:
		err := errors.New("wrong data type:default in collapsedNode")
		return nil, err
	}
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

func collapsedInternalNode(nt *internalNode) (*internalNode, error) {
	var collapsed internalNode
	if nt.Id == nil {
		err := errors.New("empty node")
		return nil, err
	}
	collapsed.Id = nt.Id
	for i := 0; i < len(nt.Children); i++ {
		if nt.Children[i] == nil {
			err := errors.New("n.children is nil")
			return nil, err
		}
		switch ct := (nt.Children[i]).(type) {
		case childEncode:
			err := errors.New("child is encoded in fold function")
			return nil, err
		case child:
			if ct.Pointer != nil {
				var pet ByteNode
				switch pt := (ct.Pointer).(type) {
				case *leafNode:
					pet = pt.Id
					var cchild child
					cchild.Pointer = &pet
					cchild.Value = ct.Value
					collapsed.Children = append(collapsed.Children, cchild)
				case *internalNode:
					pet = pt.Id
					var cchild child
					cchild.Pointer = &pet
					cchild.Value = ct.Value
					collapsed.Children = append(collapsed.Children, cchild)
				default:
					err := errors.New("wrong type")
					return nil, err
				}
			}
		default:
			err := errors.New("wrong type in child")
			return nil, err
		}
	}
	return &collapsed, nil
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

*/
