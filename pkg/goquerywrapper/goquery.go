package goquerywrapper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"regexp"
)

type SubtreeExtractor func(*goquery.Selection) (string, error)

var AddText = func(s *goquery.Selection) (string, error) { return s.Text(), nil }

func ExtractURLAttrValSubstrByRegex(logger *logrus.Logger, url, query string, extractors ...SubtreeExtractor) ([][]string, error) {
	doc, err := URLToDoc(url)
	if err != nil {
		return nil, nil
	}
	return Extract(logger, doc, query, extractors...)
}

func DefaultSubtreeExtractor(logger *logrus.Logger, s *goquery.Selection, attr, regex string) (string, error) {
	re, err := regexp.Compile(regex)
	if err != nil {
		return "", err
	}

	var str string
	if attr != "" {
		attrVal, exists := s.Attr(attr)
		if !exists {
			docStr, _ := s.Html()
			logger.Warnf("attribute is not exists, attr: %s, regex: %s, doc: %s", attr, regex, docStr)
			return "", xerrors.Errorf("regex is not matched")
		}
		str = attrVal
	} else {
		str, _ = goquery.OuterHtml(s)
	}

	subMatch := re.FindStringSubmatch(str)
	if len(subMatch) != 2 {
		return "", xerrors.Errorf("regex is not matched")
	}
	return subMatch[1], nil
}

func Extract(logger *logrus.Logger, doc *goquery.Document, query string, extractors ...SubtreeExtractor) ([][]string, error) {
	result := make([][]string, 0)

	doc.Find(query).Each(func(i int, s *goquery.Selection) {
		resultElem := make([]string, len(extractors))

		for j, currExtractor := range extractors {
			currStr, err := currExtractor(s)
			if err != nil {
				selectedStr, _ := goquery.OuterHtml(s)
				logger.Warnf("extractor returned err, err:%s, query:%s, selected:%s", err, query, selectedStr)
				return
			}
			resultElem[j] = currStr
		}

		result = append(result, resultElem)
	})
	return result, nil
}
