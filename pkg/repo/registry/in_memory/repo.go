package in_memory

import (
	"context"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"personal-feed/pkg/model"
	"personal-feed/pkg/repo"
	"sync"
	"time"
)

type Repo struct {
	base            map[int][]model.DBTreeNode
	cronLastRunTime *time.Time
	mutex           sync.Mutex
}

func (r *Repo) GenerateLiquibaseProperties() (string, error) {
	return "", xerrors.Errorf("liquibase is unsupported for in-memory repo")
}

func (r *Repo) NewTx(_ context.Context) (repo.Tx, error) {
	return &TxStub{}, nil
}

func (r *Repo) GetUserInfo(_ repo.Tx, _ context.Context, _ string) (*model.User, error) {
	return nil, nil
}

func (r *Repo) UpdateUserInfo(_ repo.Tx, _ context.Context, _ string, _ *model.User) error {
	return nil
}

func (r *Repo) ListSources(_ context.Context) ([]model.Source, error) {
	return nil, nil
}

func (r *Repo) InsertNewTreeNodes(ctx context.Context, sourceID int, nodes []model.DBTreeNode) error {
	return r.InsertNewTreeNodesTx(nil, ctx, sourceID, nodes)
}

func (r *Repo) InsertNewTreeNodesTx(_ repo.Tx, _ context.Context, sourceID int, nodes []model.DBTreeNode) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.base[sourceID]; !ok {
		r.base[sourceID] = make([]model.DBTreeNode, 0)
	}

	r.base[sourceID] = append(r.base[sourceID], nodes...)
	return nil
}

func (r *Repo) ExtractTreeNodesTx(_ repo.Tx, _ context.Context, sourceID int) ([]model.DBTreeNode, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	return r.base[sourceID], nil
}

func (r *Repo) ExtractTreeNodes(ctx context.Context, sourceID int) ([]model.DBTreeNode, error) {
	return r.ExtractTreeNodesTx(nil, ctx, sourceID)
}

func (r *Repo) GetNextCronPeriod(_ context.Context) (lastRunTime *time.Time, currentTime time.Time, err error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	cronLastRunTimeCopy := *r.cronLastRunTime
	return &cronLastRunTimeCopy, time.Now(), nil
}

func (r *Repo) SetCronLastRunTime(_ context.Context, cronLastRunTime time.Time) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.cronLastRunTime = &cronLastRunTime
	return nil
}

func (r *Repo) Len() int {
	sum := 0
	for _, v := range r.base {
		sum += len(v)
	}
	return sum
}

func (r *Repo) SetState(_ context.Context, _ int, _ string) error {
	return nil
}

func (r *Repo) InsertSourceIterationTx(_ repo.Tx, _ context.Context, _ int, _, _ string) error {
	return nil
}

func (r *Repo) InsertSourceIteration(_ context.Context, _ int, _, _ string) error {
	return nil
}

func (r *Repo) TestExtractAllTreeNodes(_ repo.Tx, _ context.Context) ([]model.DBTreeNode, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	result := make([]model.DBTreeNode, 0)
	for _, v := range r.base {
		result = append(result, v...)
	}

	return result, nil
}

func (r *Repo) Close(_ context.Context) error {
	return nil
}

func NewRepo(_ context.Context, _ interface{}, _ *logrus.Logger) (repo.Repo, error) {
	return &Repo{
		base:            make(map[int][]model.DBTreeNode),
		cronLastRunTime: nil,
		mutex:           sync.Mutex{},
	}, nil
}

func init() {
	repo.Register(NewRepo, &RepoConfigInMemory{})
}
