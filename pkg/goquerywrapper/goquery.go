package goquerywrapper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"personal-feed/pkg/goquerywrapper/extractors"
	"personal-feed/pkg/goquerywrapper/extractors/util"
	"regexp"
)

type SubtreeExtractor func(*goquery.Selection) (string, error)

var AddText = func(s *goquery.Selection) (string, error) { return s.Text(), nil }

func ExtractURLAttrValSubstrByRegex(logger *logrus.Logger, url, query string, extractors ...SubtreeExtractor) ([][]string, error) {
	doc, err := util.URLToDoc(url)
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

	if regex == "" {
		return str, nil
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

func ExtractItemsByProgram(logger *logrus.Logger, doc *goquery.Document, query string, extractors ...*extractors.GoQueryProgram) ([][]string, error) {
	result := make([][]string, 0)

	doc.Find(query).Each(func(i int, s *goquery.Selection) {
		resultElem := make([]string, 0, len(extractors))

		for j, currExtractor := range extractors {
			currStr, err := currExtractor.Do(s)
			if err != nil {
				selectedStr, _ := goquery.OuterHtml(s)
				logger.Warnf("extractor #%d returned err, err:%s, query:%s, selected:%s", j, err, query, selectedStr)
				return
			}
			resultElem = append(resultElem, currStr)
		}

		result = append(result, resultElem)
	})
	return result, nil
}

func ExtractByProgram(doc *goquery.Document, program *extractors.GoQueryProgram) (string, error) {
	return program.Do(doc.Selection)
}
