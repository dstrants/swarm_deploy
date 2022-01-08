package config

import env "github.com/Netflix/go-env"

type Config struct {
	// Gin Webserver configuration
	WebServer struct {
		Host string `env:"WEBSERVER_HOST,default="`
		Port int    `env:"WEBSERVER_PORT,default=8080"`
	}
	// Github Related Configuration
	Github struct {
		WebhookSecret string `env:"GITHUB_WEBHOOK_SECRET,required=true"`
	}
	Extras env.EnvSet
}
