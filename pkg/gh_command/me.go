package gh_command

import (
	"log"
	"os/exec"
	"strings"
)

func GetMyUserLogin() string {
	cmd := exec.Command(
		"gh",
		"api",
		"user",
		"--jq",
		".login",
	)
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to execute gh command: %v", err)
	}

	return strings.TrimSpace(string(output))
}
