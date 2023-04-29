package tree

import (
	"golang.org/x/xerrors"
	"personal-feed/pkg/model"
)

func BuildDiffTreeAndSerialize(sourceID int, known, existing *Tree) ([]model.DBTreeNode, []model.DBTreeNode, error) {
	diffTree, newInternalNodes, err := BuildDiffTree(known, existing)
	if err != nil {
		return nil, nil, xerrors.Errorf("unable to build diff tree, err: %w", err)
	}

	dbNewInternalNodes := make([]model.DBTreeNode, 0, len(newInternalNodes))
	for fullKey, key := range newInternalNodes {
		complexKey, _ := model.ParseComplexKey(fullKey)
		dbNewInternalNodes = append(dbNewInternalNodes, *serializeKey(sourceID, complexKey.ParentKey(), key))
	}

	newDocs := diffTree.ExtractDocsUnwrapped()
	dbNewDocs := make([]model.DBTreeNode, 0, len(newDocs))
	for fullKey, key := range newDocs {
		dbNewDocs = append(dbNewDocs, *serializeDoc(sourceID, fullKey, key))
	}

	return dbNewInternalNodes, dbNewDocs, nil
}
