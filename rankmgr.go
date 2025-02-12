package rank

type RankMgr struct {
	m_RankLists map[*RankListKey]IRankList
}

func NewRankMgr() *RankMgr {
	return &RankMgr{m_RankLists: make(map[*RankListKey]IRankList)}
}