package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/lsmltesting/MicroBlog/services/engagement/internal/repository"
)

type StatsHandler struct {
	repo repository.StatsRepository
}

func NewStatsHandler(repo repository.StatsRepository) *StatsHandler {
	return &StatsHandler{repo: repo}
}

type PostStatsResponse struct {
	PostID    int   `json:"post_id"`
	LikeCount int64 `json:"like_count"`
}

func (h *StatsHandler) GetPostStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	postID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid post id", http.StatusBadRequest)
		return
	}

	count, err := h.repo.GetPostLikeCount(r.Context(), postID)
	if err != nil {
		http.Error(w, "failed to get stats", http.StatusInternalServerError)
		return
	}

	resp := PostStatsResponse{
		PostID:    postID,
		LikeCount: count,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
