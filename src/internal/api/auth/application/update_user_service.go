package application

import (
	"context"

	"github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"
	"github.com/AngelVlc/todos_backend/src/internal/api/auth/domain/passgen"
	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
)

type UpdateUserService struct {
	usersRepo domain.UsersRepository
	passGen   passgen.PasswordGenerator
}

func NewUpdateUserService(usersRepo domain.UsersRepository, passGen passgen.PasswordGenerator) *UpdateUserService {
	return &UpdateUserService{usersRepo, passGen}
}

func (s *UpdateUserService) UpdateUser(ctx context.Context, userID int32, userName domain.UserNameValueObject, password string, isAdmin bool) (*domain.UserEntity, error) {
	foundUser, err := s.usersRepo.FindUser(ctx, domain.UserRecord{ID: userID})
	if err != nil {
		return nil, err
	}

	entity := foundUser.ToUserEntity()

	if entity.IsTheAdminUser() {
		if userName.String() != "admin" {
			return nil, &appErrors.BadRequestError{Msg: "It is not possible to change the admin user name"}
		}

		if !isAdmin {
			return nil, &appErrors.BadRequestError{Msg: "The admin user must be an admin"}
		}
	}

	if entity.Name.String() != userName.String() {
		if existsUser, err := s.usersRepo.ExistsUser(ctx, domain.UserRecord{Name: userName.String()}); err != nil {
			return nil, &appErrors.UnexpectedError{Msg: "Error checking if a user with the same name already exists", InternalError: err}
		} else if existsUser {
			return nil, &appErrors.BadRequestError{Msg: "A user with the same user name already exists", InternalError: nil}
		}
	}

	if len(password) > 0 {
		hasshedPass, err := s.passGen.GenerateFromPassword(string(password))
		if err != nil {
			return nil, &appErrors.UnexpectedError{Msg: "Error encrypting password", InternalError: err}
		}

		foundUser.PasswordHash = hasshedPass
	}

	foundUser.Name = userName.String()
	foundUser.IsAdmin = isAdmin

	err = s.usersRepo.Update(ctx, foundUser)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error updating the user", InternalError: err}
	}

	return foundUser.ToUserEntity(), nil
}
