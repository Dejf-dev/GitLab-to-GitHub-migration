# GitLab → GitHub Migration Tool

A simple Go CLI tool to migrate repositories from **GitLab to GitHub** while preserving full history (commits, branches, tags).  
GitLab cloning and GitHub pushing are done via **SSH**.

---

## Features

- Migrates all GitLab projects you are a member of
- Preserves full git history (branches, tags, commits)
- Uses GitLab API + GitHub API
- Optional project name filtering
- Clear progress logging

---

## Requirements

- Go 1.24.4+
- Git installed
- SSH access to GitLab and GitHub

---

## SSH Setup

Verify SSH access:

```bash
ssh -T git@gitlab.fit.cvut.cz
ssh -T git@github.com
```

## GitLab configuration
* GITLAB_URL=https://gitlab.fit.cvut.cz
* GITLAB_TOKEN=glpat_xxxxxxxxxxxxxxxxx

## GitHub configuration
* GITHUB_USERNAME=your-github-username
* GITHUB_TOKEN=ghp_xxxxxxxxxxxxxxxxx

## Optional: migrate only projects whose name contains this string
* FILTER=

## Usage
```go
go run cmd/migrate/main.go
```

or

```go
go build -o migrate cmd/migrate/main.go
./migrate
```

