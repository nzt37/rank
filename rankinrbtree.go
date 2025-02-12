package rank

type RankRBNodeKey struct {
	m_Value     int
	m_Elems     []*RankElem
	m_NodeTotal int
	m_ElemTotal int
}

type RankInRBtree struct {
	RBTree
	RankConfig
	RankList
	m_RankElems map[int]*RankElem
}

func NewRankInRBtree() *RankInRBtree {
	return &RankInRBtree{m_RankElems: make(map[int]*RankElem)}
}

func (elems *RankRBNodeKey) GetValue() int {
	return elems.m_Value
}

func (elems *RankRBNodeKey) GetElemsLen() int {
	return len(elems.m_Elems)
}

func (elems *RankRBNodeKey) AddElemTotal(v int) {
	elems.m_ElemTotal += v
}

func (elems *RankRBNodeKey) AddNodeTotal(v int) {
	elems.m_NodeTotal += v
}

func (elems *RankRBNodeKey) AddElems(v *RankElem) {
	elems.m_Elems = append(elems.m_Elems, v)
}

func (elems *RankRBNodeKey) GetElems() []*RankElem {
	return elems.m_Elems
}

func (elems *RankRBNodeKey) GetNodeTotal() int {
	return elems.m_NodeTotal
}

func (elems *RankRBNodeKey) Compare(other IRankElem) ECompareType {
	t := other.(*RankRBNodeKey)
	if elems.m_Value < t.m_Value {
		return ECompareType_Less
	} else if elems.m_Value > t.m_Value {
		return ECompareType_More
	} else {
		return ECompareType_Equal
	}
}
func (elems *RankRBNodeKey) GetElemTotal() int {
	return elems.m_ElemTotal
}
func (rankRB *RankInRBtree) Init(orderType EOrderType) {
	rankRB.RankConfig.orderType = orderType
}

func (rankRB *RankInRBtree) SearchRank(elem *RankElem) int {
	elem.m_Position = 0
	tmpRBNodeKey := &RankRBNodeKey{m_Value: elem.m_Value}
	rbNode := rankRB.m_Root
	iRankType := rankRB.RankConfig.orderType
	for rbNode != nil {
		if tmpRBNodeKey.m_Value < rbNode.m_key.GetValue() {
			rbNode = rbNode.m_Left
		} else if tmpRBNodeKey.m_Value > rbNode.m_key.GetValue() {
			if iRankType == EOrderType_Normal || iRankType == EOrderType_Parallel {
				if rbNode.m_Left != nil {
					elem.m_Position += rbNode.m_Left.m_key.GetElemTotal()
				}
				elem.m_Position += rbNode.m_key.GetElemsLen()
			} else if iRankType == EOrderType_ParallelIncrease {
				if rbNode.m_Left != nil {
					elem.m_Position += rbNode.m_Left.m_key.GetNodeTotal()
				}
				elem.m_Position += 1
			}
			rbNode = rbNode.m_Right
		} else {
			if iRankType == EOrderType_Normal || iRankType == EOrderType_Parallel {
				if rbNode.m_Left != nil {
					elem.m_Position += rbNode.m_Left.m_key.GetElemTotal()
				}
				if iRankType == EOrderType_Normal {
					for i, v := range rbNode.m_key.GetElems() {
						if v.m_ID == elem.m_ID {
							elem.m_Position += i
							break
						}
					}
				}
			} else if iRankType == EOrderType_ParallelIncrease && rbNode.m_Left != nil {
				elem.m_Position += rbNode.m_Left.m_key.GetNodeTotal()
				break
			}
			break
		}
	}
	elem.m_Position += 1
	return elem.m_Position
}

func (rankRB *RankInRBtree) RankUpdate(rankID int, rankValue *RankElem) *RankElem {
	elem := rankRB.RankDelete(rankID)
	if elem == nil {
		elem = &RankElem{m_Value: rankValue.m_Value, m_ID: rankValue.m_ID}
	}
	// fmt.Println(elem.m_ID, rankID)
	rankRB.m_RankElems[rankID] = elem
	newRBNodeKey := &RankRBNodeKey{m_Value: rankValue.m_Value}
	rbNode := rankRB.m_Root
	var rbNodeParent *RBNode
	for rbNode != nil {
		if newRBNodeKey.m_Value < rbNode.m_key.GetValue() {
			rbNodeParent, rbNode = rbNode, rbNode.m_Left
		} else if newRBNodeKey.m_Value > rbNode.m_key.GetValue() {
			rbNodeParent, rbNode = rbNode, rbNode.m_Right
		} else {
			break
		}
	}
	for rbNodeParent != nil {
		rbNodeParent.m_key.AddElemTotal(1)
		if rbNode == nil {
			rbNodeParent.m_key.AddNodeTotal(1)
		}
		rbNodeParent = rbNodeParent.m_Parent
	}
	if rbNode != nil {
		rbNode.m_key.AddElemTotal(1)
		rbNode.m_key.AddElems(elem)
	} else {
		newRBNodeKey.m_NodeTotal, newRBNodeKey.m_ElemTotal = 1, 1
		newRBNodeKey.m_Elems = append(newRBNodeKey.m_Elems, elem)
		rankRB.Insert(newRBNodeKey)
	}
	rankRB.SearchRank(elem)
	return elem
}

func (rankRB *RankInRBtree) RankFind(rankID int) *RankElem {
	if elem, ok := rankRB.m_RankElems[rankID]; ok {
		rankRB.SearchRank(elem)
		return elem
	}
	return nil
}

func (rankRB *RankInRBtree) RankDelete(rankID int) *RankElem {
	// return nil
	elem, ok := rankRB.m_RankElems[rankID]
	if !ok {
		return nil
	}
	delete(rankRB.m_RankElems, rankID)
	// fmt.Println("删除后 ", rankRB.m_RankElems)
	rbNode := rankRB.Find(&RankRBNodeKey{m_Value: elem.m_Value})
	if rbNode == nil {
		return elem
	}
	rbNodeKey := rbNode.m_key
	needDelNode := false
	if rbNodeKey.GetElemsLen() < 2 {
		needDelNode = true
	}
	tmpRBNode := rbNode.m_Parent
	for tmpRBNode != nil {
		tmpRBNode.m_key.AddElemTotal(-1)
		if needDelNode {
			tmpRBNode.m_key.AddNodeTotal(-1)
		}
		tmpRBNode = tmpRBNode.m_Parent
	}
	if needDelNode {
		rankRB.Delete(rbNodeKey)
	} else {
		rbNodeKey.AddElemTotal(-1)
		key := rbNodeKey.(*RankRBNodeKey)
		for i, e := range key.m_Elems {
			if e.m_ID == elem.m_ID {
				key.m_Elems = append(key.m_Elems[:i], key.m_Elems[i+1:]...)
				break
			}
		}
	}
	return elem
}

func (rankRB *RankInRBtree) PositionRange(posA, posB int) []*RankElem {
	if posA > posB || posA <= 0 {
		return nil
	}
	l := posA
	rbNode := rankRB.m_Root
	iRankType := rankRB.RankConfig.orderType
	for rbNode != nil {
		nowNum := 0
		lTotal := 0
		if iRankType == EOrderType_Normal || iRankType == EOrderType_Parallel {
			nowNum = rbNode.m_key.GetElemsLen()
			if rbNode.m_Left != nil {
				lTotal += rbNode.m_Left.m_key.GetElemsLen()
			}
		} else if iRankType == EOrderType_ParallelIncrease {
			nowNum = 1
			if rbNode.m_Left != nil {
				lTotal += rbNode.m_Left.m_key.GetNodeTotal()
			}
		}
		if l <= lTotal {
			rbNode = rbNode.m_Left
		} else if l <= lTotal+nowNum {
			l -= lTotal + 1
			break
		} else {
			rbNode = rbNode.m_Right
			l -= lTotal + nowNum
		}
	}
	if rbNode == nil {
		return nil
	}
	if iRankType == EOrderType_Parallel && l > 0 {
		posA += len(rbNode.m_key.GetElems()) - 1
		l = 0
		rbNode = rankRB.Next(rbNode)
		if posA > posB {
			return nil
		}
	}
	rangeElems := []*RankElem{}
	for rbNode != nil {
		key := rbNode.m_key.(*RankRBNodeKey)
		for i := l; i < key.GetElemsLen(); i++ {
			if posA > posB {
				return rangeElems
			}
			elem := key.m_Elems[i]
			rangeElems = append(rangeElems, elem)
			if iRankType == EOrderType_Normal {
				posA += 1
			}
		}
		l = 0
		if iRankType == EOrderType_Parallel {
			posA += key.GetElemsLen()
		} else if iRankType == EOrderType_ParallelIncrease {
			posA += 1
		}
		rbNode = rankRB.Next(rbNode)
	}
	return rangeElems
}