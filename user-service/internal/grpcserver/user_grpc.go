package grpcserver

import (
	"context"

	"github.com/RibunLoc/microservices-learn/user-service/proto/userpb"
	repository "github.com/RibunLoc/microservices-learn/user-service/repository/user"
)

type UserGRPCHandler struct {
	userpb.UnimplementedUserServiceServer
	Repo *repository.RedisMongo
}

func (h *UserGRPCHandler) GetUserByID(ctx context.Context, req *userpb.GetUserByIDRequest) (*userpb.GetUserByIDResponse, error) {
	user, err := h.Repo.FindByID(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &userpb.GetUserByIDResponse{
		UserId:   user.ID.Hex(),
		Email:    user.Email,
		Fullname: user.Fullname,
		Role:     user.Role,
	}, nil
}
