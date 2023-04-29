package model

type IDable interface {
	ID() string
}

//-----------------------------------

type RootKey struct {
}

func (k *RootKey) ID() string {
	return "ROOT"
}

func NewRootKey() *RootKey {
	return &RootKey{}
}

//-----------------------------------

type Node interface {
	Depth() int
	IsLeaf() bool
	ID() string
	ComplexKey() *ComplexKey
	Key() IDable
	ChildrenKeysLen() int
	CreateOrGetChildNode(nextKey IDable) (interface{}, error)
}

//-----------------------------------
