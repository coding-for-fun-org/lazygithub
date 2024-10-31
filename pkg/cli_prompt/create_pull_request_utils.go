package cli_prompt

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
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

func getReviewersFilePath() (string, error) {
	fileName := ".__reviewers.csv"
	home, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Failed to get user home directory: %s", err)

		return "", err
	}

	return filepath.Join(home, fileName), nil
}

func writeLatestReviewers(repo string, reviewers []string) error {
	filePath, err := getReviewersFilePath()
	if err != nil {
		return err
	}

	// Open the CSV file
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	// Initialize the CSV reader
	reader := csv.NewReader(file)

	// Read all records from the file
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read CSV file: %s", err)
	}

	// Track if repo was found
	repoFound := false

	// first column is repo name, second column is reviewers
	// if repo name already exists, update the reviewers
	// if repo name does not exist, append the repo name and reviewers
	for i, record := range records {
		if record[0] == repo {
			records[i][1] = strings.Join(reviewers, ",")
			repoFound = true
			break
		}
	}

	// If repo was not found, append a new record
	if !repoFound {
		newRecord := []string{repo, strings.Join(reviewers, ",")}
		records = append(records, newRecord)
	}

	// Move the file pointer back to the beginning to overwrite the file
	if _, err := file.Seek(0, 0); err != nil {
		log.Fatalf("Failed to seek file: %s", err)
		return err
	}

	// write the updated records to the file
	writer := csv.NewWriter(file)
	err = writer.WriteAll(records)
	if err != nil {
		log.Fatalf("Failed to write CSV file: %s", err)
		return err
	}

	writer.Flush()

	return writer.Error()
}

func getLatestReviewers(repo string) ([]string, error) {
	filePath, err := getReviewersFilePath()
	if err != nil {
		return nil, err
	}

	// Open the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Failed to open file: %s", err)
		return nil, err
	}
	defer file.Close()

	// Initialize the CSV reader
	reader := csv.NewReader(file)

	// Read all records from the file
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("Failed to read CSV file: %s", err)
		return nil, err
	}

	reviewers := []string{}
	// first column is repo name, second column is reviewers

	for _, record := range records {
		if record[0] == repo {
			reviewers = strings.Split(record[1], ",")
			break
		}
	}

	return reviewers, nil
}
