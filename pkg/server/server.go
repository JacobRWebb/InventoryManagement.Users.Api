package server

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/JacobRWebb/InventoryManagement.Users.Api/pkg/config"
	"github.com/JacobRWebb/InventoryManagement.Users.Api/pkg/consul"
	pb "github.com/JacobRWebb/InventoryManagement.Users.Api/pkg/proto/v1/user"
	"github.com/JacobRWebb/InventoryManagement.Users.Api/pkg/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthPB "google.golang.org/grpc/health/grpc_health_v1"
)

type Server struct {
	cfg          *config.Config
	grpcServer   *grpc.Server
	consulClient *consul.Client
}

func NewServer(cfg *config.Config) (*Server, error) {
	consulClient, err := consul.NewClient(cfg.ConsulAddr)

	if err != nil {
		return nil, fmt.Errorf("failed to create Consul client: %v", err)
	}

	s := &Server{
		cfg:          cfg,
		grpcServer:   grpc.NewServer(),
		consulClient: consulClient,
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

	healthServer := health.NewServer()
	healthPB.RegisterHealthServer(s.grpcServer, healthServer)
	pb.RegisterServiceServer(s.grpcServer, &service.UserService{})

	healthServer.SetServingStatus("", healthPB.HealthCheckResponse_SERVING)

	if err := s.consulClient.Register(s.cfg); err != nil {
		return fmt.Errorf("failed to register with Consul: %v", err)
	}

	log.Printf("gRPC server is running")

	return s.grpcServer.Serve(lis)
}
