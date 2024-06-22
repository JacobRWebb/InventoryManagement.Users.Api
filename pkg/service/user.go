package service

import (
	"context"
	"fmt"
	"log"

	pb "github.com/JacobRWebb/InventoryManagement.Users.Api/pkg/proto/v1/user"
)

type UserService struct {
	pb.UnimplementedServiceServer
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	log.Printf("Logged incoming create user request")
	return &pb.CreateUserResponse{
		Result: &pb.CreateUserResponse_UserId{
			UserId: fmt.Sprintf("%s-%s", req.Email, req.Password),
		},
	}, nil
}
