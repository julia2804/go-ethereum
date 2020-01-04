package ebtree_v2

import "testing"

func TestCreateCache(t *testing.T) {
	var nt *LeafNode
	var ebt *EBTree

	pool := CreatPoolAndRun(ebt, 10, 10)
	pool.CacheChan <- nt

	close(pool.CacheChan)
}
