package gh_command

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// CreatePullRequestParams struct to represent the parameters for creating a pull request
type CreatePullRequestParams struct {
	BaseBranch string
	HeadBranch string
	Title      string
	Body       string
	Reviewers  []string
	IsDraft    bool
}

// CreatePullRequest function to create a pull request on GitHub
func CreatePullRequest(owner string, repo string, options CreatePullRequestParams) (string, error) {
	if options.BaseBranch == "" || options.HeadBranch == "" || options.Title == "" {
		log.Fatalf("BaseBranch, HeadBranch, and Title are required")
	}

	args := []string{
		"pr",
		"create",
		"--repo", fmt.Sprintf("%s/%s", owner, repo),
		"--base", options.BaseBranch,
		"--head", options.HeadBranch,
		"--title", options.Title,
	}

	// Append "--body" and options.body if options.body is available
	if options.Body != "" {
		args = append(args, "--body", options.Body)
	}

	// Append draft flag if options.IsDraft is true
	if options.IsDraft == true {
		args = append(args, "--draft")
	}

	// Append reviewers if options.Reviewers is available
	if len(options.Reviewers) > 0 {
		args = append(args, "--reviewer", strings.Join(options.Reviewers, ","))
	}

	// Run the GitHub CLI command and capture the output
	cmd := exec.Command("gh", args...)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Failed to execute gh command: %v", err)

		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}
