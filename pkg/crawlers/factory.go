package crawlers

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"personal-feed/pkg/clients"
	commongoparse "personal-feed/pkg/crawlers/common_goparse"
	youtube "personal-feed/pkg/crawlers/youtube"
	"personal-feed/pkg/model"
)

func NewCrawler(source model.Source, logger *logrus.Logger) (model.Crawler, error) {
	switch source.CrawlerID {
	case model.CrawlerTypeYoutube:
		var youtubeSource model.YoutubeSource
		err := json.Unmarshal([]byte(source.CrawlerMeta), &youtubeSource)
		if err != nil {
			return nil, err
		}
		youtubeClient, err := clients.NewYoutubeGoparseClient(youtubeSource)
		if err != nil {
			return nil, xerrors.Errorf("unable to create youtube client: %w", err)
		}
		return youtube.NewCrawler(youtubeSource, logger, youtubeClient)
	case model.CrawlerTypeCommonGoparse:
		var commonGoparse model.CommonGoparseSource
		err := json.Unmarshal([]byte(source.CrawlerMeta), &commonGoparse)
		if err != nil {
			return nil, err
		}
		return commongoparse.NewCrawler(commonGoparse, logger)
	default:
		return nil, xerrors.Errorf("unknown crawler type: %d", source.CrawlerID)
	}
}
