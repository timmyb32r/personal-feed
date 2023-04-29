package model

import (
	"context"
)

type Tx interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// > mockgen -source ./database_client.go -package db -destination ./database_client_mock.go

type DatabaseClient interface {
	NewTx() (Tx, error)

	GetUserInfo(tx Tx, userEmail string) (*User, error) // returns nil if user not found
	UpdateUserInfo(tx Tx, userEmail string, user *User) error

	//WithCrawlerProjectLock(ctx context.Context, projectName, crawlerType string, f func(ctx context.Context) error) error

	ListSources(tx Tx) ([]Source, error)

	InsertNewTreeNodes(tx Tx, sourceID int, nodes []DBTreeNode) error
	ExtractTreeNodes(tx Tx, sourceID int) ([]DBTreeNode, error)
}
