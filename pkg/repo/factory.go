package repo

import (
	"golang.org/x/xerrors"
	"personal-feed/pkg/config"
	"personal-feed/pkg/util"
)

func NewRepo(repoConfig config.RepoConfig) (Repo, error) {
	configName := util.GetStructName(repoConfig)
	if currRepoFactory, ok := configNameToRepoFactory[configName]; ok {
		return currRepoFactory(repoConfig)
	} else {
		return nil, xerrors.Errorf("unknown configName: %s", configName)
	}
}
