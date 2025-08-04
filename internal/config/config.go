package config

import (
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/rs/zerolog/log"
)

type Service struct {
	URL string `koanf:"url"`
}

type Config struct {
	mu       sync.RWMutex // everything within this config should be thread safe, therefore RWMutex is used.
	apiKey   string
	services map[string]Service
}

// LoadConfig loads the config and sets up file watching for hot reloads
func LoadConfig() *Config {
	k := koanf.New(".")

	cfg := &Config{}
	cfg.load(k)

	go cfg.watch(k)

	return cfg
}

func (c *Config) load(k *koanf.Koanf) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := k.Load(file.Provider("config/config.yaml"), yaml.Parser()); err != nil {
		log.Error().Err(err).Msg("could not load config")
		return
	}

	var raw struct {
		APIKey   string             `koanf:"api_key"`
		Services map[string]Service `koanf:"services"`
	}

	if err := k.Unmarshal("", &raw); err != nil {
		log.Error().Err(err).Msg("could not unmarshal config")
		return
	}

	c.apiKey = raw.APIKey
	c.services = raw.Services
	log.Info().Msg("configuration reloaded")
}

func (c *Config) watch(k *koanf.Koanf) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error().Err(err).Msg("failed to create watcher")
		return
	}
	defer watcher.Close()

	if err := watcher.Add("config/config.yaml"); err != nil {
		log.Error().Err(err).Msg("failed to watch config file")
		return
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Info().Msg("detected config change")
				c.load(k)
			}
		case err := <-watcher.Errors:
			log.Error().Err(err).Msg("watcher error")
		}
	}
}

// GetServiceURL returns the service URL for a given name (thread-safe)
func (c *Config) GetServiceURL(name string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	svc, ok := c.services[name]
	if !ok {
		return "", false
	}
	return svc.URL, true
}

// GetAPIKey returns the current API key (thread-safe)
func (c *Config) GetAPIKey() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.apiKey
}

func (c *Config) GetAllServices() map[string]Service {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Create a copy to avoid concurrent access or external modifications
	servicesCopy := make(map[string]Service, len(c.services))
	for name, svc := range c.services {
		servicesCopy[name] = svc
	}

	return servicesCopy
}
