package db

import (
	"context"
)

type Querier interface {
	CreateSnapshot(ctx context.Context, arg CreateSnapshotParams) (Snapshot, error)
	ListSnapshots(ctx context.Context) ([]Snapshot, error)
}

var _ Querier = (*Queries)(nil)
