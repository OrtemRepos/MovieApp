package consul

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/consul/api"
	"movieexample.com/pkg/discovery"
)

type Registry struct {
	client *api.Client
}

func NewRegistry(addr string) (*Registry, error) {
	config := api.DefaultConfig()
	config.Address = addr
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &Registry{client}, nil
}

func (r *Registry) Register(
	ctx context.Context, id string,
	name string, hostPort string,
) error {
	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 {
		return errors.New("hostPort must be in a form of <host>:<port>, example: localhost:8080")
	}
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}
	agent := api.AgentServiceRegistration{
		Address: parts[0],
		ID:      id,
		Name:    name,
		Port:    port,
		Check:   &api.AgentServiceCheck{CheckID: id, TTL: "5s"},
	}
	return r.client.Agent().ServiceRegister(&agent)
}

func (r *Registry) Deregister(ctx context.Context, id string, _ string) error {
	return r.client.Agent().ServiceDeregister(id)
}

func (r *Registry) ServiceAddresses(ctx context.Context, name string) ([]string, error) {
	entries, _, err := r.client.Health().Service(name, "", true, nil)
	if err != nil {
		return nil, err
	} else if len(entries) == 0 {
		return nil, discovery.ErrNotFound
	}

	res := make([]string, len(entries))
	for _, e := range entries {
		res = append(res, fmt.Sprintf("%s:%d", e.Service.Address, e.Service.Port))
	}
	return res, nil
}

func (r *Registry) ReportHealthyState(id string, _ string) error {
	return r.client.Agent().PassTTL(id, "")
}