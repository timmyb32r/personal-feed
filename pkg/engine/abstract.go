package engine

import "context"

type AbstractEngine interface {
	RunOnce(ctx context.Context) error
}
