package configsengine

import (
	"fmt"
	"reflect"
)

type typeTag string
type typeTaggedInterface reflect.Type
type typeTaggedImplementation reflect.Type

type TypeTagged interface {
	IsTypeTagged()
}

type registryEntry struct {
	implType typeTaggedImplementation
}

var typeTagRegistry = map[typeTaggedInterface]map[typeTag]registryEntry{}

func isTypeTaggedInterface(typ reflect.Type) bool {
	if typ.Kind() != reflect.Interface {
		return false
	}
	currTypeTaggedInterface := reflect.TypeOf((*TypeTagged)(nil)).Elem()
	return typ.Implements(currTypeTaggedInterface)
}

func RegisterTypeTagged(iface interface{}, impl TypeTagged, tag string) {
	ifaceType := reflect.TypeOf(iface).Elem()
	implType := reflect.TypeOf(impl).Elem()
	tagMap, ok := typeTagRegistry[ifaceType]
	if !ok {
		tagMap = map[typeTag]registryEntry{}
		typeTagRegistry[ifaceType] = tagMap
	}
	if existingEntry, ok := tagMap[typeTag(tag)]; ok {
		panic(fmt.Sprintf(
			"tag %s for interface %s is already registered: conflicting implementations are %s and %s",
			tag,
			ifaceType.Name(),
			implType.Name(),
			existingEntry.implType.Name(),
		))
	}
	if implType.Kind() != reflect.Struct {
		panic("type-tagged interface implementation must be a pointer to a struct")
	}

	tagMap[typeTag(tag)] = registryEntry{
		implType: implType,
	}
}
