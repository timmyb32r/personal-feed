package pg

import (
	"github.com/stretchr/testify/require"
	"personal-feed/pkg/config"
	"personal-feed/pkg/util"
	"strings"
	"testing"
)

var testYaml = `
repo:
  type: pg
  db_host: my_host
`

func TestConfig(t *testing.T) {
	yamlReader := strings.NewReader(testYaml)
	currConfig, err := config.Load(yamlReader)
	require.NoError(t, err)
	require.Equal(t, "RepoConfigPG", util.GetStructName(currConfig.Repo))
	require.Equal(t, "my_host", currConfig.Repo.(*RepoConfigPG).Host)
}
