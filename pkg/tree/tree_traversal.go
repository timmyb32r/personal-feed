package tree

import "personal-feed/pkg/model"

func extractInternalNodes(in *node) map[string]*node {
	if in.IsLeaf() {
		return map[string]*node{in.ComplexKey().FullKey(): in}
	}
	result := make(map[string]*node)
	childrenKeys := in.ChildrenKeys()
	for _, childrenKey := range childrenKeys {
		childNodeObj, _ := in.GetChildNodeByKey(childrenKey)
		childNode := childNodeObj.(*node)
		currMap := extractInternalNodes(childNode)
		for k, v := range currMap {
			result[k] = v
		}
		result[childNode.ComplexKey().FullKey()] = childNode
	}
	return result
}

func extractDocs(in *node) map[string]doc {
	result := make(map[string]doc)
	childrenKeys := in.ChildrenKeys()
	if in.IsLeaf() {
		for _, childrenKey := range childrenKeys {
			currDoc := doc{
				parentNode: in,
				key:        childrenKey,
			}
			result[in.ComplexKey().MakeSubkey(childrenKey.ID()).FullKey()] = currDoc
		}
	} else {
		for _, childrenKey := range childrenKeys {
			childNode, _ := in.GetChildNodeByKey(childrenKey)
			currMap := extractDocs(childNode.(*node))
			for k, v := range currMap {
				result[k] = v
			}
		}
	}
	return result
}

func extractDocsUnwrapped(in *node) map[string]model.IDable {
	result := make(map[string]model.IDable)
	childrenKeys := in.ChildrenKeys()
	if in.IsLeaf() {
		for _, childrenKey := range childrenKeys {
			result[in.ComplexKey().MakeSubkey(childrenKey.ID()).FullKey()] = childrenKey
		}
	} else {
		for _, childrenKey := range childrenKeys {
			childNode, _ := in.GetChildNodeByKey(childrenKey)
			currMap := extractDocsUnwrapped(childNode.(*node))
			for k, v := range currMap {
				result[k] = v
			}
		}
	}
	return result
}
