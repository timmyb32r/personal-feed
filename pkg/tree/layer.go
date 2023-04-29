package tree

import (
	"golang.org/x/xerrors"
	"personal-feed/pkg/model"
)

type layer struct {
	// types stuff
	masterObjCurrLayerKey model.IDable

	// layer info
	nextLayer *layer

	// flag isLeaf
	isLeaf bool
}

func (l *layer) checkKey(in interface{}) error {
	if !model.IsSameType(in, l.masterObjCurrLayerKey) {
		return xerrors.Errorf("key has wrong type: %T vs expected %T", in, l.masterObjCurrLayerKey)
	}
	return nil
}

func (l *layer) keyType() interface{} {
	return l.masterObjCurrLayerKey
}

func newLayer(masterObjCurrLayerKey model.IDable, masterObjNextLayerKey interface{}) *layer {
	var isLeaf bool
	var MasterObjNextLayer *layer
	switch t := masterObjNextLayerKey.(type) {
	case *layer:
		isLeaf = false
		MasterObjNextLayer = t
	default:
		isLeaf = true
		MasterObjNextLayer = nil
	}

	return &layer{
		masterObjCurrLayerKey: masterObjCurrLayerKey,
		nextLayer:             MasterObjNextLayer,
		isLeaf:                isLeaf,
	}
}
