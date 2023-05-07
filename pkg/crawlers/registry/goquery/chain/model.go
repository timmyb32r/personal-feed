package goquery

const (
	CrawlerTypeCommonGoparseChain = 3
)

type QueryIntoSelected struct {
	Attr  string
	Regex string
}

type QueryIntoDoc struct {
	Query string
	Attr  string
	Regex string
}

type CommonGoparseSourceItem struct {
	Query  string
	Header QueryIntoSelected
	Link   QueryIntoSelected
}

type CommonGoparseSource struct {
	URL  string
	Item CommonGoparseSourceItem
	Next QueryIntoDoc
}
