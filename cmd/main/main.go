package main

import (
	"log"

	"github.com/JacobRWebb/InventoryManagement.Users.Api/pkg/config"
	"github.com/JacobRWebb/InventoryManagement.Users.Api/pkg/server"
	"github.com/JacobRWebb/InventoryManagement.Users.Api/pkg/store"
)

func main() {
	cfg, err := config.NewConfig()

	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db := store.MustOpen(cfg)

	srv, err := server.NewServer(cfg, db)

	if err != nil {
		log.Fatalf("There was a problem with getting a new server: %v", err)
	}

	if err := srv.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
