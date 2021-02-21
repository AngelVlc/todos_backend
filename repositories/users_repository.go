package repositories

import (
	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/AngelVlc/todos/models"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
)

type UsersRepository interface {
	GetAll() ([]*models.User, error)
	FindByID(id int32) (*models.User, error)
	FindByName(name string) (*models.User, error)
	Insert(user *models.User) (int32, error)
	Remove(id int32) error
	Update(user *models.User) error
}

type MockedUsersRepository struct {
	mock.Mock
}

func NewMockedUsersRepository() *MockedUsersRepository {
	return &MockedUsersRepository{}
}

func (m *MockedUsersRepository) GetAll() ([]*models.User, error) {
	args := m.Called()
	got := args.Get(0)
	if got == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockedUsersRepository) FindByID(id int32) (*models.User, error) {
	args := m.Called(id)
	got := args.Get(0)
	if got == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockedUsersRepository) FindByName(name string) (*models.User, error) {
	args := m.Called(name)
	got := args.Get(0)
	if got == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockedUsersRepository) Insert(user *models.User) (int32, error) {
	args := m.Called(user)
	got := args.Get(0)
	if got == nil {
		return -1, args.Error(1)
	}
	return args.Get(0).(int32), args.Error(1)
}

func (m *MockedUsersRepository) Remove(id int32) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockedUsersRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

type DefaultUsersRepository struct {
	db *gorm.DB
}

func NewDefaultUsersRepository(db *gorm.DB) *DefaultUsersRepository {
	return &DefaultUsersRepository{db}
}

// FindByID returns a single user from its id
func (r *DefaultUsersRepository) FindByID(id int32) (*models.User, error) {
	foundUser := models.User{}
	err := r.db.Where(models.User{ID: id}).Table("users").First(&foundUser).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error getting user by user id", InternalError: err}
	}

	return &foundUser, nil
}

// FindByName returns a single user by its name
func (r *DefaultUsersRepository) FindByName(name string) (*models.User, error) {
	foundUser := models.User{}
	err := r.db.Where(models.User{Name: name}).Table("users").First(&foundUser).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error getting user by user name", InternalError: err}
	}

	return &foundUser, nil
}

// Insert adds a new user
func (r *DefaultUsersRepository) Insert(user *models.User) (int32, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		return -1, &appErrors.UnexpectedError{Msg: "Error inserting in the database", InternalError: err}
	}

	return user.ID, nil
}

// Remove removes a user
func (r *DefaultUsersRepository) Remove(id int32) error {
	if err := r.db.Where(models.User{ID: id}).Delete(models.User{}).Error; err != nil {
		return &appErrors.UnexpectedError{Msg: "Error deleting user", InternalError: err}
	}
	return nil
}

func (r *DefaultUsersRepository) Update(user *models.User) error {
	if err := r.db.Save(&user).Error; err != nil {
		return &appErrors.UnexpectedError{Msg: "Error updating user", InternalError: err}
	}

	return nil
}
