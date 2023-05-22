package config

import "personal-feed/pkg/config/configsengine"

type typeTagged configsengine.TypeTagged

//---

type Config struct {
	Repo      RepoConfig `mapstructure:"repo"`
	AllowCORS bool
}

type RepoConfig interface {
	typeTagged
	IsRepoConfig()
}
