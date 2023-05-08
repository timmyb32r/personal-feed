package chain

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	chaingoquery "personal-feed/pkg/crawlers/registry/goquery/chain"
	"personal-feed/pkg/model"
	"personal-feed/pkg/repo/registry/in_memory"
	"testing"
)

//---

type MockedHTMLGetter struct {
	index int
}

func (g *MockedHTMLGetter) Get(_ string) (string, error) {
	if g.index > 0 {
		return "", nil
	}
	g.index++
	return `
<!DOCTYPE html>
<html>
<head>
    <title>QQQ Blog</title>
    <meta charset="utf-8">
</head>
<body class="blog">
    <div id>
        <div class="container" id="content">
            <div class="row post-text-padding row-no-expand">
                <div class="col-md-9">
                    <h1 class="section full">QQQ Blog</h1>
                    <div class="component-wrapper">
                        <div class="grid__item width-10-12 width-12-12-m">
                            <div class="blog-list-item grid-wrapper">
                                <div class="row" style="margin-left: 0; margin-right: 0; margin-bottom: 10px;">
                                    <div class="col-sm-12" style="padding-left: 0px;">
                                        <div style="display: table-cell; vertical-align: top;">
                                            <div style="width: 72px; border: 1px solid #ccc; padding: 3px; display: inline-block;"> <img src="/assets/images/author.jpg" style="width: 64px;"> </div>
                                        </div>
                                        <div style="display: table-cell; vertical-align: top;">
                                            <div style="margin-left: 8px;">
                                            <span class="hidden-sm hidden-xs" style="font-size: 2.75rem; line-height: 1;"> 
                                                <a href="/blog/2023/05/02/blablabla/">blablabla-blablabla</a> 
                                            </span>
                                            <span class="hidden-md hidden-lg" style="font-size: 2rem; line-height: 1;">
                                                <a href="/blog/2023/05/02/blablabla/">blablabla-blablabla</a> 
                                            </span>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                    <ul class="pager">
                        <li class="previous"> <a href="/blog/page/2/" class="previous">Older</a> </li>
                        <li class="pages "> Page: 1 of 41 </li>
                        <li class="next disabled"><a href="">Newer</a></li>
                    </ul>
                </div>
            </div>
        </div>
        <footer class="container">
            <div>
            </div>
        </footer>
    </div>
    <script async src="https://www.googletagmanager.com/gtag/js?id=UA-76464546-1"></script>
    <script>
        function gtag() {
            dataLayer.push(arguments)
        }
        window.dataLayer = window.dataLayer || [], gtag("js", new Date), gtag("config", "UA-76464546-1");
    </script>
</body>
</html>
`, nil
}

//---

func TestChain(t *testing.T) {
	sourceCrawlerMeta := chaingoquery.CommonGoparseSource{
		URL: "https://test-blog.io/blog/",
		Item: chaingoquery.CommonGoparseSourceItem{
			Query: ".blog-list-item",
			Header: chaingoquery.QueryIntoSelected{
				Attr:  "",
				Regex: `.*?<a href=[^>]+>(.*?)</a>.*`,
			},
			Link: chaingoquery.QueryIntoSelected{
				Attr:  "",
				Regex: `.*?<a href="([^"]+)".*`,
			},
		},
		Next: chaingoquery.QueryIntoDoc{
			Query: ".previous",
			Attr:  "",
			Regex: `.*href=\"([^"]+)\".*`,
		},
	}
	sourceCrawlerMetaArr, _ := json.Marshal(sourceCrawlerMeta)

	source := &model.Source{
		ID:          1,
		Description: "blablabla",
		CrawlerID:   1,
		CrawlerMeta: string(sourceCrawlerMetaArr),
		Schedule:    "",
	}
	var log = logrus.New()

	crawlerImpl, err := chaingoquery.NewCrawlerImpl(*source, log, &MockedHTMLGetter{})
	require.NoError(t, err)

	stubNotifier := func(crawlerDescr string, expected *int, real int) {}

	inMemoryRepoWrapped, _ := in_memory.NewRepo(struct{}{}, nil)
	inMemoryRepo := inMemoryRepoWrapped.(*in_memory.Repo)

	ctx := context.TODO()

	engine := NewEngine(source, stubNotifier, crawlerImpl, inMemoryRepo)
	err = engine.RunOnce(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, inMemoryRepo.Len())
}
