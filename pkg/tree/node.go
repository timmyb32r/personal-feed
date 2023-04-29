package tree

import (
	"golang.org/x/xerrors"
	"personal-feed/pkg/model"
)

type node struct {
	// current node key
	key        model.IDable
	complexKey *model.ComplexKey

	// layer info
	layerDescr *layer

	// parent
	parentNode *node

	// links to children - it either *node (if not leaf node) or some value (if leaf node)
	keyIDToChildNode map[string]interface{}
}

func (n *node) Key() model.IDable {
	return n.key
}

func (n *node) ComplexKey() *model.ComplexKey {
	return n.complexKey
}

func (n *node) ID() string {
	return n.key.ID()
}

func (n *node) IsLeaf() bool {
	return n.layerDescr.isLeaf
}

func (n *node) Depth() int {
	return n.complexKey.Depth()
}

//---------------------------------------------------------------------------------------------------------------------
// create child nodes

func (n *node) CreateOrGetChildNode(nextKey model.IDable) (interface{}, error) {
	nextID := nextKey.ID()
	if n.layerDescr.isLeaf {
		n.keyIDToChildNode[nextID] = nextKey
	} else {
		if err := n.layerDescr.nextLayer.checkKey(nextKey); err != nil {
			return nil, err
		}
		if _, ok := n.keyIDToChildNode[nextID]; !ok {
			node, err := newNode(nextKey, n.layerDescr.nextLayer, n)
			if err != nil {
				return nil, err
			}
			n.keyIDToChildNode[nextID] = node
		}
	}
	return n.keyIDToChildNode[nextID], nil
}

func (n *node) CreateOrGetChildNodes(nextKeys []model.IDable) (interface{}, error) {
	var currNode interface{} = n
	for _, currKey := range nextKeys {
		childNode, err := currNode.(*node).CreateOrGetChildNode(currKey)
		if err != nil {
			return nil, xerrors.Errorf("unable to create child node by ID: %s", currKey.ID())
		}
		currNode = childNode
	}
	return currNode, nil
}

//---------------------------------------------------------------------------------------------------------------------
// navigation between layers (by concrete values)

func (n *node) GetChildNodeByKeyID(id string) (interface{}, error) {
	if result, ok := n.keyIDToChildNode[id]; ok {
		return result, nil
	} else {
		return nil, xerrors.Errorf("unable to find id on current level: %s", id)
	}
}

func (n *node) GetChildNodeByKey(curr interface{}) (interface{}, error) {
	if err := n.layerDescr.nextLayer.checkKey(curr); err != nil {
		return nil, err
	}
	id := curr.(model.IDable).ID()
	return n.GetChildNodeByKeyID(id)
}

func (n *node) GetChildNodeByComplexKey(complexKey *model.ComplexKey) (*node, error) {
	currNode := n
	keys := complexKey.Keys()
	for _, el := range keys {
		childNode, err := currNode.GetChildNodeByKeyID(el)
		if err != nil {
			return nil, xerrors.Errorf("unable to get child node by id: %s", el)
		}
		currNode = childNode.(*node)
	}
	return currNode, nil
}

//---------------------------------------------------------------------------------------------------------------------
// enumerating over node

func (n *node) ChildrenKeys() []model.IDable {
	result := make([]model.IDable, 0, len(n.keyIDToChildNode))
	if n.IsLeaf() {
		for _, v := range n.keyIDToChildNode {
			result = append(result, v.(model.IDable))
		}
	} else {
		for _, v := range n.keyIDToChildNode {
			result = append(result, v.(*node).Key())
		}
	}
	return result
}

func (n *node) ChildrenKeysLen() int {
	return len(n.keyIDToChildNode)
}

//---------------------------------------------------------------------------------------------------------------------
// ctor

func newNode(currNodeKey model.IDable, layer *layer, parentNode *node) (*node, error) {
	if err := layer.checkKey(currNodeKey); err != nil {
		return nil, err
	}
	var currComplexKey *model.ComplexKey
	if parentNode == nil {
		currComplexKey = model.NewComplexKey(currNodeKey.ID())
	} else {
		currComplexKey = parentNode.ComplexKey().MakeSubkey(currNodeKey.ID())
	}
	return &node{
		key:              currNodeKey,
		complexKey:       currComplexKey,
		layerDescr:       layer,
		parentNode:       parentNode,
		keyIDToChildNode: make(map[string]interface{}),
	}, nil
}
