package extractors

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/require"
	"personal-feed/pkg/goquerywrapper/extractors/abstract"
	"personal-feed/pkg/goquerywrapper/extractors/registry"
	"personal-feed/pkg/goquerywrapper/extractors/util"
	"testing"
)

func newQueryTestWrapper(t *testing.T, query string) *registry.Query {
	result, err := registry.NewQuery(query)
	require.NoError(t, err)
	return result
}

func newAttrTestWrapper(t *testing.T, attr string) *registry.Attr {
	result, err := registry.NewAttr(attr)
	require.NoError(t, err)
	return result
}

func newRegexTestWrapper(t *testing.T, regex string) *registry.Regex {
	result, err := registry.NewRegex(regex)
	require.NoError(t, err)
	return result
}

func TestExec(t *testing.T) {
	doc, err := util.HTMLToDoc(`<body><div><a href="blablabla!ururu">lol</a></div></body>`)
	require.NoError(t, err)

	var selection *goquery.Selection = nil
	doc.Find(`body`).Each(func(_ int, s *goquery.Selection) {
		selection = s
	})
	require.NotNil(t, selection)

	//---1

	result, err := Exec(selection, []abstract.Extractor{
		newQueryTestWrapper(t, "a"),
		newAttrTestWrapper(t, "href"),
		newRegexTestWrapper(t, `.*?!(.*)`),
	})
	require.NoError(t, err)
	require.Equal(t, "ururu", result)

	//---2

	result2, err := Exec(selection, []abstract.Extractor{
		newQueryTestWrapper(t, "a"),
		registry.NewText(),
	})
	require.NoError(t, err)
	require.Equal(t, "lol", result2)
}
