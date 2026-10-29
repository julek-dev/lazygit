package git_commands

import (
	"errors"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/stretchr/testify/assert"
)

func TestGetCommitFilesFromFilenames(t *testing.T) {
	tests := []struct {
		testName string
		input    string
		output   []*models.CommitFile
	}{
		{
			testName: "no files",
			input:    "",
			output:   []*models.CommitFile{},
		},
		{
			testName: "one file",
			input:    "MM\x00Myfile\x00",
			output: []*models.CommitFile{
				{
					Path:         "Myfile",
					ChangeStatus: "MM",
				},
			},
		},
		{
			testName: "two files",
			input:    "MM\x00Myfile\x00M \x00MyOtherFile\x00",
			output: []*models.CommitFile{
				{
					Path:         "Myfile",
					ChangeStatus: "MM",
				},
				{
					Path:         "MyOtherFile",
					ChangeStatus: "M ",
				},
			},
		},
		{
			testName: "three files",
			input:    "MM\x00Myfile\x00M \x00MyOtherFile\x00 M\x00YetAnother\x00",
			output: []*models.CommitFile{
				{
					Path:         "Myfile",
					ChangeStatus: "MM",
				},
				{
					Path:         "MyOtherFile",
					ChangeStatus: "M ",
				},
				{
					Path:         "YetAnother",
					ChangeStatus: " M",
				},
			},
		},
		{
			testName: "a rename among regular files",
			input:    "M\x00Myfile\x00R100\x00before\x00after\x00A\x00Added\x00",
			output: []*models.CommitFile{
				{
					Path:         "Myfile",
					ChangeStatus: "M",
				},
				{
					Path:         "after",
					PreviousPath: "before",
					ChangeStatus: "R",
				},
				{
					Path:         "Added",
					ChangeStatus: "A",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			result := getCommitFilesFromFilenames(test.input)
			assert.Equal(t, test.output, result)
		})
	}
}

func TestGetFilesInDiff(t *testing.T) {
	nameStatusArgs := []string{"-c", "diff.dstPrefix=b/", "-c", "diff.srcPrefix=a/", "-c", "diff.noprefix=false", "diff", "--submodule", "--no-ext-diff", "--name-status", "-z", "--find-renames=50%", "from", "to"}
	numstatArgs := []string{"-c", "diff.dstPrefix=b/", "-c", "diff.srcPrefix=a/", "-c", "diff.noprefix=false", "diff", "--submodule", "--no-ext-diff", "--numstat", "-z", "--find-renames=50%", "from", "to"}
	nameStatus := "M\x00file\x00R090\x00before\x00after\x00M\x00image.png\x00"

	scenarios := []struct {
		testName    string
		withNumstat bool
		runner      *oscommands.FakeCmdObjRunner
		expected    []*models.CommitFile
	}{
		{
			testName:    "without numstat",
			withNumstat: false,
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs(nameStatusArgs, nameStatus, nil),
			expected: []*models.CommitFile{
				{Path: "file", ChangeStatus: "M"},
				{Path: "after", PreviousPath: "before", ChangeStatus: "R"},
				{Path: "image.png", ChangeStatus: "M"},
			},
		},
		{
			testName:    "with numstat",
			withNumstat: true,
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs(nameStatusArgs, nameStatus, nil).
				ExpectGitArgs(numstatArgs, "4\t1\tfile\x002\t0\t\x00before\x00after\x00-\t-\timage.png\x00", nil),
			expected: []*models.CommitFile{
				{Path: "file", ChangeStatus: "M", LinesAdded: 4, LinesDeleted: 1},
				{Path: "after", PreviousPath: "before", ChangeStatus: "R", LinesAdded: 2},
				{Path: "image.png", ChangeStatus: "M"},
			},
		},
		{
			testName:    "with numstat failing",
			withNumstat: true,
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs(nameStatusArgs, nameStatus, nil).
				ExpectGitArgs(numstatArgs, "", errors.New("fatal: could not fetch blob from promisor remote")),
			expected: []*models.CommitFile{
				{Path: "file", ChangeStatus: "M"},
				{Path: "after", PreviousPath: "before", ChangeStatus: "R"},
				{Path: "image.png", ChangeStatus: "M"},
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			loader := NewCommitFileLoader(common.NewDummyCommon(), oscommands.NewDummyCmdObjBuilder(s.runner))

			files, err := loader.GetFilesInDiff("from", "to", false, s.withNumstat)
			assert.NoError(t, err)
			assert.Equal(t, s.expected, files)
			s.runner.CheckForMissingCalls()
		})
	}
}
