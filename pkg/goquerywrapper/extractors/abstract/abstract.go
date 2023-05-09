package abstract

import "github.com/PuerkitoBio/goquery"

type Extractor interface {
	isExtractor()
}

type ExtractorSelectionToString interface {
	Extractor
	Do(in *goquery.Selection) *string
}

type ExtractorSelectionToSelectionMaybeError interface {
	Extractor
	Do(in *goquery.Selection) (*goquery.Selection, error)
}

type ExtractorStringToString interface {
	Extractor
	Do(in string) *string
}
