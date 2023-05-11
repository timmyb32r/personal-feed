package tree

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"personal-feed/pkg/model"
	"sort"
	"testing"
)

type stNt string // serializer test node type

func (n stNt) ID() string {
	return string(n)
}

func normalizeForCanonizing(t *testing.T, in []model.DBTreeNode) []string {
	type nodeWithoutBusinessTime struct {
		SourceID        int    `db:"source_id"`
		Depth           int    `db:"depth"` // 0 means 'root'
		ParentFullKey   string `db:"parent_full_key"`
		CurrentNodeJSON string `db:"current_node_json"` // here are serialized object of current depth type
	}

	resultStrings := make([]string, 0)
	for _, el := range in {
		jsonBytes, _ := json.Marshal(el)

		// remove BusinessTime
		var tmp nodeWithoutBusinessTime
		err := json.Unmarshal(jsonBytes, &tmp)
		require.NoError(t, err)
		jsonBytes, _ = json.Marshal(tmp)

		resultStrings = append(resultStrings, string(jsonBytes))
	}
	sort.Strings(resultStrings)
	return resultStrings
}

func TestSerde(t *testing.T) {
	sourceID := 1

	layers := []model.IDable{
		stNt(""),
		stNt(""),
		stNt(""),
	}
	knownTree, err := NewTree(layers)
	require.NoError(t, err)

	//-----------------------------
	// https://tree.nathanfriend.io
	//-----------------------------
	//    Root/
	//    ├── A/
	//    │   ├── C/
	//    │   │   ├── F
	//    │   │   ├── G
	//    │   │   └── H
	//    │   └── D/
	//    │       └── I
	//    └── B/
	//        └── E/
	//            ├── D
	//            └── K
	//-----------------------------

	paths := [][]model.IDable{
		{stNt("A"), stNt("C"), stNt("F")},
		{stNt("A"), stNt("C"), stNt("G")},
		{stNt("A"), stNt("C"), stNt("H")},
		{stNt("A"), stNt("D"), stNt("I")},
		{stNt("B"), stNt("E"), stNt("J")},
		{stNt("B"), stNt("E"), stNt("K")},
	}
	for _, currPath := range paths {
		_, err := knownTree.Root().(*node).CreateOrGetChildNodes(currPath)
		require.NoError(t, err)
	}

	resultTreeNodes := serialize(sourceID, knownTree)
	resultStrings := normalizeForCanonizing(t, resultTreeNodes)
	canonized := []string{
		`{"SourceID":1,"Depth":1,"ParentFullKey":"ROOT","CurrentNodeJSON":"\"A\""}`,
		`{"SourceID":1,"Depth":1,"ParentFullKey":"ROOT","CurrentNodeJSON":"\"B\""}`,
		`{"SourceID":1,"Depth":2,"ParentFullKey":"ROOT!A","CurrentNodeJSON":"\"C\""}`,
		`{"SourceID":1,"Depth":2,"ParentFullKey":"ROOT!A","CurrentNodeJSON":"\"D\""}`,
		`{"SourceID":1,"Depth":2,"ParentFullKey":"ROOT!B","CurrentNodeJSON":"\"E\""}`,
		`{"SourceID":1,"Depth":3,"ParentFullKey":"ROOT!A!C","CurrentNodeJSON":"\"F\""}`,
		`{"SourceID":1,"Depth":3,"ParentFullKey":"ROOT!A!C","CurrentNodeJSON":"\"G\""}`,
		`{"SourceID":1,"Depth":3,"ParentFullKey":"ROOT!A!C","CurrentNodeJSON":"\"H\""}`,
		`{"SourceID":1,"Depth":3,"ParentFullKey":"ROOT!A!D","CurrentNodeJSON":"\"I\""}`,
		`{"SourceID":1,"Depth":3,"ParentFullKey":"ROOT!B!E","CurrentNodeJSON":"\"J\""}`,
		`{"SourceID":1,"Depth":3,"ParentFullKey":"ROOT!B!E","CurrentNodeJSON":"\"K\""}`,
	}
	require.Equal(t, canonized, resultStrings)

	deserializedTree, err := Deserialize(resultTreeNodes, layers)
	require.NoError(t, err)

	resultTreeNodes2 := serialize(sourceID, deserializedTree)
	resultStrings2 := normalizeForCanonizing(t, resultTreeNodes2)
	require.Equal(t, canonized, resultStrings2)
}
