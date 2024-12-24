package config

type ServiceNameEnum string

// Service names
const (
	ProfileService ServiceNameEnum = "profile"
	StatsService   ServiceNameEnum = "stats"
	CameraService  ServiceNameEnum = "camera"
)

// Create a set of valid service names
var validServiceNames = map[ServiceNameEnum]struct{}{
	ProfileService: {},
	StatsService:   {},
	CameraService:  {},
}

// Function to check if the service name is valid
func IsValidServiceName(serviceName string) bool {
	_, exists := validServiceNames[ServiceNameEnum(serviceName)]
	return exists
}
