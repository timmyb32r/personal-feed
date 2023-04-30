package config

import "personal-feed/pkg/config/configsengine"

type typeTagged configsengine.TypeTagged

//---

type Config struct {
	Repo RepoConfig `mapstructure:"repo"`
}

type RepoConfig interface {
	typeTagged
	IsRepoConfig()
}
