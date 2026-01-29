package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gitlab-to-github-migration/internal/config"
)

type createRepoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
}

type Client struct {
	username string
	token    string
	private  bool
	client   *http.Client
}

func NewClient(cfg config.Config) *Client {
	return &Client{
		username: cfg.GitHubUsername,
		token:    cfg.GitHubToken,
		private:  cfg.PrivateRepos,
		client:   &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) CreateRepo(name, desc string) error {
	body, _ := json.Marshal(createRepoRequest{
		Name:        name,
		Description: desc,
		Private:     c.private,
	})

	req, _ := http.NewRequest(
		"POST",
		"https://api.github.com/user/repos",
		bytes.NewBuffer(body),
	)

	req.Header.Set("Authorization", "token "+c.token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated ||
		resp.StatusCode == http.StatusUnprocessableEntity {
		return nil
	}

	return fmt.Errorf("GitHub API error: %s", resp.Status)
}

func (c *Client) RemoteURL(repo string) string {
	return fmt.Sprintf("git@github.com:%s/%s.git", c.username, repo)
}
