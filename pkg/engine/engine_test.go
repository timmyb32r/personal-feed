package engine

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	youtube "personal-feed/pkg/crawlers/youtube"
	"personal-feed/pkg/db/in_memory"
	"personal-feed/pkg/model"
	"testing"
)

//---------------------------------------------------------------------------------------------------------------------

type mockYoutubeClientTime1 struct {
}

func (c *mockYoutubeClientTime1) ListPlaylists(channelID string) ([]model.IDable, error) {
	return []model.IDable{
		&model.ContentSourceYoutubePlaylist{
			YoutubePlaylistID:    "my_playlist_id_1",
			YoutubePlaylistTitle: "my_playlist_title_1",
		},
	}, nil
}

func (c *mockYoutubeClientTime1) ListPlaylist(playlistID string) ([]model.IDable, error) {
	return []model.IDable{
		&model.ContentSourceYoutubeVideo{
			YoutubeVideoID:          "my_video_id_1",
			YoutubeVideoTitle:       "my_video_title_1",
			YoutubeVideoDescription: "my_video_description_1",
		},
		&model.ContentSourceYoutubeVideo{
			YoutubeVideoID:          "my_video_id_2",
			YoutubeVideoTitle:       "my_video_title_2",
			YoutubeVideoDescription: "my_video_description_2",
		},
	}, nil
}

//---

type mockYoutubeClientTime2 struct {
}

func (c *mockYoutubeClientTime2) ListPlaylists(channelID string) ([]model.IDable, error) {
	return []model.IDable{
		&model.ContentSourceYoutubePlaylist{
			YoutubePlaylistID:    "my_playlist_id_1",
			YoutubePlaylistTitle: "my_playlist_title_1",
		},
		&model.ContentSourceYoutubePlaylist{
			YoutubePlaylistID:    "my_playlist_id_2",
			YoutubePlaylistTitle: "my_playlist_title_2",
		},
	}, nil
}

func (c *mockYoutubeClientTime2) ListPlaylist(playlistID string) ([]model.IDable, error) {
	if playlistID == "my_playlist_id_1" {
		return []model.IDable{
			&model.ContentSourceYoutubeVideo{
				YoutubeVideoID:          "my_video_id_1",
				YoutubeVideoTitle:       "my_video_title_1",
				YoutubeVideoDescription: "my_video_description_1",
			},
			&model.ContentSourceYoutubeVideo{
				YoutubeVideoID:          "my_video_id_2",
				YoutubeVideoTitle:       "my_video_title_2",
				YoutubeVideoDescription: "my_video_description_2",
			},
			&model.ContentSourceYoutubeVideo{
				YoutubeVideoID:          "my_video_id_3",
				YoutubeVideoTitle:       "my_video_title_3",
				YoutubeVideoDescription: "my_video_description_3",
			},
		}, nil
	} else {
		return []model.IDable{
			&model.ContentSourceYoutubeVideo{
				YoutubeVideoID:          "my_video_id_4",
				YoutubeVideoTitle:       "my_video_title_4",
				YoutubeVideoDescription: "my_video_description_4",
			},
		}, nil
	}
}

//---------------------------------------------------------------------------------------------------------------------

func TestEngine(t *testing.T) {
	source := model.Source{
		ID:          1,
		Description: "blablabla",
		CrawlerID:   1,
		CrawlerMeta: `{"ChannelURL": "https://www.youtube.com/c/blablabla"}`,
		Schedule:    "",
	}

	inMemoryDBClient := in_memory.NewInMemoryDatabaseClient()
	var log = logrus.New()
	var err error

	var youtubeSource model.YoutubeSource
	err = json.Unmarshal([]byte(source.CrawlerMeta), &youtubeSource)
	require.NoError(t, err)

	crawler1, err := youtube.NewCrawler(youtubeSource, log, &mockYoutubeClientTime1{})
	require.NoError(t, err)
	engine1 := NewEngine(source, crawler1, inMemoryDBClient)
	err = engine1.RunOnce()
	require.NoError(t, err)
	require.Equal(t, 3, inMemoryDBClient.Len())

	crawler2, err := youtube.NewCrawler(youtubeSource, log, &mockYoutubeClientTime1{})
	require.NoError(t, err)
	engine2 := NewEngine(source, crawler2, inMemoryDBClient)
	err = engine2.RunOnce()
	require.NoError(t, err)
	require.Equal(t, 3, inMemoryDBClient.Len())

	crawler3, err := youtube.NewCrawler(youtubeSource, log, &mockYoutubeClientTime2{})
	require.NoError(t, err)
	engine3 := NewEngine(source, crawler3, inMemoryDBClient)
	err = engine3.RunOnce()
	require.NoError(t, err)
	require.Equal(t, 6, inMemoryDBClient.Len())
}
