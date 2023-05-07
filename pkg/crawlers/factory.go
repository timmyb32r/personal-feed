package crawlers

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"personal-feed/pkg/model"
)

func NewCrawler(source model.Source, logger *logrus.Logger) (AbstractCrawler, error) {
	if currCrawlerFactory, ok := crawlerTreeIDToFactory[source.CrawlerID]; ok {
		return currCrawlerFactory(source, logger)
	} else if currCrawlerFactory2, ok := crawlerChainIDToFactory[source.CrawlerID]; ok {
		return currCrawlerFactory2(source, logger)
	} else {
		return nil, xerrors.New("unknown crawler type")
	}
}
