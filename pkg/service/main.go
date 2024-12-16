package service

import (
	"github.com/vandi37/TgLogger/pkg/db"
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

// A service for working with database
type Service struct {
	db     *db.DB
	logger *logger.Logger
}

// Creates a new service
func New(db *db.DB, logger *logger.Logger) *Service {
	return &Service{db, logger}
}

func (s *Service) NewUser(id int64) error {
	return s.db.NewUser(id)
}

// Creates a new token
func (s *Service) AddToken(id int64) (string, error) {
	var err error

	var token string

	for exists := true; exists; {
		token = tokens.NewToken()
		exists, err = s.db.CheckToken(token)
		if err != nil {
			return "", err
		}
	}

	err = s.db.NewToken(token, id)
	if err != nil {
		return "", err
	}

	s.logger.Printf("user %d add token %s", id, token)
	return token, nil
}

// Deletes token
func (s *Service) DeleteToken(token string, id int64) error {
	ok := tokens.Ok(token)
	if !ok {
		return vanerrors.NewSimple(InvalidToken)
	}

	ok, err := s.db.CheckToken(token)
	if err != nil {
		return err
	}
	if !ok {
		return vanerrors.NewSimple(TokenNotExist)
	}

	usr, err := s.db.GetOwner(token)
	if err != nil {
		return err
	}
	if usr != id {
		return vanerrors.NewSimple(NotAllowed)
	}

	return s.db.DeleteToken(token)
}

// Gets all tokens of user
func (s *Service) GetTokens(id int64) ([]string, error) {
	return s.db.SelectTokens(id)
}

// Checks token existence
func (s *Service) CheckToken(token string) (bool, error) {
	return s.db.CheckToken(token)
}
