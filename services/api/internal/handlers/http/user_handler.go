package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lsmltesting/MicroBlog/services/api/internal/dto"
	"github.com/lsmltesting/MicroBlog/services/api/internal/events"
	"github.com/lsmltesting/MicroBlog/services/api/internal/service/user"
	"github.com/lsmltesting/MicroBlog/services/api/internal/trace"
)

type UserHTTPHandler struct {
	UserService user.UserService
	Producer    events.Producer
}

func NewUserHTTPHandler(userService user.UserService, producer events.Producer) *UserHTTPHandler {
	return &UserHTTPHandler{
		UserService: userService,
		Producer:    producer,
	}
}

func (h *UserHTTPHandler) sendError(w http.ResponseWriter, message string, statusCode int) {
	errUserDTO := dto.ErrorDTO{
		Message: message,
		Time:    time.Now(),
	}

	http.Error(w, errUserDTO.ToString(), statusCode)
}

func (h *UserHTTPHandler) HandlerCreate(w http.ResponseWriter, r *http.Request) {
	var userDTO dto.UserDTO
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := h.UserService.Create(ctx, userDTO.Username, userDTO.Email, userDTO.Password)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusConflict)
		return
	}

	user, err := h.UserService.FindByID(ctx, userID)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// publish an event
	traceID := trace.IDFromContext(ctx)
	env := events.Envelope{
		ID:         uuid.NewString(),
		OccurredAt: time.Now().UTC(),
		TraceID:    traceID,
	}

	_ = h.Producer.PublishUserRegistered(ctx, env, events.UserCreatedPayload{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
	})

	b, err := json.MarshalIndent(user, "", "    ")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}
