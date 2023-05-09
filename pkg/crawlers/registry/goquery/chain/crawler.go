package goquery

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"log"
	"net/url"
	"personal-feed/pkg/crawlers"
	"personal-feed/pkg/goquerywrapper"
	"personal-feed/pkg/goquerywrapper/extractors"
	"personal-feed/pkg/goquerywrapper/extractors/util"
	"personal-feed/pkg/model"
)

type stNt struct {
	Link         string
	HeaderText   string
	BusinessTime string
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

	itemHeaderExtractor       *extractors.GoQueryProgram
	itemLinkExtractor         *extractors.GoQueryProgram
	itemBusinessTimeExtractor *extractors.GoQueryProgram
	nextLinkExtractor         *extractors.GoQueryProgram
	contentExtractor          *extractors.GoQueryProgram
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
		result = append(result, stNt{HeaderText: el[0], Link: c.MakeLink(el[1]), BusinessTime: el[2]})
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
	content, err := goquerywrapper.ExtractByProgram(doc, c.contentExtractor)
	if err != nil {
		return nil, "", "", xerrors.Errorf("unable to extract content from link %s, err: %w", link, err)
	}
	result := []model.IDable{
		stContent{Link: link, Content: content},
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
