package model

type PaginationOpts struct {
	Skip  int
	Limit int
}

type Filter struct {
	Value    string
	IsSended bool
}

type FieldsOpts struct {
	Fields []string
}
