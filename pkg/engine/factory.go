package engine

import (
	"golang.org/x/xerrors"
	"personal-feed/pkg/crawlers"
	enginechain "personal-feed/pkg/engine/chain"
	enginetree "personal-feed/pkg/engine/tree"
	"personal-feed/pkg/model"
	"personal-feed/pkg/repo"
)

func NewEngine(source *model.Source, numMatchedNotifier model.NumMatchedNotifier, crawler interface{}, db repo.Repo) (AbstractEngine, error) {
	if _, ok := crawlers.CrawlerTreeIDToName[source.CrawlerID]; ok {
		return enginetree.NewEngine(source, numMatchedNotifier, crawler.(crawlers.CrawlerTree), db), nil
	}

	if _, ok := crawlers.CrawlerChainIDToName[source.CrawlerID]; ok {
		return enginechain.NewEngine(source, numMatchedNotifier, crawler.(crawlers.CrawlerChain), db), nil
	}

	return nil, xerrors.Errorf("unable to find crawlerID: %d", source.ID)
}
