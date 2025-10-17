package post

import (
	"context"
	"testing"

	"github.com/lsmltesting/MicroBlog/internal/models"
	"github.com/lsmltesting/MicroBlog/internal/service/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPostRepository struct {
	mock.Mock
}

type MockUserRepository struct {
	mock.Mock
}

// Posts methods
func (mockPost *MockPostRepository) Save(ctx context.Context, post *models.Post) (int, error) {
	args := mockPost.Called(ctx, post)
	return args.Int(0), args.Error(1)
}

func (mockPost *MockPostRepository) FindPostByID(ctx context.Context, postID int) (*models.Post, error) {
	args := mockPost.Called(ctx, postID)
	return args.Get(0).(*models.Post), args.Error(1)
}

func (mockPost *MockPostRepository) GetAllPosts(ctx context.Context) (map[int]*models.Post, error) {
	args := mockPost.Called(ctx)
	return args.Get(0).(map[int]*models.Post), args.Error(1)
}

// Users methods
func (mockUser *MockUserRepository) Save(ctx context.Context, user *models.User) (int, error) {
	args := mockUser.Called(ctx, user)
	return args.Int(0), args.Error(1)
}

func (mockUser *MockUserRepository) FindUserByID(ctx context.Context, ID int) (*models.User, error) {
	args := mockUser.Called(ctx, ID)
	return args.Get(0).(*models.User), args.Error(1)
}

func TestPostService_Save_Success(t *testing.T) {
	mockRepoPost := new(MockPostRepository)
	mockRepoUser := new(MockUserRepository)
	ctx := context.Background()

	userService := user.NewUserService(mockRepoUser)
	postService := NewPostService(mockRepoPost, userService)

	postID := 2
	textPost := "smth new test text for post is already written"

	userID := 1

	testUser := &models.User{
		Username: "test_user_name_for_post",
		Email:    "test_mail@gmail.com",
		Password: "test_password_qwerty_for_post",
		ID:       userID,
	}

	mockRepoUser.On("FindUserByID", ctx, userID).Return(testUser, nil)
	mockRepoPost.On("Save", ctx, mock.MatchedBy(func(post *models.Post) bool {
		return (post.Text == textPost &&
			post.UserID == userID)
	})).Return(postID, nil)

	createdPostId, err := postService.CreatePost(ctx, userID, textPost)

	assert.NoError(t, err)
	assert.Equal(t, postID, createdPostId)

	mockRepoPost.AssertExpectations(t)
	mockRepoUser.AssertExpectations(t)
}

func TestPostService_Save_Error(t *testing.T) {
	mockPostRepository := new(MockPostRepository)
	mockUserRepository := new(MockUserRepository)
	ctx := context.Background()

	userService := user.NewUserService(mockUserRepository)
	postService := NewPostService(mockPostRepository, userService)

	userID := 10
	postText := "test post text"

	testUser := &models.User{
		Username: "test_user_name_for_post",
		Email:    "test_mail@gmail.com",
		Password: "test_password_qwerty_for_post",
		ID:       userID,
	}

	mockUserRepository.On("FindUserByID", ctx, userID).Return(testUser, nil)
	mockPostRepository.On("Save", ctx, mock.Anything).Return(0, assert.AnError)

	postID, err := postService.CreatePost(ctx, userID, postText)

	assert.Error(t, err)
	assert.Equal(t, 0, postID)

	mockPostRepository.AssertExpectations(t)
	mockUserRepository.AssertExpectations(t)
}

func TestPostService_FindPostByID_Success(t *testing.T) {
	mockPostRepository := new(MockPostRepository)
	mockUserRepository := new(MockUserRepository)
	ctx := context.Background()

	userService := user.NewUserService(mockUserRepository)
	postService := NewPostService(mockPostRepository, userService)

	testPostID := 13
	testUserID := 7
	testPostText := "random text for post test"

	expectedPost := &models.Post{
		Text:   testPostText,
		UserID: testUserID,
	}

	mockPostRepository.On("FindPostByID", ctx, testPostID).Return(expectedPost, nil)

	post, err := postService.GetPostByID(ctx, testPostID)

	assert.NoError(t, err)
	assert.Equal(t, post, expectedPost)

	mockPostRepository.AssertExpectations(t)
	mockUserRepository.AssertExpectations(t)
}

func TestPostService_FindPostByID_NotFound(t *testing.T) {
	mockPostRepository := new(MockPostRepository)
	mockUserRepository := new(MockUserRepository)
	ctx := context.Background()

	userService := user.NewUserService(mockUserRepository)
	postService := NewPostService(mockPostRepository, userService)

	testPostID := 13

	mockPostRepository.On("FindPostByID", ctx, mock.Anything).Return((*models.Post)(nil), assert.AnError)

	post, err := postService.GetPostByID(ctx, testPostID)

	assert.Error(t, err)
	assert.Nil(t, post)

	mockPostRepository.AssertExpectations(t)
	mockUserRepository.AssertExpectations(t)
}

func TestPostService_GetAllPosts_Success(t *testing.T) {
	mockPostRepository := new(MockPostRepository)
	mockUserRepository := new(MockUserRepository)
	ctx := context.Background()

	userService := user.NewUserService(mockUserRepository)
	postService := NewPostService(mockPostRepository, userService)

	expectedPosts := map[int]*models.Post{
		1: {
			ID:     10,
			Text:   "first test text",
			UserID: 9,
		},
		2: {
			ID:     33,
			Text:   "second test text",
			UserID: 123,
		},
	}

	mockPostRepository.On("GetAllPosts", ctx).Return(expectedPosts, nil)

	posts, err := postService.GetAllPosts(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedPosts, posts)

	mockPostRepository.AssertExpectations(t)
	mockUserRepository.AssertExpectations(t)
}

func TestPostService_GetAllPosts_Error(t *testing.T) {
	mockPostRepository := new(MockPostRepository)
	mockUserRepository := new(MockUserRepository)
	ctx := context.Background()

	userService := user.NewUserService(mockUserRepository)
	postService := NewPostService(mockPostRepository, userService)

	mockPostRepository.On("GetAllPosts", ctx).Return((map[int]*models.Post)(nil), assert.AnError)

	posts, err := postService.GetAllPosts(ctx)

	assert.Error(t, err)
	assert.Nil(t, posts)

	mockPostRepository.AssertExpectations(t)
	mockUserRepository.AssertExpectations(t)
}
