package model

type PaginationOpts struct {
	Skip  int
	Limit int
}

type Filter struct {
	IsSended bool
	Value    string
}

type FieldsOpts struct {
	Fields []string
}
