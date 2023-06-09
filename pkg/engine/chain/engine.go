package chain

import (
	"context"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"personal-feed/pkg/crawlers"
	"personal-feed/pkg/model"
	"personal-feed/pkg/operation"
	"personal-feed/pkg/repo"
	"personal-feed/pkg/tree"
	"time"
)

type Engine struct {
	source             model.Source
	numMatchedNotifier model.NumMatchedNotifier
	crawler            crawlers.CrawlerChain
	db                 repo.Repo
	logger             *logrus.Logger
}

func (e *Engine) diff(knownTree *tree.Tree, existingNodes []model.IDable) ([]model.DBTreeNode, []model.IDable) {
	knownNodesMap := knownTree.ExtractFullKeysSet() // fullKey->Node
	rawResultNodes := make([]model.IDable, 0)
	for _, currNode := range existingNodes {
		currNodeFullKey := model.NewComplexKey("ROOT").MakeSubkey(currNode.ID()).FullKey()
		if _, ok := knownNodesMap[currNodeFullKey]; !ok {
			rawResultNodes = append(rawResultNodes, currNode)
		}
	}
	rawSerializedNodes := make([]model.DBTreeNode, 0)
	for _, el := range rawResultNodes {
		currKey, _ := model.ParseComplexKey("ROOT")
		rawSerializedNodes = append(rawSerializedNodes, *tree.SerializeKey(e.source.ID, currKey, el))
	}
	return rawSerializedNodes, rawResultNodes
}

func (e *Engine) RunOnce(ctx context.Context, op operation.Operation) error {
	layers := e.crawler.Layers()
	if len(layers) != 2 {
		return xerrors.Errorf("wrong configuration of chain crawler - there should be only two layer")
	}

	nextLink := e.source.HistoryState

	for {
		// TODO - extract every iteration is not effective - will optimize it then
		knownNodes, err := e.db.ExtractTreeNodes(ctx, e.source.ID)
		if err != nil {
			return xerrors.Errorf("unable to extract tree nodes: %w", err)
		}
		knownTree, err := tree.Deserialize(knownNodes, layers)
		if err != nil {
			return xerrors.Errorf("unable to deserialize tree: %w", err)
		}

		currLink := nextLink
		e.logger.Infof("handling iteration of sourceID:%d, currLink:%s", e.source.ID, currLink)

		var items []model.IDable
		var body string
		items, nextLink, body, err = e.crawler.ListItems(1, currLink)
		if err != nil {
			return xerrors.Errorf("unable to list items, err: %w", err)
		}
		e.numMatchedNotifier(e.source.ToJSON(), e.source.NumShouldBeMatched, len(items))
		newRawSerializedNodes, rawResultNodes := e.diff(knownTree, items)

		e.logger.Infof("extracted %d elements", len(newRawSerializedNodes))

		if len(newRawSerializedNodes) != 0 {
			for _, newItem := range newRawSerializedNodes {
				e.logger.Infof("    new el: %s", newItem.CurrentNodeJSON)
			}
			for i, newItem := range newRawSerializedNodes {
				e.logger.Infof("    start handling el: %s", newItem.CurrentNodeJSON)

				currDepth1Node := rawResultNodes[i]
				currID := currDepth1Node.ID()
				docs, _, _, err := e.crawler.ListItems(2, currID)
				if err != nil {
					return xerrors.Errorf("unable to get content, err: %w", err)
				}
				if len(docs) != 1 {
					return xerrors.Errorf("len(docs) != 1, err: %w", err)
				}
				currDoc := docs[0]
				if businessTimeable, ok := currDoc.(model.BusinessTimeable); ok {
					if businessTimeable1, ok := currDepth1Node.(model.BusinessTimeable); ok {
						businessTimeable.SetBusinessTime(businessTimeable1.GetBusinessTime())
					}
				}
				fullKey := model.NewComplexKey("ROOT").MakeSubkey(currID).MakeSubkey("doc").FullKey()
				dbDocNode := tree.SerializeDoc(e.source.ID, fullKey, currDoc)

				err = e.db.InsertNewTreeNodes(ctx, e.source.ID, []model.DBTreeNode{newItem, *dbDocNode})
				if err != nil {
					return xerrors.Errorf("unable to insert items, err: %w", err)
				}
				e.logger.Infof("    finish handling el: %s", newItem.CurrentNodeJSON)

				time.Sleep(2 * time.Second) // to not to ddos
			}
		}
		// if there are something new OR number of new items is not expected OR it's 'load-history' operation
		// in other words, if regular update found all known items and NumShouldBeMatched is expected - then we are now saving it
		if (len(newRawSerializedNodes) != 0 || (e.source.NumShouldBeMatched != nil && *e.source.NumShouldBeMatched != len(newRawSerializedNodes))) || op.OperationType == operation.OpTypeLoadHistory {
			e.logger.Infof("will save iteration: %s", currLink)
			err = e.db.InsertSourceIteration(ctx, e.source.ID, currLink, body)
			if err != nil {
				return xerrors.Errorf("unable to InsertSourceIteration, err: %w", err)
			}
		}

		countOfKnownItems := len(items) - len(newRawSerializedNodes)
		if op.OperationType == operation.OpTypeRegularUpdate && countOfKnownItems != 0 {
			e.logger.Infof("regular_update found at least one known element, so it's traversal finished")
			return nil
		}

		e.logger.Infof("next link: %s", nextLink)

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

func NewEngine(source *model.Source, numMatchedNotifier model.NumMatchedNotifier, crawler crawlers.CrawlerChain, db repo.Repo, logger *logrus.Logger) *Engine {
	return &Engine{
		source:             *source,
		numMatchedNotifier: numMatchedNotifier,
		crawler:            crawler,
		db:                 db,
		logger:             logger,
	}
}
