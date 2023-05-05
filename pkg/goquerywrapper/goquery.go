package goquerywrapper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"regexp"
)

type SubtreeExtractor func(*goquery.Selection) (string, bool)

var AddText = func(s *goquery.Selection) (string, bool) { return s.Text(), true }

func printWarn(warnText string, logger *logrus.Logger, s *goquery.Selection, query, attr, regex string) {
	selectedStr, _ := goquery.OuterHtml(s)
	logger.Warnf("%s, query:%s, attr:%s, regex:%s, selected:%s", warnText, query, attr, regex, selectedStr)
}

func ExtractURLAttrValSubstrByRegex(logger *logrus.Logger, url, query, attr, regex string, extractors ...SubtreeExtractor) ([][]string, error) {
	doc, err := URLToDoc(url)
	if err != nil {
		return nil, nil
	}
	return ExtractDocAttrValSubstrByRegex(logger, doc, query, attr, regex, extractors...)
}

func ExtractDocAttrValSubstrByRegex(logger *logrus.Logger, doc *goquery.Document, query, attr, regex string, extractors ...SubtreeExtractor) ([][]string, error) {
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
				docStr, _ := doc.Html()
				logger.Warnf("attribute is not exists, query: %s, attr: %s, regex: %s, doc: %s", query, attr, regex, docStr)
				return
			}
			str = attrVal
		} else {
			str, _ = goquery.OuterHtml(s)
		}

		subMatch := re.FindStringSubmatch(str)
		if len(subMatch) != 2 {
			printWarn("regex is not matched", logger, s, query, attr, regex)
			return
		}
		resultElem[0] = subMatch[1]

		for j, currExtractor := range extractors {
			currStr, ok := currExtractor(s)
			if !ok {
				printWarn("extractor returned !ok", logger, s, query, attr, regex)
				return
			}
			resultElem[j+1] = currStr
		}

		result = append(result, resultElem)
	})
	return result, nil
}
