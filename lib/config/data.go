package config

import (
	"strings"

	env "github.com/Netflix/go-env"
)

type Config struct {
	// Gin Webserver configuration
	WebServer struct {
		Host string `env:"WEBSERVER_HOST,default="`
		Port int    `env:"WEBSERVER_PORT,default=8080"`
	}
	// Github Related Configuration
	Github struct {
		WebhookSecret string `env:"GITHUB_WEBHOOK_SECRET,required=true"`
		WebhookEvents string `env:"GITHUB_WEBHOOK_EVENTS,default=package,ping"`
	}

	Slack struct {
		Token   string `env:"SLACK_TOKEN,required=true"`
		Channel string `env:"SLACK_CHANNEL,default=news"`
	}

	Extras env.EnvSet
}

func (c Config) GithubWebhookEvents() []string {
	AcceptedEvents := strings.Split(c.Github.WebhookEvents, ",")

	return AcceptedEvents
}
