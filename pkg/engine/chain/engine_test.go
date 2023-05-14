package chain

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	chaingoquery "personal-feed/pkg/crawlers/registry/goquery/chain"
	"personal-feed/pkg/goquerywrapper/extractors"
	"personal-feed/pkg/model"
	"personal-feed/pkg/operation"
	"personal-feed/pkg/repo/registry/in_memory"
	"testing"
)

//---

type MockedHTMLGetter struct {
	index int
}

var pages = []string{`
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
												<div class="byline" style="line-height: 1;"> <em> May 2, 2023 by </em> <em> name surname </em>
												</div>
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
`, `
<!DOCTYPE html>
<html>
<body class="post">
    <div id="rhbar"> <a class="jbdevlogo" href="https://blablabla"></a> <a class="rhlogo" href="https://blablabla/"></a> </div>
    <div id>
        <div class="container" id="content">
            <div class="row post-text-padding row-no-expand">
                <div class="col-md-9">
                    <div class="post">
                        <div class="row" style="margin-left: 0; margin-right: 0; margin-bottom: 10px">
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</body>
</html>
`}

func (g *MockedHTMLGetter) Get(_ string) (string, error) {
	currIndex := g.index
	g.index++
	if currIndex >= len(pages) {
		return "", nil
	}
	return pages[currIndex], nil
}

//---

func TestChain(t *testing.T) {
	sourceCrawlerMeta := chaingoquery.CommonGoparseSource{
		URL: "https://test-blog.io/blog/",
		Item: chaingoquery.Item{
			Query: ".blog-list-item",
			Header: extractors.Program{
				extractors.Instruction{Query: "span.hidden-lg"},
				extractors.Instruction{Query: "a[href]"},
				extractors.Instruction{Text: "!"},
			},
			Link: extractors.Program{
				extractors.Instruction{Query: "span.hidden-lg"},
				extractors.Instruction{Query: "a[href]"},
				extractors.Instruction{Attr: "href"},
			},
			BusinessTime: extractors.Program{
				extractors.Instruction{Query: "div.byline"},
				extractors.Instruction{Text: "!"},
				extractors.Instruction{Regex: `\s*(.*?) by`},
			},
		},
		NextLink: extractors.Program{
			extractors.Instruction{Query: "a.previous"},
			extractors.Instruction{Attr: "href"},
		},
		Content: extractors.Program{
			extractors.Instruction{Query: "div.post"},
		},
	}
	sourceCrawlerMetaArr, _ := json.Marshal(sourceCrawlerMeta)

	fmt.Println(string(sourceCrawlerMetaArr))

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

	ctx := context.TODO()

	inMemoryRepoWrapped, _ := in_memory.NewRepo(ctx, struct{}{}, nil)
	inMemoryRepo := inMemoryRepoWrapped.(*in_memory.Repo)

	op := operation.Operation{
		OperationType: operation.OpTypeRegularUpdate,
	}

	engine := NewEngine(source, stubNotifier, crawlerImpl, inMemoryRepo, logrus.New())
	err = engine.RunOnce(ctx, op)
	require.NoError(t, err)
	require.Equal(t, 2, inMemoryRepo.Len())
}
