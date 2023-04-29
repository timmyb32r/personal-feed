package in_memory

import (
	"personal-feed/pkg/model"
	"sync"
)

type InMemoryDatabaseClient struct {
	base  map[int][]model.DBTreeNode
	mutex sync.Mutex
}

func (c *InMemoryDatabaseClient) NewTx() (model.Tx, error) {
	return &TxStub{}, nil
}

//func (c *InMemoryDatabaseClient) WithCrawlerProjectLock(ctx context.Context, projectName, crawlerType string, f func(ctx context.Context) error) error {
//	return nil
//}

func (c *InMemoryDatabaseClient) GetUserInfo(tx model.Tx, userEmail string) (*model.User, error) {
	return nil, nil
}

func (c *InMemoryDatabaseClient) UpdateUserInfo(tx model.Tx, userEmail string, user *model.User) error {
	return nil
}

func (c *InMemoryDatabaseClient) ListSources(tx model.Tx) ([]model.Source, error) {
	return nil, nil
}

func (c *InMemoryDatabaseClient) InsertNewTreeNodes(tx model.Tx, sourceID int, nodes []model.DBTreeNode) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.base[sourceID]; !ok {
		c.base[sourceID] = make([]model.DBTreeNode, 0)
	}

	c.base[sourceID] = append(c.base[sourceID], nodes...)
	return nil
}

func (c *InMemoryDatabaseClient) ExtractTreeNodes(tx model.Tx, sourceID int) ([]model.DBTreeNode, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.base[sourceID], nil
}

func (c *InMemoryDatabaseClient) Len() int {
	sum := 0
	for _, v := range c.base {
		sum += len(v)
	}
	return sum
}

func NewInMemoryDatabaseClient() *InMemoryDatabaseClient {
	return &InMemoryDatabaseClient{
		base:  make(map[int][]model.DBTreeNode),
		mutex: sync.Mutex{},
	}
}
