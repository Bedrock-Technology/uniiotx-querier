package storer

import (
	"context"
	"database/sql"
	_ "embed"
	"github.com/Bedrock-Technology/uniiotx-querier/common"
	"github.com/Bedrock-Technology/uniiotx-querier/storer/sqlc"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var ddl string

type MyStorer struct {
	logger common.Logger
	db     *sql.DB
	q      *sqlc.Queries
}

func NewStorer(sqliteDSN string, logger common.Logger) *MyStorer {
	ctx := context.Background()

	db, err := sql.Open("sqlite3", sqliteDSN)
	if err != nil {
		logger.Fatal("failed to open database", err)
	}

	// Create tables if not exists
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		logger.Fatal("failed to create tables if not exits", err)
	}

	logger.Info("database opened")
	return &MyStorer{logger: logger, db: db, q: sqlc.New(db)}
}

func (s *MyStorer) DB() *sql.DB {
	return s.db
}

func (s *MyStorer) Queries() *sqlc.Queries {
	return s.q
}

func (s *MyStorer) Close() {
	err := s.db.Close()
	if err != nil {
		s.logger.Error("failed to close database", err)
		return
	}
	s.logger.Info("database closed")
}
