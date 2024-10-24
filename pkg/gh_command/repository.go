package gh_command

import (
	"encoding/json"
	"log"
	"os/exec"
)

type Repo struct {
	RepoName string
}

// RepoOwner struct to represent the owner of the repository
type RepoOwner struct {
	ID    string `json:"id"`
	Login string `json:"login"`
}

// RepoAssignableUser struct to represent an assignable user
type RepoAssignableUser struct {
	ID    string `json:"id"`
	Login string `json:"login"`
	Name  string `json:"name"`
}

// RepoDefaultBranchRef struct to represent the default branch of a repository
type RepoDefaultBranchRef struct {
	Name string `json:"name"`
}

// GetRepoOptions struct to represent the options for getting a repository
type GetRepoOptions struct {
	Json string
}

// GetRepoResponse struct to represent the response of getting a repository
type GetRepoResponse struct {
	AssignableUsers  []RepoAssignableUser `json:"assignableUsers"`
	DefaultBranchRef RepoDefaultBranchRef `json:"defaultBranchRef"`
	Owner            RepoOwner            `json:"owner"`
	Name             string               `json:"name"`
}

// GetRepo function to get the detail of a repository
func (r *Repo) Get(
	options GetRepoOptions,
) GetRepoResponse {
	args := []string{
		"repo",
		"view",
		r.RepoName,
	}
	if options.Json == "" {
		args = append(args, "--json")
		args = append(args, "assignableUsers,defaultBranchRef,owner,name")
	}
	// Run the GitHub CLI command and capture the output
	cmd := exec.Command("gh", args...)
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to execute gh command: %v", err)
	}

	// Parse the JSON output into a slice of repo detail structs
	var repo GetRepoResponse
	err = json.Unmarshal(output, &repo)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	return repo
}
