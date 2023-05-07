package tree

import (
	"context"
	"golang.org/x/xerrors"
	"personal-feed/pkg/crawlers"
	"personal-feed/pkg/model"
	"personal-feed/pkg/repo"
	"personal-feed/pkg/tree"
	"personal-feed/pkg/util"
)

type Engine struct {
	source             model.Source
	numMatchedNotifier model.NumMatchedNotifier
	crawler            crawlers.CrawlerTree
	db                 repo.Repo
}

func (e *Engine) RunOnce(ctx context.Context) error {
	rollbacks := util.Rollbacks{}
	defer rollbacks.Do()

	tx, err := e.db.NewTx()
	if err != nil {
		return xerrors.Errorf("unable to begin new transaction: %w", err)
	}

	rollbacks.Add(func() { _ = tx.Rollback(ctx) })

	knownNodes, err := e.db.ExtractTreeNodesTx(tx, e.source.ID)
	if err != nil {
		return xerrors.Errorf("unable to extract tree nodes: %w", err)
	}

	layers := e.crawler.Layers()
	maxDepth := len(layers)

	knownTree, err := tree.Deserialize(knownNodes, layers)
	if err != nil {
		return xerrors.Errorf("unable to deserialize tree: %w", err)
	}

	existingTree, err := tree.NewTree(layers)
	if err != nil {
		return xerrors.Errorf("unable to create new tree: %w", err)
	}

	nextLayerNodes := make([]interface{}, 0)
	nextLayerNodes = append(nextLayerNodes, existingTree.Root())
	for i := 0; i < maxDepth; i++ {
		currLayerNodes := nextLayerNodes
		nextLayerNodes = make([]interface{}, 0)
		for j := 0; j < len(currLayerNodes); j++ {
			currNode := currLayerNodes[j].(model.Node)
			children, err := e.crawler.ListLayer(i+1, currNode) // i+1, bcs depth=0 is Root everywhere
			if err != nil {
				return xerrors.Errorf("unable to list node: %s, err: %w", currNode.ComplexKey().FullKey(), err)
			}
			e.numMatchedNotifier(e.source.ToJSON(), e.source.NumShouldBeMatched, len(children))
			for _, child := range children {
				if !model.IsSameType(child, layers[i]) {
					return xerrors.Errorf("parser returned wrong type, on node: %s, %T vs expected %T", currNode.ComplexKey().FullKey(), child, layers[i+1])
				}
				newNode, err := currNode.CreateOrGetChildNode(child)
				if err != nil {
					return xerrors.Errorf("unable to add child node on node: %s, err: %w", currNode.ComplexKey().FullKey(), err)
				}
				nextLayerNodes = append(nextLayerNodes, newNode)
			}
		}
	}

	newDBInternalNodes, newDBDocs, err := tree.BuildDiffTreeAndSerialize(e.source.ID, knownTree, existingTree)
	if err != nil {
		return xerrors.Errorf("unable to build diff tree, err: %w", err)
	}

	err = e.db.InsertNewTreeNodesTx(tx, e.source.ID, append(newDBInternalNodes, newDBDocs...))
	if err != nil {
		return xerrors.Errorf("unable to inset new nodes, err: %w", err)
	}

	tx.Commit(context.Background())
	rollbacks.Cancel()
	return nil
}

func NewEngine(source *model.Source, numMatchedNotifier model.NumMatchedNotifier, crawler crawlers.CrawlerTree, db repo.Repo) *Engine {
	return &Engine{
		source:             *source,
		numMatchedNotifier: numMatchedNotifier,
		crawler:            crawler,
		db:                 db,
	}
}
