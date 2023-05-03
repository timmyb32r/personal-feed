package crawlers

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"personal-feed/pkg/model"
)

func NewCrawler(source model.Source, logger *logrus.Logger) (Crawler, error) {
	if currCrawlerFactory, ok := crawlerIDToFactory[source.CrawlerID]; ok {
		return currCrawlerFactory(source, logger)
	} else {
		return nil, xerrors.Errorf("unknown crawlerID: %d", source.CrawlerID)
	}
}
