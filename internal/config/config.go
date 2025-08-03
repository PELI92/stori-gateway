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
	APIKey   string             `koanf:"api_key"`
	Services map[string]Service `koanf:"services"`
}

var (
	K    = koanf.New(".")
	Lock sync.RWMutex
	Conf Config
)

func LoadConfig() *Config {
	// Initial config load
	load()

	// Watcher for config file changes
	go watch()

	return &Conf
}

func load() {
	Lock.Lock()
	defer Lock.Unlock()

	if err := K.Load(file.Provider("config/config.yaml"), yaml.Parser()); err != nil {
		log.Error().Err(err).Msg("could not load config")
		return
	}

	if err := K.Unmarshal("", &Conf); err != nil {
		log.Error().Err(err).Msg("could not unmarshal config")
		return
	}

	log.Info().Msg("configuration reloaded")
}

func watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error().Err(err).Msg("failed to create watcher")
		return
	}
	defer watcher.Close()

	err = watcher.Add("config/config.yaml")
	if err != nil {
		log.Error().Err(err).Msg("failed to watch config file")
		return
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Info().Msg("detected config change")
				load()
			}
		case err := <-watcher.Errors:
			log.Error().Err(err).Msg("watcher error")
		}
	}
}
