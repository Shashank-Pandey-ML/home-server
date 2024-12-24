package services

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/consul/api"

	"gateway/app/config"
	"gateway/app/logger"
)

var client *api.Client

const consul_host = "consul"
const consul_port = "8500"

func init() {
	logger.Logger.Debug("Initializing service package")
	if config.AppConfig.Service.DeploymentMode == "docker-compose" {
		// Initialize Consul client
		var err error
		consul_address := consul_host + ":" + consul_port
		client, err = api.NewClient(&api.Config{
			Address: consul_address, // Point to the Consul service in Docker
		})
		if err != nil {
			log.Fatalf("Failed to create Consul client: %v", err)
		}

		// Register the service in Consul
		registerServiceInConsul()
	}
}

func registerServiceInConsul() {
	serviceName := config.AppConfig.Service.Name
	servicePort := config.AppConfig.Service.Port

	// Define the service registration
	registration := &api.AgentServiceRegistration{
		ID:      serviceName + "-1", // Unique service ID
		Name:    serviceName,        // Service name
		Port:    servicePort,        // Service port
		Address: "127.0.0.1",        // Service address (replace with your service's IP)
		Tags:    []string{},         // Tags for the service (optional)
		Check: &api.AgentServiceCheck{ // Health check configuration
			HTTP:     "http://127.0.0.1:8080/health", // Replace with your health check endpoint
			Interval: "10s",                          // Health check interval
			Timeout:  "5s",                           // Timeout for the health check
		},
	}

	// Register the service
	err := client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Fatalf("Failed to register service with Consul: %v", err)
	}

	logger.Logger.Info(fmt.Sprintf("Service %s registered successfully with Consul", serviceName))
}

// Function to discover a service URL using the service name from consul
func discoverServiceFromConsul(serviceName string) (string, error) {
	// Query the catalog to get the list of services
	serviceEntries, _, err := client.Catalog().Service(serviceName, "", nil)
	if err != nil {
		return "", fmt.Errorf("error querying services: %v", err)
	}

	var portNumber int
	var address string

	// Display the service instances
	log.Printf("\nInstances of %s:\n", serviceName)
	for _, entry := range serviceEntries {
		address = entry.Address
		portNumber = entry.ServicePort
		log.Printf("Service ID: %s, Address: %s, Port: %d\n", entry.ServiceID, entry.Address, entry.ServicePort)
	}

	return address + ":" + strconv.Itoa(portNumber), nil
}

// GetServiceURL fetches the URL of a microservice from the discovery service
func GetServiceURL(serviceName string) (string, error) {
	if !config.IsValidServiceName(serviceName) {
		return "", fmt.Errorf("service name '%s' is invalid", serviceName)
	}

	return discoverServiceFromConsul(serviceName)
}

// ProxyRequest proxies a request to a specific service
func ProxyRequest(serviceName string, path string, method string) (*http.Response, error) {
	baseUrl, err := GetServiceURL(serviceName)
	if err != nil {
		return nil, err
	}

	var resp *http.Response
	url := fmt.Sprintf("%s%s", baseUrl, path)
	for attempt := 0; attempt <= config.AppConfig.Api.MaxRetries; attempt++ {
		client := &http.Client{Timeout: time.Duration(config.AppConfig.Api.Timeout) * time.Second}
		// TODO: Handle other methods PUT, POST, DELETE.
		resp, err = client.Get(url)
		if err == nil {
			return resp, nil // Successful request
		}

		// Retry only for transient errors
		if attempt < config.AppConfig.Api.MaxRetries {
			fmt.Printf("Retry %d/%d: Error - %v\n", attempt+1, config.AppConfig.Api.MaxRetries, err)
			time.Sleep(time.Duration(config.AppConfig.Api.RetryDelay)) // Delay before retrying
		}
	}

	return nil, fmt.Errorf("request failed after %d retries: %w", config.AppConfig.Api.MaxRetries, err)
}
