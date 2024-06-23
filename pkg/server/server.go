package server

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/JacobRWebb/InventoryManagement.Users.Api/pkg/config"
	"github.com/JacobRWebb/InventoryManagement.Users.Api/pkg/consul"
	UserServiceProto "github.com/JacobRWebb/InventoryManagement.Users.Api/pkg/proto/v1/user"
	"github.com/JacobRWebb/InventoryManagement.Users.Api/pkg/service"
	"github.com/JacobRWebb/InventoryManagement.Users.Api/pkg/store/user"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type Server struct {
	cfg          *config.Config
	grpcServer   *grpc.Server
	consulClient *consul.Client
	db           *gorm.DB
}

func NewServer(cfg *config.Config, db *gorm.DB) (*Server, error) {
	consulClient, err := consul.NewClient(cfg.ConsulAddr)

	if err != nil {
		return nil, fmt.Errorf("failed to create Consul client: %v", err)
	}

	grpcServer := grpc.NewServer()

	s := &Server{
		cfg:          cfg,
		grpcServer:   grpcServer,
		consulClient: consulClient,
		db:           db,
	}

	return s, nil
}

func (s *Server) Run() error {
	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()

		if err := s.runGRPCServer(); err != nil {
			log.Fatalf("Failed to run gRPC server: %v", err)
		}
	}()

	wg.Wait()
	return nil
}

func (s *Server) runGRPCServer() (err error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.GRPCPort))

	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	userStore := user.NewStore(s.db)

	userService := service.NewUserService(userStore)

	UserServiceProto.RegisterUserServiceServer(s.grpcServer, userService)

	if err := s.consulClient.Register(s.cfg); err != nil {
		return fmt.Errorf("failed to register with Consul: %v", err)
	}

	log.Printf("gRPC server is running")

	return s.grpcServer.Serve(lis)
}
