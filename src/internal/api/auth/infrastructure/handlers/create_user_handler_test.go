//go:build !e2e
// +build !e2e

package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestCreateUserHandler_Validations_Returns_An_ErrorResult_With_A_BadRequestError_If_The_Request_Does_Not_Have_Body(t *testing.T) {
	h := handler.Handler{}
	request, _ := http.NewRequest(http.MethodGet, "/", nil)

	result := CreateUserHandler(httptest.NewRecorder(), request, h)

	results.CheckBadRequestErrorResult(t, result, "Invalid body")
}

func TestCreateUserHandler_Validations_Returns_A_BadRequesError_If_The_Body_Is_Not_A_CreateUserRequest(t *testing.T) {
	h := handler.Handler{}
	request, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader("wadus"))

	result := CreateUserHandler(httptest.NewRecorder(), request, h)

	results.CheckBadRequestErrorResult(t, result, "Invalid body")
}

func TestCreateUserHandler_Validations_Returns_A_BadRequest_Error_If_Its_A_Create_Admin_Request_With_Not_Admin_UserName_(t *testing.T) {
	h := handler.Handler{}
	createReq := createUserRequest{Name: "another"}
	body, _ := json.Marshal(createReq)
	request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
	request.RequestURI = "/auth/createadmin"

	result := CreateUserHandler(httptest.NewRecorder(), request, h)

	results.CheckBadRequestErrorResult(t, result, "/auth/createadmin only can be used to create the admin user")
}

func TestCreateUserHandler_Validations_Returns_An_ErrorResult_With_A_BadRequestError_If_The_CreateUserRequest_Has_An_Empty_UserName(t *testing.T) {
	h := handler.Handler{}
	createReq := createUserRequest{Name: ""}
	body, _ := json.Marshal(createReq)
	request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

	result := CreateUserHandler(httptest.NewRecorder(), request, h)

	results.CheckBadRequestErrorResult(t, result, "The user name can not be empty")
}

func TestCreateUserHandler_Validations_Returns_An_ErrorResult_With_A_BadRequestError_If_The_CreateUserRequest_Does_Not_Have_Password(t *testing.T) {
	h := handler.Handler{}

	createReq := createUserRequest{Name: "Wadus", Password: ""}
	body, _ := json.Marshal(createReq)
	request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

	result := CreateUserHandler(httptest.NewRecorder(), request, h)

	results.CheckBadRequestErrorResult(t, result, "Password can not be empty")
}

func TestCreateUserHandler_Validations_Returns_An_ErrorResult_With_A_BadRequestError_If_The_CreateUserRequest_Passwords_Do_Not_Match(t *testing.T) {
	h := handler.Handler{}

	createReq := createUserRequest{Name: "Wadus", Password: "pass", ConfirmPassword: "othePass"}
	body, _ := json.Marshal(createReq)
	request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

	result := CreateUserHandler(httptest.NewRecorder(), request, h)

	results.CheckBadRequestErrorResult(t, result, "Passwords don't match")
}

func TestCreateUserHandler_Returns_An_Error_If_The_Query_To_Check_If_The_User_Exists_Fails(t *testing.T) {
	mockedUsersRepo := authRepository.MockedUsersRepository{}
	mockedPassGen := passgen.MockedPasswordGenerator{}
	h := handler.Handler{UsersRepository: &mockedUsersRepo, PassGen: &mockedPassGen}

	createReq := createUserRequest{Name: "wadus", Password: "pass", ConfirmPassword: "pass", IsAdmin: true}
	body, _ := json.Marshal(createReq)

	request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
	mockedUsersRepo.On("ExistsUser", request.Context(), &domain.UserEntity{Name: domain.UserNameValueObject("wadus")}).Return(false, fmt.Errorf("some error")).Once()

	result := CreateUserHandler(httptest.NewRecorder(), request, h)

	results.CheckError(t, result, "some error")
	mockedUsersRepo.AssertExpectations(t)
}

func TestCreateUserHandler_Returns_A_BadRequest_Error_If_A_User_With_The_Same_Name_Already_Exist(t *testing.T) {
	mockedUsersRepo := authRepository.MockedUsersRepository{}
	mockedPassGen := passgen.MockedPasswordGenerator{}
	h := handler.Handler{UsersRepository: &mockedUsersRepo, PassGen: &mockedPassGen}

	createReq := createUserRequest{Name: "wadus", Password: "pass", ConfirmPassword: "pass", IsAdmin: true}
	body, _ := json.Marshal(createReq)

	request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
	mockedUsersRepo.On("ExistsUser", request.Context(), &domain.UserEntity{Name: domain.UserNameValueObject("wadus")}).Return(true, nil).Once()

	result := CreateUserHandler(httptest.NewRecorder(), request, h)

	results.CheckBadRequestErrorResult(t, result, "A user with the same user name already exists")
	mockedUsersRepo.AssertExpectations(t)
}

func TestCreateUserHandler_Returns_An_ErrorResult_With_An_UnexpectedError_If_The_User_Does_Not_Exist_But_Generating_The_Password_Fails(t *testing.T) {
	mockedUsersRepo := authRepository.MockedUsersRepository{}
	mockedPassGen := passgen.MockedPasswordGenerator{}
	h := handler.Handler{UsersRepository: &mockedUsersRepo, PassGen: &mockedPassGen}

	createReq := createUserRequest{Name: "wadus", Password: "pass", ConfirmPassword: "pass", IsAdmin: true}
	body, _ := json.Marshal(createReq)

	request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
	mockedUsersRepo.On("ExistsUser", request.Context(), &domain.UserEntity{Name: domain.UserNameValueObject("wadus")}).Return(false, nil).Once()
	mockedPassGen.On("GenerateFromPassword", "pass").Return("", fmt.Errorf("some error")).Once()

	result := CreateUserHandler(httptest.NewRecorder(), request, h)

	results.CheckUnexpectedErrorResult(t, result, "Error encrypting password")
	mockedUsersRepo.AssertExpectations(t)
	mockedPassGen.AssertExpectations(t)
}

func TestCreateUserHandler_Returns_An_ErrorResult_With_An_UnexpectedError_If_The_User_Does_Not_Exist_But_Creating_The_User_Fails(t *testing.T) {
	mockedUsersRepo := authRepository.MockedUsersRepository{}
	mockedPassGen := passgen.MockedPasswordGenerator{}
	h := handler.Handler{UsersRepository: &mockedUsersRepo, PassGen: &mockedPassGen}

	createReq := createUserRequest{Name: "wadus", Password: "pass", ConfirmPassword: "pass", IsAdmin: true}
	body, _ := json.Marshal(createReq)

	request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
	mockedUsersRepo.On("ExistsUser", request.Context(), &domain.UserEntity{Name: domain.UserNameValueObject("wadus")}).Return(false, nil).Once()
	hassedPass := "hassed"
	mockedPassGen.On("GenerateFromPassword", "pass").Return(hassedPass, nil).Once()
	user := domain.UserEntity{Name: domain.UserNameValueObject("wadus"), PasswordHash: hassedPass, IsAdmin: true}
	mockedUsersRepo.On("Create", request.Context(), &user).Return(fmt.Errorf("some error")).Once()

	result := CreateUserHandler(httptest.NewRecorder(), request, h)

	results.CheckUnexpectedErrorResult(t, result, "Error creating the user")
	mockedUsersRepo.AssertExpectations(t)
	mockedPassGen.AssertExpectations(t)
}

func TestCreateUserHandler_Creates_The_User(t *testing.T) {
	mockedUsersRepo := authRepository.MockedUsersRepository{}
	mockedPassGen := passgen.MockedPasswordGenerator{}
	h := handler.Handler{UsersRepository: &mockedUsersRepo, PassGen: &mockedPassGen}

	createReq := createUserRequest{Name: "wadus", Password: "pass", ConfirmPassword: "pass", IsAdmin: true}
	body, _ := json.Marshal(createReq)

	request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
	mockedUsersRepo.On("ExistsUser", request.Context(), &domain.UserEntity{Name: domain.UserNameValueObject("wadus")}).Return(false, nil).Once()
	hassedPass := "hassed"
	mockedPassGen.On("GenerateFromPassword", "pass").Return(hassedPass, nil).Once()
	user := domain.UserEntity{Name: domain.UserNameValueObject("wadus"), PasswordHash: hassedPass, IsAdmin: true}
	mockedUsersRepo.On("Create", request.Context(), &user).Return(nil).Once().Run(func(args mock.Arguments) {
		arg := args.Get(1).(*domain.UserEntity)
		*arg = domain.UserEntity{ID: int32(1)}
	})

	result := CreateUserHandler(httptest.NewRecorder(), request, h)

	okRes := results.CheckOkResult(t, result, http.StatusCreated)
	res, isOk := okRes.Content.(infrastructure.UserResponse)
	require.Equal(t, true, isOk, "should be a UserResponse")
	assert.Equal(t, int32(1), res.ID)

	mockedUsersRepo.AssertExpectations(t)
	mockedPassGen.AssertExpectations(t)
}
