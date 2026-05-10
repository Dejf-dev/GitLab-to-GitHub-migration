package migrate

import (
	"fmt"
	"os"
	"path/filepath"

	"gitlab-to-github-migration/internal/config"
	"gitlab-to-github-migration/internal/git"
	"gitlab-to-github-migration/internal/github"
	"gitlab-to-github-migration/internal/gitlab"

	log "github.com/sirupsen/logrus"
)

type Result struct {
	Success []string
	Failed  []string
}

func Filter(projects []gitlab.Project, filter string) []gitlab.Project {
	log.Info("Filtering projects")
	if filter == "" {
		return projects
	}

	out := []gitlab.Project{}
	for _, p := range projects {
		if p.Name == filter {
			out = append(out, p)
		}
	}
	return out
}

func Run(
	projects []gitlab.Project,
	gl *gitlab.Client,
	gh *github.Client,
	cfg config.Config,
) Result {
	// Create temporary directory for cloning
	tmp := "temp_migration"
	os.MkdirAll(tmp, 0755)
	defer os.RemoveAll(tmp)

	res := Result{}

	for _, p := range projects {
		repoPath := fmt.Sprintf("%s-%d", p.Path, p.ID)
		path := filepath.Join(tmp, repoPath)

		log.Infof("Cloning project %s", p.Name)
		if err := git.MirrorClone(p.SSHURLToRepo, path); err != nil {
			log.Error(err)
			res.Failed = append(res.Failed, p.Name)
			continue
		}

		// cleanup large files
		log.Infof("Cleaning large files in %s", p.Name)
		if err := git.RemoveLargeFiles(path); err != nil {
			log.Warnf("Skipping cleanup failed repo %s: %v", p.Name, err)
			res.Failed = append(res.Failed, p.Name)
			continue
		}

		log.Infof("Creating GitHub repository %s", repoPath)
		if err := gh.CreateRepo(repoPath, p.Description); err != nil {
			log.Error(err)
			res.Failed = append(res.Failed, p.Name)
			continue
		}

		log.Infof("Pushing project %s to GitHub", p.Name)
		if err := git.MirrorPush(path, gh.RemoteURL(repoPath)); err != nil {
			log.Error(err)
			res.Failed = append(res.Failed, p.Name)
			continue
		}

		res.Success = append(res.Success, p.Name)
	}

	return res
}
