package tree

import "personal-feed/pkg/model"

type doc struct {
	parentNode *node
	key        model.IDable
}
