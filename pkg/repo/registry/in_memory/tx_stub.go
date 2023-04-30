package in_memory

import "context"

type TxStub struct {
}

func (t *TxStub) Commit(_ context.Context) error {
	return nil
}

func (t *TxStub) Rollback(_ context.Context) error {
	return nil
}
