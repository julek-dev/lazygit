package git_commands

import (
	"errors"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/stretchr/testify/assert"
)

func TestGetCommitFilesFromRawDiff(t *testing.T) {
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
			input:    ":100644 100644 1234567 89abcde M\x00Myfile\x00",
			output: []*models.CommitFile{
				{
					Path:         "Myfile",
					ChangeStatus: "M",
				},
			},
		},
		{
			testName: "several files",
			input: ":100644 100644 1234567 89abcde M\x00Myfile\x00" +
				":000000 100644 0000000 89abcde A\x00MyOtherFile\x00" +
				":100644 000000 1234567 0000000 D\x00YetAnother\x00",
			output: []*models.CommitFile{
				{
					Path:         "Myfile",
					ChangeStatus: "M",
				},
				{
					Path:         "MyOtherFile",
					ChangeStatus: "A",
				},
				{
					Path:         "YetAnother",
					ChangeStatus: "D",
				},
			},
		},
		{
			testName: "a rename and a copy among regular files",
			input: ":100644 100644 1234567 89abcde M\x00Myfile\x00" +
				":100644 100644 1234567 1234567 R100\x00before\x00after\x00" +
				":100644 100644 1234567 89abcde C075\x00original\x00copy\x00" +
				":000000 100644 0000000 89abcde A\x00Added\x00",
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
					Path:         "copy",
					PreviousPath: "original",
					ChangeStatus: "C",
				},
				{
					Path:         "Added",
					ChangeStatus: "A",
				},
			},
		},
		{
			testName: "a path that starts with a colon",
			input:    ":100644 100644 1234567 89abcde M\x00:file\x00",
			output: []*models.CommitFile{
				{
					Path:         ":file",
					ChangeStatus: "M",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			result := getCommitFilesFromRawDiff(test.input)
			assert.Equal(t, test.output, result)
		})
	}
}

func TestGetFilesInDiff(t *testing.T) {
	diffArgs := func(withNumstat bool, reverse bool) []string {
		args := []string{"-c", "diff.dstPrefix=b/", "-c", "diff.srcPrefix=a/", "-c", "diff.noprefix=false", "diff", "--submodule", "--no-ext-diff", "--raw"}
		if withNumstat {
			args = append(args, "--numstat")
		}
		args = append(args, "-z", "--find-renames=50%")
		if reverse {
			args = append(args, "-R")
		}
		return append(args, "from", "to")
	}
	raw := ":100644 100644 1234567 89abcde M\x00file\x00" +
		":100644 100644 1234567 89abcde R090\x00before\x00after\x00" +
		":100644 100644 1234567 89abcde M\x00image.png\x00" +
		":160000 160000 1234567 89abcde M\x00sub\x00"
	numstat := "4\t1\tfile\x002\t0\t\x00before\x00after\x00-\t-\timage.png\x001\t1\tsub\x00"
	filesWithoutCounts := []*models.CommitFile{
		{Path: "file", ChangeStatus: "M"},
		{Path: "after", PreviousPath: "before", ChangeStatus: "R"},
		{Path: "image.png", ChangeStatus: "M"},
		{Path: "sub", ChangeStatus: "M"},
	}

	scenarios := []struct {
		testName    string
		withNumstat bool
		reverse     bool
		runner      *oscommands.FakeCmdObjRunner
		expected    []*models.CommitFile
	}{
		{
			testName:    "without numstat",
			withNumstat: false,
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs(diffArgs(false, false), raw, nil),
			expected: filesWithoutCounts,
		},
		{
			testName:    "with numstat",
			withNumstat: true,
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs(diffArgs(true, false), raw+numstat, nil),
			expected: []*models.CommitFile{
				{Path: "file", ChangeStatus: "M", LinesAdded: 4, LinesDeleted: 1},
				{Path: "after", PreviousPath: "before", ChangeStatus: "R", LinesAdded: 2},
				{Path: "image.png", ChangeStatus: "M"},
				{Path: "sub", ChangeStatus: "M"},
			},
		},
		{
			testName:    "with numstat failing",
			withNumstat: true,
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs(diffArgs(true, false), "", errors.New("fatal: could not fetch blob from promisor remote")).
				ExpectGitArgs(diffArgs(false, false), raw, nil),
			expected: filesWithoutCounts,
		},
		{
			testName:    "reversed",
			withNumstat: true,
			reverse:     true,
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs(diffArgs(true, true),
					":100644 100644 89abcde 1234567 R090\x00after\x00before\x000\t2\t\x00after\x00before\x00",
					nil,
				),
			expected: []*models.CommitFile{
				{Path: "before", PreviousPath: "after", ChangeStatus: "R", LinesDeleted: 2},
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			loader := NewCommitFileLoader(common.NewDummyCommon(), oscommands.NewDummyCmdObjBuilder(s.runner))

			files, err := loader.GetFilesInDiff("from", "to", s.reverse, s.withNumstat)
			assert.NoError(t, err)
			assert.Equal(t, s.expected, files)
			s.runner.CheckForMissingCalls()
		})
	}
}
