package service

import (
	"context"

	"github.com/vandi37/TgLogger/internal/database/db"
	"github.com/vandi37/TgLogger/internal/database/repo/token_repo"
	"github.com/vandi37/TgLogger/internal/database/repo/user_repo"
	"github.com/vandi37/TgLogger/pkg/logger"
	"github.com/vandi37/TgLogger/pkg/tokens"
	"github.com/vandi37/vanerrors"
)

// Errors
const (
	InvalidToken  = "invalid token"
	TokenNotExist = "token not exist"
	NotAllowed    = "deleting not allowed"
)

type Service struct {
	user_repo  *user_repo.UserRepo
	token_repo *token_repo.TokenRepo
	logger     *logger.Logger
}

func New(db *db.DB, logger *logger.Logger) *Service {
	return &Service{user_repo.New(db), token_repo.New(db), logger}
}

func (s *Service) NewUser(ctx context.Context, id int64) error {
	return s.user_repo.NewUser(ctx, id)
}

func (s *Service) AddToken(ctx context.Context, id int64) (string, error) {
	var err error
	var token string

	for exists := true; exists; {
		token = tokens.NewToken()
		exists, err = s.token_repo.Exist(ctx, token)
		if err != nil {
			return "", err
		}
	}

	err = s.token_repo.New(ctx, token, id)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) DeleteToken(ctx context.Context, token string, id int64) error {
	ok := tokens.Ok(token)
	if !ok {
		return vanerrors.NewSimple(InvalidToken)
	}

	ok, err := s.token_repo.Exist(ctx, token)
	if err != nil {
		return err
	}
	if !ok {
		return vanerrors.NewSimple(TokenNotExist)
	}

	usr, err := s.token_repo.GetOwner(ctx, token)
	if err != nil {
		return err
	}
	if usr != id {
		return vanerrors.NewSimple(NotAllowed)
	}

	return s.token_repo.Delete(ctx, token)
}

// Gets all tokens of user
func (s *Service) GetTokens(ctx context.Context, id int64) ([]string, error) {
	return s.token_repo.Select(ctx, id)
}

func (s *Service) CheckToken(ctx context.Context, token string) (bool, error) {
	return s.token_repo.Exist(ctx, token)
}
