package registry

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/require"
	"personal-feed/pkg/goquerywrapper/extractors/util"
	"testing"
)

func TestNewQuery(t *testing.T) {
	t.Run("test working case", func(t *testing.T) {
		extractor, err := NewQuery("a")
		require.NoError(t, err)

		doc, err := util.HTMLToDoc(`<body><div><a href="blablabla">lol</a></div></body>`)
		require.NoError(t, err)

		var selection *goquery.Selection = nil
		doc.Find("body").Each(func(_ int, s *goquery.Selection) {
			selection = s
		})
		require.NotNil(t, selection)

		result, err := extractor.Do(selection)
		require.NotNil(t, result)

		resultStr, err := goquery.OuterHtml(result)
		require.NoError(t, err)
		require.Equal(t, `<a href="blablabla">lol</a>`, resultStr)
	})
}
