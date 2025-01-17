package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/vandi37/TgLogger/internal/config"
	"github.com/vandi37/vanerrors"
)

const (
	ErrorOpeningDataBase = "error opining database"
	CheckingConnection   = "checking database connection failed"
	ErrorCreateTable     = "error to create table"
)

type DB struct {
	*sql.DB
}

func New(cfg config.DBConfig) (*DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("%s://%s:%s@db:%d/%s?sslmode=disable", cfg.Host, cfg.Username, cfg.Password, cfg.Port, cfg.Name))
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorOpeningDataBase, err, vanerrors.EmptyHandler)
	}
	err = db.Ping()
	if err != nil {
		return nil, vanerrors.NewWrap(CheckingConnection, err, vanerrors.EmptyHandler)
	}
	return &DB{db}, nil
}

func (db *DB) Init() error {

	query := `
CREATE TABLE  IF NOT EXISTS users (
    id BIGINT NOT NULL,
    created TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS tokens  (
    token CHAR(25) PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
`

	_, err := db.Exec(query)
	if err != nil {
		return vanerrors.NewWrap(ErrorCreateTable, err, vanerrors.EmptyHandler)
	}

	return nil
}
