package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lithiferous/cmd-etl/util"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func createRandomSnapshot(t *testing.T) Snapshot {
	rand_ts := util.RandomTimestamp()
	arg := CreateSnapshotParams{
		UserName:    util.RandomUser(),
		StoreName:   util.RandomStore(),
		CreditLimit: decimal.NewFromFloat(util.RandomMoney()),
		SnapshotAt: pgtype.Timestamp{
			Time:  rand_ts,
			Valid: !rand_ts.IsZero(),
		},
	}

	snapshot, err := testStore.CreateSnapshot(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, snapshot)

	require.Equal(t, arg.UserName, snapshot.UserName)
	require.Equal(t, arg.StoreName, snapshot.StoreName)
	require.Equal(t, arg.CreditLimit, snapshot.CreditLimit)
	require.Equal(t, arg.SnapshotAt, snapshot.SnapshotAt)

	require.NotZero(t, snapshot.CreatedAt)

	return snapshot
}

func TestCreateSnapshot(t *testing.T) {
	createRandomSnapshot(t)
}

func TestListSnapshots(t *testing.T) {
	for i := 0; i < 10; i++ {
		_ = createRandomSnapshot(t)
	}

	snapshots, err := testStore.ListSnapshots(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, snapshots)

	for _, snapshot := range snapshots {
		require.NotEmpty(t, snapshot)
	}
}
