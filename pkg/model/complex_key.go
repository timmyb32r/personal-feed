package model

import (
	"net/url"
	"strings"
)

type ComplexKey struct {
	keys []string // stored in original form
}

func (k *ComplexKey) Empty() bool {
	return len(k.keys) == 0
}

func (k *ComplexKey) Keys() []string {
	return k.keys
}

func (k *ComplexKey) FullKey() string {
	result := make([]string, 0, len(k.keys))
	for _, el := range k.keys {
		result = append(result, url.QueryEscape(el))
	}
	return strings.Join(result, "!")
}

func (k *ComplexKey) ParentKey() *ComplexKey {
	return &ComplexKey{
		keys: k.keys[0 : len(k.keys)-1],
	}
}

func (k *ComplexKey) CutFirstSubkey() *ComplexKey {
	return &ComplexKey{
		keys: k.keys[1:],
	}
}

func (k *ComplexKey) ShortKey() string {
	return k.keys[len(k.keys)-1]
}

func (k *ComplexKey) Depth() int {
	return len(k.keys) - 1
}

func (k *ComplexKey) MakeSubkey(in string) *ComplexKey {
	var keysCopy []string
	copy(keysCopy, k.keys)
	keysCopy = append(keysCopy, in)

	return &ComplexKey{
		keys: keysCopy,
	}
}

func ParseComplexKey(fullKey string) (*ComplexKey, error) {
	keys := strings.Split(fullKey, "!")
	result := make([]string, 0, len(keys))
	for _, el := range keys {
		currKey, err := url.QueryUnescape(el)
		if err != nil {
			return nil, err
		}
		result = append(result, currKey)
	}
	return &ComplexKey{
		keys: result,
	}, nil
}

func NewComplexKey(topLevelKey string) *ComplexKey {
	return &ComplexKey{
		keys: []string{topLevelKey},
	}
}
