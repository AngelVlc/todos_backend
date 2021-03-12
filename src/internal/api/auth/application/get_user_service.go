package application

import (
	"github.com/AngelVlc/todos/internal/api/auth/domain"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/domain/errors"
)

type GetUserService struct {
	repo domain.AuthRepository
}

func NewGetUserService(repo domain.AuthRepository) *GetUserService {
	return &GetUserService{repo}
}

func (s *GetUserService) GetUser(userID int32) (*domain.User, error) {
	foundUser, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error getting user by id", InternalError: err}
	}

	if foundUser == nil {
		return nil, &appErrors.BadRequestError{Msg: "The user does not exist"}
	}

	return foundUser, nil
}
