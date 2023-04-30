package in_memory

import (
	"personal-feed/pkg/model"
	"personal-feed/pkg/repo"
	"sync"
)

type Repo struct {
	base  map[int][]model.DBTreeNode
	mutex sync.Mutex
}

func (c *Repo) NewTx() (repo.Tx, error) {
	return &TxStub{}, nil
}

func (c *Repo) GetUserInfo(_ repo.Tx, _ string) (*model.User, error) {
	return nil, nil
}

func (c *Repo) UpdateUserInfo(_ repo.Tx, _ string, _ *model.User) error {
	return nil
}

func (c *Repo) ListSources(_ repo.Tx) ([]model.Source, error) {
	return nil, nil
}

func (c *Repo) InsertNewTreeNodes(_ repo.Tx, sourceID int, nodes []model.DBTreeNode) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.base[sourceID]; !ok {
		c.base[sourceID] = make([]model.DBTreeNode, 0)
	}

	c.base[sourceID] = append(c.base[sourceID], nodes...)
	return nil
}

func (c *Repo) ExtractTreeNodes(_ repo.Tx, sourceID int) ([]model.DBTreeNode, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.base[sourceID], nil
}

func (c *Repo) Len() int {
	sum := 0
	for _, v := range c.base {
		sum += len(v)
	}
	return sum
}

func (c *Repo) TestExtractAllTreeNodes(_ repo.Tx) ([]model.DBTreeNode, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	result := make([]model.DBTreeNode, 0)
	for _, v := range c.base {
		result = append(result, v...)
	}
	return result, nil
}

func NewRepo(_ interface{}) (repo.Repo, error) {
	return &Repo{
		base:  make(map[int][]model.DBTreeNode),
		mutex: sync.Mutex{},
	}, nil
}

func init() {
	repo.Register(NewRepo, &RepoConfigInMemory{})
}
