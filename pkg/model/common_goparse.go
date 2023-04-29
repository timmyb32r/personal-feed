package model

type CommonGoparseSourceItem struct {
	Query string
	Attr  string
	Regex string
}

type CommonGoparseSource struct {
	URL    string
	Layers []CommonGoparseSourceItem
}

type CommonGoparseSourceV struct {
	Val string
}

func (p *CommonGoparseSourceV) ID() string {
	return p.Val
}
