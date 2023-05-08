package engine

import (
	"context"
	"personal-feed/pkg/operation"
)

type AbstractEngine interface {
	RunOnce(ctx context.Context, op operation.Operation) error
}
