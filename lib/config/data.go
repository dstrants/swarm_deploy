package config

import env "github.com/Netflix/go-env"

type Config struct {
	GithubWebhookSecret string `env:"GITHUB_WEBHOOK_SECRET,required=true"`
	Extras              env.EnvSet
}
