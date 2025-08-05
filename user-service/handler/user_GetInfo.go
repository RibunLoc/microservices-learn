package handler

import (
	"encoding/json"
	"net/http"

	repository "github.com/RibunLoc/microservices-learn/user-service/repository/user"
	"github.com/RibunLoc/microservices-learn/user-service/util"

	"github.com/go-chi/chi/v5"
)

type UserGetInfo struct {
	Repo *repository.RedisMongo
}

func (h *UserGetInfo) GetInfoHandler(w http.ResponseWriter, r *http.Request) {
	targetUserID := chi.URLParam(r, "id")

	// 1. Lấy userID từ JWT
	userID, ok := util.GetUserIDFromRequest(r, h.Repo.JwtSecret)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if userID != targetUserID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var repBody struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		Fullname  string `json:"full_name"`
		Role      string `json:"role"`
		IsActive  bool   `json:"is_active"`
		CreatedAt string `json:"created_at"`
	}

	user, err := h.Repo.FindByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	repBody.ID = user.ID.Hex()
	repBody.Email = user.Email
	repBody.Fullname = user.Fullname
	repBody.Role = user.Role
	repBody.IsActive = user.IsActive
	repBody.CreatedAt = user.CreatedAt.String()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(repBody); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
