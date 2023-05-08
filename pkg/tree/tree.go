package tree

import (
	"golang.org/x/xerrors"
	"personal-feed/pkg/model"
)

type Tree struct {
	rootKey     model.IDable
	layersTypes []model.IDable // without Root one
	layers      []*layer
	root        *node
}

func (t *Tree) ProjectName() string {
	return t.rootKey.ID()
}

func (t *Tree) LayersTypes() []model.IDable { // without Root one
	return t.layersTypes
}

func (t *Tree) Layers() []*layer {
	return t.layers
}

func (t *Tree) Root() interface{} {
	return t.root
}

func (t *Tree) insertNode(inNode *node) (*node, error) {
	depth := inNode.Depth()
	pathNodes := make([]*node, depth)
	currNode := inNode
	for i := 0; i < depth; i++ {
		pathNodes[depth-1-i] = currNode
		currNode = inNode.parentNode
	}
	currNode = t.root
	for i := 0; i < depth; i++ {
		newNode, err := currNode.CreateOrGetChildNode(pathNodes[i].Key())
		if err != nil {
			return nil, err
		}
		currNode = newNode.(*node)
	}
	return currNode, nil
}

func (t *Tree) InsertDoc(inDoc doc) error {
	parentNode, err := t.insertNode(inDoc.parentNode)
	if err != nil {
		return xerrors.Errorf("unable to insert node: %w", err)
	}
	_, err = parentNode.CreateOrGetChildNode(inDoc.key)
	if err != nil {
		return xerrors.Errorf("unable to insert inDoc: %w", err)
	}
	return nil
}

func (t *Tree) ExtractInternalNodes() map[string]*node {
	return extractInternalNodes(t.root)
}

func (t *Tree) ExtractDocs() map[string]doc {
	return extractDocs(t.root)
}

func (t *Tree) ExtractDocsUnwrapped() map[string]model.IDable {
	return extractDocsUnwrapped(t.root)
}

func (t *Tree) Serialize(sourceID int) []model.DBTreeNode {
	return serialize(sourceID, t)
}

func (t *Tree) SerializeKey(sourceID int, fullKey string, key model.IDable) *model.DBTreeNode {
	currComplexKey, _ := model.ParseComplexKey(fullKey)
	return SerializeKey(sourceID, currComplexKey.ParentKey(), key)
}

func NewTree(layersTypes []model.IDable) (*Tree, error) {
	realLayersTypes := append([]model.IDable{model.NewRootKey()}, layersTypes...)
	if len(realLayersTypes) == 0 || len(realLayersTypes) == 1 {
		return nil, xerrors.Errorf("unable to create tree, when levels<=1")
	}
	if !model.IsSameType(realLayersTypes[0], model.NewRootKey()) {
		return nil, xerrors.Errorf("root should be type of model.RootKey, got: %T", realLayersTypes[0])
	}

	layers := make([]*layer, len(realLayersTypes)-1)
	lastLayer := newLayer(realLayersTypes[len(realLayersTypes)-2], realLayersTypes[len(realLayersTypes)-1])
	layers[len(realLayersTypes)-2] = lastLayer

	if len(realLayersTypes) > 2 {
		i := len(realLayersTypes) - 3
		for {
			currLevel := newLayer(realLayersTypes[i], layers[i+1])
			layers[i] = currLevel

			if i == 0 {
				break
			}
			i--
		}
	}

	rootKey := model.NewRootKey()
	root, err := newNode(rootKey, layers[0], nil)
	if err != nil {
		return nil, xerrors.Errorf("unable to create node: %w", err)
	}

	return &Tree{
		rootKey:     rootKey,
		layersTypes: realLayersTypes[1:],
		layers:      layers,
		root:        root,
	}, nil
}
