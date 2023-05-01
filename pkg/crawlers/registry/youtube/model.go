package youtube

import "personal-feed/pkg/model"

const (
	CrawlerTypeYoutube = 1
)

type ContentSourceYoutubePlaylist struct {
	YoutubePlaylistID    string
	YoutubePlaylistTitle string
}

func (p *ContentSourceYoutubePlaylist) ID() string {
	return p.YoutubePlaylistID
}

//---

type ContentSourceYoutubeVideo struct {
	YoutubeVideoID          string
	YoutubeVideoTitle       string
	YoutubeVideoDescription string
}

func (p *ContentSourceYoutubeVideo) ID() string {
	return p.YoutubeVideoID
}

//---

type YoutubeClient interface {
	ListPlaylists(channelID string) ([]model.IDable, error) // channel -> []ContentSourceYoutubePlaylist
	ListPlaylist(playlistID string) ([]model.IDable, error) // ContentSourceYoutubePlaylist -> []ContentSourceYoutubeVideo
}

//---

type YoutubeSource struct {
	ChannelURL string
	ChannelID  string // if ChannelID is empty - will be resolved
}
