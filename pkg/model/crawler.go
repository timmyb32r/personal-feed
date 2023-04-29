package model

const (
	CrawlerTypeYoutube       = 1
	CrawlerTypeCommonGoparse = 2
)

type Crawler interface {
	CrawlerType() int
	Layers() []IDable // without root layer
	ListLayer(depth int, node Node) ([]IDable, error)
	// ListLayer handler can implement max_children_on_level, and not to parse useless paths!
	//
	// And here can be probabilistic approach
	// On every impossible mode start parsing with 1% probability
	// This way you always eventually get full info, and every time will parse a little
}
