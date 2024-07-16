package consul

import (
	"fmt"
	"log"
	"time"

	"github.com/JacobRWebb/InventoryManagement.Users.Api/internal/config"
	"github.com/hashicorp/consul/api"
)

type Client struct {
	consulClient *api.Client
	cfg          *config.Config
}

func NewClient(cfg *config.Config) (*api.Client, error) {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = cfg.Consul.Addr

	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Consul client: %w", err)
	}

	client := &Client{
		consulClient: consulClient,
		cfg:          cfg,
	}

	err = client.register()

	if err != nil {
		return nil, err
	}

	return consulClient, nil
}

func (c *Client) register() error {
	log.Printf("Registering service %s at %s:%d", c.cfg.Service.Name, c.cfg.Service.GRPCAddr, c.cfg.Service.GRPCPort)

	grpcReg := &api.AgentServiceRegistration{
		ID:      "User_Service",
		Name:    c.cfg.Service.Name,
		Address: c.cfg.Service.GRPCAddr,
		Port:    c.cfg.Service.GRPCPort,
		Check: &api.AgentServiceCheck{
			TTL:                            "8s",
			DeregisterCriticalServiceAfter: c.cfg.Consul.DeregisterTime,
			CheckID:                        "CheckAlive",
		},
	}

	if err := c.consulClient.Agent().ServiceRegister(grpcReg); err != nil {
		return fmt.Errorf("failed to register service with Consul: %w", err)
	}

	go c.updateHealthChecker()

	return nil
}

func (c *Client) updateHealthChecker() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		err := c.consulClient.Agent().UpdateTTL("CheckAlive", "online", api.HealthPassing)
		if err != nil {
			log.Printf("failed to update TTL: %v", err)
		}
	}
}
