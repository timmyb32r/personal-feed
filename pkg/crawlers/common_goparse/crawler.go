package commongoparse

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"personal-feed/pkg/goquerywrapper"
	"personal-feed/pkg/model"
)

type stNt string

func (n stNt) ID() string {
	return string(n)
}

type Crawler struct {
	logger              *logrus.Logger
	commonGoparseSource model.CommonGoparseSource
}

func (c *Crawler) CrawlerType() int {
	return model.CrawlerTypeCommonGoparse
}

func (c *Crawler) Layers() []model.IDable {
	result := make([]model.IDable, 0)
	for range c.commonGoparseSource.Layers {
		result = append(result, stNt(""))
	}
	return result
}

func (c *Crawler) ListLayer(depth int, _ model.Node) ([]model.IDable, error) {
	currDepth := depth - 1 // TODO - fix
	if depth-1 < len(c.commonGoparseSource.Layers) {
		return c.listLayer(currDepth)
	} else {
		return nil, xerrors.Errorf("")
	}
}

//---

func (c *Crawler) listLayer(depth int) ([]model.IDable, error) {
	result := make([]model.IDable, 0)
	currLayer := c.commonGoparseSource.Layers[depth]
	res, err := goquerywrapper.ExtractURLAttrValSubstrByRegex(c.commonGoparseSource.URL, currLayer.Query, currLayer.Attr, currLayer.Regex, goquerywrapper.AddText)
	if err != nil {
		return nil, nil
	}
	for _, el := range res {
		result = append(result, stNt(el[0]))
	}
	return result, nil
}

//---

func NewCrawler(commonGoparseSource model.CommonGoparseSource, logger *logrus.Logger) (*Crawler, error) {
	return &Crawler{
		logger:              logger,
		commonGoparseSource: commonGoparseSource,
	}, nil
}
