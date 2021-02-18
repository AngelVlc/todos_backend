package repositories

import (
	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/AngelVlc/todos/models"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
)

type UsersRepository interface {
	FindByID(id int32) (*models.User, error)
}

type MockedUsersRepository struct {
	mock.Mock
}

func NewMockedUsersRepository() *MockedUsersRepository {
	return &MockedUsersRepository{}
}

func (m *MockedUsersRepository) FindByID(id int32) (*models.User, error) {
	args := m.Called(id)
	got := args.Get(0)
	if got == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

type DefaultUsersRepository struct {
	db *gorm.DB
}

func NewDefaultUsersRepository(db *gorm.DB) *DefaultUsersRepository {
	return &DefaultUsersRepository{db}
}

// FindByID returns a single user from its id
func (s *DefaultUsersRepository) FindByID(id int32) (*models.User, error) {
	foundUser := models.User{}
	err := s.db.Where(models.User{ID: id}).Table("users").First(&foundUser).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error getting user by user id", InternalError: err}
	}

	return &foundUser, nil
}
