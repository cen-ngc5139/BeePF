package utils

type Query struct {
	PageSize   int
	PageNum    int
	IsAdmin    bool
	Authorized []string // 已经授权的集群列表
}

func NewQueryParma(size, num int, isAdmin bool, authorized []string) *Query {
	return &Query{
		PageSize:   size,
		PageNum:    num,
		IsAdmin:    isAdmin,
		Authorized: authorized,
	}
}
