//go:build !e2e
// +build !e2e

package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"
	"github.com/AngelVlc/todos_backend/src/internal/api/auth/domain/passgen"
	"github.com/AngelVlc/todos_backend/src/internal/api/auth/infrastructure"
	authRepository "github.com/AngelVlc/todos_backend/src/internal/api/auth/infrastructure/repository"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateUserHandler_Validations_Returns_An_ErrorResult_With_A_BadRequestError_If_The_CreateUserInput_Passwords_Do_Not_Match(t *testing.T) {
	userName, _ := domain.NewUserNameValueObject("wadus")
	userPassword, _ := domain.NewUserPasswordValueObject("pass")
	h := handler.Handler{
		RequestInput: &infrastructure.CreateUserInput{Name: userName, Password: userPassword, ConfirmPassword: "othePass"},
	}

	request, _ := http.NewRequest(http.MethodPost, "/", nil)

	result := CreateUserHandler(httptest.NewRecorder(), request, h)

	results.CheckBadRequestErrorResult(t, result, "Passwords don't match")
}

func TestCreateUserHandler_Returns_An_Error_If_The_Query_To_Check_If_The_User_Exists_Fails(t *testing.T) {
	mockedUsersRepo := authRepository.MockedUsersRepository{}
	mockedPassGen := passgen.MockedPasswordGenerator{}
	userName, _ := domain.NewUserNameValueObject("wadus")
	userPassword, _ := domain.NewUserPasswordValueObject("pass")
	h := handler.Handler{
		UsersRepository: &mockedUsersRepo,
		PassGen:         &mockedPassGen,
		RequestInput:    &infrastructure.CreateUserInput{Name: userName, Password: userPassword, ConfirmPassword: "pass", IsAdmin: true},
	}

	request, _ := http.NewRequest(http.MethodPost, "/", nil)
	mockedUsersRepo.On("ExistsUser", request.Context(), &domain.UserRecord{Name: "wadus"}).Return(false, fmt.Errorf("some error")).Once()

	result := CreateUserHandler(httptest.NewRecorder(), request, h)

	results.CheckError(t, result, "Error checking if a user with the same name already exists")
	mockedUsersRepo.AssertExpectations(t)
}

func TestCreateUserHandler_Returns_A_BadRequest_Error_If_A_User_With_The_Same_Name_Already_Exist(t *testing.T) {
	mockedUsersRepo := authRepository.MockedUsersRepository{}
	mockedPassGen := passgen.MockedPasswordGenerator{}
	userName, _ := domain.NewUserNameValueObject("wadus")
	userPassword, _ := domain.NewUserPasswordValueObject("pass")
	h := handler.Handler{
		UsersRepository: &mockedUsersRepo,
		PassGen:         &mockedPassGen,
		RequestInput:    &infrastructure.CreateUserInput{Name: userName, Password: userPassword, ConfirmPassword: "pass", IsAdmin: true},
	}

	request, _ := http.NewRequest(http.MethodPost, "/", nil)
	mockedUsersRepo.On("ExistsUser", request.Context(), &domain.UserRecord{Name: "wadus"}).Return(true, nil).Once()

	result := CreateUserHandler(httptest.NewRecorder(), request, h)

	results.CheckBadRequestErrorResult(t, result, "A user with the same user name already exists")
	mockedUsersRepo.AssertExpectations(t)
}

func TestCreateUserHandler_Returns_An_ErrorResult_With_An_UnexpectedError_If_The_User_Does_Not_Exist_But_Generating_The_Password_Fails(t *testing.T) {
	mockedUsersRepo := authRepository.MockedUsersRepository{}
	mockedPassGen := passgen.MockedPasswordGenerator{}
	userName, _ := domain.NewUserNameValueObject("wadus")
	userPassword, _ := domain.NewUserPasswordValueObject("pass")
	h := handler.Handler{
		UsersRepository: &mockedUsersRepo,
		PassGen:         &mockedPassGen,
		RequestInput:    &infrastructure.CreateUserInput{Name: userName, Password: userPassword, ConfirmPassword: "pass", IsAdmin: true},
	}

	request, _ := http.NewRequest(http.MethodPost, "/", nil)
	mockedUsersRepo.On("ExistsUser", request.Context(), &domain.UserRecord{Name: "wadus"}).Return(false, nil).Once()
	mockedPassGen.On("GenerateFromPassword", "pass").Return("", fmt.Errorf("some error")).Once()

	result := CreateUserHandler(httptest.NewRecorder(), request, h)

	results.CheckUnexpectedErrorResult(t, result, "Error encrypting password")
	mockedUsersRepo.AssertExpectations(t)
	mockedPassGen.AssertExpectations(t)
}

func TestCreateUserHandler_Returns_An_ErrorResult_With_An_UnexpectedError_If_The_User_Does_Not_Exist_But_Creating_The_User_Fails(t *testing.T) {
	mockedUsersRepo := authRepository.MockedUsersRepository{}
	mockedPassGen := passgen.MockedPasswordGenerator{}
	userName, _ := domain.NewUserNameValueObject("wadus")
	userPassword, _ := domain.NewUserPasswordValueObject("pass")
	h := handler.Handler{
		UsersRepository: &mockedUsersRepo,
		PassGen:         &mockedPassGen,
		RequestInput:    &infrastructure.CreateUserInput{Name: userName, Password: userPassword, ConfirmPassword: "pass", IsAdmin: true},
	}

	request, _ := http.NewRequest(http.MethodPost, "/", nil)
	mockedUsersRepo.On("ExistsUser", request.Context(), &domain.UserRecord{Name: "wadus"}).Return(false, nil).Once()
	hassedPass := "hassed"
	mockedPassGen.On("GenerateFromPassword", "pass").Return(hassedPass, nil).Once()
	user := domain.UserRecord{Name: "wadus", PasswordHash: hassedPass, IsAdmin: true}
	mockedUsersRepo.On("Create", request.Context(), &user).Return(fmt.Errorf("some error")).Once()

	result := CreateUserHandler(httptest.NewRecorder(), request, h)

	results.CheckUnexpectedErrorResult(t, result, "Error creating the user")
	mockedUsersRepo.AssertExpectations(t)
	mockedPassGen.AssertExpectations(t)
}

func TestCreateUserHandler_Creates_The_User(t *testing.T) {
	mockedUsersRepo := authRepository.MockedUsersRepository{}
	mockedPassGen := passgen.MockedPasswordGenerator{}
	userName, _ := domain.NewUserNameValueObject("wadus")
	userPassword, _ := domain.NewUserPasswordValueObject("pass")
	h := handler.Handler{
		UsersRepository: &mockedUsersRepo,
		PassGen:         &mockedPassGen,
		RequestInput:    &infrastructure.CreateUserInput{Name: userName, Password: userPassword, ConfirmPassword: "pass", IsAdmin: true},
	}

	request, _ := http.NewRequest(http.MethodPost, "/", nil)
	mockedUsersRepo.On("ExistsUser", request.Context(), &domain.UserRecord{Name: "wadus"}).Return(false, nil).Once()
	hassedPass := "hassed"
	mockedPassGen.On("GenerateFromPassword", "pass").Return(hassedPass, nil).Once()
	user := domain.UserRecord{Name: "wadus", PasswordHash: hassedPass, IsAdmin: true}
	mockedUsersRepo.On("Create", request.Context(), &user).Return(nil).Once().Run(func(args mock.Arguments) {
		arg := args.Get(1).(*domain.UserRecord)
		*arg = domain.UserRecord{ID: int32(1)}
	})

	result := CreateUserHandler(httptest.NewRecorder(), request, h)

	okRes := results.CheckOkResult(t, result, http.StatusCreated)
	res, isOk := okRes.Content.(*domain.UserRecord)
	require.Equal(t, true, isOk, "should be a UserRecord")
	assert.Equal(t, int32(1), res.ID)

	mockedUsersRepo.AssertExpectations(t)
	mockedPassGen.AssertExpectations(t)
}
