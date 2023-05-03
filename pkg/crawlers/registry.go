package crawlers

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"personal-feed/pkg/model"
	"personal-feed/pkg/util"
)

type crawlerFactory func(model.Source, *logrus.Logger) (Crawler, error)

var crawlerIDToFactory = make(map[int]crawlerFactory)
var CrawlerIDToName = make(map[int]string)

func Register(foo crawlerFactory, crawlerID int) {
	if existingFactory, ok := crawlerIDToFactory[crawlerID]; ok {
		panic(fmt.Sprintf("this crawlerID (id=%d) is already registered. old: %s, new: %s", crawlerID, util.GetFuncName(existingFactory), util.GetFuncName(foo)))
	}
	crawlerIDToFactory[crawlerID] = foo
	CrawlerIDToName[crawlerID] = util.GetPackageNameOfFunc(foo)
}
