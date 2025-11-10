package user

import (
	"context"
	"strconv"

	"github.com/lsmltesting/MicroBlog/internal/logger"
	"github.com/lsmltesting/MicroBlog/services/api/internal/models"
)

type userServiceDecorator struct {
	userService UserService
	lg          logger.Logger
}

func NewUserServiceDecorator(userService UserService, lg logger.Logger) UserService {
	return &userServiceDecorator{
		userService: userService,
		lg:          lg,
	}
}

func (u *userServiceDecorator) Create(ctx context.Context, username string, email string, password string) (int, error) {
	u.lg.AddLog(
		logger.LevelInfo,
		logger.SourceService,
		map[string]string{
			"username": username,
			"email":    email,
			"password": password,
		},
		"CreateUser from like service is called",
	)

	userID, err := u.userService.Create(ctx, username, email, password)

	if err != nil {
		u.lg.AddLog(
			logger.LevelError,
			logger.SourceService,
			map[string]string{
				"username": username,
				"email":    email,
				"password": password,
			},
			err.Error(),
		)
	} else {
		u.lg.AddLog(
			logger.LevelInfo,
			logger.SourceService,
			map[string]string{
				"username": username,
				"email":    email,
				"password": password,
			},
			"CreateUser called successfully",
		)
	}

	return userID, err
}
func (u *userServiceDecorator) FindByID(ctx context.Context, ID int) (*models.User, error) {
	u.lg.AddLog(
		logger.LevelInfo,
		logger.SourceService,
		map[string]string{
			"ID": strconv.Itoa(ID),
		},
		"GetUserByID from like service is called",
	)

	userModel, err := u.userService.FindByID(ctx, ID)

	if err != nil {
		u.lg.AddLog(
			logger.LevelError,
			logger.SourceService,
			map[string]string{
				"ID": strconv.Itoa(ID),
			},
			err.Error(),
		)
	} else {
		u.lg.AddLog(
			logger.LevelInfo,
			logger.SourceService,
			map[string]string{
				"ID": strconv.Itoa(ID),
			},
			"GetUserByID called successfully",
		)
	}

	return userModel, err
}
