package cli_prompt

import (
	"fmt"
	"log"
	"sort"
	"strings"

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
	title           string
	body            string
	reviewers       []string
}

// initializeBaseInfo method to initialize the base information for creating a pull request
func (p *CreatePullRequest) initializeBaseInfo() {
	r := gh_command.Repo{RepoName: ""}
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

// initializePullRequestTitleAndBody method to initialize the pull request title and body
func (p *CreatePullRequest) initializePullRequestTitleAndBody() {
	commits := gh_command.GetBranchCommits(p.repoOwner, p.repoName, p.baseBranch, p.headBranch)

	title, body := getPrePopulatedTitleAndBody(commits)
	p.title = title
	p.body = body
}

func (p *CreatePullRequest) restForm() *huh.Form {
	restForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter the pull request title").
				Value(&p.title),

			huh.NewText().
				Title("Enter the pull request body").
				ShowLineNumbers(true).
				Value(&p.body).
				// If I pass 0 to WithCharLimit, it will not limit the number of characters.
				CharLimit(0).
				// Calculate the line count of current text and add 5 to it as the height of the text box.
				WithHeight(strings.Count(p.body, "\n")+5),

			huh.NewMultiSelect[string]().
				Title("Select reviewers").
				Options((func() []huh.Option[string] {
					myUserLogin := gh_command.GetMyUserLogin()
					users := make([]huh.Option[string], 0)
					userLoginMap := make(map[string]string)

					// Step 1: Map each login to the corresponding name
					for _, user := range p.assignableUsers {
						if user.Login != myUserLogin {
							userLoginMap[user.Login] = user.Name
						}
					}

					// Step 2: Collect and sort the logins
					sortedLogins := make([]string, 0, len(userLoginMap))
					for login := range userLoginMap {
						sortedLogins = append(sortedLogins, login)
					}
					sort.Strings(sortedLogins)

					// Step 3: Construct `users` in alphabetical order
					for _, login := range sortedLogins {
						name := userLoginMap[login]

						format := "%s"
						if name != "" {
							format += " (%s)"
						}
						spread := []interface{}{login}
						if name != "" {
							spread = append(spread, name)
						}

						users = append(users, huh.NewOption(fmt.Sprintf(format, spread...), login))
					}

					return users
				})()...).
				Value(&p.reviewers),
		),
	)

	return restForm
}

// Run method to run the create pull request prompt
func (p *CreatePullRequest) Run() {
	initializeBaseInfo := p.initializeBaseInfo
	spinner.New().
		Title("Loading base information to create a pull request...").
		Action(initializeBaseInfo).
		Run()

	branchForm := p.branchForm()
	errBranchForm := branchForm.Run()
	if errBranchForm != nil {
		log.Fatal(errBranchForm)
	}

	// If the user stops the program, we don't want to go to the next form
	if branchForm.State == huh.StateAborted {
		fmt.Println("Aborted")
		return
	}

	initializePullRequestTitleAndBody := p.initializePullRequestTitleAndBody

	spinner.New().
		Title("Loading").
		Action(initializePullRequestTitleAndBody).
		Run()

	restForm := p.restForm()
	errRestForm := restForm.Run()
	if errRestForm != nil {
		log.Fatal(errRestForm)
	}

	fmt.Println("===REPO OWNER===")
	fmt.Println(p.repoOwner)
	fmt.Println("===REPO NAME===")
	fmt.Println(p.repoName)

	fmt.Println("===HEAD BRANCH===")
	fmt.Println(p.headBranch)
	fmt.Println("===BASE BRANCH===")
	fmt.Println(p.baseBranch)

	fmt.Println("===TITLE===")
	fmt.Println(p.title)
	fmt.Println("===BODY===")
	fmt.Println(p.body)
	fmt.Println("===REVIEWERS===")
	for _, reviewer := range p.reviewers {
		fmt.Println(reviewer)
	}
}
