package handler

import (
	"encoding/json"
	"net/http"

	repository "github.com/RibunLoc/microservices-learn/user-service/repository/user"
	"github.com/RibunLoc/microservices-learn/user-service/util"

	"github.com/go-chi/chi/v5"
)

type UserChangePassword struct {
	Repo *repository.RedisMongo
}

func (h *UserChangePassword) ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	targetUserID := chi.URLParam(r, "id")

	// Lấy userid từ request
	userID, ok := util.GetUserIDFromRequest(r, h.Repo.JwtSecret)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if userID != targetUserID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var body struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	user, err := h.Repo.FindByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	if !util.CheckPasswordHash(body.OldPassword, user.Password) {
		http.Error(w, "Old password incorrect", http.StatusUnauthorized)
		return
	}
	if len(body.NewPassword) < 6 {
		http.Error(w, "Password too short", http.StatusBadRequest)
		return
	}

	hash, _ := util.HashPassword(body.NewPassword)
	err = h.Repo.UpdatePassword(r.Context(), userID, hash)
	if err != nil {
		http.Error(w, "failed to update password", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Password changed successfully"}`))
}
