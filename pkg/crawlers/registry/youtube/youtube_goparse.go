package youtube

import (
	"golang.org/x/xerrors"
	"personal-feed/pkg/goquerywrapper"
	"personal-feed/pkg/model"
)

type YoutubeGoparseClient struct {
	youtubeSource YoutubeSource
}

func (c *YoutubeGoparseClient) ListPlaylists(_ string) ([]model.IDable, error) {
	result := make([]model.IDable, 0)
	url := c.youtubeSource.ChannelURL + "/playlists?view=1"
	res, err := goquerywrapper.ExtractURLAttrValSubstrByRegex(url, "a[id=video-title]", "href", `list=(.*)`, goquerywrapper.AddText)
	if err != nil {
		return nil, nil
	}
	for _, el := range res {
		result = append(result, &ContentSourceYoutubePlaylist{YoutubePlaylistID: el[0], YoutubePlaylistTitle: el[1]})
	}
	return result, nil
}

func (c *YoutubeGoparseClient) ListPlaylist(playlistID string) ([]model.IDable, error) {
	result := make([]model.IDable, 0)
	url := c.youtubeSource.ChannelURL + "/playlist?list=" + playlistID
	res, err := goquerywrapper.ExtractURLAttrValSubstrByRegex(url, "a[id=video-title]", "href", `/watch?v=(.*?)&`, goquerywrapper.AddText)
	if err != nil {
		return nil, nil
	}
	for _, el := range res {
		result = append(result, &ContentSourceYoutubeVideo{YoutubeVideoID: el[0], YoutubeVideoTitle: el[1], YoutubeVideoDescription: ""})
	}
	return result, nil
}

func newYoutubeGoparseClient(youtubeSource YoutubeSource) (*YoutubeGoparseClient, error) {
	if !ValidateLinkToChannel(youtubeSource.ChannelURL) {
		return nil, xerrors.Errorf("invalid youtube URL: %s", youtubeSource.ChannelURL)
	}
	return &YoutubeGoparseClient{
		youtubeSource: youtubeSource,
	}, nil
}
