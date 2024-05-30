package db

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/lithiferous/cmd-etl/util"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
)

var testStore Store

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	_, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	pgconf, err := pgxpool.ParseConfig(config.DBSource)
	if err != nil {
		log.Fatal("cannot parse URI db:", err)
	}

	pgconf.MaxConns = 40

	pgconf.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxdecimal.Register(conn.TypeMap())
		return nil
	}

	connPool, err := pgxpool.NewWithConfig(context.Background(), pgconf)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testStore = NewStore(connPool)
	os.Exit(m.Run())
}
