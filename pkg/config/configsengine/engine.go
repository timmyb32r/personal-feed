package configsengine

import (
	"github.com/mitchellh/mapstructure"
	"golang.org/x/xerrors"
	"gopkg.in/yaml.v3"
	"io"
	"reflect"
)

func makeConfigDecodeHook() mapstructure.DecodeHookFuncValue {
	return func(sourceValue, destinationValue reflect.Value) (interface{}, error) {
		destinationType := destinationValue.Type()
		if isTypeTaggedInterface(destinationType) {
			value, err := decodeTypeTaggedValue(destinationType, sourceValue.Interface())
			if err != nil {
				return nil, xerrors.Errorf("unable to decode value of type-tagged type %s, err: %w", destinationType.String(), err)
			}
			return value, nil
		}
		return sourceValue.Interface(), nil
	}
}

func decodeTypeTaggedValue(ifaceType reflect.Type, value interface{}) (interface{}, error) {
	typeTaggedMap, ok := value.(map[string]interface{})
	if !ok {
		return nil, xerrors.Errorf("unable to case value to map, value type: %T", value)
	}
	typeTagTypeErased, ok := typeTaggedMap["type"]
	if !ok {
		return nil, xerrors.New("unable to find into typeTag field 'type'")
	}
	delete(typeTaggedMap, "type")

	tag, ok := typeTagTypeErased.(string)
	if !ok {
		return nil, xerrors.Errorf("unable to convert 'type' into string, it's: %T", typeTagTypeErased)
	}

	tagMap, ok := typeTagRegistry[ifaceType]
	if !ok {
		return nil, xerrors.New("unable to find structure in the registry")
	}

	currRegistryEntry, ok := tagMap[typeTag(tag)]
	if !ok {
		return nil, xerrors.Errorf("unable to find typeTag into the registry, typeTag: %s", tag)
	}

	implValue := reflect.New(currRegistryEntry.implType).Interface()
	if err := decodeMap(typeTaggedMap, implValue, makeConfigDecodeHook()); err != nil {
		return nil, xerrors.Errorf("unable to decode type-tagged map, err: %w", err)
	}
	for key := range typeTaggedMap {
		return nil, xerrors.Errorf("found extra field: %s", key)
	}
	if !reflect.TypeOf(implValue).Implements(ifaceType) {
		return nil, xerrors.Errorf("value of type %T does not implement interface type %s", implValue, ifaceType.Name())
	}

	return implValue, nil
}

func decodeMap(configMap map[string]interface{}, result interface{}, decodeHook mapstructure.DecodeHookFuncValue) error {
	var err error
	var decoderMetadata mapstructure.Metadata
	configDecoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:     result,
		Metadata:   &decoderMetadata,
		DecodeHook: decodeHook,
	})
	if err != nil {
		return xerrors.Errorf("unable to create map decoder, err: %w", err)
	}
	if err := configDecoder.Decode(configMap); err != nil {
		return xerrors.Errorf("unable to decode configMap, err: %w", err)
	}
	for _, key := range decoderMetadata.Keys {
		delete(configMap, key)
	}
	return nil
}

func FillConfigStruct(reader io.Reader, outputStruct interface{}) error {
	configMap := map[string]interface{}{}
	yamlDecoder := yaml.NewDecoder(reader)
	if err := yamlDecoder.Decode(&configMap); err != nil {
		return xerrors.Errorf("unable to decode yaml, err: %w", err)
	}

	common := reflect.New(reflect.TypeOf(outputStruct).Elem().Elem())
	if err := decodeMap(configMap, common.Interface(), makeConfigDecodeHook()); err != nil {
		return xerrors.Errorf("unable to decode config, err: %w", err)
	}
	reflect.ValueOf(outputStruct).Elem().Set(common)
	return nil
}
