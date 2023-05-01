package engine

import (
	"context"
	"golang.org/x/xerrors"
	"personal-feed/pkg/crawlers"
	"personal-feed/pkg/model"
	"personal-feed/pkg/repo"
	"personal-feed/pkg/tree"
)

type Engine struct {
	source  model.Source
	crawler crawlers.Crawler
	db      repo.Repo
}

func (e *Engine) RunOnce() error {
	tx, err := e.db.NewTx()
	if err != nil {
		return xerrors.Errorf("unable to begin new transaction: %w", err)
	}

	knownNodes, err := e.db.ExtractTreeNodes(tx, e.source.ID)
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

	err = e.db.InsertNewTreeNodes(tx, e.source.ID, append(newDBInternalNodes, newDBDocs...))
	if err != nil {
		return xerrors.Errorf("unable to inset new nodes, err: %w", err)
	}

	// TODO - insert writing to 'feed' table

	return tx.Commit(context.Background())
}

func NewEngine(source *model.Source, crawler crawlers.Crawler, db repo.Repo) *Engine {
	return &Engine{
		source:  *source,
		crawler: crawler,
		db:      db,
	}
}
