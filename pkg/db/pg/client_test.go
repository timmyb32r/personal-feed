package pg

import (
	"context"
	"github.com/stretchr/testify/require"
	"personal-feed/pkg/model"
	"testing"
)

func TestPgClient(t *testing.T) {
	t.Skip()

	sourceID := 1

	cfg := NewConfig("?", "?", "?", 6432, "?", true)
	client, err := NewPgClient(cfg)
	require.NoError(t, err)

	//---
	// prepare db

	_, err = client.conn.Exec(context.Background(), "delete from events;")
	require.NoError(t, err)

	//---
	// insert 1 row

	tx, err := client.NewTx()
	require.NoError(t, err)

	nodes, err := client.ExtractTreeNodes(tx, sourceID)
	require.NoError(t, err)
	require.Equal(t, 0, len(nodes))

	newNodes := []model.DBTreeNode{
		{
			Depth:           1,
			ParentFullKey:   "parent",
			CurrentNodeJSON: "{}",
		},
	}
	err = client.InsertNewTreeNodes(tx, sourceID, newNodes)
	require.NoError(t, err)

	err = tx.Commit(context.Background())
	require.NoError(t, err)

	//---
	// list table from tx

	tx2, err := client.NewTx()
	require.NoError(t, err)

	nodes2, err := client.ExtractTreeNodes(tx2, sourceID)
	require.NoError(t, err)
	require.Equal(t, 1, len(nodes2))
}
