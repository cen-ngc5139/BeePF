package utils

type Query struct {
	PageSize int
	PageNum  int
}

func NewQueryParma(size, num int) *Query {
	return &Query{
		PageSize: size,
		PageNum:  num,
	}
}
