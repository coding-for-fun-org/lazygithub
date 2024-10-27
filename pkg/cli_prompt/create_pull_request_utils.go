package cli_prompt

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/coding-for-fun-org/lazygithub/pkg/gh_command"
)

// extractPatterns function extracts all occurrences of the pattern [A-Z]{1,}-\d{1,} from the input string.
func extractPatterns(input string) ([]string, error) {
	// Define the regular expression
	pattern := `[A-Z]+-\d+`

	// Compile the regular expression
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	// Find all matches
	matches := re.FindAllString(input, -1)
	return matches, nil
}

// concatenateAndRemoveDuplicates function to concatenate two slices and remove duplicates
func concatenateAndRemoveDuplicates(slice1, slice2 []string) []string {
	// Use a map to track unique items
	uniqueMap := make(map[string]struct{})
	result := make([]string, 0)

	// Helper function to add items to the result while checking uniqueness
	addUnique := func(item string) {
		if _, exists := uniqueMap[item]; !exists {
			uniqueMap[item] = struct{}{}
			result = append(result, item)
		}
	}

	// Add elements from the first slice
	for _, item := range slice1 {
		addUnique(item)
	}

	// Add elements from the second slice
	for _, item := range slice2 {
		addUnique(item)
	}

	return result
}

// splitCommitSummaryAndDescription function to split the commit message into summary and description
func splitCommitSummaryAndDescription(commitMessage string) (string, string) {
	parts := strings.SplitN(commitMessage, "\n\n", 2)

	if len(parts) == 2 {
		return parts[0], parts[1]
	}

	return commitMessage, ""
}

func getPrePopulatedTitleAndBody(
	commits []gh_command.Commit,
) (string, string) {
	if len(commits) == 1 {
		commitTitle, commitBody := splitCommitSummaryAndDescription(commits[0].Message)
		issueNumbersFromTitle, _ := extractPatterns(commitTitle)
		issueNumbersFromBody, _ := extractPatterns(commitBody)
		issueNumbers := concatenateAndRemoveDuplicates(
			issueNumbersFromTitle,
			issueNumbersFromBody,
		)
		commitBody = commitBody + "\n\n### Jira Link\n\n"
		for _, issueNumber := range issueNumbers {
			commitBody = commitBody + fmt.Sprintf(
				"[%s](https://keends.atlassian.net/browse/%s)\n",
				issueNumber,
				issueNumber,
			)
		}
		return commitTitle, commitBody
	}

	linkBody := "### Jira Link\n\n"
	commitFullBody := ""
	for _, commit := range commits {
		commitTitle, commitBody := splitCommitSummaryAndDescription(commit.Message)
		issueNumbersFromTitle, _ := extractPatterns(commitTitle)
		issueNumbersFromBody, _ := extractPatterns(commitBody)
		issueNumbers := concatenateAndRemoveDuplicates(
			issueNumbersFromTitle,
			issueNumbersFromBody,
		)
		for _, issueNumber := range issueNumbers {
			linkBody = linkBody + fmt.Sprintf(
				"[%s](https://keends.atlassian.net/browse/%s)\n",
				issueNumber,
				issueNumber,
			)
		}

		commitFullBody = commitFullBody + commitTitle + "\n\n" + commitBody + "\n---\n"
	}

	commitFullBody = commitFullBody + linkBody

	return "", commitFullBody
}
