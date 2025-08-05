package handler

import (
	"encoding/json"
	"net/http"
	repository "user-service/repository/user"
	"user-service/util"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
)

type UserUpdateInfo struct {
	Repo *repository.RedisMongo
}

func (h *UserUpdateInfo) UpdateInfoHandler(w http.ResponseWriter, r *http.Request) {
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

	// 2. Parse body
	var reqBody struct {
		Email    string `json:"email"`
		Fullname string `json:"full_name"`
		// sau nay có thể thêm các trường khác nếu cần
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	bsonM := bson.M{
		"full_name": reqBody.Fullname,
		"email":     reqBody.Email,
	}

	err := h.Repo.UpdateUserFields(r.Context(), userID, bsonM)
	if err != nil {
		http.Error(w, "failed to update user info", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"User info updated successfully"}`))

}
