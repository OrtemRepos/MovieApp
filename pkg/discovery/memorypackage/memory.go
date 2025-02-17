package memory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"movieexample.com/pkg/discovery"
)

type serviceName string
type instanceID string

type Registry struct {
	sync.RWMutex
	serviceAddrs map[serviceName]map[instanceID]*serviceInstance
}

type serviceInstance struct {
	hostPort   string
	lastActive time.Time
}

func NewRegistry() *Registry {
	return &Registry{serviceAddrs: map[serviceName]map[instanceID]*serviceInstance{}}
}

func (r *Registry) Register(
	ctx context.Context, id string,
	name string, hostPort string,
) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddrs[serviceName(name)]; !ok {
		r.serviceAddrs[serviceName(name)] = make(map[instanceID]*serviceInstance, 100)
	}
	r.serviceAddrs[serviceName(name)][instanceID(id)] = &serviceInstance{hostPort: hostPort, lastActive: time.Now()}
	return nil
}

func (r *Registry) Deregister(ctx context.Context, id string, name string) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddrs[serviceName(name)]; !ok {
		return discovery.ErrNotFound
	}
	delete(r.serviceAddrs[serviceName(name)], instanceID(id))
	return nil
}

func (r *Registry) ReportHealthyState(id string, name string) error {
	r.Lock()
	if _, ok := r.serviceAddrs[serviceName(name)]; !ok {
		return fmt.Errorf("service is not registered yet: %w", discovery.ErrNotFound)
	}
	if _, ok := r.serviceAddrs[serviceName(name)][instanceID(id)]; !ok {
		return fmt.Errorf("service instance is not registered yet: %w", discovery.ErrNotFound)
	}
	r.serviceAddrs[serviceName(name)][instanceID(id)].lastActive = time.Now()
	return nil
}

func (r *Registry) ServiceAddresses(ctx context.Context, name string) ([]string, error) {
	r.RLock()
	defer r.RUnlock()
	if len(r.serviceAddrs[serviceName(name)]) == 0 {
		return nil, discovery.ErrNotFound
	}
	res := make([]string, len(r.serviceAddrs[serviceName(name)]))
	for _, i := range r.serviceAddrs[serviceName(name)] {
		if !i.lastActive.Before(time.Now().Add(-5 * time.Second)) {
			res = append(res, i.hostPort)
		}
	}
	return res, nil
}
