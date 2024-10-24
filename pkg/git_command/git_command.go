package git_command

import (
	"encoding/json"
	"log"
	"os/exec"
	"strings"
)

type ListLatestBranchesResponse struct {
	Ref    string `json:"ref"`
	Commit string `json:"commit"`
	Date   string `json:"date"`
}

func ListLatestBranches() []ListLatestBranchesResponse {
	cmd := exec.Command(
		"git",
		"for-each-ref",
		"refs/heads/",
		"--sort=-committerdate",
		"--format={\"ref\": \"%(refname:short)\", \"commit\": \"%(objectname)\", \"date\": \"%(authordate:iso8601)\"}",
	)
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to execute git command: %v", err)
	}

	// Convert output to string and split into lines
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	// Join lines with "," and wrap with "[" and "]"
	jsonArray := "[" + strings.Join(lines, ",") + "]"

	// Parse the JSON output into a slice of Branch structs
	var bs []ListLatestBranchesResponse
	err = json.Unmarshal([]byte(jsonArray), &bs)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	return bs
}
