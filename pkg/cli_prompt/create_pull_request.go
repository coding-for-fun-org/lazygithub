package cli_prompt

import (
	"fmt"
	"log"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/coding-for-fun-org/lazygithub/pkg/gh_command"
	"github.com/coding-for-fun-org/lazygithub/pkg/git_command"
)

type CreatePullRequest struct {
	repoOwner       string
	repoName        string
	assignableUsers []gh_command.RepoAssignableUser
	defaultBranch   string
	latestBranches  []git_command.ListLatestBranchesResponse
	baseBranch      string
	headBranch      string
}

// initializeBaseInfo method to initialize the base information for creating a pull request
func (p *CreatePullRequest) initializeBaseInfo() {
	r := gh_command.Repo{RepoName: "RockRabbit-ai/rockrabbit-web"}
	repo := r.Get(gh_command.GetRepoOptions{})

	p.repoOwner = repo.Owner.Login
	p.repoName = repo.Name
	p.assignableUsers = repo.AssignableUsers
	p.defaultBranch = repo.DefaultBranchRef.Name
	p.latestBranches = git_command.ListLatestBranches()
}

// branchForm method to create a form for selecting the base and head branches
func (p *CreatePullRequest) branchForm() *huh.Form {
	branchForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select the head branch").
				Options((func() []huh.Option[string] {
					// filter out the default branch
					branches := make([]huh.Option[string], 0)
					for _, branch := range p.latestBranches {
						if branch.Ref != p.defaultBranch {
							branches = append(
								branches,
								huh.NewOption(branch.Ref, branch.Ref),
							)
						}
					}

					return branches
				})()...).
				Value(&p.headBranch),

			// I'd like to put this in another group but
			// because of the current bug, I can not do it.
			// https://github.com/charmbracelet/huh/issues/419
			huh.NewSelect[string]().
				Title("Select the base branch").
				OptionsFunc(func() []huh.Option[string] {
					// filter out the default branch
					branches := make([]huh.Option[string], 0)
					branches = append(
						branches,
						huh.NewOption(p.defaultBranch, p.defaultBranch),
					)
					for _, branch := range p.latestBranches {
						if branch.Ref != p.defaultBranch && branch.Ref != p.headBranch {
							branches = append(
								branches,
								huh.NewOption(branch.Ref, branch.Ref),
							)
						}
					}

					return branches
				}, &p.headBranch).
				Value(&p.baseBranch),
		),
	)

	return branchForm
}

func (p *CreatePullRequest) Run() {
	initializeBaseInfo := p.initializeBaseInfo
	spinner.New().
		Title("Loading base information to create a pull request...").
		Action(initializeBaseInfo).
		Run()

	branchForm := p.branchForm()
	err := branchForm.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Head branch: %s\n", p.headBranch)
	fmt.Printf("Base branch: %s\n", p.baseBranch)
}
