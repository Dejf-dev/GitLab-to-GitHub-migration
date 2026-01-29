package config

import (
	"os"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	GitLabURL      string
	GitLabToken    string
	GitHubUsername string
	GitHubToken    string
	Filter         string
	PrivateRepos   bool
}

func MustLoad() Config {
	log.Info("Loading configuration")
	cfg := Config{
		GitLabURL:      os.Getenv("GITLAB_URL"),
		GitLabToken:    os.Getenv("GITLAB_TOKEN"),
		GitHubUsername: os.Getenv("GITHUB_USERNAME"),
		GitHubToken:    os.Getenv("GITHUB_TOKEN"),
		Filter:         os.Getenv("FILTER"),
		PrivateRepos:   true,
	}

	if cfg.GitLabToken == "" || cfg.GitHubToken == "" {
		log.Fatal("missing tokens")
	}

	return cfg
}
