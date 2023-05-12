package pg

import (
	"context"
	"github.com/stretchr/testify/require"
	"personal-feed/pkg/model"
	"testing"
)

func TestRepo(t *testing.T) {
	t.Skip()

	sourceID := 1

	cfg := RepoConfigPG{}
	client, err := NewRepo(cfg, nil)
	require.NoError(t, err)

	////---
	//// prepare db
	//
	//_, err = client.conn.Exec(context.Background(), "delete from events;")
	//require.NoError(t, err)

	//---
	// insert 1 row

	tx, err := client.NewTx()
	require.NoError(t, err)

	nodes, err := client.ExtractTreeNodesTx(tx, sourceID)
	require.NoError(t, err)
	require.Equal(t, 0, len(nodes))

	newNodes := []model.DBTreeNode{
		{
			Depth:           1,
			CurrentFullKey:  "parent!abc",
			CurrentNodeJSON: "{}",
		},
	}
	err = client.InsertNewTreeNodesTx(tx, sourceID, newNodes)
	require.NoError(t, err)

	err = tx.Commit(context.Background())
	require.NoError(t, err)

	//---
	// list table from tx

	tx2, err := client.NewTx()
	require.NoError(t, err)

	nodes2, err := client.ExtractTreeNodesTx(tx2, sourceID)
	require.NoError(t, err)
	require.Equal(t, 1, len(nodes2))
}
