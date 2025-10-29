package like

import (
	"context"
	"strconv"

	"github.com/lsmltesting/MicroBlog/internal/logger"
	"github.com/lsmltesting/MicroBlog/internal/models"
)

type likeServiceDecorator struct {
	likeService LikeService
	lg          logger.Logger
}

func NewLikeServiceDecorator(likeService LikeService, lg logger.Logger) LikeService {
	return &likeServiceDecorator{
		likeService: likeService,
		lg:          lg,
	}
}

func (l *likeServiceDecorator) Create(ctx context.Context, userID int, postID int) (int, error) {
	l.lg.AddLog(
		logger.LevelInfo,
		logger.SourceService,
		map[string]string{
			"userID": strconv.Itoa(userID),
			"postID": strconv.Itoa(postID),
		},
		"CreateLike from like service is called",
	)

	likeID, err := l.likeService.Create(ctx, userID, postID)

	if err != nil {
		l.lg.AddLog(
			logger.LevelError,
			logger.SourceService,
			map[string]string{
				"userID": strconv.Itoa(userID),
				"postID": strconv.Itoa(postID),
			},
			err.Error(),
		)
	} else {
		l.lg.AddLog(
			logger.LevelInfo,
			logger.SourceService,
			map[string]string{
				"userID": strconv.Itoa(userID),
				"postID": strconv.Itoa(postID),
			},
			"CreateLike called successfully",
		)
	}

	return likeID, err
}

func (l *likeServiceDecorator) FindById(ctx context.Context, likeID int) (*models.Like, error) {
	likeModel, err := l.likeService.FindById(ctx, likeID)

	l.lg.AddLog(
		logger.LevelInfo,
		logger.SourceService,
		map[string]string{
			"likeID": strconv.Itoa(likeID),
		},
		"GetLikeById from like service is called",
	)

	if err != nil {
		l.lg.AddLog(
			logger.LevelError,
			logger.SourceService,
			map[string]string{
				"likeID": strconv.Itoa(likeID),
			},
			err.Error(),
		)
	} else {
		l.lg.AddLog(
			logger.LevelInfo,
			logger.SourceService,
			map[string]string{
				"likeID": strconv.Itoa(likeID),
			},
			"GetLikeById called successfully",
		)
	}

	return likeModel, err
}

func (l *likeServiceDecorator) GetAll(ctx context.Context) (map[int]*models.Like, error) {
	mapLike, err := l.likeService.GetAll(ctx)

	l.lg.AddLog(
		logger.LevelInfo,
		logger.SourceService,
		make(map[string]string),
		"GetAllLikes from like service is called",
	)

	if err != nil {
		l.lg.AddLog(
			logger.LevelError,
			logger.SourceService,
			make(map[string]string),
			err.Error(),
		)
	} else {
		l.lg.AddLog(
			logger.LevelInfo,
			logger.SourceService,
			make(map[string]string),
			"GetAllLikes called successfully",
		)
	}

	return mapLike, err
}
