package cli_prompt

import (
	"reflect"
	"testing"

	"github.com/coding-for-fun-org/lazygithub/pkg/gh_command"
	"github.com/coding-for-fun-org/lazygithub/pkg/git_command"
)

func Test_extractPatterns(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "base test case",
			args: args{
				input: "ABC-123",
			},
			want:    []string{"ABC-123"},
			wantErr: false,
		},
		{
			name: "multiple patterns",
			args: args{
				input: "fix(ABC-456,DEF-146): this is test commit.\n\nThis is related to XYZ-7890",
			},
			want:    []string{"ABC-456", "DEF-146", "XYZ-7890"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractPatterns(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractPatterns() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractPatterns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_concatenateAndRemoveDuplicates(t *testing.T) {
	type args struct {
		slice1 []string
		slice2 []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "just concatenate if one of the slices is empty",
			args: args{
				slice1: []string{},
				slice2: []string{"b", "c", "d"},
			},
			want: []string{"b", "c", "d"},
		},
		{
			name: "concatenate and remove duplicates",
			args: args{
				slice1: []string{"a", "b", "c"},
				slice2: []string{"b", "c", "d"},
			},
			want: []string{"a", "b", "c", "d"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := concatenateAndRemoveDuplicates(tt.args.slice1, tt.args.slice2); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("concatenateAndRemoveDuplicates() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_splitCommitSummaryAndDescription(t *testing.T) {
	type args struct {
		commitMessage string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{
			name: "base test case",
			args: args{
				commitMessage: "fix(ABC-456,DEF-146): this is test commit.\n\nThis is related to XYZ-7890",
			},
			want:  "fix(ABC-456,DEF-146): this is test commit.",
			want1: "This is related to XYZ-7890",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := splitCommitSummaryAndDescription(tt.args.commitMessage)
			if got != tt.want {
				t.Errorf("splitCommitSummaryAndDescription() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("splitCommitSummaryAndDescription() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestCreatePullRequest_getPrePopulatedTitleAndBody(t *testing.T) {
	type fields struct {
		repoOwner       string
		repoName        string
		assignableUsers []gh_command.RepoAssignableUser
		defaultBranch   string
		latestBranches  []git_command.ListLatestBranchesResponse
		title           string
		body            string
		baseBranch      string
		headBranch      string
		reviewers       []string
	}
	type args struct {
		commits []gh_command.Commit
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		want1  string
	}{
		{
			name: "base test case",
			fields: fields{
				title: "test",
				body:  "test",
			},
			args: args{
				commits: []gh_command.Commit{
					{
						Message: "feat(ABC-123): this is test commit.\n\nThis is related to XYZ-7890",
					},
				},
			},
			want:  "feat(ABC-123): this is test commit.",
			want1: "This is related to XYZ-7890\n\n### Jira Link\n\n[ABC-123](https://keends.atlassian.net/browse/ABC-123)\n[XYZ-7890](https://keends.atlassian.net/browse/XYZ-7890)\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getPrePopulatedTitleAndBody(tt.args.commits)
			if got != tt.want {
				t.Errorf(
					"CreatePullRequest.getPrePopulatedTitleAndBody() got = %v, want %v",
					got,
					tt.want,
				)
			}
			if got1 != tt.want1 {
				t.Errorf(
					"CreatePullRequest.getPrePopulatedTitleAndBody() got1 = %v, want %v",
					got1,
					tt.want1,
				)
			}
		})
	}
}
