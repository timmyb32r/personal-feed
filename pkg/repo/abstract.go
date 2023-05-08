package repo

import (
	"context"
	"personal-feed/pkg/model"
	"time"
)

type Tx interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// > mockgen -source ./abstract.go -package db -destination ./abstract_mock.go

type Repo interface {
	GenerateLiquibaseProperties() (string, error)

	NewTx() (Tx, error)

	GetUserInfo(tx Tx, userEmail string) (*model.User, error) // returns nil if user not found
	UpdateUserInfo(tx Tx, userEmail string, user *model.User) error

	ListSources() ([]model.Source, error)

	InsertNewTreeNodesTx(tx Tx, sourceID int, nodes []model.DBTreeNode) error
	InsertNewTreeNodes(ctx context.Context, sourceID int, nodes []model.DBTreeNode) error

	ExtractTreeNodesTx(tx Tx, sourceID int) ([]model.DBTreeNode, error)
	ExtractTreeNodes(ctx context.Context, sourceID int) ([]model.DBTreeNode, error)

	GetNextCronPeriod(ctx context.Context) (lastRunTime *time.Time, currentTime time.Time, err error)
	SetCronLastRunTime(ctx context.Context, cronLastRunTime time.Time) error

	SetState(ctx context.Context, sourceID int, state string) error

	InsertSourceIterationTx(tx Tx, ctx context.Context, sourceID int, body string) error
	InsertSourceIteration(ctx context.Context, sourceID int, body string) error

	// temporary things

	TestExtractAllTreeNodes(tx Tx) ([]model.DBTreeNode, error)
}
