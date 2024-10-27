package gh_command

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

type Commit struct {
	Sha     string `json:"sha"`
	Message string `json:"message"`
	Author  string `json:"author"`
}

// GetBranchCommits function to get the commits between two branches from GitHub
func GetBranchCommits(owner string, repo string, baseBranch string, headBranch string) []Commit {
	// Run the GitHub CLI command and capture the output
	cmd := exec.Command(
		"gh",
		"api",
		fmt.Sprintf("repos/%s/%s/compare/%s...%s", owner, repo, baseBranch, headBranch),
		"--jq",
		"[.commits[] | {sha: .sha, message: .commit.message, author: .commit.author.name}]",
	)
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to execute gh command: %v", err)
	}

	// Parse the JSON output into a slice of strings
	var commits []Commit
	err = json.Unmarshal(output, &commits)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	return commits
}
