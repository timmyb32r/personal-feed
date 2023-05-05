package youtube

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"personal-feed/pkg/crawlers"
	"personal-feed/pkg/model"
)

type Crawler struct {
	source        model.Source
	logger        *logrus.Logger
	youtubeClient YoutubeClient
	youtubeSource *YoutubeSource
	channelID     string
}

func (c *Crawler) CrawlerType() int {
	return CrawlerTypeYoutube
}

func (c *Crawler) CrawlerDescr() string {
	serializedSource, _ := json.Marshal(c.source)
	return string(serializedSource)
}

func (c *Crawler) Layers() []model.IDable {
	return []model.IDable{
		&ContentSourceYoutubePlaylist{},
		&ContentSourceYoutubeVideo{},
	}
}

func (c *Crawler) ListLayer(depth int, node model.Node) ([]model.IDable, error) {
	if depth == 1 {
		return c.listPlaylists()
	} else if depth == 2 {
		return c.listPlaylist(node)
	} else {
		return nil, xerrors.Errorf("")
	}
}

//---

func (c *Crawler) listPlaylists() ([]model.IDable, error) {
	playlists, err := c.youtubeClient.ListPlaylists(c.channelID)
	if err != nil {
		return nil, xerrors.Errorf("unable to list playlists, channelID: %s, err: %w", c.channelID, err)
	}
	return playlists, nil
}

func (c *Crawler) listPlaylist(playlistNode model.Node) ([]model.IDable, error) {
	playlistID := playlistNode.ID()
	videos, err := c.youtubeClient.ListPlaylist(playlistID)
	if err != nil {
		return nil, xerrors.Errorf("unable to list playlist, playlistID: %s, err: %w", playlistID, err)
	}
	return videos, nil
}

//---

func NewCrawlerImpl(source model.Source, youtubeSource *YoutubeSource, logger *logrus.Logger, youtubeClient YoutubeClient) (crawlers.Crawler, error) {
	if !ValidateLinkToChannel(youtubeSource.ChannelURL) {
		return nil, xerrors.Errorf("invalid youtube URL: %s", youtubeSource.ChannelURL)
	}
	channelID := youtubeSource.ChannelID
	if channelID == "" {
		var err error
		channelID, err = ChannelIDByURL(youtubeSource.ChannelURL)
		if err != nil {
			return nil, err
		}
	}
	return &Crawler{
		source:        source,
		logger:        logger,
		youtubeClient: youtubeClient,
		youtubeSource: youtubeSource,
		channelID:     channelID,
	}, nil
}

func NewCrawler(source model.Source, logger *logrus.Logger) (crawlers.Crawler, error) {
	youtubeSource := YoutubeSource{}
	err := json.Unmarshal([]byte(source.CrawlerMeta), &youtubeSource)
	if err != nil {
		return nil, xerrors.Errorf("unable to unmarshal crawlerMetaStr, crawlerMeta: %s, err: %w", source.CrawlerMeta, err)
	}
	youtubeClient, err := newYoutubeGoparseClient(logger, youtubeSource)
	if err != nil {
		return nil, xerrors.Errorf("unable to create youtube client: %w", err)
	}
	return NewCrawlerImpl(source, &youtubeSource, logger, youtubeClient)
}

func init() {
	crawlers.Register(NewCrawler, CrawlerTypeYoutube)
}
