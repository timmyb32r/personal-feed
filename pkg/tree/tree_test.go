package tree

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"personal-feed/pkg/model"
	"testing"
)

type treeTestStringObj string

func (o treeTestStringObj) ID() string {
	return string(o)
}

type treeTestIntObj int

func (o treeTestIntObj) ID() string {
	return fmt.Sprintf("%d", o)
}

type treeTestBoolObj bool

func (o treeTestBoolObj) ID() string {
	return fmt.Sprintf("%v", o)
}

func TestBaseCases(t *testing.T) {
	_, err := NewTree([]model.IDable{
		treeTestStringObj(""),
	})
	require.NoError(t, err)

	_, err = NewTree([]model.IDable{
		treeTestIntObj(2),
		treeTestIntObj(3),
	})
	require.NoError(t, err)

	_, err = NewTree([]model.IDable{
		treeTestIntObj(2),
		treeTestIntObj(3),
		treeTestIntObj(4),
	})
	require.NoError(t, err)

	_, err = NewTree([]model.IDable{
		treeTestStringObj(""),
		treeTestBoolObj(true),
	})
	require.NoError(t, err)

	_, err = NewTree(nil) // levels=0
	require.Error(t, err)
	_, err = NewTree([]model.IDable{}) // levels=0
	require.Error(t, err)
}

func TestYoutubeExample(t *testing.T) {
	tree, err := NewTree([]model.IDable{
		&model.ContentSourceYoutubePlaylist{},
		&model.ContentSourceYoutubeVideo{},
	})
	require.NoError(t, err)

	testPlaylists := tree.Root().(*node).ChildrenKeys()
	require.Equal(t, 0, len(testPlaylists))

	playlist := &model.ContentSourceYoutubePlaylist{
		YoutubePlaylistID:    "playlist_id1",
		YoutubePlaylistTitle: "playlist_name1",
	}
	playlistNodeObj, err := tree.Root().(*node).CreateOrGetChildNode(playlist)
	require.NoError(t, err)
	playlistNode := playlistNodeObj.(*node)

	require.Equal(t, `ROOT!playlist_id1`, playlistNode.ComplexKey().FullKey())

	video := &model.ContentSourceYoutubeVideo{
		YoutubeVideoID:          "video_id_1",
		YoutubeVideoTitle:       "video_title_1",
		YoutubeVideoDescription: "video_description_1",
	}
	_, err = playlistNode.CreateOrGetChildNode(video)
	require.NoError(t, err)

	videos, err := playlistNode.GetChildNodeByKeyID("video_id_1")
	require.NoError(t, err)
	videosTest := videos.(*model.ContentSourceYoutubeVideo)
	require.Equal(t, "video_id_1", videosTest.YoutubeVideoID)
	require.Equal(t, "video_title_1", videosTest.YoutubeVideoTitle)
	require.Equal(t, "video_description_1", videosTest.YoutubeVideoDescription)

	//

	fullKeyToInternalNode := tree.ExtractInternalNodes()
	fullKeysInternalNodes := make([]string, 0, len(fullKeyToInternalNode))
	for k := range fullKeyToInternalNode {
		fullKeysInternalNodes = append(fullKeysInternalNodes, k)
	}
	require.Equal(t, []string{`ROOT!playlist_id1`}, fullKeysInternalNodes)

	fullKeyToDoc := tree.ExtractDocs()
	docFullKeysArr := make([]string, 0, len(fullKeyToDoc))
	for k := range fullKeyToDoc {
		docFullKeysArr = append(docFullKeysArr, k)
	}
	require.Equal(t, []string{`ROOT!playlist_id1!video_id_1`}, docFullKeysArr)
}
