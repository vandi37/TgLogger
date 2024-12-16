package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/vandi37/TgLogger/config"
	"github.com/vandi37/vanerrors"
)

// The errors
const (
	ErrorOpeningDataBase = "error opining database"
	ErrorCreateTable     = "error creating table"
	// ErrorPreparingQuery  = "error preparing query"
	ErrorInserting    = "error inserting"
	ErrorSelecting    = "error selecting"
	ErrorScanningRows = "error scanning rows"
	ErrorDeleting     = "error deleting"
	// NotFound             = "not found"
)

// The data base
type DB struct {
	db *sql.DB
}

// Creates a new data base connection
func New(cfg config.DBConfig) (*DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=localhost port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Port, cfg.Username, cfg.Password, cfg.Name))
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorOpeningDataBase, err, vanerrors.EmptyHandler)
	}
	return &DB{db: db}, nil
}

// Creates table if not exists
func (db *DB) Init() error {

	query := `
CREATE TABLE users (
    id BIGINT NOT NULL,
    created TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE TABLE tokens (
    token CHAR(15) PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
`

	_, err := db.db.Exec(query)
	if err != nil {
		return vanerrors.NewWrap(ErrorCreateTable, err, vanerrors.EmptyHandler)
	}

	return nil
}

// Close database
func (db *DB) Close() error {
	return db.db.Close()
}

// type User struct {
// 	Id      int64     `json:"id"`
// 	Created time.Time `json:"created"`
// }

// type Token struct {
// 	Token   string    `json:"token"`
// 	UserId  int64     `json:"user_id"`
// 	Created time.Time `json:"created"`
// }

// Creating new user
func (db *DB) NewUser(id int64) error {
	query := `insert into users (id) values ($1) on conflict (Id) do nothing;`

	_, err := db.db.Exec(query, id)
	if err != nil {
		return vanerrors.NewWrap(ErrorInserting, err, vanerrors.EmptyHandler)
	}

	return nil
}

func (db *DB) NewToken(token string, id int64) error {
	query := `insert into tokens (token, user_id) values ($1, $2)`

	_, err := db.db.Exec(query, token, id)
	if err != nil {
		return vanerrors.NewWrap(ErrorInserting, err, vanerrors.EmptyHandler)
	}

	return nil
}

// Checks token existence
func (db *DB) CheckToken(token string) (bool, error) {
	query := `select count(*) from tokens where token = $1`
	rows, err := db.db.Query(query, token)
	if err != nil {
		return false, vanerrors.NewWrap(ErrorSelecting, err, vanerrors.EmptyHandler)
	}
	rows.Next()

	var count int64
	err = rows.Scan(&count)
	if err != nil {
		return false, vanerrors.NewWrap(ErrorScanningRows, err, vanerrors.EmptyHandler)
	}

	return count > 0, nil
}

// Delete token
func (db *DB) DeleteToken(token string) error {
	query := `delete from tokens where token = $1`

	_, err := db.db.Exec(query, token)
	if err != nil {
		return vanerrors.NewWrap(ErrorDeleting, err, vanerrors.EmptyHandler)
	}

	return nil
}

// Getting tokens
func (db *DB) SelectTokens(id int64) ([]string, error) {
	query := `select token from tokens where user_id = $1`

	rows, err := db.db.Query(query, id)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorSelecting, err, vanerrors.EmptyHandler)
	}

	defer rows.Close()

	var tokens []string = []string{}
	for rows.Next() {
		var token string
		err = rows.Scan(&token)
		if err != nil {
			return nil, vanerrors.NewWrap(ErrorScanningRows, err, vanerrors.EmptyHandler)
		}

		tokens = append(tokens, token)
	}

	return tokens, err
}

func (db *DB) GetOwner(token string) (int64, error) {
	query := `select user_id from tokens where token = $1`

	rows, err := db.db.Query(query, token)
	if err != nil {
		return 0, vanerrors.NewWrap(ErrorSelecting, err, vanerrors.EmptyHandler)
	}

	defer rows.Close()

	rows.Next()

	var id int64
	err = rows.Scan(&id)
	if err != nil {
		return 0, vanerrors.NewWrap(ErrorScanningRows, err, vanerrors.EmptyHandler)
	}

	return id, nil
}
