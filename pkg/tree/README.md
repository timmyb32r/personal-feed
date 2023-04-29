There are main exported structure: `Tree`

```
type Tree struct {
	rootKey     model.IDable
	layersTypes []model.IDable // without Root one
	layers      []*layer
	root        *node
}

type IDable interface {
	ID() string
}

type leaf struct {
	parentNode *node
	key        model.IDable
}

type layer struct {
	// types stuff
	masterObjCurrLayerKey model.IDable

	// layer info
	nextLayer *layer

	// flag isLeaf
	isLeaf bool
}
```

ID - it's 'ProjectName' string

SourceID - it's project ID (integer)

ComplexKey - it's path (keys sequence) to node via '!'. Every key is url-encoded. For example: "a!b!c'

FullKey - is the same as 'ComplexKey', but ComplexKey is an object, and 'FullKey' is a string

There is 'key' & 'key_id'. On 'key' you can call .ID() to get 'key_id'. Key can be complex object, and 'key_id' is just identifier of this object

If the document complexKey: a!b!c, 'b' will be the leaf! 'c' will be the document. So, 'b' will be 'node' with isLeaf:true, and map 'keyIDToChildNode' will map from 'c' to this 'a!b!c' document - it's 'IDable' but not a node

```
root->internal_node->...->internal_node(leaf)->doc
```

Serialize - serializes objects into model.DB* things, not into string

Node fields:
- 'key' (like short filename)
- 'complexKey' (like full filename)
- layerDescr - pointer to description of current level
- parentNode - pointer to parent node
- keyIDToChildNode - map from child keyName to child

Node methods:
- func (n *node) CreateOrGetChildNode(nextKey model.IDable) (interface{}, error) // create child node by 'key'
- func (n *node) CreateOrGetChildNodes(nextKeys []model.IDable) (interface{}, error) // create child node by 'complexKey'
- func (n *node) GetChildNodeByComplexKey(complexKey *model.ComplexKey) (*node, error)
- func (n *node) ChildrenKeys() []model.IDable

Tree methods:
- func (t *Tree) ProjectName() string
- func (t *Tree) LayersTypes() []model.IDable
- func (t *Tree) Layers() []*layer
- func (t *Tree) Root() *node
- func (t *Tree) InsertNode(inNode *node) (*node, error)
- func (t *Tree) InsertLeaf(leaf leaf) error
- func (t *Tree) ExtractInternalNodes() map[string]*node
- func (t *Tree) ExtractDocs() map[string]doc
- func (t *Tree) ExtractDocsUnwrapped() map[string]model.IDable
- func (t *Tree) Serialize(sourceID int) []model.DBTreeNode
- func (t *Tree) SerializeKey(sourceID int, fullKey string, key model.IDable) *model.DBTreeNode
- ctor: NewTree(layersTypes []model.IDable) (*Tree, error)

---

There are main exported function `BuildDiffTreeAndSerialize`:

```
func BuildDiffTreeAndSerialize

    (
        sourceID int,
        known  *Tree,
        existing *Tree
    )

    (
        []model.DBTreeNode, // dbNewInternalNodes
        []model.DBTreeNode, // dbNewDocs
        error
    )
```

So, this function determines new nodes, which can be inserted into the database then.
