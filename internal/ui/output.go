package ui

import (
	"fmt"
	"strings"

	"gitlab-to-github-migration/internal/gitlab"
	"gitlab-to-github-migration/internal/migrate"
)

func PrintBanner() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println(" GitLab → GitHub Migration Tool")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

func ConfirmProjects(projects []gitlab.Project) bool {
	fmt.Println("Projects:")
	for i, p := range projects {
		fmt.Printf("  %d. %s\n", i+1, p.Name)
	}

	fmt.Print("Proceed? (yes/no): ")
	var in string
	fmt.Scanln(&in)

	return strings.HasPrefix(strings.ToLower(in), "y")
}

func PrintResults(r migrate.Result) {
	fmt.Printf("Successful: %d\n", len(r.Success))
	fmt.Printf("Failed: %d\n", len(r.Failed))
}
