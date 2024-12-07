package wrapify

func NewPagination() *pagination {
	p := &pagination{}
	return p
}

func NewMeta() *meta {
	m := &meta{
		CustomFields: map[string]interface{}{},
	}
	return m
}

func NewHeader() *header {
	h := &header{}
	return h
}

func NewWrap() *wrapper {
	w := &wrapper{
		Debug: &map[string]interface{}{},
	}
	return w
}
