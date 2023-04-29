package in_memory

import "context"

type TxStub struct {
}

func (t *TxStub) Commit(ctx context.Context) error {
	return nil
}

func (t *TxStub) Rollback(ctx context.Context) error {
	return nil
}
