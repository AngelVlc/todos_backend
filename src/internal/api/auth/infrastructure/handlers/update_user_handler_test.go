//go:build !e2e
// +build !e2e

package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	authDomain "github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"
	"github.com/AngelVlc/todos_backend/src/internal/api/auth/domain/passgen"
	"github.com/AngelVlc/todos_backend/src/internal/api/auth/infrastructure"
	authRepository "github.com/AngelVlc/todos_backend/src/internal/api/auth/infrastructure/repository"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateUserHandler_Validations_Returns_An_ErrorResult_With_A_BadRequestError_If_The_Request_Has_Passwords_But_They_Do_Not_Match(t *testing.T) {
	userName, _ := authDomain.NewUserNameValueObject("wadus")
	h := handler.Handler{
		RequestInput: &infrastructure.UpdateUserInput{Name: userName, Password: "one", ConfirmPassword: "another"},
	}

	request, _ := http.NewRequest(http.MethodPost, "/", nil)

	result := UpdateUserHandler(httptest.NewRecorder(), request, h)

	results.CheckBadRequestErrorResult(t, result, "Passwords don't match")
}

func TestUpdateUserHandler_Returns_An_Error_If_The_Query_To_Find_The_User_Fails(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "1",
		})
		return request
	}

	mockedUsersRepo := authRepository.MockedUsersRepository{}
	mockedPassGen := passgen.MockedPasswordGenerator{}
	userName, _ := authDomain.NewUserNameValueObject("wadus")
	h := handler.Handler{
		UsersRepository: &mockedUsersRepo,
		PassGen:         &mockedPassGen,
		RequestInput:    &infrastructure.UpdateUserInput{Name: userName},
	}

	req := request()
	mockedUsersRepo.On("FindUser", req.Context(), &authDomain.UserRecord{ID: int32(1)}).Return(nil, fmt.Errorf("some error")).Once()

	result := UpdateUserHandler(httptest.NewRecorder(), req, h)

	results.CheckError(t, result, "some error")
	mockedUsersRepo.AssertExpectations(t)
}

func TestUpdateUserHandler_Returns_An_ErrorResult_With_A_BadRequestError_If_Tries_To_Update_The_Admin_User_UserName(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "1",
		})
		return request
	}

	mockedUsersRepo := authRepository.MockedUsersRepository{}
	mockedPassGen := passgen.MockedPasswordGenerator{}
	userName, _ := authDomain.NewUserNameValueObject("newAdmin")
	h := handler.Handler{
		UsersRepository: &mockedUsersRepo,
		PassGen:         &mockedPassGen,
		RequestInput:    &infrastructure.UpdateUserInput{Name: userName},
	}

	req := request()
	foundUser := authDomain.UserRecord{Name: "admin"}
	mockedUsersRepo.On("FindUser", req.Context(), &authDomain.UserRecord{ID: int32(1)}).Return(&foundUser, nil).Once()

	result := UpdateUserHandler(httptest.NewRecorder(), req, h)

	results.CheckBadRequestErrorResult(t, result, "It is not possible to change the admin user name")
	mockedUsersRepo.AssertExpectations(t)
}

func TestUpdateUserHandler_Returns_An_ErrorResult_With_A_BadRequestError_If_Tries_To_Update_The_Admin_User_IsAdmin(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "1",
		})
		return request
	}

	mockedUsersRepo := authRepository.MockedUsersRepository{}
	mockedPassGen := passgen.MockedPasswordGenerator{}
	userName, _ := authDomain.NewUserNameValueObject("admin")
	h := handler.Handler{
		UsersRepository: &mockedUsersRepo,
		PassGen:         &mockedPassGen,
		RequestInput:    &infrastructure.UpdateUserInput{Name: userName, IsAdmin: false},
	}

	req := request()
	foundUser := authDomain.UserRecord{Name: "admin"}
	mockedUsersRepo.On("FindUser", req.Context(), &authDomain.UserRecord{ID: int32(1)}).Return(&foundUser, nil).Once()

	result := UpdateUserHandler(httptest.NewRecorder(), req, h)

	results.CheckBadRequestErrorResult(t, result, "The admin user must be an admin")
	mockedUsersRepo.AssertExpectations(t)
}

func TestUpdateUserHandler_Returns_An_Error_If_The_Query_To_Check_If_A_User_With_The_Same_UserName_Already_Exists(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "1",
		})
		return request
	}

	mockedUsersRepo := authRepository.MockedUsersRepository{}
	mockedPassGen := passgen.MockedPasswordGenerator{}
	userName, _ := authDomain.NewUserNameValueObject("wadusR")
	h := handler.Handler{
		UsersRepository: &mockedUsersRepo,
		PassGen:         &mockedPassGen,
		RequestInput:    &infrastructure.UpdateUserInput{Name: userName, Password: "newPass", ConfirmPassword: "newPass"},
	}

	req := request()
	foundUser := authDomain.UserRecord{Name: "wadus"}
	mockedUsersRepo.On("FindUser", req.Context(), &authDomain.UserRecord{ID: int32(1)}).Return(&foundUser, nil).Once()
	mockedUsersRepo.On("ExistsUser", req.Context(), &authDomain.UserRecord{Name: "wadusR"}).Return(false, fmt.Errorf("some error")).Once()

	result := UpdateUserHandler(httptest.NewRecorder(), req, h)

	results.CheckError(t, result, "Error checking if a user with the same name already exists")
	mockedUsersRepo.AssertExpectations(t)
	mockedPassGen.AssertExpectations(t)
}

func TestUpdateUserHandler_Returns_An_ErrorResult_With_A_BadRequestError_If_The_UserName_Already_Exists(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "1",
		})
		return request
	}

	mockedUsersRepo := authRepository.MockedUsersRepository{}
	mockedPassGen := passgen.MockedPasswordGenerator{}
	userName, _ := authDomain.NewUserNameValueObject("wadusR")
	h := handler.Handler{
		UsersRepository: &mockedUsersRepo,
		PassGen:         &mockedPassGen,
		RequestInput:    &infrastructure.UpdateUserInput{Name: userName, Password: "newPass", ConfirmPassword: "newPass"},
	}

	req := request()
	foundUser := authDomain.UserRecord{Name: "wadus"}
	mockedUsersRepo.On("FindUser", req.Context(), &authDomain.UserRecord{ID: int32(1)}).Return(&foundUser, nil).Once()
	mockedUsersRepo.On("ExistsUser", req.Context(), &authDomain.UserRecord{Name: "wadusR"}).Return(true, nil).Once()

	result := UpdateUserHandler(httptest.NewRecorder(), req, h)

	results.CheckBadRequestErrorResult(t, result, "A user with the same user name already exists")
	mockedUsersRepo.AssertExpectations(t)
	mockedPassGen.AssertExpectations(t)
}

func TestUpdateUserHandler_Returns_An_ErrorResult_With_An_UnexpectedError_If_Generating_The_Password_Fails(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "1",
		})
		return request
	}

	mockedUsersRepo := authRepository.MockedUsersRepository{}
	mockedPassGen := passgen.MockedPasswordGenerator{}
	userName, _ := authDomain.NewUserNameValueObject("wadus")
	h := handler.Handler{
		UsersRepository: &mockedUsersRepo,
		PassGen:         &mockedPassGen,
		RequestInput:    &infrastructure.UpdateUserInput{Name: userName, Password: "newPass", ConfirmPassword: "newPass"},
	}

	req := request()
	foundUser := authDomain.UserRecord{Name: "wadus"}
	mockedUsersRepo.On("FindUser", req.Context(), &authDomain.UserRecord{ID: int32(1)}).Return(&foundUser, nil).Once()
	mockedPassGen.On("GenerateFromPassword", "newPass").Return("", fmt.Errorf("some error")).Once()

	result := UpdateUserHandler(httptest.NewRecorder(), req, h)

	results.CheckUnexpectedErrorResult(t, result, "Error encrypting password")
	mockedUsersRepo.AssertExpectations(t)
	mockedPassGen.AssertExpectations(t)
}

func TestUpdateUserHandler_Returns_An_ErrorResult_With_An_UnexpectedError_If_The_Update_Fails(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "1",
		})
		return request
	}

	mockedUsersRepo := authRepository.MockedUsersRepository{}
	mockedPassGen := passgen.MockedPasswordGenerator{}
	userName, _ := authDomain.NewUserNameValueObject("updated")
	h := handler.Handler{
		UsersRepository: &mockedUsersRepo,
		PassGen:         &mockedPassGen,
		RequestInput:    &infrastructure.UpdateUserInput{Name: userName},
	}

	req := request()
	foundUser := authDomain.UserRecord{Name: "wadus"}
	mockedUsersRepo.On("FindUser", req.Context(), &authDomain.UserRecord{ID: int32(1)}).Return(&foundUser, nil).Once()
	foundUser.Name = "updated"
	mockedUsersRepo.On("Update", req.Context(), &foundUser).Return(fmt.Errorf("some error")).Once()

	result := UpdateUserHandler(httptest.NewRecorder(), req, h)

	results.CheckUnexpectedErrorResult(t, result, "Error updating the user")
	mockedUsersRepo.AssertExpectations(t)
	mockedPassGen.AssertExpectations(t)
}

func TestUpdateUserHandler_Updates_The_UserName(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "1",
		})
		return request
	}

	mockedUsersRepo := authRepository.MockedUsersRepository{}
	mockedPassGen := passgen.MockedPasswordGenerator{}
	userName, _ := authDomain.NewUserNameValueObject("updated")
	h := handler.Handler{
		UsersRepository: &mockedUsersRepo,
		PassGen:         &mockedPassGen,
		RequestInput:    &infrastructure.UpdateUserInput{Name: userName},
	}

	req := request()
	foundUser := authDomain.UserRecord{ID: int32(1), Name: "wadus", IsAdmin: false}
	mockedUsersRepo.On("FindUser", req.Context(), &authDomain.UserRecord{ID: int32(1)}).Return(&foundUser, nil).Once()
	mockedUsersRepo.On("ExistsUser", req.Context(), &authDomain.UserRecord{Name: "updated"}).Return(false, nil).Once()
	foundUser2 := authDomain.UserRecord{ID: int32(1), Name: "wadus", IsAdmin: false}
	foundUser2.Name = "updated"
	mockedUsersRepo.On("Update", req.Context(), &foundUser2).Return(nil).Once()

	result := UpdateUserHandler(httptest.NewRecorder(), req, h)

	okRes := results.CheckOkResult(t, result, http.StatusOK)
	userRes, isOk := okRes.Content.(infrastructure.UserResponse)
	require.Equal(t, true, isOk, "should be a user response")

	assert.Equal(t, int32(1), userRes.ID)
	assert.Equal(t, "updated", userRes.Name)
	assert.False(t, userRes.IsAdmin)

	mockedUsersRepo.AssertExpectations(t)
	mockedPassGen.AssertExpectations(t)
}

func TestUpdateUserHandler_Updates_The_Password(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "1",
		})
		return request
	}

	mockedUsersRepo := authRepository.MockedUsersRepository{}
	mockedPassGen := passgen.MockedPasswordGenerator{}
	userName, _ := authDomain.NewUserNameValueObject("wadus")
	h := handler.Handler{
		UsersRepository: &mockedUsersRepo,
		PassGen:         &mockedPassGen,
		RequestInput:    &infrastructure.UpdateUserInput{Name: userName, Password: "newPass", ConfirmPassword: "newPass"},
	}

	req := request()
	foundUser := authDomain.UserRecord{ID: int32(1), Name: "wadus", IsAdmin: false}
	mockedUsersRepo.On("FindUser", req.Context(), &authDomain.UserRecord{ID: int32(1)}).Return(&foundUser, nil).Once()
	mockedPassGen.On("GenerateFromPassword", "newPass").Return("hassedPass", nil).Once()
	mockedUsersRepo.On("Update", req.Context(), &foundUser).Return(nil).Once()

	result := UpdateUserHandler(httptest.NewRecorder(), req, h)

	okRes := results.CheckOkResult(t, result, http.StatusOK)
	userRes, isOk := okRes.Content.(infrastructure.UserResponse)
	require.Equal(t, true, isOk, "should be a user response")

	assert.Equal(t, int32(1), userRes.ID)
	assert.Equal(t, "wadus", userRes.Name)
	assert.False(t, userRes.IsAdmin)

	mockedUsersRepo.AssertExpectations(t)
	mockedPassGen.AssertExpectations(t)
}

func TestUpdateUserHandler_Updates_The_IsAmin(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "1",
		})
		return request
	}

	mockedUsersRepo := authRepository.MockedUsersRepository{}
	mockedPassGen := passgen.MockedPasswordGenerator{}
	userName, _ := authDomain.NewUserNameValueObject("wadus")
	h := handler.Handler{
		UsersRepository: &mockedUsersRepo,
		PassGen:         &mockedPassGen,
		RequestInput:    &infrastructure.UpdateUserInput{Name: userName, IsAdmin: true},
	}

	req := request()
	foundUser := authDomain.UserRecord{ID: int32(1), Name: "wadus", IsAdmin: false}
	mockedUsersRepo.On("FindUser", req.Context(), &authDomain.UserRecord{ID: int32(1)}).Return(&foundUser, nil).Once()
	foundUser.IsAdmin = true
	mockedUsersRepo.On("Update", req.Context(), &foundUser).Return(nil).Once()

	result := UpdateUserHandler(httptest.NewRecorder(), req, h)

	okRes := results.CheckOkResult(t, result, http.StatusOK)
	userRes, isOk := okRes.Content.(infrastructure.UserResponse)
	require.Equal(t, true, isOk, "should be a user response")

	assert.Equal(t, int32(1), userRes.ID)
	assert.Equal(t, "wadus", userRes.Name)
	assert.True(t, userRes.IsAdmin)

	mockedUsersRepo.AssertExpectations(t)
	mockedPassGen.AssertExpectations(t)
}
