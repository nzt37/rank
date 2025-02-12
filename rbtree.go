package rank

type RBNode struct {
	m_key    IRankElem
	m_Parent *RBNode
	m_Left   *RBNode
	m_Right  *RBNode
	m_Red    bool
}

type RBTree struct {
	m_Root *RBNode
	m_Size int
}

type RBTreeElem struct {
	m_Value int
}

func (elem *RBTreeElem) Compare(t IRankElem) ECompareType {
	if elem.m_Value < t.GetValue() {
		return ECompareType_Less
	} else if elem.m_Value > t.GetValue() {
		return ECompareType_More
	} else {
		return ECompareType_Equal
	}
}

func (elem *RBTreeElem) GetValue() int {
	return elem.m_Value
}

func (elem *RBTreeElem) GetElemTotal() int {
	return 0
}

func (elem *RBTreeElem) GetNodeTotal() int {
	return 0
}

func (elem *RBTreeElem) GetElemsLen() int {
	return 0
}

func (elem *RBTreeElem) AddElemTotal(v int) {
	// return 0
}
func (elem *RBTreeElem) AddNodeTotal(v int) {
	// return 0
}
func (elem *RBTreeElem) AddElems(v *RankElem) {
	// return nil
}
func (elem *RBTreeElem) GetElems() []*RankElem {
	return nil
}

func NewRBTree() *RBTree {
	return &RBTree{}
}

func (tree *RBTree) First() *RBNode {
	x := tree.m_Root
	for x != nil && x.m_Left != nil {
		x = x.m_Left
	}
	return x
}

func (tree *RBTree) Last() *RBNode {
	x := tree.m_Root
	for x != nil && x.m_Right != nil {
		x = x.m_Right
	}
	return x
}

func (tree *RBTree) Find(key IRankElem) *RBNode {
	x := tree.m_Root
	for x != nil {
		switch v := x.m_key.Compare(key); v {
		case ECompareType_Less:
			x = x.m_Right
		case ECompareType_More:
			x = x.m_Left
		default:
			return x
		}
	}
	return x
}

func (tree *RBTree) Next(x *RBNode) *RBNode {
	if x.m_Right != nil {
		y := x.m_Right
		for y.m_Left != nil {
			y = y.m_Left
		}
		return y
	}
	y := x.m_Parent
	for y != nil && x == y.m_Right {
		x = y
		y = y.m_Parent
	}
	return y
}

func (tree *RBTree) Insert(key IRankElem) *RBNode {
	// 插入节点，已有则返回已有节点
	// 核心思想：在叶子节点插入红色，再在路径黑数晶相等规则下，解决可能产生的连续红问题
	z := tree.m_Root
	var p *RBNode
	for z != nil {
		p = z
		switch v := z.m_key.Compare(key); v {
		case ECompareType_Less:
			z = z.m_Right
		case ECompareType_More:
			z = z.m_Left
		case ECompareType_Equal:
			return z
		}
	}
	z = &RBNode{m_key: key, m_Parent: p, m_Red: true}
	if p == nil {
		tree.m_Root = z
	} else if key.Compare(p.m_key) == ECompareType_Less {
		p.m_Left = z
	} else {
		p.m_Right = z
	}
	tree.m_Size += 1

	// 红黑修复逻辑
	x := z
	for x.m_Parent != nil && x.m_Parent.m_Red {
		if x.m_Parent == x.m_Parent.m_Parent.m_Left { // 父亲L型
			y := x.m_Parent.m_Parent.m_Right
			if y != nil && y.m_Red { // 子红父红叔红
				x.m_Parent.m_Red = false
				y.m_Red = false
				x.m_Parent.m_Parent.m_Red = true
				x = x.m_Parent.m_Parent
				continue
			}
			if x == x.m_Parent.m_Right { // 父亲R型
				x = x.m_Parent
				LeftRotate(tree, x)
			}
			x.m_Parent.m_Red = false // LL型子红父红叔黑祖黑->父叔变黑，祖变红，对祖右旋，即将多
			x.m_Parent.m_Parent.m_Red = true
			RightRotate(tree, x.m_Parent.m_Parent)
		} else {
			y := x.m_Parent.m_Parent.m_Left
			if y != nil && y.m_Red == true {
				// 子红父红叔红，则祖黑->父叔层变黑，祖变红，红向祖传递
				x.m_Parent.m_Red = false
				y.m_Red = false
				x.m_Parent.m_Parent.m_Red = true
				x = x.m_Parent.m_Parent
				continue
			}
			if x == x.m_Parent.m_Left {
				x = x.m_Parent
				RightRotate(tree, x)
			}
			x.m_Parent.m_Red = false
			x.m_Parent.m_Parent.m_Red = true
			LeftRotate(tree, x.m_Parent.m_Parent)
		}
	}
	tree.m_Root.m_Red = false
	return z
}

func (tree *RBTree) Delete(key IRankElem) *RBNode {
	// 删除节点
	// 核心思想：后继节点替换删除节点，再解决删除黑色节点可能产生的连续红、路径黑数量相等规则问题
	z := tree.Find(key)
	if z == nil {
		return z
	}
	var child *RBNode
	var parent *RBNode
	red := false

	if z.m_Left != nil && z.m_Right != nil { // 删除节点同时拥有两个儿子，以右子树首节点为替换节点
		replace := tree.Next(z)
		if z.m_Parent == nil {
			tree.m_Root = replace
		} else if z == z.m_Parent.m_Left {
			z.m_Parent.m_Left = replace
		} else {
			z.m_Parent.m_Right = replace
		}
		child = replace.m_Right // replace是后继节点，无左儿子
		parent = replace.m_Parent
		red = replace.m_Red
		if z == parent { // z的右子树只有replace一个节点
			parent = replace
		} else {
			if child != nil {
				child.m_Parent = parent
			}
			parent.m_Left = child // z的右子树有多个节点，则replace必然是左儿子
			replace.m_Right = z.m_Right
			z.m_Right.m_Parent = replace
		}
		replace.m_Parent = z.m_Parent
		replace.m_Red = z.m_Red
		replace.m_Left = z.m_Left
		z.m_Left.m_Parent = replace
	} else {
		// 删除节点最多一个儿子，可以直接作为替换节点删除
		if z.m_Left != nil {
			child = z.m_Left
		} else {
			child = z.m_Right
		}
		parent = z.m_Parent
		red = z.m_Red
		if child != nil {
			child.m_Parent = parent
		}
		if parent == nil {
			if z == parent.m_Left {
				parent.m_Left = child
			} else {
				parent.m_Right = child
			}
		} else {
			tree.m_Root = child
		}

		parent.m_Right = child
	}
	tree.m_Size -= 1
	if red == true {
		z.m_Parent, z.m_Left, z.m_Right = nil, nil, nil
		return z
	}

	for (child == nil || !child.m_Red) && child != tree.m_Root {
		if child == parent.m_Left {
			w := parent.m_Right
			if w.m_Red {
				w.m_Red = false
				parent.m_Red = true
				LeftRotate(tree, parent)
				w = parent.m_Right
			}
			if (w.m_Left == nil || !w.m_Left.m_Red) && (w.m_Right == nil || !w.m_Right.m_Red) {
				w.m_Red = true
				child = parent
				parent = child.m_Parent
			} else {
				if w.m_Right == nil || !w.m_Right.m_Red {
					if w.m_Left != nil {
						w.m_Left.m_Red = false
					}
					w.m_Red = true
					RightRotate(tree, w)
					w = parent.m_Right
				}
				w.m_Red = parent.m_Red
				parent.m_Red = false
				if w.m_Right != nil {
					w.m_Right.m_Red = false
				}
				
				LeftRotate(tree, parent)
				child = tree.m_Root
			}
		} else {
			w := parent.m_Left
			if w.m_Red {
				w.m_Red = false
				parent.m_Red = true
				RightRotate(tree, parent)
				w = parent.m_Left
			}
			if (w.m_Left == nil || !w.m_Left.m_Red) && (w.m_Right == nil || !w.m_Right.m_Red) {
				w.m_Red = true
				child = parent
				parent = child.m_Parent
			} else {
				if w.m_Left == nil || !w.m_Left.m_Red {
					if w.m_Right != nil {
						w.m_Right.m_Red = false
					}
					w.m_Red = true
					LeftRotate(tree, w)
					w = parent.m_Left
				}
				w.m_Red = parent.m_Red
				parent.m_Red = false
				if w.m_Left != nil {
					w.m_Left.m_Red = false
				}
				
				RightRotate(tree, parent)
				child = tree.m_Root
			}
		}
	}
	if child != nil {
		child.m_Red = false
	}
	z.m_Parent, z.m_Left, z.m_Right = nil, nil, nil
	return z
}

func LeftRotate(t *RBTree, x *RBNode) {
	// 左旋， 右儿子变父，自己变左儿子，右儿子的左儿子变自己的右儿子
	y := x.m_Right
	x.m_Right = y.m_Left
	if y.m_Left != nil {
		y.m_Left.m_Parent = x
	}
	y.m_Parent = x.m_Parent
	if x.m_Parent == nil {
		t.m_Root = y
	} else if x == x.m_Parent.m_Left {
		x.m_Parent.m_Left = y
	} else {
		x.m_Parent.m_Right = y
	}
	y.m_Left = x
	x.m_Parent = y
}

func RightRotate(t *RBTree, y *RBNode) {
	// 右旋， 左儿子变父，自己变右儿子，左儿子的右儿子变自己左儿子
	x := y.m_Left
	y.m_Left = x.m_Right
	if x.m_Right != nil {
		x.m_Right.m_Parent = y
	}
	x.m_Parent = y.m_Parent
	if y.m_Parent == nil {
		t.m_Root = x
	} else if y == y.m_Parent.m_Left {
		y.m_Parent.m_Left = x
	} else {
		y.m_Parent.m_Right = x
	}
	x.m_Right = y
	y.m_Parent = x
}