package goquery

import "personal-feed/pkg/goquerywrapper/extractors"

const (
	CrawlerTypeCommonGoparseChain = 3
)

type Item struct {
	Query        string
	Header       extractors.Program
	Link         extractors.Program
	BusinessTime extractors.Program
}

type CommonGoparseSource struct {
	URL      string
	Item     Item
	NextLink extractors.Program
	Content  extractors.Program
}
