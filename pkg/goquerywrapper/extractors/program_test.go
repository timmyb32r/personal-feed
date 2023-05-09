package extractors

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/require"
	"personal-feed/pkg/goquerywrapper/extractors/util"
	"testing"
)

func TestProgram(t *testing.T) {
	doc, err := util.HTMLToDoc(`<body><div><a href="blablabla!ururu">lol</a></div></body>`)
	require.NoError(t, err)

	var selection *goquery.Selection = nil
	doc.Find(`body`).Each(func(_ int, s *goquery.Selection) {
		selection = s
	})
	require.NotNil(t, selection)

	//---1

	program1Text := `
	[
		{"Query": "a"},
		{"Attr": "href"},
		{"Regex": ".*?!(.*)"}
	]
	`
	program1, err := NewProgramFromString(program1Text)
	require.NoError(t, err)
	result1, err := program1.Do(selection)
	require.NoError(t, err)
	require.Equal(t, "ururu", result1)

	//---2

	program2Text := `
	[
		{"Query": "a"},
		{"Text": "!"}
	]
	`
	program2, err := NewProgramFromString(program2Text)
	require.NoError(t, err)
	result2, err := program2.Do(selection)
	require.NoError(t, err)
	require.Equal(t, "lol", result2)
}
