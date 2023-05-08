package chain

import (
	"context"
	"golang.org/x/xerrors"
	"personal-feed/pkg/crawlers"
	"personal-feed/pkg/model"
	"personal-feed/pkg/repo"
	"personal-feed/pkg/tree"
	"time"
)

type Engine struct {
	source             model.Source
	numMatchedNotifier model.NumMatchedNotifier
	crawler            crawlers.CrawlerChain
	db                 repo.Repo
}

func (e *Engine) diff(knownTree *tree.Tree, existingNodes []model.IDable) []model.DBTreeNode {
	knownNodesMap := knownTree.ExtractInternalNodes()
	resultNodes := make([]model.IDable, 0)
	for _, currNode := range existingNodes {
		if _, ok := knownNodesMap[currNode.ID()]; !ok {
			resultNodes = append(resultNodes, currNode)
		}
	}
	result := make([]model.DBTreeNode, 0)
	for _, el := range resultNodes {
		currKey, _ := model.ParseComplexKey("ROOT")
		result = append(result, *tree.SerializeKey(e.source.ID, currKey, el))
	}
	return result
}

func (e *Engine) RunOnce(ctx context.Context) error {
	knownNodes, err := e.db.ExtractTreeNodes(ctx, e.source.ID)
	if err != nil {
		return xerrors.Errorf("unable to extract tree nodes: %w", err)
	}

	layers := e.crawler.Layers()
	if len(layers) != 1 {
		return xerrors.Errorf("wrong configuration of chain crawler - there should be only one layer")
	}

	knownTree, err := tree.Deserialize(knownNodes, layers)
	if err != nil {
		return xerrors.Errorf("unable to deserialize tree: %w", err)
	}

	nextLink := e.source.HistoryState

	for {
		var items []model.IDable
		var body string
		items, nextLink, body, err = e.crawler.ListItems(nextLink)
		if err != nil {
			return xerrors.Errorf("unable to list items, err: %w", err)
		}
		newItems := e.diff(knownTree, items)
		if len(newItems) != 0 {
			err = e.db.InsertNewTreeNodes(ctx, e.source.ID, newItems)
			if err != nil {
				return xerrors.Errorf("unable to insert items, err: %w", err)
			}
		}

		// if there are something new OR number of items is not expected!
		if len(newItems) != 0 || (e.source.NumShouldBeMatched != nil && *e.source.NumShouldBeMatched != len(newItems)) {
			err = e.db.InsertSourceIteration(ctx, e.source.ID, body)
			if err != nil {
				return xerrors.Errorf("unable to InsertSourceIteration, err: %w", err)
			}
		}
		err = e.db.SetState(ctx, e.source.ID, nextLink)
		if err != nil {
			return xerrors.Errorf("unable to set state, err: %w", err)
		}
		if nextLink == "" {
			// TODO - log there is the end of history
			break
		}
		time.Sleep(2 * time.Second) // to not to ddos
	}

	return nil
}

func NewEngine(source *model.Source, numMatchedNotifier model.NumMatchedNotifier, crawler crawlers.CrawlerChain, db repo.Repo) *Engine {
	return &Engine{
		source:             *source,
		numMatchedNotifier: numMatchedNotifier,
		crawler:            crawler,
		db:                 db,
	}
}
