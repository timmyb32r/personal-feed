package tree

import (
	"github.com/stretchr/testify/require"
	"personal-feed/pkg/model"
	"testing"
)

type isNt string // intersector test node type

func (n isNt) ID() string {
	return string(n)
}

func TestFillAllInternalNodes(t *testing.T) {
	layers := []model.IDable{
		isNt(""),
		isNt(""),
		isNt(""),
	}

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

	knownTree, err := NewTree(layers)
	require.NoError(t, err)

	knownPaths := [][]model.IDable{
		{isNt("A"), isNt("C"), isNt("F")},
		{isNt("A"), isNt("C"), isNt("G")},
		{isNt("A"), isNt("C"), isNt("H")},
		{isNt("A"), isNt("D"), isNt("I")},
		{isNt("B"), isNt("E"), isNt("J")},
		{isNt("B"), isNt("E"), isNt("K")},
	}
	for _, currPath := range knownPaths {
		_, err := knownTree.Root().(*node).CreateOrGetChildNodes(currPath)
		require.NoError(t, err)
	}

	//-----------------------------
	// https://tree.nathanfriend.io
	//-----------------------------
	//    Root/
	//    ├── A/
	//    │   ├── C/
	//    │   │   ├── F
	//    │   │   ├── G
	//    │   │   ├── H
	//    │   │   └── N
	//    │   ├── D/
	//    │   │   └── I
	//    │   └── M/
	//    │       └── O
	//    ├── B/
	//    │   └── E/
	//    │       ├── J
	//    │       └── K
	//    └── L/
	//        ├── P/
	//        │   ├── R
	//        │   └── S
	//        └── Q/
	//            └── T
	//-----------------------------

	existingTree, err := NewTree(layers)
	require.NoError(t, err)

	newPaths := [][]model.IDable{
		{isNt("A"), isNt("C"), isNt("F")},
		{isNt("A"), isNt("C"), isNt("G")},
		{isNt("A"), isNt("C"), isNt("H")},
		{isNt("A"), isNt("C"), isNt("N")},
		{isNt("A"), isNt("D"), isNt("I")},
		{isNt("A"), isNt("M"), isNt("O")},

		{isNt("B"), isNt("E"), isNt("J")},
		{isNt("B"), isNt("E"), isNt("K")},

		{isNt("L"), isNt("P"), isNt("R")},
		{isNt("L"), isNt("P"), isNt("S")},
		{isNt("L"), isNt("Q"), isNt("T")},
	}
	for _, currPath := range newPaths {
		_, err := existingTree.Root().(*node).CreateOrGetChildNodes(currPath)
		require.NoError(t, err)
	}

	//-----------------------------

	diffTree, newInternalNodes, err := BuildDiffTree(knownTree, existingTree)
	require.NoError(t, err)
	require.Equal(t, 4, len(newInternalNodes))
	newDocs := diffTree.ExtractDocsUnwrapped()
	require.Equal(t, 5, len(newDocs))
}
