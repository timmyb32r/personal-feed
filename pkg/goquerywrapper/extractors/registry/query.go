package registry

import (
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/xerrors"
	"personal-feed/pkg/goquerywrapper/extractors/abstract"
)

type Query struct {
	abstract.Extractor
	Query string
}

var _ abstract.ExtractorSelectionToSelectionMaybeError = (*Query)(nil)

func (q *Query) Do(in *goquery.Selection) (*goquery.Selection, error) {
	resultArr := make([]*goquery.Selection, 0)
	in.Find(q.Query).Each(func(_ int, s *goquery.Selection) {
		resultArr = append(resultArr, s)
	})
	if len(resultArr) == 0 {
		return nil, nil
	} else if len(resultArr) == 1 {
		return resultArr[0], nil
	} else {
		doc, _ := goquery.OuterHtml(in)
		return nil, xerrors.Errorf("Query %s is matched %d times, doc: %s", q.Query, doc)
	}
}

func NewQuery(query string) (*Query, error) {
	if query == "" {
		return nil, xerrors.New("Query shouldn't be empty. it's optional entity - if you don't need it - just don't use it")
	}
	return &Query{
		Query: query,
	}, nil
}
