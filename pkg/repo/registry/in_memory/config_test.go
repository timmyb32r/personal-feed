package in_memory

import (
	"github.com/stretchr/testify/require"
	"personal-feed/pkg/config"
	"personal-feed/pkg/util"
	"strings"
	"testing"
)

var testYaml = `
repo:
  type: in-memory
`

func TestConfig(t *testing.T) {
	yamlReader := strings.NewReader(testYaml)
	currConfig, err := config.Load(yamlReader)
	require.NoError(t, err)
	require.Equal(t, "RepoConfigInMemory", util.GetStructName(currConfig.Repo))
}
