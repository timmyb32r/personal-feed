package repo

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"personal-feed/pkg/config"
	"personal-feed/pkg/util"
)

func NewRepo(repoConfig config.RepoConfig, logger *logrus.Logger) (Repo, error) {
	configName := util.GetStructName(repoConfig)
	if currRepoFactory, ok := configNameToRepoFactory[configName]; ok {
		return currRepoFactory(repoConfig, logger)
	} else {
		return nil, xerrors.Errorf("unknown configName: %s", configName)
	}
}
