package server

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/JacobRWebb/InventoryManagement.Users.Api/internal/config"
	"github.com/JacobRWebb/InventoryManagement.Users.Api/internal/service"
	UserProto "github.com/JacobRWebb/InventoryManagement.Users.Api/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	Healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type Server struct {
	cfg          *config.Config
	grpcServer   *grpc.Server
	userService  *service.UserService
	healthServer *health.Server
}

func NewServer(cfg *config.Config, userService *service.UserService) *Server {
	grpcServer := grpc.NewServer()
	healthServer := health.NewServer()

	s := &Server{
		cfg:          cfg,
		grpcServer:   grpcServer,
		userService:  userService,
		healthServer: healthServer,
	}

	return s
}

func (s *Server) Run() {
	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()

		if err := s.startgrpcServer(); err != nil {
			log.Fatalf("Failed to run gRPC server: %v", err)
		}
	}()

	wg.Wait()
}

func (s *Server) Shutdown() {
	s.grpcServer.Stop()
}

func (s *Server) startgrpcServer() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.Service.GRPCPort))

	if err != nil {
		return err
	}

	UserProto.RegisterUserServiceServer(s.grpcServer, s.userService)
	Healthpb.RegisterHealthServer(s.grpcServer, s.healthServer)

	log.Printf("gRPC server is running")

	s.healthServer.SetServingStatus("", Healthpb.HealthCheckResponse_SERVING)

	return s.grpcServer.Serve(lis)
}

func (s *Server) SetServingStatus(serving bool) {
	if serving {
		s.healthServer.SetServingStatus("", Healthpb.HealthCheckResponse_SERVING)
	} else {
		s.healthServer.SetServingStatus("", Healthpb.HealthCheckResponse_NOT_SERVING)
	}
}
