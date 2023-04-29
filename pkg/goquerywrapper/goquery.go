package goquerywrapper

import (
	"github.com/PuerkitoBio/goquery"
	"regexp"
)

type SubtreeExtractor func(*goquery.Selection) (string, bool)

var AddText = func(s *goquery.Selection) (string, bool) { return s.Text(), true }

func ExtractURLAttrValSubstrByRegex(url, query, attr, regex string, extractors ...SubtreeExtractor) ([][]string, error) {
	doc, err := URLToDoc(url)
	if err != nil {
		return nil, nil
	}
	return ExtractDocAttrValSubstrByRegex(doc, query, attr, regex, extractors...)
}

func ExtractDocAttrValSubstrByRegex(doc *goquery.Document, query, attr, regex string, extractors ...SubtreeExtractor) ([][]string, error) {
	re, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}

	result := make([][]string, 0)

	doc.Find(query).Each(func(i int, s *goquery.Selection) {
		resultElem := make([]string, 1+len(extractors))

		var str string
		if attr != "" {
			attrVal, exists := s.Attr(attr)
			if !exists {
				// TODO - print warning
				return
			}
			str = attrVal
		} else {
			str, _ = goquery.OuterHtml(s)
		}

		subMatch := re.FindStringSubmatch(str)
		if len(subMatch) != 2 {
			// TODO - print warning
			return
		}
		resultElem[0] = subMatch[1]

		for j, currExtractor := range extractors {
			currStr, ok := currExtractor(s)
			if !ok {
				// TODO - print warning
				return
			}
			resultElem[j+1] = currStr
		}

		result = append(result, resultElem)
	})
	return result, nil
}
