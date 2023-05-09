package registry

import (
	"github.com/PuerkitoBio/goquery"
	"personal-feed/pkg/goquerywrapper/extractors/abstract"
)

type Text struct {
	abstract.Extractor
}

var _ abstract.ExtractorSelectionToString = (*Text)(nil)

func (a *Text) Do(in *goquery.Selection) *string {
	result := in.Text()
	return &result
}

func NewText() *Text {
	return &Text{}
}
