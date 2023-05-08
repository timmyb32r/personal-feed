package crawlers

import (
	"personal-feed/pkg/model"
)

type AbstractCrawler interface {
	CrawlerType() int
}

type CrawlerTree interface {
	CrawlerType() int
	Layers() []model.IDable // without root layer
	ListLayer(depth int, node model.Node) ([]model.IDable, error)
	// ListLayer handler can implement max_children_on_level, and not to parse useless paths!
	//
	// And here can be probabilistic approach
	// On every impossible mode start parsing with 1% probability
	// This way you always eventually get full info, and every time will parse a little
}

type CrawlerChain interface {
	CrawlerType() int
	Layers() []model.IDable
	ListItems(link string) ([]model.IDable, string, string, error) // nodes, nextLink, body, err
}
