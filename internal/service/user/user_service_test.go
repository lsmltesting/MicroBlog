package user

import (
	"context"
	"testing"

	"github.com/lsmltesting/MicroBlog/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Save(ctx context.Context, user *models.User) (int, error) {
	args := m.Called(ctx, user)
	return args.Int(0), args.Error(1)
}

func (m *MockUserRepository) FindUserByID(ctx context.Context, ID int) (*models.User, error) {
	args := m.Called(ctx, ID)
	return args.Get(0).(*models.User), args.Error(1)
}

// func (m *MockUserRepository) UpdatePostHistory(userID int, postID int) error {
// 	args := m.Called(userID, postID)
// 	return args.Error(0)
// }

func TestUserService_CreateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	testID := 42
	// Set up mock for Save()
	mockRepo.On("Save", ctx, mock.MatchedBy(func(user *models.User) bool {
		return user.Username == "testingUser" && user.Email == "test@gmail.com"
	})).Return(testID, nil)

	userId, err := service.CreateUser(ctx, "testingUser", "test@gmail.com", "sdfmdsfmsdkfm123")

	assert.NoError(t, err)
	assert.Equal(t, testID, userId)

	// Check that method Save is called only 1 time
	mockRepo.AssertExpectations(t)
}

// Check that service work correctly with errors
func TestUserService_CreateUser_Error(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	//Set up mock for Save() that Save() return error
	mockRepo.On("Save", ctx, mock.Anything).Return(0, assert.AnError)

	userId, err := service.CreateUser(ctx, "testgingUser", "test@gmail.com", "password123")

	assert.Error(t, err)
	assert.Equal(t, 0, userId)

	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUserByID_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	userID := 33

	expectedUser := &models.User{
		Username: "test_dev_user",
		Email:    "test_dev_user@gmail.com",
		Password: "password_qwerty",
		ID:       userID,
	}

	//Set up mock for FindUserByID()
	mockRepo.On("FindUserByID", ctx, userID).Return(expectedUser, nil)

	user, err := service.GetUserByID(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, user, expectedUser)

	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUserByID_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	userID := 1231

	mockRepo.On("FindUserByID", ctx, mock.Anything).Return((*models.User)(nil), assert.AnError)

	user, err := service.GetUserByID(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, user)

	mockRepo.AssertExpectations(t)
}
