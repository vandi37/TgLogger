package user_repo

import (
	"context"

	"github.com/vandi37/TgLogger/internal/database/db"
	"github.com/vandi37/TgLogger/internal/database/repo"
	"github.com/vandi37/vanerrors"
)

type UserRepo struct {
	db *db.DB
}

func New(db *db.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) NewUser(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `insert into users (id) values ($1) on conflict (Id) do nothing;`, id)
	if err != nil {
		return vanerrors.NewWrap(repo.ErrorExecuting, err, vanerrors.EmptyHandler)
	}

	return repo.ReturnByRes(res, repo.Equals(1))

}
