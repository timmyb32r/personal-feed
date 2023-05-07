package crawlers

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"personal-feed/pkg/model"
	"personal-feed/pkg/util"
)

type crawlerTreeFactory func(model.Source, *logrus.Logger) (CrawlerTree, error)

var crawlerTreeIDToFactory = make(map[int]crawlerTreeFactory)
var CrawlerTreeIDToName = make(map[int]string)

func RegisterTree(foo crawlerTreeFactory, crawlerID int) {
	if existingFactory, ok := crawlerTreeIDToFactory[crawlerID]; ok {
		panic(fmt.Sprintf("this crawlerID (id=%d) is already registered. old: %s, new: %s", crawlerID, util.GetFuncName(existingFactory), util.GetFuncName(foo)))
	}
	crawlerTreeIDToFactory[crawlerID] = foo
	CrawlerTreeIDToName[crawlerID] = util.GetPackageNameOfFunc(foo)
}

//---

type crawlerChainFactory func(model.Source, *logrus.Logger) (CrawlerChain, error)

var crawlerChainIDToFactory = make(map[int]crawlerChainFactory)
var CrawlerChainIDToName = make(map[int]string)

func RegisterChain(foo crawlerChainFactory, crawlerID int) {
	if existingFactory, ok := crawlerChainIDToFactory[crawlerID]; ok {
		panic(fmt.Sprintf("this crawlerID (id=%d) is already registered. old: %s, new: %s", crawlerID, util.GetFuncName(existingFactory), util.GetFuncName(foo)))
	}
	crawlerChainIDToFactory[crawlerID] = foo
	CrawlerChainIDToName[crawlerID] = util.GetPackageNameOfFunc(foo)
}
