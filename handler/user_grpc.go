package handler

import (
	"context"
	pb "github.com/ReStorePUC/protobucket/user"
	"strconv"
)

// UserServer is used to implement User.
type UserServer struct {
	pb.UnimplementedUserServer

	controller controller
}

func NewUserServer(c controller) *UserServer {
	return &UserServer{
		controller: c,
	}
}

func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := s.controller.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserResponse{
		Id:      strconv.Itoa(user.ID),
		IsAdmin: user.IsAdmin,
	}, nil
}
