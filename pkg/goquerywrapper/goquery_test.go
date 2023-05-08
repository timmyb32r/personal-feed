package goquerywrapper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestExtractAttrValSubstrByRegex(t *testing.T) {
	t.Run("0 extractors", func(t *testing.T) {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(`<html><a id="video-title" href="/watch?v=blablabla&amp;list=ururu">qwert</a></html>`))
		require.NoError(t, err)
		vals, err := Extract(logrus.New(), doc, "a[id=video-title]", func(s *goquery.Selection) (string, error) {
			return DefaultSubtreeExtractor(logrus.New(), s, "href", `list=(.*)`)
		})
		require.NoError(t, err)
		require.Equal(t, 1, len(vals))
		require.Equal(t, "ururu", vals[0][0])
	})

	t.Run("1 extractor", func(t *testing.T) {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(`<html><a id="video-title" href="/watch?v=blablabla&amp;list=ururu">qwert</a></html>`))
		require.NoError(t, err)
		vals, err := Extract(
			logrus.New(),
			doc,
			"a[id=video-title]",
			func(s *goquery.Selection) (string, error) {
				return DefaultSubtreeExtractor(logrus.New(), s, "href", `list=(.*)`)
			},
			func(s *goquery.Selection) (string, error) { return s.Text(), nil })
		require.NoError(t, err)
		require.Equal(t, 1, len(vals))
		require.Equal(t, 2, len(vals[0]))
		require.Equal(t, "ururu", vals[0][0])
		require.Equal(t, "qwert", vals[0][1])
	})
}
