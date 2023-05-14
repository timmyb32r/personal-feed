package repo

import (
	"context"
	"github.com/sirupsen/logrus"
	"personal-feed/pkg/config"
	"personal-feed/pkg/config/configsengine"
	"personal-feed/pkg/util"
	"strings"
)

type repoFactory func(context.Context, interface{}, *logrus.Logger) (Repo, error)

var configNameToRepoFactory = make(map[string]repoFactory)

func Register(foo repoFactory, repoConfig config.RepoConfig) {
	configNameToRepoFactory[util.GetStructName(repoConfig)] = foo
	repoName := util.ToKebabCase(strings.TrimPrefix(util.GetStructName(repoConfig), "RepoConfig"))
	var tmpVal *config.RepoConfig = nil
	configsengine.RegisterTypeTagged(tmpVal, repoConfig, repoName)
}
