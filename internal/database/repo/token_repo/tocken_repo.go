package token_repo

import (
	"context"

	"github.com/vandi37/TgLogger/internal/database/db"
	"github.com/vandi37/TgLogger/internal/database/repo"
	"github.com/vandi37/vanerrors"
)

type TokenRepo struct {
	db *db.DB
}

func New(db *db.DB) *TokenRepo {
	return &TokenRepo{db: db}
}

func (r *TokenRepo) New(ctx context.Context, token string, id int64) error {
	res, err := r.db.ExecContext(ctx, `insert into tokens (token, user_id) values ($1, $2)`, token, id)
	if err != nil {
		return vanerrors.NewWrap(repo.ErrorExecuting, err, vanerrors.EmptyHandler)
	}

	return repo.ReturnByRes(res, repo.Equals(1))
}

func (r *TokenRepo) Exist(ctx context.Context, token string) (bool, error) {
	stmt, err := r.db.PrepareContext(ctx, `select coalesce( (select 1 from tokens where token = $1), 0 );`)
	if err != nil {
		return false, vanerrors.NewWrap(repo.ErrorPreparing, err, vanerrors.EmptyHandler)
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, token)
	if err != nil {
		return false, vanerrors.NewWrap(repo.ErrorExecuting, err, vanerrors.EmptyHandler)
	}

	defer rows.Close()

	rows.Next()

	var res bool

	err = rows.Scan(&res)
	if err != nil {
		return false, vanerrors.NewWrap(repo.ErrorScanning, err, vanerrors.EmptyHandler)
	}

	return res, nil
}

func (r *TokenRepo) Delete(ctx context.Context, token string) error {
	res, err := r.db.ExecContext(ctx, `delete from tokens where token = $1`, token)
	if err != nil {
		return vanerrors.NewWrap(repo.ErrorExecuting, err, vanerrors.EmptyHandler)
	}

	return repo.ReturnByRes(res, repo.Equals(1))
}

func (r *TokenRepo) Select(ctx context.Context, id int64) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `select token from tokens where user_id = $1`, id)
	if err != nil {
		return nil, vanerrors.NewWrap(repo.ErrorExecuting, err, vanerrors.EmptyHandler)
	}

	defer rows.Close()

	var tokens []string = []string{}
	for rows.Next() {
		var token string
		err = rows.Scan(&token)
		if err != nil {
			return nil, vanerrors.NewWrap(repo.ErrorScanning, err, vanerrors.EmptyHandler)
		}

		tokens = append(tokens, token)
	}

	return tokens, err
}

func (r *TokenRepo) GetOwner(ctx context.Context, token string) (int64, error) {
	rows, err := r.db.QueryContext(ctx, `select user_id from tokens where token = $1`, token)
	if err != nil {
		return 0, vanerrors.NewWrap(repo.ErrorExecuting, err, vanerrors.EmptyHandler)
	}

	defer rows.Close()

	rows.Next()

	var id int64
	err = rows.Scan(&id)
	if err != nil {
		return 0, vanerrors.NewWrap(repo.ErrorScanning, err, vanerrors.EmptyHandler)
	}

	return id, nil
}
