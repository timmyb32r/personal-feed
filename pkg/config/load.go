package config

import (
	"golang.org/x/xerrors"
	"io"
	"personal-feed/pkg/config/configsengine"
)

func Load(reader io.Reader) (*Config, error) {
	result := &Config{}
	err := configsengine.FillConfigStruct(reader, &result)
	if err != nil {
		return nil, xerrors.Errorf("unable to decode config, err: %w", err)
	}
	return result, nil
}
