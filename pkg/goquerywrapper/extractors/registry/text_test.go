package registry

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/require"
	"personal-feed/pkg/goquerywrapper/extractors/util"
	"testing"
)

func TestNewText(t *testing.T) {
	t.Run("test working case", func(t *testing.T) {
		extractor := NewText()

		doc, err := util.HTMLToDoc(`<a href="blablabla">lol</a>`)
		require.NoError(t, err)

		var selection *goquery.Selection = nil
		doc.Find("a").Each(func(_ int, s *goquery.Selection) {
			selection = s
		})
		require.NotNil(t, selection)

		result := extractor.Do(selection)
		require.NotNil(t, result)
		require.Equal(t, "lol", *result)
	})
}
