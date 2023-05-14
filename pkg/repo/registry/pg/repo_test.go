package pg

import (
	"context"
	"github.com/stretchr/testify/require"
	"os"
	"personal-feed/pkg/model"
	"strings"
	"testing"
)

func TestContext(t *testing.T) {
	buf, err := os.ReadFile("./repo.go")
	require.NoError(t, err)
	repoStr := string(buf)
	require.False(t, strings.Contains(repoStr, "context.Background"))
	require.False(t, strings.Contains(repoStr, "context.TODO"))
}

func TestRepo(t *testing.T) {
	t.Skip()

	sourceID := 1

	ctx := context.Background()

	cfg := RepoConfigPG{}
	client, err := NewRepo(ctx, cfg, nil)
	require.NoError(t, err)

	////---
	//// prepare db
	//
	//_, err = client.conn.Exec(context.Background(), "delete from events;")
	//require.NoError(t, err)

	//---
	// insert 1 row

	tx, err := client.NewTx(ctx)
	require.NoError(t, err)

	nodes, err := client.ExtractTreeNodesTx(tx, ctx, sourceID)
	require.NoError(t, err)
	require.Equal(t, 0, len(nodes))

	newNodes := []model.DBTreeNode{
		{
			Depth:           1,
			CurrentFullKey:  "parent!abc",
			CurrentNodeJSON: "{}",
		},
	}
	err = client.InsertNewTreeNodesTx(tx, ctx, sourceID, newNodes)
	require.NoError(t, err)

	err = tx.Commit(context.Background())
	require.NoError(t, err)

	//---
	// list table from tx

	tx2, err := client.NewTx(ctx)
	require.NoError(t, err)

	nodes2, err := client.ExtractTreeNodesTx(tx2, ctx, sourceID)
	require.NoError(t, err)
	require.Equal(t, 1, len(nodes2))
}
