package youtube_old

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"personal-feed/pkg/clients"
	"personal-feed/pkg/model"
)

type Crawler struct {
	logger        *logrus.Logger
	youtubeClient model.YoutubeClient
	youtubeSource model.YoutubeSource
	channelID     string
}

func (c *Crawler) CrawlerType() int {
	return model.CrawlerTypeYoutube
}

func (c *Crawler) Layers() []model.IDable {
	return []model.IDable{
		&model.ContentSourceYoutubePlaylist{},
		&model.ContentSourceYoutubeVideo{},
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

func NewCrawler(youtubeSource model.YoutubeSource, logger *logrus.Logger, youtubeClient model.YoutubeClient) (model.Crawler, error) {
	if !clients.ValidateLinkToChannel(youtubeSource.ChannelURL) {
		return nil, xerrors.Errorf("invalid youtube URL: %s", youtubeSource.ChannelURL)
	}
	channelID := youtubeSource.ChannelID
	if channelID == "" {
		var err error
		channelID, err = clients.ChannelIDByURL(youtubeSource.ChannelURL)
		if err != nil {
			return nil, err
		}
	}
	return &Crawler{
		logger:        logger,
		youtubeClient: youtubeClient,
		youtubeSource: youtubeSource,
		channelID:     channelID,
	}, nil
}
