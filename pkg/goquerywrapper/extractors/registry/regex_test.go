package registry

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/require"
	"personal-feed/pkg/goquerywrapper/extractors/util"
	"testing"
)

func TestNewRegex(t *testing.T) {
	t.Run("check number of capturing groups", func(t *testing.T) {
		var err error

		_, err = NewRegex("a")
		require.Error(t, err)

		_, err = NewRegex("(a)(b)")
		require.Error(t, err)

		_, err = NewRegex("(a)")
		require.NoError(t, err)
	})

	t.Run("test working case", func(t *testing.T) {
		extractor, err := NewRegex(`.*?<a href=[^>]+>(.*?)</a>.*`)
		require.NoError(t, err)

		doc, err := util.HTMLToDoc(`<body><div><a href="blablabla">lol</a></div></body>`)
		require.NoError(t, err)

		var selection *goquery.Selection = nil
		doc.Find(`body`).Each(func(_ int, s *goquery.Selection) {
			selection = s
		})
		require.NotNil(t, selection)

		selectionStr, err := goquery.OuterHtml(selection)
		require.NoError(t, err)

		result := extractor.Do(selectionStr)
		require.NotNil(t, result)
		require.Equal(t, `lol`, *result)
	})
}
