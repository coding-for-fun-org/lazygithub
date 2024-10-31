package main

import "github.com/coding-for-fun-org/lazygithub/pkg/cli_prompt"

func main() {
	c := cli_prompt.CreatePullRequest{}

	c.Run()
}
