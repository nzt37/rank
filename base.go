package rank

/*
写在前面
游戏排行榜多以 redis zset,  红黑树， 跳表实现
考虑到需要并列排名的需求， 这里以红黑树做底层数据结构实现排行榜

排行榜接口 如 IRankList 示例
定义了EOrderType 排序类型， 可根据情况初始化排行榜，普通排名或并列排名


系统设计部分回答
可靠性要求讨论
若要保证排行榜的可靠性，可以分布式多节点部署排行榜，
多节点排行榜的数据一致性， 可以使用分布式一致性算法如raft来保证

性能要求
玩家数量达到百万级，考虑到内存情况， 可以多节点部署排行榜，由多个节点分区排行
类似一主多从的形式
主排行榜维护各个从排行榜当前维护的最高和最低名次，从排行榜维护真实排名

选做
已完成




*/


type (
	EOrderType   int
	ECompareType int
	RankElem     struct {
		m_ID       int // ID
		m_Value    int // 值
		m_Position int // 排名
	}
	RankListKey struct {
		m_SrvNum   int
		m_RankType int // 排行榜类型
	}
	RankConfig struct {
		orderType   EOrderType // 排行榜类型
		maxCapacity int        // 最大容量
	}
	RankList struct {
		m_RankListKey RankListKey
		m_RankConfig  RankConfig
	}
	IRankElem interface {
		Compare(other IRankElem) ECompareType
		GetValue() int
		GetNodeTotal() int
		GetElemTotal() int
		GetElemsLen() int
		AddElemTotal(v int)
		AddNodeTotal(v int)
		AddElems(elems *RankElem)
		GetElems() []*RankElem
	}
	IRankList interface {
		RankUpdate(rankID int, rankValue *RankElem) *RankElem
		RankDelete(rankID int) *RankElem
		RankFind(rankID int) *RankElem      // 目标排名
		PositionRange(a, b int) []*RankElem // a,b 区间的元素
	}
)

const (
	EOrderType_None             EOrderType = iota
	EOrderType_Normal                      // 普通排名
	EOrderType_Parallel                    // 并列排名，名次按人数递增[(value:1,pos:1),(value:2,pos:2)(value:2,pos:2)(value:3,pos:4)]
	EOrderType_ParallelIncrease            // 并列排名，名次按排名递增[(value:1,pos:1),(value:2,pos:2)(value:2,pos:2)(value:3,pos:3)]
)

const (
	ECompareType_None  ECompareType = iota
	ECompareType_Less               // 小于
	ECompareType_More               // 大于
	ECompareType_Equal              // 等于
)