package registry

import (
	"fmt"
	"golang.org/x/exp/maps"
	"personal-feed/pkg/crawlers"
	"sort"
	"testing"
)

func TestPrintCrawlersList(t *testing.T) {
	ids := maps.Keys(crawlers.CrawlerTreeIDToName)
	sort.Ints(ids)
	for _, id := range ids {
		fmt.Printf("%d - %s\n", id, crawlers.CrawlerTreeIDToName[id])
	}
}
