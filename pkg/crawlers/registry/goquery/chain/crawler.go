package goquery

import (
	"encoding/json"
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"log"
	"net/url"
	"personal-feed/pkg/crawlers"
	"personal-feed/pkg/goquerywrapper"
	"personal-feed/pkg/goquerywrapper/extractors"
	"personal-feed/pkg/goquerywrapper/extractors/util"
	"personal-feed/pkg/model"
	"time"
)

type stNt struct {
	Link         string
	HeaderText   string
	BusinessTime *time.Time
}

func (n stNt) ID() string {
	return n.Link
}

func (n stNt) GetBusinessTime() *time.Time {
	return n.BusinessTime
}

type StContent struct {
	Link         string
	BusinessTime *time.Time
	Author       string
	HeaderText   string
	Content      string
}

func (n StContent) ID() string {
	return n.Link
}

func (n StContent) GetBusinessTime() *time.Time {
	return n.BusinessTime
}

func (n StContent) SetBusinessTime(in *time.Time) {
	inCopy := *in
	n.BusinessTime = &inCopy
}

type Crawler struct {
	source              model.Source
	commonGoparseSource CommonGoparseSource
	urlGetter           URLGetter
	logger              *logrus.Logger

	itemHeaderExtractor       *extractors.GoQueryProgram
	itemLinkExtractor         *extractors.GoQueryProgram
	itemBusinessTimeExtractor *extractors.GoQueryProgram
	itemAuthorExtractor       *extractors.GoQueryProgram
	nextLinkExtractor         *extractors.GoQueryProgram
	contentExtractor          *extractors.GoQueryProgram
}

func (c *Crawler) CrawlerType() int {
	return CrawlerTypeCommonGoparseChain
}

func (c *Crawler) Layers() []model.IDable {
	return []model.IDable{
		stNt{Link: "", HeaderText: ""},
		StContent{Link: "", Content: ""},
	}
}

func (c *Crawler) MakeLink(URL string) string {
	u, err := url.Parse(c.commonGoparseSource.URL)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("https://%s%s", u.Hostname(), URL)
}

func (c *Crawler) listPage(page, link string) ([]model.IDable, string, string, error) {
	doc, err := util.HTMLToDoc(page)
	if err != nil {
		return nil, "", "", xerrors.Errorf("unable to convert html page to doc, link: %s, err: %w", link, err)
	}
	res, err := goquerywrapper.ExtractItemsByProgram(c.logger, doc, c.commonGoparseSource.Item.Query, c.itemHeaderExtractor, c.itemLinkExtractor, c.itemBusinessTimeExtractor)
	if err != nil {
		return nil, "", "", xerrors.Errorf("unable to extract from link %s, err: %w", link, err)
	}
	result := make([]model.IDable, 0)
	for _, el := range res {
		var businessTime *time.Time = nil
		if c.commonGoparseSource.Item.BusinessTime != nil {
			businessTimeVal, err := dateparse.ParseAny(el[2])
			if err != nil {
				return nil, "", "", xerrors.Errorf("unable to parse business time, str: %s, err: %s", el[2], err)
			}
			businessTime = &businessTimeVal
		}
		result = append(
			result,
			stNt{
				HeaderText:   el[0],
				Link:         c.MakeLink(el[1]),
				BusinessTime: businessTime,
			},
		)
	}

	// next_link

	rawNextLink, err := goquerywrapper.ExtractByProgram(doc, c.nextLinkExtractor)
	if err != nil {
		//return nil, "", "", xerrors.Errorf("unable to extract from link %s, err: %w", link, err)
		c.logger.Infof("look like i've found the last page")
	}
	nextLink := ""
	if rawNextLink != "" {
		nextLink = c.MakeLink(rawNextLink)
	}

	return result, nextLink, page, nil
}

func (c *Crawler) getPost(page, link string) ([]model.IDable, string, string, error) {
	doc, err := util.HTMLToDoc(page)
	if err != nil {
		return nil, "", "", xerrors.Errorf("unable to convert html page to doc, link: %s, err: %w", link, err)
	}
	res, err := goquerywrapper.ExtractItemsByProgram(c.logger, doc, ":root", c.itemHeaderExtractor, c.contentExtractor, c.itemBusinessTimeExtractor, c.itemAuthorExtractor)
	if err != nil {
		return nil, "", "", xerrors.Errorf("unable to extract content from link %s, err: %w", link, err)
	}
	if len(res) != 1 {
		return nil, "", "", xerrors.Errorf("len(res) != 1. unable to extract content/time/author from doc, link: %s, err: %w", link, err)
	}
	currRes := res[0]
	if len(currRes) != 4 {
		return nil, "", "", xerrors.Errorf("len(currRes) != 3. unable to extract content/time/author from doc, link: %s, err: %w", link, err)
	}

	businessTimeVal, err := dateparse.ParseAny(currRes[2])
	if err != nil {
		return nil, "", "", xerrors.Errorf("unable to parse business time, str: %s, err: %s", currRes[2], err)
	}
	result := []model.IDable{
		StContent{
			Link:         link,
			BusinessTime: &businessTimeVal,
			Author:       currRes[3],
			HeaderText:   currRes[0],
			Content:      currRes[1],
		},
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

	itemHeaderExtractor, err := extractors.NewProgramFromProgram(commonGoparseSource.Item.Header)
	if err != nil {
		return nil, xerrors.Errorf("unable to create program for item.header, err: %w", err)
	}
	itemLinkExtractor, err := extractors.NewProgramFromProgram(commonGoparseSource.Item.Link)
	if err != nil {
		return nil, xerrors.Errorf("unable to create program for item.link, err: %w", err)
	}
	itemBusinessTimeExtractor, err := extractors.NewProgramFromProgram(commonGoparseSource.Item.BusinessTime)
	if err != nil {
		return nil, xerrors.Errorf("unable to create program for item.business_time, err: %w", err)
	}
	itemAuthorExtractor, err := extractors.NewProgramFromProgram(commonGoparseSource.Item.Author)
	if err != nil {
		return nil, xerrors.Errorf("unable to create program for item.author, err: %w", err)
	}
	nextLinkExtractor, err := extractors.NewProgramFromProgram(commonGoparseSource.NextLink)
	if err != nil {
		return nil, xerrors.Errorf("unable to create program for item.next, err: %w", err)
	}
	contentExtractor, err := extractors.NewProgramFromProgram(commonGoparseSource.Content)
	if err != nil {
		return nil, xerrors.Errorf("unable to create program for item.content, err: %w", err)
	}

	return &Crawler{
		source:              source,
		commonGoparseSource: commonGoparseSource,
		urlGetter:           htmlGetter,
		logger:              logger,

		itemHeaderExtractor:       itemHeaderExtractor,
		itemLinkExtractor:         itemLinkExtractor,
		itemAuthorExtractor:       itemAuthorExtractor,
		itemBusinessTimeExtractor: itemBusinessTimeExtractor,
		nextLinkExtractor:         nextLinkExtractor,
		contentExtractor:          contentExtractor,
	}, nil
}

func NewCrawler(source model.Source, logger *logrus.Logger) (crawlers.CrawlerChain, error) {
	return NewCrawlerImpl(source, logger, &DefaultURLGetter{})
}

func init() {
	crawlers.RegisterChain(NewCrawler, CrawlerTypeCommonGoparseChain)
}
