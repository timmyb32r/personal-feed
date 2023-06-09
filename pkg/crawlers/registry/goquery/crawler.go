package goquery

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"personal-feed/pkg/crawlers"
	"personal-feed/pkg/goquerywrapper"
	"personal-feed/pkg/model"
)

type stNt string

func (n stNt) ID() string {
	return string(n)
}

type Crawler struct {
	source              model.Source
	commonGoparseSource CommonGoparseSource
	logger              *logrus.Logger
}

func (c *Crawler) CrawlerType() int {
	return CrawlerTypeCommonGoparse
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
	res, err := goquerywrapper.ExtractURLAttrValSubstrByRegex(c.logger, c.commonGoparseSource.URL, currLayer.Query, func(s *goquery.Selection) (string, error) {
		return goquerywrapper.DefaultSubtreeExtractor(logrus.New(), s, currLayer.Attr, currLayer.Regex)
	})
	if err != nil {
		return nil, nil
	}
	for _, el := range res {
		result = append(result, stNt(el[0]))
	}
	return result, nil
}

//---

func NewCrawler(source model.Source, logger *logrus.Logger) (crawlers.CrawlerTree, error) {
	commonGoparseSource := CommonGoparseSource{}
	err := json.Unmarshal([]byte(source.CrawlerMeta), &commonGoparseSource)
	if err != nil {
		return nil, xerrors.Errorf("unable to unmarshal crawlerMetaStr, crawlerMeta: %s, err: %w", source.CrawlerMeta, err)
	}
	return &Crawler{
		source:              source,
		commonGoparseSource: commonGoparseSource,
		logger:              logger,
	}, nil
}

func init() {
	crawlers.RegisterTree(NewCrawler, CrawlerTypeCommonGoparse)
}
