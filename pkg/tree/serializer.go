package tree

import (
	"encoding/json"
	"golang.org/x/xerrors"
	"personal-feed/pkg/model"
	"reflect"
	"sort"
)

func SerializeKey(sourceID int, parentComplexKey *model.ComplexKey, key model.IDable) *model.DBTreeNode {
	JSONObj, _ := json.Marshal(key)
	return &model.DBTreeNode{
		SourceID:        sourceID,
		Depth:           parentComplexKey.Depth() + 1,
		ParentFullKey:   parentComplexKey.FullKey(),
		CurrentNodeJSON: string(JSONObj),
	}
}

func serializeInternalNode(sourceID int, el *node) *model.DBTreeNode {
	return SerializeKey(sourceID, el.ComplexKey().ParentKey(), el.Key())
}

func serializeDoc(sourceID int, fullKey string, key model.IDable) *model.DBTreeNode {
	complexKey, _ := model.ParseComplexKey(fullKey)
	return SerializeKey(sourceID, complexKey.ParentKey(), key)
}

func serialize(sourceID int, in *Tree) []model.DBTreeNode {
	internalNodes := in.ExtractInternalNodes()
	docs := in.ExtractDocsUnwrapped()

	result := make([]model.DBTreeNode, 0)
	for _, el := range internalNodes {
		result = append(result, *serializeInternalNode(sourceID, el))
	}
	for fullKey, el := range docs {
		result = append(result, *serializeDoc(sourceID, fullKey, el))
	}
	return result
}

func Deserialize(in []model.DBTreeNode, layersTypes []model.IDable) (*Tree, error) {
	tree, err := NewTree(layersTypes)
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}

	if len(in) == 0 {
		return tree, nil
	}

	depthToNodes := make(map[int][]model.DBTreeNode)
	for _, currNode := range in {
		if _, ok := depthToNodes[currNode.Depth]; !ok {
			depthToNodes[currNode.Depth] = make([]model.DBTreeNode, 0)
		}
		depthToNodes[currNode.Depth] = append(depthToNodes[currNode.Depth], currNode)
	}

	realLayersTypes := append([]model.IDable{model.NewRootKey()}, layersTypes...)

	// validate

	if len(depthToNodes) > len(realLayersTypes) {
		return nil, xerrors.Errorf("len(depthToNodes) > len(realLayersTypes)")
	}
	for i := 1; i < len(realLayersTypes); i++ {
		if i > len(depthToNodes) { // if there are only top levels of tree
			break
		}
		if arr, ok := depthToNodes[i]; ok {
			if i == 0 {
				if len(arr) != 1 {
					return nil, xerrors.Errorf("len(arr) != 1")
				}
			}
		} else {
			return nil, xerrors.Errorf("absent depth: %d", i)
		}
	}

	// deserialize
	// performance here can be significantly improved - but later :)

	for i := 1; i < len(realLayersTypes); i++ {
		if i > len(depthToNodes) { // if there are only top levels of tree
			break
		}
		serializedNodes := depthToNodes[i]

		for _, el := range serializedNodes {
			// hack
			// the only way I can get unmarshal the type which I need - to pointer (which I got via Interface())
			// and then we need to dereference pointer into interface
			pointerToRealObj := reflect.New(reflect.TypeOf(realLayersTypes[i])).Interface()
			err = json.Unmarshal([]byte(el.CurrentNodeJSON), pointerToRealObj)
			if err != nil {
				return nil, xerrors.Errorf("%w", err)
			}
			realObj := reflect.ValueOf(pointerToRealObj).Elem().Interface()

			parentKey, err := model.ParseComplexKey(el.ParentFullKey)
			if err != nil {
				return nil, xerrors.Errorf("%w", err)
			}

			parentNode, err := tree.Root().(*node).GetChildNodeByComplexKey(parentKey.CutFirstSubkey())
			if err != nil {
				return nil, xerrors.Errorf("%w", err)
			}

			_, err = parentNode.CreateOrGetChildNode(realObj.(model.IDable))
			if err != nil {
				return nil, xerrors.Errorf("%w", err)
			}
		}
	}

	return tree, nil
}

func normalizeForCanonizing(in []model.DBTreeNode) []string {
	resultStrings := make([]string, 0)
	for _, el := range in {
		jsonBytes, _ := json.Marshal(el)
		resultStrings = append(resultStrings, string(jsonBytes))
	}
	sort.Strings(resultStrings)
	return resultStrings
}
