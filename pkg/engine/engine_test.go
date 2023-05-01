package engine

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"personal-feed/pkg/crawlers/registry/youtube"
	"personal-feed/pkg/model"
	"personal-feed/pkg/repo/registry/in_memory"
	"testing"
)

//---------------------------------------------------------------------------------------------------------------------

type mockYoutubeClientTime1 struct {
}

func (c *mockYoutubeClientTime1) ListPlaylists(channelID string) ([]model.IDable, error) {
	return []model.IDable{
		&youtube.ContentSourceYoutubePlaylist{
			YoutubePlaylistID:    "my_playlist_id_1",
			YoutubePlaylistTitle: "my_playlist_title_1",
		},
	}, nil
}

func (c *mockYoutubeClientTime1) ListPlaylist(playlistID string) ([]model.IDable, error) {
	return []model.IDable{
		&youtube.ContentSourceYoutubeVideo{
			YoutubeVideoID:          "my_video_id_1",
			YoutubeVideoTitle:       "my_video_title_1",
			YoutubeVideoDescription: "my_video_description_1",
		},
		&youtube.ContentSourceYoutubeVideo{
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
		&youtube.ContentSourceYoutubePlaylist{
			YoutubePlaylistID:    "my_playlist_id_1",
			YoutubePlaylistTitle: "my_playlist_title_1",
		},
		&youtube.ContentSourceYoutubePlaylist{
			YoutubePlaylistID:    "my_playlist_id_2",
			YoutubePlaylistTitle: "my_playlist_title_2",
		},
	}, nil
}

func (c *mockYoutubeClientTime2) ListPlaylist(playlistID string) ([]model.IDable, error) {
	if playlistID == "my_playlist_id_1" {
		return []model.IDable{
			&youtube.ContentSourceYoutubeVideo{
				YoutubeVideoID:          "my_video_id_1",
				YoutubeVideoTitle:       "my_video_title_1",
				YoutubeVideoDescription: "my_video_description_1",
			},
			&youtube.ContentSourceYoutubeVideo{
				YoutubeVideoID:          "my_video_id_2",
				YoutubeVideoTitle:       "my_video_title_2",
				YoutubeVideoDescription: "my_video_description_2",
			},
			&youtube.ContentSourceYoutubeVideo{
				YoutubeVideoID:          "my_video_id_3",
				YoutubeVideoTitle:       "my_video_title_3",
				YoutubeVideoDescription: "my_video_description_3",
			},
		}, nil
	} else {
		return []model.IDable{
			&youtube.ContentSourceYoutubeVideo{
				YoutubeVideoID:          "my_video_id_4",
				YoutubeVideoTitle:       "my_video_title_4",
				YoutubeVideoDescription: "my_video_description_4",
			},
		}, nil
	}
}

//---------------------------------------------------------------------------------------------------------------------

func TestEngine(t *testing.T) {
	source := &model.Source{
		ID:          1,
		Description: "blablabla",
		CrawlerID:   1,
		CrawlerMeta: `{"ChannelURL": "https://www.youtube.com/blablabla"}`,
		Schedule:    "",
	}

	inMemoryRepoWrapped, _ := in_memory.NewRepo(struct{}{}, nil)
	inMemoryRepo := inMemoryRepoWrapped.(*in_memory.Repo)
	var log = logrus.New()
	var err error

	var youtubeSource youtube.YoutubeSource
	err = json.Unmarshal([]byte(source.CrawlerMeta), &youtubeSource)
	require.NoError(t, err)

	crawler1, err := youtube.NewCrawlerImpl(&youtubeSource, log, &mockYoutubeClientTime1{})
	require.NoError(t, err)
	engine1 := NewEngine(source, crawler1, inMemoryRepo)
	err = engine1.RunOnce()
	require.NoError(t, err)
	require.Equal(t, 3, inMemoryRepo.Len())

	crawler2, err := youtube.NewCrawlerImpl(&youtubeSource, log, &mockYoutubeClientTime1{})
	require.NoError(t, err)
	engine2 := NewEngine(source, crawler2, inMemoryRepo)
	err = engine2.RunOnce()
	require.NoError(t, err)
	require.Equal(t, 3, inMemoryRepo.Len())

	crawler3, err := youtube.NewCrawlerImpl(&youtubeSource, log, &mockYoutubeClientTime2{})
	require.NoError(t, err)
	engine3 := NewEngine(source, crawler3, inMemoryRepo)
	err = engine3.RunOnce()
	require.NoError(t, err)
	require.Equal(t, 6, inMemoryRepo.Len())
}
