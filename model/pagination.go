package model

type PaginationOpts struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
}

type Filter struct {
	Value    string
	IsSended bool
}

type FieldsOpts struct {
	Fields []string
}
