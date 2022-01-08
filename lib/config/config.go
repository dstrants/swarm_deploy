package config

import (
	env "github.com/Netflix/go-env"
	log "github.com/sirupsen/logrus"
)

func LoadConfig() Config {

	var config Config
	es, err := env.UnmarshalFromEnviron(&config)

	if err != nil {
		log.Fatalf("Loading environment configuration failed with: %v", err)
	}

	config.Extras = es

	return config
}
