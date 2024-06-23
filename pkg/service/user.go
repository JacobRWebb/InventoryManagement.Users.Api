package service

import (
	"context"
	"log"
	"time"

	UserServiceProto "github.com/JacobRWebb/InventoryManagement.Users.Api/pkg/proto/v1/user"
	"github.com/JacobRWebb/InventoryManagement.Users.Api/pkg/store/user"
)

type UserService struct {
	UserServiceProto.UnimplementedUserServiceServer
	userStore *user.Store
}

func NewUserService(us *user.Store) *UserService {
	s := &UserService{userStore: us}

	return s
}

func (s *UserService) CreateUser(ctx context.Context, req *UserServiceProto.CreateUserRequest) (*UserServiceProto.CreateUserResponse, error) {
	go func(begin time.Time) {
		latency := time.Since(begin)
		log.Printf("Request processed. Latency: %v", latency)
	}(time.Now())

	user, err := s.userStore.CreateUser(&user.CreateUser{Email: req.Email, Password: req.Password})

	if err != nil {
		return &UserServiceProto.CreateUserResponse{
			Result: &UserServiceProto.CreateUserResponse_Error{
				Error: err.Error(),
			},
		}, err
	}

	return &UserServiceProto.CreateUserResponse{
		Result: &UserServiceProto.CreateUserResponse_UserId{
			UserId: user.Id.String(),
		},
	}, nil
}
