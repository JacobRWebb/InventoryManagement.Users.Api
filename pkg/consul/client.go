package consul

import (
	"fmt"
	"log"
	"time"

	"github.com/JacobRWebb/InventoryManagement.Users.Api/pkg/config"
	"github.com/hashicorp/consul/api"
)

type Client struct {
	client *api.Client
}

func NewClient(addr string) (*Client, error) {
	config := api.DefaultConfig()
	config.Address = addr

	client, err := api.NewClient(config)

	if err != nil {
		return nil, fmt.Errorf("failed to create Consule client: %v", err)
	}

	return &Client{client: client}, nil
}

func (c *Client) updateHealthChecker() {
	ticker := time.NewTicker(time.Second * 5)
	for {
		c.client.Agent().UpdateTTL("CheckAlive", "online", api.HealthPassing)
		<-ticker.C
	}
}

func (c *Client) Register(cfg *config.Config) (err error) {
	log.Printf("%s:%d", cfg.GRPCAddr, cfg.GRPCPort)
	// GRPC Service Consul
	grpcReg := &api.AgentServiceRegistration{
		ID:      "User_Service",
		Name:    cfg.ServiceName,
		Address: cfg.ConsulAddr,
		Port:    cfg.GRPCPort,
		Check: &api.AgentServiceCheck{
			TTL:                            fmt.Sprintf("%ds", 8),
			DeregisterCriticalServiceAfter: cfg.ConsulDeregisterTime,
			CheckID:                        "CheckAlive",
		},
	}

	err = c.client.Agent().ServiceRegister(grpcReg)

	go c.updateHealthChecker()

	if err != nil {
		return err
	}

	return nil
}
