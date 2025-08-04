package config

type Provider interface {
	GetAPIKey() string
	GetServiceURL(name string) (string, bool)
	GetAllServices() map[string]Service
}
