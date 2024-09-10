package model

type PaginationOpts struct {
	Skip  uint
	Limit uint
}

type Filter struct {
	IsSended bool
	Value    string
}
