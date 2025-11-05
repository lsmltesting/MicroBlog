package post

import (
	"context"
	"strconv"

	"github.com/lsmltesting/MicroBlog/services/api/internal/logger"
	"github.com/lsmltesting/MicroBlog/services/api/internal/models"
)

type postServiceDecorator struct {
	postService PostService
	lg          logger.Logger
}

func NewPostServiceDecorator(postService PostService, lg logger.Logger) PostService {
	return &postServiceDecorator{
		postService: postService,
		lg:          lg,
	}
}

func (p *postServiceDecorator) Create(ctx context.Context, userID int, text string) (int, error) {
	p.lg.AddLog(
		logger.LevelInfo,
		logger.SourceService,
		map[string]string{
			"userID": strconv.Itoa(userID),
			"text":   text,
		},
		"CreatePost from post service is called",
	)

	postID, err := p.postService.Create(ctx, userID, text)
	if err != nil {
		p.lg.AddLog(
			logger.LevelError,
			logger.SourceService,
			map[string]string{
				"userID": strconv.Itoa(userID),
				"text":   text,
				"postID": strconv.Itoa(postID),
			},
			err.Error(),
		)
	} else {
		p.lg.AddLog(
			logger.LevelInfo,
			logger.SourceService,
			map[string]string{
				"userID": strconv.Itoa(userID),
				"text":   text,
				"postID": strconv.Itoa(postID),
			},
			"CreatePost called successfully",
		)
	}

	return postID, err
}

func (p *postServiceDecorator) FindByID(ctx context.Context, postID int) (*models.Post, error) {
	p.lg.AddLog(
		logger.LevelInfo,
		logger.SourceService,
		map[string]string{
			"postID": strconv.Itoa(postID),
		},
		"GetPostByID from post service is called",
	)

	postModel, err := p.postService.FindByID(ctx, postID)

	if err != nil {
		p.lg.AddLog(
			logger.LevelError,
			logger.SourceService,
			map[string]string{
				"postID": strconv.Itoa(postID),
			},
			err.Error(),
		)
	} else {
		p.lg.AddLog(
			logger.LevelInfo,
			logger.SourceService,
			map[string]string{
				"postID": strconv.Itoa(postID),
			},
			"GetPostByID called successfully",
		)
	}

	return postModel, err
}

func (p *postServiceDecorator) GetAll(ctx context.Context) (map[int]*models.Post, error) {
	p.lg.AddLog(
		logger.LevelInfo,
		logger.SourceService,
		make(map[string]string),
		"GetAllPosts from post service is called",
	)

	postModel, err := p.postService.GetAll(ctx)

	if err != nil {
		p.lg.AddLog(
			logger.LevelError,
			logger.SourceService,
			make(map[string]string),
			err.Error(),
		)
	} else {
		p.lg.AddLog(
			logger.LevelInfo,
			logger.SourceService,
			make(map[string]string),
			"GetAllPosts called successfully",
		)
	}

	return postModel, err
}
