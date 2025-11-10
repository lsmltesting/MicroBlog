package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lsmltesting/MicroBlog/internal/events"
	"github.com/lsmltesting/MicroBlog/internal/trace"
	"github.com/lsmltesting/MicroBlog/services/api/internal/dto"
	"github.com/lsmltesting/MicroBlog/services/api/internal/service/post"
	"github.com/lsmltesting/MicroBlog/services/api/internal/service/user"
)

type PostHTTPHandler struct {
	PostService post.PostService
	UserService user.UserService
	Producer    events.Producer
}

func NewPostHTTPHandler(
	postService post.PostService,
	userService user.UserService,
	producer events.Producer,
) *PostHTTPHandler {
	return &PostHTTPHandler{
		PostService: postService,
		UserService: userService,
		Producer:    producer,
	}
}

func (h *PostHTTPHandler) sendError(w http.ResponseWriter, message string, statusCode int) {
	errUserDTO := dto.ErrorDTO{
		Message: message,
		Time:    time.Now(),
	}

	http.Error(w, errUserDTO.ToString(), statusCode)
}

func (p *PostHTTPHandler) HandlerCreate(w http.ResponseWriter, r *http.Request) {
	var postDTO dto.PostDTO
	ctx := r.Context()

	// Check to get PostDTO from request body
	if err := json.NewDecoder(r.Body).Decode(&postDTO); err != nil {
		p.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	postID, err := p.PostService.Create(ctx, postDTO.UserID, postDTO.Text)
	if err != nil {
		p.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	post, err := p.PostService.FindByID(ctx, postID)
	if err != nil {
		p.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// publish an event
	traceID := trace.IDFromContext(ctx)
	env := events.Envelope{
		ID:         uuid.NewString(),
		OccurredAt: time.Now().UTC(),
		TraceID:    traceID,
	}

	_ = p.Producer.PublishPostCreated(ctx, env, events.PostCreatedPayload{
		PostID: post.ID,
		UserID: post.UserID,
		Text:   post.Text,
	})

	b, err := json.MarshalIndent(post, "", "    ")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}

func (p *PostHTTPHandler) HandlerGetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	posts, err := p.PostService.GetAll(ctx)
	if err != nil {
		p.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	b, err := json.MarshalIndent(posts, "", "    ")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
