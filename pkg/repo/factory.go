package repo

import (
	"context"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"personal-feed/pkg/config"
	"personal-feed/pkg/util"
)

func NewRepo(ctx context.Context, repoConfig config.RepoConfig, logger *logrus.Logger) (Repo, error) {
	configName := util.GetStructName(repoConfig)
	if currRepoFactory, ok := configNameToRepoFactory[configName]; ok {
		return currRepoFactory(ctx, repoConfig, logger)
	} else {
		return nil, xerrors.Errorf("unknown configName: %s", configName)
	}
}
