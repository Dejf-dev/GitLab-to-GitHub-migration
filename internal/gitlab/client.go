package gitlab

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gitlab-to-github-migration/internal/config"
)

type Project struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Path          string `json:"path"`
	Description   string `json:"description"`
	SSHURLToRepo  string `json:"ssh_url_to_repo"`
	HTTPURLToRepo string `json:"http_url_to_repo"`
}

type Client struct {
	baseURL string
	token   string
	client  *http.Client
}

func NewClient(cfg config.Config) *Client {
	return &Client{
		baseURL: cfg.GitLabURL,
		token:   cfg.GitLabToken,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) ListProjects() ([]Project, error) {
	var all []Project
	page := 1

	for {
		url := fmt.Sprintf(
			"%s/api/v4/projects?membership=true&per_page=100&page=%d",
			c.baseURL, page,
		)

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("PRIVATE-TOKEN", c.token)

		resp, err := c.client.Do(req)
		if err != nil {
			return nil, err
		}

		var batch []Project
		if err := json.NewDecoder(resp.Body).Decode(&batch); err != nil {
			resp.Body.Close()
			return nil, err
		}
		resp.Body.Close()

		if len(batch) == 0 {
			break
		}

		all = append(all, batch...)
		page++
	}

	return all, nil
}
