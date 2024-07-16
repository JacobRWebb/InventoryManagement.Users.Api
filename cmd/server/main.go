package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/JacobRWebb/InventoryManagement.Users.Api/internal/config"
	"github.com/JacobRWebb/InventoryManagement.Users.Api/internal/consul"
	"github.com/JacobRWebb/InventoryManagement.Users.Api/internal/database"
	"github.com/JacobRWebb/InventoryManagement.Users.Api/internal/server"
	"github.com/JacobRWebb/InventoryManagement.Users.Api/internal/service"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

func run() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	logger := log.New(os.Stdout, "", log.LstdFlags)

	_, err = consul.NewClient(cfg)

	if err != nil {
		return fmt.Errorf("error while starting consul client: %v", err)
	}

	db := database.MustOpen(cfg.DB.DSN)

	userService := service.NewUserService(db)

	srv := server.NewServer(cfg, userService)

	go func() {
		logger.Println("Starting gRPC server...")
		srv.Run()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Println("Shutting down server...")

	srv.SetServingStatus(false)
	srv.Shutdown()

	logger.Println("Server exiting")
	return nil
}
