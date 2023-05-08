package goquery

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"log"
	"net/url"
	"personal-feed/pkg/crawlers"
	"personal-feed/pkg/goquerywrapper"
	"personal-feed/pkg/model"
)

type stNt struct {
	Link       string
	HeaderText string
}

func (n stNt) ID() string {
	return n.Link
}

type stContent struct {
	Link    string
	Content string
}

func (n stContent) ID() string {
	return n.Link
}

type Crawler struct {
	source              model.Source
	commonGoparseSource CommonGoparseSource
	urlGetter           URLGetter
	logger              *logrus.Logger
}

func (c *Crawler) CrawlerType() int {
	return CrawlerTypeCommonGoparseChain
}

func (c *Crawler) Layers() []model.IDable {
	return []model.IDable{
		stNt{Link: "", HeaderText: ""},
		stContent{Link: "", Content: ""},
	}
}

func (c *Crawler) listPage(page, link string) ([]model.IDable, string, string, error) {
	doc, err := goquerywrapper.HTMLToDoc(page)
	if err != nil {
		return nil, "", "", xerrors.Errorf("unable to convert html page to doc, link: %s, err: %w", link, err)
	}
	res, err := goquerywrapper.Extract(c.logger, doc, c.commonGoparseSource.Item.Query, func(s *goquery.Selection) (string, error) {
		return goquerywrapper.DefaultSubtreeExtractor(c.logger, s, c.commonGoparseSource.Item.Header.Attr, c.commonGoparseSource.Item.Header.Regex)
	}, func(s *goquery.Selection) (string, error) {
		return goquerywrapper.DefaultSubtreeExtractor(c.logger, s, c.commonGoparseSource.Item.Link.Attr, c.commonGoparseSource.Item.Link.Regex)
	})
	if err != nil {
		return nil, "", "", xerrors.Errorf("unable to extract from link %s, err: %w", link, err)
	}
	result := make([]model.IDable, 0)
	for _, el := range res {
		result = append(result, stNt{HeaderText: el[0], Link: el[1]})
	}

	// next_link

	res, err = goquerywrapper.Extract(c.logger, doc, c.commonGoparseSource.Next.Query, func(s *goquery.Selection) (string, error) {
		return goquerywrapper.DefaultSubtreeExtractor(c.logger, s, c.commonGoparseSource.Next.Attr, c.commonGoparseSource.Next.Regex)
	})
	if err != nil {
		return nil, "", "", xerrors.Errorf("unable to extract from link %s, err: %w", link, err)
	}
	nextLink := ""
	if len(res) != 0 {
		u, err := url.Parse(c.commonGoparseSource.URL)
		if err != nil {
			log.Fatal(err)
		}
		nextLink = fmt.Sprintf("https://%s%s", u.Hostname(), res[0][0])
	}

	return result, nextLink, page, nil
}

func (c *Crawler) getPost(page, link string) ([]model.IDable, string, string, error) {
	doc, err := goquerywrapper.HTMLToDoc(page)
	if err != nil {
		return nil, "", "", xerrors.Errorf("unable to convert html page to doc, link: %s, err: %w", link, err)
	}
	res, err := goquerywrapper.Extract(c.logger, doc, c.commonGoparseSource.Content.Query, func(s *goquery.Selection) (string, error) {
		return goquerywrapper.DefaultSubtreeExtractor(c.logger, s, c.commonGoparseSource.Content.Attr, c.commonGoparseSource.Content.Regex)
	})
	if err != nil {
		return nil, "", "", xerrors.Errorf("unable to extract content from link %s, err: %w", link, err)
	}
	result := make([]model.IDable, 0)
	for _, el := range res {
		result = append(result, stContent{Link: link, Content: el[0]})
	}
	return result, "", "", nil
}

func (c *Crawler) ListItems(depth int, link string) ([]model.IDable, string, string, error) {
	if link == "" {
		link = c.commonGoparseSource.URL
	}
	page, err := c.urlGetter.Get(link)
	if err != nil {
		return nil, "", "", xerrors.Errorf("unable to extract from link %s, err: %w", link, err)
	}

	if depth == 1 {
		return c.listPage(page, link)
	} else {
		return c.getPost(page, link)
	}
}

//---

func NewCrawlerImpl(source model.Source, logger *logrus.Logger, htmlGetter URLGetter) (crawlers.CrawlerChain, error) {
	commonGoparseSource := CommonGoparseSource{}
	err := json.Unmarshal([]byte(source.CrawlerMeta), &commonGoparseSource)
	if err != nil {
		return nil, xerrors.Errorf("unable to unmarshal crawlerMetaStr, crawlerMeta: %s, err: %w", source.CrawlerMeta, err)
	}
	return &Crawler{
		source:              source,
		commonGoparseSource: commonGoparseSource,
		urlGetter:           htmlGetter,
		logger:              logger,
	}, nil
}

func NewCrawler(source model.Source, logger *logrus.Logger) (crawlers.CrawlerChain, error) {
	return NewCrawlerImpl(source, logger, &DefaultURLGetter{})
}

func init() {
	crawlers.RegisterChain(NewCrawler, CrawlerTypeCommonGoparseChain)
}
