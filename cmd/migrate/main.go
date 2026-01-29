package main

import (
	"gitlab-to-github-migration/internal/config"
	"gitlab-to-github-migration/internal/github"
	"gitlab-to-github-migration/internal/gitlab"
	"gitlab-to-github-migration/internal/migrate"
	"gitlab-to-github-migration/internal/ui"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Warn("no .env file found, relying on environment variables")
	}

	ui.PrintBanner()

	cfg := config.MustLoad()

	gl := gitlab.NewClient(cfg)
	gh := github.NewClient(cfg)

	log.Info("Fetching projects from GitLab")
	projects, err := gl.ListProjects()
	if err != nil {
		log.Fatal(err)
	}

	projects = migrate.Filter(projects, cfg.Filter)

	if !ui.ConfirmProjects(projects) {
		log.Warn("Migration cancelled")
		return
	}

	results := migrate.Run(projects, gl, gh, cfg)

	ui.PrintResults(results)
}
