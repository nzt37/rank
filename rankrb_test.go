package rank

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestRBTree(t *testing.T) {
	tree := NewRBTree()
	tree.Insert(&RBTreeElem{m_Value: 1})
	tree.Insert(&RBTreeElem{m_Value: -1})
	tree.Insert(&RBTreeElem{m_Value: 2})
	tree.Insert(&RBTreeElem{m_Value: 6})
	tree.Insert(&RBTreeElem{m_Value: 3})

	dfs(tree.m_Root)

	tree.Delete(&RBTreeElem{m_Value: 3})
	fmt.Println()
	dfs(tree.m_Root)
	a := tree.Find(&RBTreeElem{m_Value: 2})
	fmt.Println(a.m_key)

}

func dfs(root *RBNode) {
	if root == nil {
		return 
	}

	dfs(root.m_Left)

	fmt.Print(root.m_key.GetValue(), " ")


	dfs(root.m_Right)
}

func TestRankMgr(t *testing.T) {
	RankMgr := NewRankMgr()
	RankList := &RankListKey{m_SrvNum: 101, m_RankType: 1}
	rankRb := NewRankInRBtree()
	RankMgr.m_RankLists[RankList] = rankRb
	rankRb.Init(EOrderType_Normal)

	elems := []*RankElem{}

	for i := 0; i <= 10; i++ {
		elem := &RankElem{m_ID: i}
		elems = append(elems, elem)
	}

	
	// t.Log("开始了")
	mi := 1000001
	for i := 0; i < 1; i++ {
		for _, elem := range elems {
			score := rand.Intn(100000)
			elem.m_Value = score
			if score < mi {
				mi = score
			}
			rankRb.RankUpdate(elem.m_ID, elem)
			fmt.Println("添加了元素", elem.m_ID)
		}
	}
	a := rankRb.PositionRange(1,5)
	fmt.Println(mi)
	for _, v  := range a {
		// t.Log(v.m_Value)
		fmt.Print(v.m_Value, "  ")
	}
	
}

func BenchmarkRankUpdate(b *testing.B) {
	RankMgr := NewRankMgr()
	RankList := &RankListKey{m_SrvNum: 101, m_RankType: 1}
	rankRb := NewRankInRBtree()
	RankMgr.m_RankLists[RankList] = rankRb
	rankRb.Init(EOrderType_Normal)

	elems := []*RankElem{}

	for i := 0; i <= 1000; i++ {
		elem := &RankElem{m_ID: i}
		elems = append(elems, elem)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for  _, elem := range elems {
			rankRb.RankUpdate(elem.m_ID, elem)
		}
	}
}