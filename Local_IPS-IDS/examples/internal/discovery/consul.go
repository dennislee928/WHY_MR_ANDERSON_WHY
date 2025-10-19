package discovery

import (
	"fmt"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
)

// ConsulDiscovery implements service discovery using Consul
type ConsulDiscovery struct {
	client *api.Client
	logger *logrus.Logger
}

// ServiceConfig contains service registration configuration
type ServiceConfig struct {
	ID      string
	Name    string
	Address string
	Port    int
	Tags    []string
	Meta    map[string]string
	Check   *HealthCheck
}

// HealthCheck defines health check configuration
type HealthCheck struct {
	HTTP                           string
	Interval                       time.Duration
	Timeout                        time.Duration
	DeregisterCriticalServiceAfter time.Duration
}

// NewConsulDiscovery creates a new Consul service discovery client
func NewConsulDiscovery(address string, logger *logrus.Logger) (*ConsulDiscovery, error) {
	if logger == nil {
		logger = logrus.New()
	}

	config := api.DefaultConfig()
	config.Address = address

	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Consul client: %w", err)
	}

	return &ConsulDiscovery{
		client: client,
		logger: logger,
	}, nil
}

// Register registers a service with Consul
func (cd *ConsulDiscovery) Register(config *ServiceConfig) error {
	registration := &api.AgentServiceRegistration{
		ID:      config.ID,
		Name:    config.Name,
		Address: config.Address,
		Port:    config.Port,
		Tags:    config.Tags,
		Meta:    config.Meta,
	}

	if config.Check != nil {
		registration.Check = &api.AgentServiceCheck{
			HTTP:                           config.Check.HTTP,
			Interval:                       config.Check.Interval.String(),
			Timeout:                        config.Check.Timeout.String(),
			DeregisterCriticalServiceAfter: config.Check.DeregisterCriticalServiceAfter.String(),
		}
	}

	if err := cd.client.Agent().ServiceRegister(registration); err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}

	cd.logger.Infof("Service registered: %s (ID: %s)", config.Name, config.ID)
	return nil
}

// Deregister removes a service from Consul
func (cd *ConsulDiscovery) Deregister(serviceID string) error {
	if err := cd.client.Agent().ServiceDeregister(serviceID); err != nil {
		return fmt.Errorf("failed to deregister service: %w", err)
	}

	cd.logger.Infof("Service deregistered: %s", serviceID)
	return nil
}

// Discover finds all healthy instances of a service
func (cd *ConsulDiscovery) Discover(serviceName string) ([]*ServiceInstance, error) {
	services, _, err := cd.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to discover service: %w", err)
	}

	instances := make([]*ServiceInstance, 0, len(services))
	for _, service := range services {
		instances = append(instances, &ServiceInstance{
			ID:      service.Service.ID,
			Name:    service.Service.Service,
			Address: service.Service.Address,
			Port:    service.Service.Port,
			Tags:    service.Service.Tags,
			Meta:    service.Service.Meta,
		})
	}

	cd.logger.Debugf("Discovered %d instances of service %s", len(instances), serviceName)
	return instances, nil
}

// Watch watches for changes to a service
func (cd *ConsulDiscovery) Watch(serviceName string, callback func([]*ServiceInstance)) error {
	plan, err := api.WatchPlan{
		Type: "service",
		Service: serviceName,
		Handler: func(idx uint64, data interface{}) {
			if entries, ok := data.([]*api.ServiceEntry); ok {
				instances := make([]*ServiceInstance, 0, len(entries))
				for _, entry := range entries {
					instances = append(instances, &ServiceInstance{
						ID:      entry.Service.ID,
						Name:    entry.Service.Service,
						Address: entry.Service.Address,
						Port:    entry.Service.Port,
						Tags:    entry.Service.Tags,
						Meta:    entry.Service.Meta,
					})
				}
				callback(instances)
			}
		},
	}.Run(cd.client.Config().Address)

	if err != nil {
		return fmt.Errorf("failed to watch service: %w", err)
	}

	return plan
}

// ServiceInstance represents a discovered service instance
type ServiceInstance struct {
	ID      string
	Name    string
	Address string
	Port    int
	Tags    []string
	Meta    map[string]string
}

// GetAddress returns the full address of the service instance
func (si *ServiceInstance) GetAddress() string {
	return fmt.Sprintf("%s:%d", si.Address, si.Port)
}

