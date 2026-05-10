package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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
		Description: cleanDescription(desc),
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

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)

		return fmt.Errorf(
			"GitHub API error: %s | body: %s",
			resp.Status,
			string(bodyBytes),
		)
	}

	return nil
}

func (c *Client) RemoteURL(repo string) string {
	return fmt.Sprintf("git@github.com:%s/%s.git", c.username, repo)
}

func cleanDescription(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")

	return strings.TrimSpace(s)
}
