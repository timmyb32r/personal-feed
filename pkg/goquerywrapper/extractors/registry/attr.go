package registry

import (
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/xerrors"
	"personal-feed/pkg/goquerywrapper/extractors/abstract"
)

type Attr struct {
	abstract.Extractor
	Attr string
}

var _ abstract.ExtractorSelectionToString = (*Attr)(nil)

func (a *Attr) Do(in *goquery.Selection) *string {
	attrVal, exists := in.Attr(a.Attr)
	if !exists {
		return nil
	}
	return &attrVal
}

func NewAttr(attr string) (*Attr, error) {
	if attr == "" {
		return nil, xerrors.New("Attr shouldn't be empty. it's optional entity - if you don't need it - just don't use it")
	}
	return &Attr{
		Attr: attr,
	}, nil
}
