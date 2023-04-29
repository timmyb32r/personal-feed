package tree

import (
	"golang.org/x/xerrors"
	"personal-feed/pkg/model"
)

func isEqualTreeSchema(known, existing *Tree) error {
	lTypes := known.Layers()
	rTypes := existing.Layers()
	if len(lTypes) != len(rTypes) {
		return xerrors.Errorf("len(types) differs: %d vs %d", len(lTypes), len(rTypes))
	}
	for i := range lTypes {
		if !model.IsSameType(lTypes[i], rTypes[i]) {
			return xerrors.Errorf("types on index %d differs: %T vs %T", i, lTypes[i], rTypes[i])
		}
	}
	return nil
}

func BuildDiffTree(known, existing *Tree) (*Tree, map[string]model.IDable, error) {
	if known.ProjectName() != existing.ProjectName() {
		return nil, nil, xerrors.Errorf("unable to compare trees from different projects: %s vs %s", known.ProjectName(), existing.ProjectName())
	}
	err := isEqualTreeSchema(known, existing)
	if err != nil {
		return nil, nil, xerrors.Errorf("unable to compare trees with different schemas: %w", err)
	}

	knownInternalNodes := known.ExtractInternalNodes()
	knownDocs := known.ExtractDocs()

	existingInternalNodes := existing.ExtractInternalNodes()
	existingDocs := existing.ExtractDocs()

	newInternalNodes := make(map[string]*node)
	for k, v := range existingInternalNodes {
		if _, ok := knownInternalNodes[k]; !ok {
			newInternalNodes[k] = v
		}
	}

	newDocs := make(map[string]doc)
	for k, v := range existingDocs {
		if _, ok := knownDocs[k]; !ok {
			newDocs[k] = v
		}
	}

	newTree, err := NewTree(
		known.LayersTypes(),
	)
	if err != nil {
		return nil, nil, xerrors.Errorf("unable to create new Tree: %w", err)
	}

	newInternalNodesSerialized := make(map[string]model.IDable)
	for k, el := range newInternalNodes {
		_, err = newTree.insertNode(el)
		if err != nil {
			return nil, nil, xerrors.Errorf("unable to insert internal node: %w", err)
		}
		newInternalNodesSerialized[k] = el.Key()
	}

	for _, el := range newDocs {
		err = newTree.InsertDoc(el)
		if err != nil {
			return nil, nil, xerrors.Errorf("unable to insert doc: %w", err)
		}
	}

	return newTree, newInternalNodesSerialized, nil
}
