package discovery

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type Registry interface {
	Register(ctx context.Context, id string, name string, hostPort string) error
	Deregister(ctx context.Context, id string, name string) error
	ServiceAddresses(ctx context.Context, name string) ([]string, error)
	ReportHealthyState(id string, name string) error
}

var ErrNotFound = errors.New("no service addresses found")

func GenerateInstanceID(serviceName string) string {
	return fmt.Sprintf("%s-%s", serviceName, uuid.New())
}
